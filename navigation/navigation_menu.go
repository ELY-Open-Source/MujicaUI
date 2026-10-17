package navigation

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// NavigationLink is a page of a NavigationMenu, or with Items a group of pages.
type NavigationLink struct {
	Label string
	// Path is pushed on the router; groups need none.
	Path        string
	Icon        *ui.SVG
	Description string
	Items       []NavigationLink
}

// NavigationMenuOptions configure a NavigationMenu.
type NavigationMenuOptions struct {
	// History shows back and forward buttons for the router.
	History bool
	// Label names the menu for assistive technology.
	Label string
}

// NavigationMenu shows links pushing pages on router; the router's path marks the current one.
func NavigationMenu(c *ui.Context, router *ui.Router, links []NavigationLink, opts NavigationMenuOptions) ui.Element {
	if router == nil {
		panic("mujica: NavigationMenu needs a router")
	}
	if len(links) == 0 {
		panic("mujica: NavigationMenu needs links")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	path := router.Path()
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "navigation_menu.label")
	}
	push := func(p string) {
		if p != path {
			slog.Debug("mujica: navigation menu push", "from", path, "to", p)
			router.Push(p)
		}
	}
	row := ui.Row(c).AlignItems(ui.Center).Gap(u * 2)
	row.Children(func() {
		if opts.History {
			ui.Row(c).Gap(u).Children(func() {
				navigationMenuHistory(c, ui.BackButton(c, router), core.Msg(c, "navigation_menu.back"))
				navigationMenuHistory(c, ui.ForwardButton(c, router), core.Msg(c, "navigation_menu.forward"))
			})
			ui.Box(c).Size(1, 20).Background(k.Ornament).Role(ui.RoleNone)
		}
		ui.Row(c).AlignItems(ui.Center).Gap(u).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).Label(label).Children(func() {
			for _, l := range links {
				navigationMenuTop(c, l, path, push)
			}
		})
	})
	return row
}

func navigationMenuHistory(c *ui.Context, b ui.Element, name string) {
	k := core.Tokens(c)
	d := core.Density(c)
	b.Label(name).MinHeight(d.ControlHeight()).MinWidth(d.ControlHeight()).Radius(theme.ControlRadius).Background(k.Surface)
	internal.Ring(b, k.Focus, theme.ControlRadius)
}

// navigationMenuActive reports whether l or a link inside it is the page at path.
func navigationMenuActive(l NavigationLink, path string) bool {
	// The root page is current only on itself, not under every path.
	if l.Path != "" && (path == l.Path || l.Path != "/" && strings.HasPrefix(path, strings.TrimSuffix(l.Path, "/")+"/")) {
		return true
	}
	for _, it := range l.Items {
		if navigationMenuActive(it, path) {
			return true
		}
	}
	return false
}

func navigationMenuTop(c *ui.Context, l NavigationLink, path string, push func(string)) {
	if l.Label == "" {
		panic("mujica: NavigationLink needs a Label")
	}
	if l.Items == nil && l.Path == "" {
		panic("mujica: NavigationLink needs a Path or Items")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	on := navigationMenuActive(l, path)
	b := ui.ButtonBase(c).MinHeight(d.ControlHeight()).Padding(0, u*3).Gap(u * 1.5).Radius(theme.ControlRadius).
		Label(l.Label).Cursor(ui.CursorPointer)
	open := ui.Local(b, "open", func() bool { return false })
	if l.Items == nil {
		b.Role(ui.RoleLink)
	} else {
		b.Role(ui.RoleMenuButton)
	}
	fg := k.TextMuted
	if on {
		b.FontWeight(600)
		fg = k.Text
	}
	if b.Hovered() || *open {
		b.Background(k.SurfaceHover)
		fg = k.Text
	}
	internal.Ease(b, core.Motion(c, theme.StateDuration), fg)
	b.FocusRing(false).DrawOver(func(p *ui.Painter, r ui.Rect) {
		if on {
			p.Fill(ui.Rect{X: r.X + u*3, Y: r.Y + r.H - 2, W: r.W - u*6, H: 2}, k.AccentText, 1)
		}
		if b.FocusVisible() {
			internal.FocusRing(p, r, theme.ControlRadius, k.Focus)
		}
	})
	b.Children(func() {
		if l.Icon != nil {
			ui.Icon(c, l.Icon).Size(16, 16)
		}
		ui.Text(c, l.Label).SingleLine()
		if l.Items != nil {
			ui.Icon(c, icons.Must("chevron-down")).Size(12, 12).TextColor(k.TextMuted)
		}
	})
	if l.Items == nil {
		if b.Clicked() {
			push(l.Path)
		}
		return
	}
	if b.Clicked() {
		*open = !*open
	} else if b.Shortcut(0, ui.KeyDown) {
		*open = true
	}
	ui.PopoverBase(c, b, open, func(panel ui.Element) {
		core.OverlayStyle(c, panel, core.OverlayEnter(c, panel))
		// Capped to the room beside the anchor; the links scroll inside.
		_, wh := c.Size()
		r := b.Bounds()
		panel.Margin(u, 0, 0, 0).Padding(u * 2).MaxHeight(max(r.Y, wh-r.Y-r.H) - 16).MinWidth(240).MaxWidth(360)
		ui.Scroll(c).Shrink(1).Gap(u).FocusGroup(ui.Vertical).Children(func() {
			navigationMenuPanel(c, l.Items, path, open, push)
		})
	})
}

// navigationMenuPanel lists a group's links, its subgroups under serif headings.
func navigationMenuPanel(c *ui.Context, items []NavigationLink, path string, open *bool, push func(string)) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	first := true
	var list func(items []NavigationLink)
	list = func(items []NavigationLink) {
		for _, it := range items {
			if it.Label == "" {
				panic("mujica: NavigationLink needs a Label")
			}
			if it.Items != nil {
				ui.Text(c, it.Label).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Ornament.Mix(k.Text, 0.35)).
					Padding(u*2, u*2, u, u*2).SingleLine()
				list(it.Items)
				continue
			}
			if it.Path == "" {
				panic("mujica: NavigationLink needs a Path or Items")
			}
			if navigationMenuItem(c, it, it.Path == path, first).Clicked() {
				*open = false
				push(it.Path)
			}
			first = false
		}
	}
	list(items)
}

func navigationMenuItem(c *ui.Context, l NavigationLink, on, first bool) ui.Element {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	e := ui.ButtonBase(c).Role(ui.RoleLink).Justify(ui.Start).AlignItems(ui.Start).Gap(u*3).
		Padding(u*2, u*3).Radius(theme.ControlRadius).Label(l.Label).TextColor(k.Text).Cursor(ui.CursorPointer)
	if first {
		core.AutoFocus(c, e)
	}
	internal.Ring(e, k.Focus, theme.ControlRadius)
	switch {
	case on:
		e.Background(k.Selection).FontWeight(600)
		e.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + 4, W: 3, H: r.H - 8}, k.AccentText, 1.5)
		})
	case e.Hovered():
		e.Background(k.SurfaceHover)
	}
	internal.Ease(e, core.Motion(c, theme.StateDuration), k.Text)
	e.Children(func() {
		if l.Icon != nil {
			ui.Icon(c, l.Icon).Size(16, 16).Margin(2, 0, 0, 0)
		}
		ui.Column(c).Gap(2).MinWidth(0).Children(func() {
			ui.Text(c, l.Label).SingleLine()
			if l.Description != "" {
				ui.Text(c, l.Description).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).FontWeight(400)
			}
		})
	})
	return e
}
