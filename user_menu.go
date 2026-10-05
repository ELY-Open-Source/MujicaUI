package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// UserStatus is the presence a user shows to others.
type UserStatus int

const (
	UserStatusOnline UserStatus = iota
	UserStatusAway
	UserStatusBusy
	UserStatusInvisible
)

// UserMenuOptions configure a UserMenu. Name is required.
type UserMenuOptions struct {
	Name, Email string
	Avatar      *ui.Bitmap
	Disabled    bool
}

// UserMenuResult reports presence, settings and sign-out events.
type UserMenuResult struct {
	Element                  *ui.Element
	changed, settings, leave bool
}

// StatusChanged reports a new presence in *status.
func (r UserMenuResult) StatusChanged() bool { return r.changed }

// Settings reports a request to open settings.
func (r UserMenuResult) Settings() bool { return r.settings }

// SignOut reports a request to sign out.
func (r UserMenuResult) SignOut() bool { return r.leave }

// UserMenu opens presence, settings and sign-out actions.
func UserMenu(c *ui.Context, status *UserStatus, opts UserMenuOptions) UserMenuResult {
	if opts.Name == "" {
		panic("mujica: UserMenu needs a Name")
	}
	if *status < UserStatusOnline || *status > UserStatusInvisible {
		panic(fmt.Sprintf("mujica: UserMenu status %d is unknown", *status))
	}
	k := Tokens(c)
	var r UserMenuResult
	word := msg(c, fmt.Sprintf("user_menu.status%d", *status))
	b := ui.ButtonBase(c).Role(ui.RoleMenuButton).Size(36, 36).Radius(18).Disabled(opts.Disabled).
		Label(fmt.Sprintf(msg(c, "user_menu.label"), opts.Name, word))
	r.Element = b
	internal.Ring(b, k.Focus, 18)
	open := ui.Local(b, "open", func() bool { return false })
	disabled := b.IsDisabled()
	if disabled {
		*open = false
	} else if b.Clicked() || b.Shortcut(0, ui.KeyDown) {
		*open = true
	}
	st := *status
	b.Children(func() { Avatar(c, opts.Name, AvatarOptions{Image: opts.Avatar, Size: 32}) })
	b.DrawOver(func(p *ui.Painter, rect ui.Rect) { userMenuDot(p, k, rect.X+rect.W-6, rect.Y+rect.H-6, st) })
	popupMenuOpen(c, b, open, func(m *PopupMenu) {
		head := opts.Name
		if opts.Email != "" {
			head += " · " + opts.Email
		}
		m.Group(head, func() {
			for s := UserStatusOnline; s <= UserStatusInvisible; s++ {
				on := s == *status
				if m.Check(msg(c, fmt.Sprintf("user_menu.status%d", s)), &on, PopupMenuItemOptions{}) && s != *status && !disabled {
					*status, r.changed = s, true
					slog.Debug("mujica: user_menu status", "status", s)
				}
			}
		})
		m.Separator()
		r.settings = m.Item(msg(c, "user_menu.settings"), PopupMenuItemOptions{}) && !disabled
		r.leave = m.Item(msg(c, "user_menu.sign_out"), PopupMenuItemOptions{}) && !disabled
	})
	if r.settings || r.leave {
		slog.Debug("mujica: user_menu chose", "settings", r.settings, "sign_out", r.leave)
	}
	return r
}

// userMenuDot distinguishes presence using shapes and colors.
func userMenuDot(p *ui.Painter, k theme.Tokens, cx, cy float32, s UserStatus) {
	const d = 11
	r := ui.Rect{X: cx - d/2, Y: cy - d/2, W: d, H: d}
	p.Fill(ui.Rect{X: r.X - 2, Y: r.Y - 2, W: d + 4, H: d + 4}, k.Surface, d)
	switch s {
	case UserStatusOnline:
		p.Fill(r, k.Success, d)
	case UserStatusAway:
		p.Stroke(r, k.Warning, d, 3)
	case UserStatusBusy:
		p.Fill(r, k.Danger, d)
		p.Fill(ui.Rect{X: r.X + 2.5, Y: cy - 1, W: d - 5, H: 2}, k.Surface, 0)
	default:
		p.Stroke(r, k.TextMuted, d, 1.5)
	}
}
