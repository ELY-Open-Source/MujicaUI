package mujica

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
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
	Element *ui.Element
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
	u := Density(c).Unit()
	res.Element = Dialog(c, open, DialogOptions{Title: fmt.Sprintf(msg(c, "share_dialog.title"), opts.Title), Width: 480}, func() {
		ui.Column(c).Gap(u * 3).Disabled(opts.Disabled).Children(func() {
			ui.Row(c).Gap(u * 2).AlignItems(ui.Start).Children(func() {
				in := TextInput(c, &draft.email, TextInputOptions{Placeholder: msg(c, "share_dialog.placeholder"), Label: msg(c, "share_dialog.invite")})
				in.Element.Grow(1).MinWidth(0)
				PermissionSelect(c, &draft.perm, PermissionSelectOptions{Label: msg(c, "share_dialog.invitePermission")})
				send := Button(c, msg(c, "share_dialog.send"), ButtonOptions{Disabled: strings.TrimSpace(draft.email) == ""})
				if (send.Clicked() || in.Submitted()) && strings.TrimSpace(draft.email) != "" && !send.IsDisabled() {
					res.Invited = &ShareInvite{Email: strings.TrimSpace(draft.email), Permission: draft.perm}
					draft.email = ""
					slog.Debug("mujica: share_dialog invite")
				}
			})
			ui.Text(c, msg(c, "share_dialog.people")).Font(theme.SerifFont).FontWeight(600)
			shareMembers(c, members, &res, &draft.rows)
			if opts.Link != "" {
				Divider(c, DividerOptions{})
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icons.Must("external-link")).Size(16, 16).Shrink(0).TextColor(Tokens(c).TextMuted)
					ui.Text(c, opts.Link).SingleLine().Grow(1).MinWidth(0).Selectable()
					b := Button(c, msg(c, "share_dialog.copy"), ButtonOptions{Variant: Secondary, Icon: icons.Must("copy")})
					if b.Clicked() && !b.IsDisabled() {
						if opts.OnCopy != nil {
							opts.OnCopy(opts.Link)
						} else {
							c.WriteClipboard(opts.Link)
						}
						res.Copied = true
						c.Announce(msg(c, "share_dialog.copied"))
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
	k := Tokens(c)
	u := Density(c).Unit()
	remove := -1
	rows.Key = func(i int) any { return (*members)[i].ID }
	_, h := c.Size()
	ui.List(c, rows, len(*members), func(i int) {
		m := &(*members)[i]
		if m.ID == "" {
			panic("mujica: ShareMember needs an ID")
		}
		ui.Row(c).Key(m.ID).Gap(u * 2).AlignItems(ui.Center).MinHeight(Density(c).ControlHeight()).Children(func() {
			Avatar(c, m.Name, AvatarOptions{Size: 28})
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				ui.Text(c, m.Name).SingleLine()
				if m.Email != "" {
					ui.Text(c, m.Email).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
				}
			})
			if m.Owner {
				ui.Text(c, msg(c, "share_dialog.owner")).TextColor(k.TextMuted)
				return
			}
			if PermissionSelect(c, &m.Permission, PermissionSelectOptions{Label: m.Name}).Changed() {
				res.Changed = m.ID
			}
			if b := iconAction(c, icons.Must("x"), fmt.Sprintf(msg(c, "share_dialog.remove"), m.Name)); b.Clicked() && !b.IsDisabled() {
				remove = i
			}
		})
	}).Height(min(float32(len(*members))*(Density(c).ControlHeight()+u), max(28, h*0.4))).FillWidth()
	if remove >= 0 {
		res.Removed = (*members)[remove].ID
		*members = append((*members)[:remove:remove], (*members)[remove+1:]...)
		slog.Debug("mujica: share_dialog remove")
	}
}
