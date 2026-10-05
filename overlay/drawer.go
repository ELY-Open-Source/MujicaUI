package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DrawerSide is the window edge a Drawer comes from. The zero value is DrawerRight.
type DrawerSide int

const (
	DrawerRight DrawerSide = iota
	DrawerLeft
	DrawerTop
	DrawerBottom
)

// DrawerOptions configure a Drawer. Title is required.
type DrawerOptions struct {
	Title string
	Side  DrawerSide
	// Size is the width or height (360 or 280 when zero), kept within 160 DIPs and 90% of the window.
	Size float32
	// NonModal leaves the window usable and the focus free.
	NonModal bool
	// Actions builds the footer, aligned to the end.
	Actions func()
}

// Drawer shows a panel along a window edge while *open is true; Escape closes it.
func Drawer(c *ui.Context, open *bool, opts DrawerOptions, content func()) *ui.Element {
	if opts.Title == "" {
		panic("mujica: Drawer needs a Title")
	}
	if opts.Side < DrawerRight || opts.Side > DrawerBottom {
		panic("mujica: Drawer Side out of range")
	}
	if opts.Size < 0 {
		panic("mujica: Drawer Size is negative")
	}
	site, focus := core.OverlaySite(c)
	if !*open {
		focus.Closed(c, site)
		return nil
	}
	k := core.Tokens(c)
	size := drawerSize(c, opts)
	closed := false
	body := func(panel *ui.Element, p float32) {
		drawerPanel(c, panel, opts.Side, size, p)
		panel.Label(opts.Title)
		closed = dialogHeader(c, opts.Title, false, true, "drawer.close")
		if content != nil {
			ui.Scroll(c).Grow(1).Shrink(1).Children(content)
		}
		if opts.Actions != nil {
			dialogActions(c, opts.Actions)
		}
	}
	var panel *ui.Element
	escaped := false
	if opts.NonModal {
		ui.Overlay(c, func() {
			panel = ui.Column(c).Absolute().Role(ui.RoleDialog)
			p := core.OverlayEnter(c, panel)
			drawerPlace(c, panel, opts.Side, size*p)
			panel.Children(func() { body(panel, p) })
		})
		escaped = core.OverlayEscape(c, open)
	} else {
		core.OverlayLayer(c, open)
		shown := true
		panel = ui.DialogBase(c, &shown, func(back, panel *ui.Element) {
			p := core.OverlayEnter(c, panel)
			back.Background(OverlayScrim(k).Alpha(p))
			drawerBackdrop(back, opts.Side)
			body(panel, p)
		})
		escaped = !shown
	}
	focus.Shown(panel)
	if closed || escaped {
		slog.Debug("mujica: drawer close", "escape", escaped)
		*open = false
	}
	return panel
}

func drawerSize(c *ui.Context, opts DrawerOptions) float32 {
	w, h := c.Size()
	extent, size := w, opts.Size
	if opts.Side == DrawerTop || opts.Side == DrawerBottom {
		extent = h
		if size == 0 {
			size = 280
		}
	} else if size == 0 {
		size = 360
	}
	return max(min(size, extent*0.9), min(160, extent))
}

func drawerHorizontal(s DrawerSide) bool { return s == DrawerLeft || s == DrawerRight }

// drawerBackdrop lines the panel up against its edge inside a modal backdrop.
func drawerBackdrop(back *ui.Element, s DrawerSide) {
	if drawerHorizontal(s) {
		back.Row()
	}
	back.AlignItems(ui.Stretch).Justify(ui.Start)
	if s == DrawerRight || s == DrawerBottom {
		back.Justify(ui.End)
	}
}

// drawerPlace puts a non-modal panel against its edge, shown by size.
func drawerPlace(c *ui.Context, panel *ui.Element, s DrawerSide, size float32) {
	w, h := c.Size()
	switch s {
	case DrawerRight:
		panel.Left(w - size).Top(0).Height(h)
	case DrawerLeft:
		panel.Left(0).Top(0).Height(h)
	case DrawerTop:
		panel.Left(0).Top(0).Width(w)
	case DrawerBottom:
		panel.Left(0).Top(h - size).Width(w)
	}
}

// drawerPanel styles the panel at entrance progress p, sliding in from its edge.
func drawerPanel(c *ui.Context, panel *ui.Element, s DrawerSide, size, p float32) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	off := -(1 - p) * size
	panel.Background(k.Surface).TextColor(k.Text).Opacity(p).Padding(u*6).Gap(u*3).
		Shadow(0, 0, 24, 0, ui.RGBA(0, 0, 0, 0.18*p))
	bw := theme.BorderWidth
	switch s {
	case DrawerRight:
		panel.Width(size).BorderWidth(0, 0, 0, bw).Margin(0, off, 0, 0)
	case DrawerLeft:
		panel.Width(size).BorderWidth(0, bw, 0, 0).Margin(0, 0, 0, off)
	case DrawerTop:
		panel.Height(size).BorderWidth(0, 0, bw, 0).Margin(off, 0, 0, 0)
	case DrawerBottom:
		panel.Height(size).BorderWidth(bw, 0, 0, 0).Margin(0, 0, off, 0)
	}
	panel.BorderColor(k.Border)
}
