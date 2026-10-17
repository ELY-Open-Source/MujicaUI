package files

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ShareMember is someone the item is shared with.
type ShareMember struct {
	ID, Name, Email string
	Permission      Permission
	// Owner cannot be changed or removed.
	Owner bool
}

// ShareInvite is a new person to share with.
type ShareInvite struct {
	Email      string
	Permission Permission
}

// ShareDialogOptions configure a ShareDialog.
type ShareDialogOptions struct {
	// Title names what is shared.
	Title string
	// Link is the item's address; empty hides the link row.
	Link     string
	Disabled bool
	// OnCopy handles copying the link; nil uses the clipboard.
	OnCopy func(link string)
}

// ShareDialogResult is what happened in the dialog this frame.
type ShareDialogResult struct {
	Element ui.Element
	// Invited requests membership; Changed and Removed identify applied changes.
	Invited          *ShareInvite
	Changed, Removed string
	// Copied reports the link was written to the clipboard.
	Copied bool
}

type shareDraft struct {
	email string
	perm  Permission
	rows  ui.ListState
}

// ShareDialog manages invitations, permissions, members and copying links.
func ShareDialog(c *ui.Context, open *bool, members *[]ShareMember, opts ShareDialogOptions) ShareDialogResult {
	if open == nil || members == nil || opts.Title == "" {
		panic("mujica: ShareDialog needs open, members and a Title")
	}
	var res ShareDialogResult
	site := ui.Box(c).Disabled(opts.Disabled)
	if site.IsDisabled() {
		*open = false
	}
	draft := ui.Local(site, "share_dialog.draft", func() shareDraft { return shareDraft{} })
	u := core.Density(c).Unit()
	res.Element = overlay.Dialog(c, open, overlay.DialogOptions{Title: fmt.Sprintf(core.Msg(c, "share_dialog.title"), opts.Title), Width: 480}, func() {
		ui.Column(c).Gap(u * 3).Disabled(opts.Disabled).Children(func() {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Start).Children(func() {
				in := input.TextInput(c, &draft.email, input.TextInputOptions{Placeholder: core.Msg(c, "share_dialog.placeholder"), Label: core.Msg(c, "share_dialog.invite")})
				in.Element.Grow(1).MinWidth(0)
				PermissionSelect(c, &draft.perm, PermissionSelectOptions{Label: core.Msg(c, "share_dialog.invitePermission")})
				send := input.Button(c, core.Msg(c, "share_dialog.send"), input.ButtonOptions{Disabled: strings.TrimSpace(draft.email) == ""})
				if (send.Clicked() || in.Submitted()) && strings.TrimSpace(draft.email) != "" && !send.IsDisabled() {
					res.Invited = &ShareInvite{Email: strings.TrimSpace(draft.email), Permission: draft.perm}
					draft.email = ""
					slog.Debug("mujica: share_dialog invite")
				}
			})
			ui.Text(c, core.Msg(c, "share_dialog.people")).Font(theme.SerifFont).FontWeight(600)
			shareMembers(c, members, &res, &draft.rows)
			if opts.Link != "" {
				layout.Divider(c, layout.DividerOptions{})
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icons.Must("external-link")).Size(16, 16).Shrink(0).TextColor(core.Tokens(c).TextMuted)
					ui.Text(c, opts.Link).SingleLine().Grow(1).MinWidth(0).Selectable()
					b := input.Button(c, core.Msg(c, "share_dialog.copy"), input.ButtonOptions{Variant: input.Secondary, Icon: icons.Must("copy")})
					if b.Clicked() && !b.IsDisabled() {
						if opts.OnCopy != nil {
							opts.OnCopy(opts.Link)
						} else {
							c.WriteClipboard(opts.Link)
						}
						res.Copied = true
						c.Announce(core.Msg(c, "share_dialog.copied"))
						slog.Debug("mujica: share_dialog copy link")
					}
				})
			}
		})
	})
	return res
}

// shareMembers virtualizes membership permissions and removal actions.
func shareMembers(c *ui.Context, members *[]ShareMember, res *ShareDialogResult, rows *ui.ListState) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	remove := -1
	rows.Key = func(i int) any { return (*members)[i].ID }
	_, h := c.Size()
	ui.List(c, rows, len(*members), func(i int) {
		m := &(*members)[i]
		if m.ID == "" {
			panic("mujica: ShareMember needs an ID")
		}
		ui.Row(c.Key(m.ID)).Gap(u * 2).AlignItems(ui.Center).MinHeight(core.Density(c).ControlHeight()).Children(func() {
			display.Avatar(c, m.Name, display.AvatarOptions{Size: 28})
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				ui.Text(c, m.Name).SingleLine()
				if m.Email != "" {
					ui.Text(c, m.Email).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
				}
			})
			if m.Owner {
				ui.Text(c, core.Msg(c, "share_dialog.owner")).TextColor(k.TextMuted)
				return
			}
			if PermissionSelect(c, &m.Permission, PermissionSelectOptions{Label: m.Name}).Changed() {
				res.Changed = m.ID
			}
			if b := core.IconAction(c, icons.Must("x"), fmt.Sprintf(core.Msg(c, "share_dialog.remove"), m.Name)); b.Clicked() && !b.IsDisabled() {
				remove = i
			}
		})
	}).Height(min(float32(len(*members))*(core.Density(c).ControlHeight()+u), max(28, h*0.4))).FillWidth()
	if remove >= 0 {
		res.Removed = (*members)[remove].ID
		*members = append((*members)[:remove:remove], (*members)[remove+1:]...)
		slog.Debug("mujica: share_dialog remove")
	}
}
