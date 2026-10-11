package account

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// Account is a signed-in account a user can switch to.
type Account struct {
	ID, Name, Email string
	Avatar          *ui.Bitmap
}

// AccountSwitcherOptions configure an AccountSwitcher.
type AccountSwitcherOptions struct {
	Disabled bool
}

// AccountSwitcherResult reports account changes and addition requests.
type AccountSwitcherResult struct {
	Element *ui.Element
	changed bool
	add     bool
}

// Changed reports a switch to another account.
func (r AccountSwitcherResult) Changed() bool { return r.changed }

// Add reports a request to sign in another account.
func (r AccountSwitcherResult) Add() bool { return r.add }

// AccountSwitcher offers account switching and an addition action.
func AccountSwitcher(c *ui.Context, current *string, accounts []Account, opts AccountSwitcherOptions) AccountSwitcherResult {
	var cur *Account
	for i := range accounts {
		if accounts[i].ID == "" || accounts[i].Name == "" {
			panic("mujica: Account needs an ID and a Name")
		}
		if accounts[i].ID == *current {
			cur = &accounts[i]
		}
	}
	if cur == nil {
		panic(fmt.Sprintf("mujica: AccountSwitcher current %q is not an account", *current))
	}
	var r AccountSwitcherResult
	b := accountTrigger(c, cur.Name, cur.Email, cur.Avatar, opts.Disabled)
	r.Element = b
	open := ui.Local(b, "open", func() bool { return false })
	if b.IsDisabled() {
		*open = false
	} else if b.Clicked() || b.Shortcut(0, ui.KeyDown) {
		*open = true
	}
	if len(accounts) > 8 {
		accountSwitcherLarge(c, b, open, current, accounts, &r)
		return r
	}
	overlay.PopupMenuOpen(c, b, open, func(m *overlay.PopupMenu) {
		m.Group(core.Msg(c, "account_switcher.accounts"), func() {
			for _, a := range accounts {
				on := a.ID == *current
				label := a.Name
				if a.Email != "" {
					label += " · " + a.Email
				}
				ui.Column(c).Key(a.ID).Shrink(0).Children(func() {
					if m.Check(label, &on, overlay.PopupMenuItemOptions{}) && a.ID != *current && !b.IsDisabled() {
						*current, r.changed = a.ID, true
						slog.Debug("mujica: account_switcher switched")
					}
				})
			}
		})
		m.Separator()
		if m.Item(core.Msg(c, "account_switcher.add"), overlay.PopupMenuItemOptions{}) && !b.IsDisabled() {
			r.add = true
			slog.Debug("mujica: account_switcher add")
		}
	})
	return r
}

func accountSwitcherLarge(c *ui.Context, b *ui.Element, open *bool, current *string, accounts []Account, r *AccountSwitcherResult) {
	overlay.Popover(c, b, open, overlay.PopoverOptions{Label: core.Msg(c, "account_switcher.accounts"), Width: 320}, func() {
		cursor := ui.Local(b, "cursor", func() int { return -1 })
		rows := ui.Local(b, "list", func() ui.ListState { return ui.ListState{} })
		rows.Key = func(i int) any { return accounts[i].ID }
		rows.Label = func(i int) string { return accounts[i].Name }
		rows.Selected = cursor
		restore := data.DataSelectTheme(c, 0)
		list := ui.List(c, rows, len(accounts), func(i int) {
			a := accounts[i]
			row := workspaceRow(c, Workspace{ID: a.ID, Name: a.Name, Detail: a.Email}, a.ID == *current, i == *cursor)
			if row.Focused() && !b.IsDisabled() {
				*cursor = i
			}
			if row.Clicked() && !row.IsDisabled() && !b.IsDisabled() {
				r.changed, *current, *open = a.ID != *current, a.ID, false
			}
		}).Height(280).MinHeight(0).AutoFocus()
		restore()
		internal.Ring(list, core.Tokens(c).Focus, theme.ControlRadius)
		if list.Changed() && !b.IsDisabled() {
			list.Focus()
		}
		if list.Submitted() && !b.IsDisabled() && *cursor >= 0 && *cursor < len(accounts) {
			id := accounts[*cursor].ID
			r.changed, *current, *open = id != *current, id, false
		}
		if input.Button(c, core.Msg(c, "account_switcher.add"), input.ButtonOptions{Variant: input.Ghost}).Clicked() && !b.IsDisabled() {
			r.add, *open = true, false
		}
	})
	if r.changed || r.add {
		slog.Debug("mujica: account_switcher action", "changed", r.changed, "add", r.add)
	}
}

// accountTrigger displays an avatar, name and detail.
func accountTrigger(c *ui.Context, name, detail string, avatar *ui.Bitmap, disabled bool) *ui.Element {
	k := core.Tokens(c)
	d := core.Density(c)
	b := ui.ButtonBase(c).Role(ui.RoleMenuButton).MinHeight(d.ControlHeight()+d.Unit()*2).Padding(d.Unit(), d.Unit()*3, d.Unit(), d.Unit()*2).
		Gap(d.Unit() * 2).Radius(theme.ControlRadius).Justify(ui.Start).Disabled(disabled).Label(name).Transition(core.ColorFade)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case b.Pressed():
		b.Background(k.Selection)
	case b.Hovered():
		b.Background(k.SurfaceHover)
	}
	b.Children(func() {
		display.Avatar(c, name, display.AvatarOptions{Image: avatar, Size: 28})
		ui.Column(c).Grow(1).MinWidth(0).Children(func() {
			ui.Text(c, name).FontWeight(600).SingleLine().Ellipsis("…")
			if detail != "" {
				ui.Text(c, detail).TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine().Ellipsis("…")
			}
		})
		ui.Icon(c, icons.Must("chevron-down")).Size(14, 14).TextColor(k.TextMuted)
	})
	return b
}
