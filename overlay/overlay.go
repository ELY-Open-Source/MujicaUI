package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
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
	site, focus := core.OverlaySite(c)
	if !*open {
		focus.Closed(c, site)
		return nil
	}
	k := core.Tokens(c)
	inset := opts.Inset
	if inset == 0 {
		inset = 16
	}
	fx, fy := overlayFractions(opts.Place)
	if opts.Modal {
		core.OverlayLayer(c, open)
		var backClicked bool
		shown := true
		panel := ui.DialogBase(c, &shown, func(back, panel *ui.Element) {
			backClicked = back.Clicked()
			p := core.OverlayEnter(c, panel)
			back.Padding(inset).Justify(overlayAlign(fy)).AlignItems(overlayAlign(fx))
			if opts.Dim {
				back.Background(OverlayScrim(k).Alpha(p))
			}
			core.OverlayStyle(c, panel, p)
			w, h := c.Size()
			panel.Label(opts.Label).Padding(core.Density(c).Unit() * 4).MaxWidth(w - 2*inset).MaxHeight(h - 2*inset)
			ui.Scroll(c).Shrink(1).Children(content)
		})
		focus.Shown(panel)
		if !shown {
			slog.Debug("mujica: overlay close", "backdrop", backClicked)
			*open = false
		}
		return panel
	}
	var panel *ui.Element
	ui.Overlay(c, func() {
		w, h := c.Size()
		panel = ui.Column(c).Absolute().Role(ui.RoleDialog).Label(opts.Label).Padding(core.Density(c).Unit() * 4).
			MaxWidth(w - 2*inset).MaxHeight(h - 2*inset)
		b := panel.Bounds()
		panel.Left(inset + (w-2*inset-b.W)*fx).Top(inset + (h-2*inset-b.H)*fy)
		p := core.OverlayEnter(c, panel)
		if b.W == 0 {
			p = 0 // not measured yet
			c.AnimationFrame()
		}
		core.OverlayStyle(c, panel, p)
		panel.Children(func() { ui.Scroll(c).Shrink(1).Children(content) })
	})
	focus.Shown(panel)
	if core.OverlayEscape(c, open) {
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

// OverlayScrim is the color dimming the window behind a modal layer.
func OverlayScrim(k theme.Tokens) ui.Color {
	if k.Dark {
		return ui.RGBA(0, 0, 0, 0.6)
	}
	return k.Text.Alpha(0.42)
}

// OverlayAnchored is overlayEnter, held at 0 until the anchor b has been laid out.
func OverlayAnchored(c *ui.Context, panel *ui.Element, b ui.Rect) float32 {
	if b.W == 0 && b.H == 0 {
		c.AnimationFrame()
		return 0
	}
	return core.OverlayEnter(c, panel)
}
