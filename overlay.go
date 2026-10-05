package mujica

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// OverlayPlace is where an Overlay sits in the window. The zero value is OverlayCenter.
type OverlayPlace int

const (
	OverlayCenter OverlayPlace = iota
	OverlayTop
	OverlayBottom
	OverlayLeft
	OverlayRight
	OverlayTopLeft
	OverlayTopRight
	OverlayBottomLeft
	OverlayBottomRight
)

// OverlayOptions configure an Overlay.
type OverlayOptions struct {
	Place OverlayPlace
	// Inset keeps the panel from the window edges; 16 DIPs when zero.
	Inset float32
	// Modal traps the focus, makes the window inert and closes on a backdrop click.
	Modal bool
	// Dim darkens the window behind a modal overlay. It requires Modal.
	Dim bool
	// Label names the layer for assistive technology.
	Label string
}

// Overlay shows content in a layer above the window while *open is true; Escape closes it.
func Overlay(c *ui.Context, open *bool, opts OverlayOptions, content func()) *ui.Element {
	if opts.Dim && !opts.Modal {
		panic("mujica: Overlay Dim requires Modal")
	}
	if opts.Place < OverlayCenter || opts.Place > OverlayBottomRight {
		panic("mujica: Overlay Place out of range")
	}
	site, focus := overlaySite(c)
	if !*open {
		focus.closed(c, site)
		return nil
	}
	k := Tokens(c)
	inset := opts.Inset
	if inset == 0 {
		inset = 16
	}
	fx, fy := overlayFractions(opts.Place)
	if opts.Modal {
		overlayLayer(c, open)
		var backClicked bool
		shown := true
		panel := ui.DialogBase(c, &shown, func(back, panel *ui.Element) {
			backClicked = back.Clicked()
			p := overlayEnter(c, panel)
			back.Padding(inset).Justify(overlayAlign(fy)).AlignItems(overlayAlign(fx))
			if opts.Dim {
				back.Background(overlayScrim(k).Alpha(p))
			}
			overlayStyle(c, panel, p)
			w, h := c.Size()
			panel.Label(opts.Label).Padding(Density(c).Unit() * 4).MaxWidth(w - 2*inset).MaxHeight(h - 2*inset)
			ui.Scroll(c).Shrink(1).Children(content)
		})
		focus.shown(panel)
		if !shown {
			slog.Debug("mujica: overlay close", "backdrop", backClicked)
			*open = false
		}
		return panel
	}
	var panel *ui.Element
	ui.Overlay(c, func() {
		w, h := c.Size()
		panel = ui.Column(c).Absolute().Role(ui.RoleDialog).Label(opts.Label).Padding(Density(c).Unit() * 4).
			MaxWidth(w - 2*inset).MaxHeight(h - 2*inset)
		b := panel.Bounds()
		panel.Left(inset + (w-2*inset-b.W)*fx).Top(inset + (h-2*inset-b.H)*fy)
		p := overlayEnter(c, panel)
		if b.W == 0 {
			p = 0 // not measured yet
			c.AnimationFrame()
		}
		overlayStyle(c, panel, p)
		panel.Children(func() { ui.Scroll(c).Shrink(1).Children(content) })
	})
	focus.shown(panel)
	if overlayEscape(c, open) {
		slog.Debug("mujica: overlay close", "escape", true)
		*open = false
	}
	return panel
}

func overlayFractions(p OverlayPlace) (fx, fy float32) {
	switch p {
	case OverlayTop:
		return 0.5, 0
	case OverlayBottom:
		return 0.5, 1
	case OverlayLeft:
		return 0, 0.5
	case OverlayRight:
		return 1, 0.5
	case OverlayTopLeft:
		return 0, 0
	case OverlayTopRight:
		return 1, 0
	case OverlayBottomLeft:
		return 0, 1
	case OverlayBottomRight:
		return 1, 1
	}
	return 0.5, 0.5
}

func overlayAlign(f float32) ui.Align {
	switch f {
	case 0:
		return ui.Start
	case 1:
		return ui.End
	}
	return ui.Center
}

// overlayEnter returns how far, 0 to 1, the popup born with e has entered.
func overlayEnter(c *ui.Context, e *ui.Element) float32 {
	born := ui.Local(e, "mujica.overlay.born", func() time.Time { return c.Now() })
	d := Motion(c, theme.PopupDuration)
	if d == 0 {
		return 1
	}
	t := float32(c.Now().Sub(*born)) / float32(d)
	if t >= 1 {
		return 1
	}
	c.AnimationFrame()
	return ui.EaseOut(t)
}

// overlayScrim is the color dimming the window behind a modal layer.
func overlayScrim(k theme.Tokens) ui.Color {
	if k.Dark {
		return ui.RGBA(0, 0, 0, 0.6)
	}
	return k.Text.Alpha(0.42)
}

// overlayStyle gives a popup panel the surface, edge and shadow of the theme.
func overlayStyle(c *ui.Context, panel *ui.Element, p float32) {
	k := Tokens(c)
	panel.Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.PopupRadius).
		TextColor(k.Text).Opacity(p).Shadow(0, 6, 18, 0, ui.RGBA(0, 0, 0, 0.16*p))
}

// overlayPlace puts a w×h box below anchor a, above when there is no room, inside the window.
func overlayPlace(c *ui.Context, a ui.Rect, w, h, gap float32) (x, y float32) {
	ww, wh := c.Size()
	const edge = 4
	x, y = a.X, a.Y+a.H+gap
	if y+h > wh-edge && a.Y-gap-h >= edge {
		y = a.Y - gap - h
	}
	x = max(edge, min(x, ww-edge-w))
	y = max(edge, min(y, wh-edge-h))
	return x, y
}

// overlayAnchored is overlayEnter, held at 0 until the anchor b has been laid out.
func overlayAnchored(c *ui.Context, panel *ui.Element, b ui.Rect) float32 {
	if b.W == 0 && b.H == 0 {
		c.AnimationFrame()
		return 0
	}
	return overlayEnter(c, panel)
}
