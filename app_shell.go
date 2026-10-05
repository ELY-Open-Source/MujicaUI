package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// AppShellOptions are the slots of an AppShell. Content is required; the others are optional.
type AppShellOptions struct {
	// TitleBar goes across the top, as a TitleBar.
	TitleBar func()
	// Sidebar scrolls on its own beside the content and collapses to a rail.
	Sidebar func()
	// Content scrolls on its own.
	Content func()
	// StatusBar goes across the bottom, as a StatusBar.
	StatusBar func()
	// SidebarWidth is the open sidebar's width in DIPs; zero is 240.
	SidebarWidth float32
	// ContentScroll keeps the content's offset; nil keeps it by the shell's place.
	ContentScroll *ui.ScrollState
}

// AppShellResult is an AppShell and its event in this frame.
type AppShellResult struct {
	*ui.Element
	// Toggled reports the sidebar opened or collapsed.
	Toggled bool
}

// AppShell lays out a window: title bar, sidebar, content and status bar.
// *collapsed folds the sidebar to a rail; its toggle button and Cmd/Ctrl+B flip it.
func AppShell(c *ui.Context, collapsed *bool, opts AppShellOptions) AppShellResult {
	if collapsed == nil {
		panic("mujica: AppShell needs a collapsed state")
	}
	if opts.Content == nil {
		panic("mujica: AppShell needs Content")
	}
	if opts.SidebarWidth < 0 {
		panic("mujica: AppShell SidebarWidth must not be negative")
	}
	k := Tokens(c)
	var res AppShellResult
	toggle := func() {
		*collapsed = !*collapsed
		res.Toggled = true
		slog.Debug("mujica: app_shell toggled", "collapsed", *collapsed)
	}
	if opts.Sidebar != nil && c.Shortcut(ui.Cmd, ui.KeyB) {
		toggle()
	}
	e := ui.Column(c).Fill().AlignItems(ui.Stretch).Background(k.Background)
	e.Children(func() {
		if opts.TitleBar != nil {
			opts.TitleBar()
		}
		ui.Row(c).Grow(1).MinHeight(0).AlignItems(ui.Stretch).Children(func() {
			if opts.Sidebar != nil && appShellSidebar(c, *collapsed, opts) {
				toggle()
			}
			ScrollArea(c, opts.ContentScroll, ScrollAreaOptions{}, opts.Content).Grow(1).MinWidth(0)
		})
		if opts.StatusBar != nil {
			opts.StatusBar()
		}
	})
	res.Element = e
	return res
}

// appShellSidebar builds the sidebar or its rail and reports whether its toggle was pressed.
func appShellSidebar(c *ui.Context, collapsed bool, opts AppShellOptions) bool {
	k := Tokens(c)
	u := Density(c).Unit()
	open := opts.SidebarWidth
	if open == 0 {
		open = 240
	}
	rail := theme.MinHitSize + 4*u
	target := open
	if collapsed {
		target = rail
	}
	side := ui.Column(c).Shrink(0).Clip().AlignItems(ui.Stretch).Background(k.Surface).
		BorderWidth(0, theme.BorderWidth, 0, 0).BorderColor(k.Border)
	side.Width(side.AnimateWith("app_shell.width", target, Motion(c, theme.StateDuration), ui.EaseOut))
	var pressed bool
	side.Children(func() {
		head := ui.Row(c).Shrink(0).Padding(u * 2).Justify(ui.End)
		if collapsed {
			head.Justify(ui.Center)
		}
		head.Children(func() {
			label := msg(c, "app_shell.collapse")
			if collapsed {
				label = msg(c, "app_shell.expand")
			}
			b := ui.ButtonBase(c).Size(theme.MinHitSize, theme.MinHitSize).Radius(theme.ControlRadius).
				Label(label).Tooltip(label).TextColor(k.TextMuted)
			internal.Ring(b, k.Focus, theme.ControlRadius)
			if b.Pressed() {
				b.Background(k.Selection).TextColor(k.Text)
			} else if b.Hovered() {
				b.Background(k.SurfaceHover).TextColor(k.Text)
			}
			b.Children(func() { ui.Icon(c, icons.Must("panel-left")).Size(16, 16) })
			pressed = b.Clicked()
		})
		if !collapsed {
			ScrollArea(c, nil, ScrollAreaOptions{}, opts.Sidebar).Grow(1).MinHeight(0)
		}
	})
	return pressed
}
