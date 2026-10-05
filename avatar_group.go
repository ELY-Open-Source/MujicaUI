package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// AvatarItem is one person of an AvatarGroup.
type AvatarItem struct {
	Name  string
	Image *ui.Bitmap
}

// AvatarGroupOptions configure an AvatarGroup.
type AvatarGroupOptions struct {
	// Max is how many avatars show before the "+N" button; zero means 4.
	Max int
	// Size is each avatar's diameter; zero means 32.
	Size float32
}

type avatarGroupKey struct{}

// AvatarGroupView is a shown AvatarGroup and its full list's state.
type AvatarGroupView struct {
	Element *ui.Element
	// Open reports whether the full list shows.
	Open bool
	// Toggled reports the full list opened or closed this frame.
	Toggled bool
}

// AvatarGroup shows overlapping avatars. Past Max, a "+N" button opens a
// popover with the full list of names.
func AvatarGroup(c *ui.Context, items []AvatarItem, opts AvatarGroupOptions) AvatarGroupView {
	limit := opts.Max
	if limit < 0 {
		panic("mujica: AvatarGroup Max must not be negative")
	}
	if limit == 0 {
		limit = 4
	}
	size := opts.Size
	if size < 0 {
		panic("mujica: AvatarGroup Size must not be negative")
	}
	if size == 0 {
		size = 32
	}
	k := Tokens(c)
	shown := min(limit, len(items))
	ring := float32(2)
	row := ui.Row(c).AlignItems(ui.Center).Role(ui.RoleGroup)
	open := ui.Local(row, avatarGroupKey{}, func() bool { return false })
	was := *open
	row.Children(func() {
		for i, it := range items[:shown] {
			cell := ui.Box(c).Key(i).Padding(ring).Radius(size/2 + ring).Background(k.Background).Shrink(0)
			if i > 0 {
				cell.Margin(0, 0, 0, -size/4)
			}
			cell.Children(func() { Avatar(c, it.Name, AvatarOptions{Image: it.Image, Size: size}) })
		}
		rest := len(items) - shown
		if rest == 0 {
			*open = false
			return
		}
		more := avatarGroupMore(c, k, size, ring, rest, len(items))
		if more.Clicked() {
			*open = !*open
		}
		ui.PopoverBase(c, more, open, func(panel *ui.Element) {
			avatarGroupList(c, k, panel, items)
		})
	})
	v := AvatarGroupView{Element: row, Open: *open, Toggled: *open != was}
	if v.Toggled {
		slog.Debug("mujica: avatargroup toggle", "open", v.Open)
	}
	return v
}

func avatarGroupMore(c *ui.Context, k theme.Tokens, size, ring float32, rest, total int) *ui.Element {
	cell := ui.Box(c).Key("more").Padding(ring).Radius(size/2+ring).Background(k.Background).
		Margin(0, 0, 0, -size/4).Shrink(0)
	var b *ui.Element
	cell.Children(func() {
		b = ui.ButtonBase(c).Size(size, size).Radius(size / 2).Label(fmt.Sprintf(msg(c, "avatar_group.more"), total))
		bg := k.Surface
		if b.Pressed() {
			bg = k.Selection
		} else if b.Hovered() {
			bg = k.SurfaceHover
		}
		b.Background(bg).Border(theme.BorderWidth, k.ControlBorder).TextColor(k.Text)
		internal.Ring(b, k.Focus, size/2)
		internal.Fade(b, Motion(c, theme.StateDuration))
		b.Children(func() {
			ui.Text(c, fmt.Sprintf("+%d", rest)).FontSize(size * 0.36).FontWeight(600).SingleLine()
		})
	})
	return b
}

func avatarGroupList(c *ui.Context, k theme.Tokens, panel *ui.Element, items []AvatarItem) {
	u := Density(c).Unit()
	panel.Margin(u, 0, 0, 0).Padding(u*2).Radius(theme.PopupRadius).Background(k.Surface).
		Border(theme.BorderWidth, k.Border).MinWidth(200)
	if d := Motion(c, theme.PopupDuration); d > 0 {
		panel.Transition(ui.ElementTransition{Duration: d, Enter: &ui.Motion{Y: -4}})
	}
	ui.Scroll(c).MaxHeight(240).Children(func() {
		ui.Column(c).Gap(u * 2).Role(ui.RoleList).Children(func() {
			for i, it := range items {
				ui.Row(c).Key(i).Gap(u * 2).AlignItems(ui.Center).Role(ui.RoleListItem).Children(func() {
					Avatar(c, it.Name, AvatarOptions{Image: it.Image, Size: 24})
					ui.Text(c, it.Name).SingleLine()
				})
			}
		})
	})
}
