package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// DialogOptions configure a Dialog. Title is required.
type DialogOptions struct {
	Title       string
	Description string
	// Width is 440 DIPs when zero; the window always bounds it.
	Width float32
	// Decorated draws the double border, corner ornaments and title rule.
	Decorated bool
	// Actions builds the footer, usually buttons, aligned to the end.
	Actions func()
	// KeepOnBackdrop ignores clicks on the backdrop.
	KeepOnBackdrop bool
	// KeepOnEscape ignores Escape.
	KeepOnEscape bool
	// HideClose removes the close button from the header.
	HideClose bool
}

// Dialog shows a modal dialog while *open is true, trapping and then restoring the focus.
func Dialog(c *ui.Context, open *bool, opts DialogOptions, content func()) *ui.Element {
	if opts.Title == "" {
		panic("mujica: Dialog needs a Title")
	}
	site, focus := overlaySite(c)
	if !*open {
		focus.closed(c, site)
		return nil
	}
	overlayLayer(c, open)
	k := Tokens(c)
	shown, backClicked, closed := true, false, false
	panel := ui.DialogBase(c, &shown, func(back, panel *ui.Element) {
		backClicked = back.Clicked()
		p := overlayEnter(c, panel)
		back.Background(overlayScrim(k).Alpha(p)).Padding(16)
		dialogPanel(c, panel, opts.Width, opts.Decorated, p)
		panel.Label(opts.Title)
		if opts.Description != "" {
			panel.Description(opts.Description)
		}
		closed = dialogHeader(c, opts.Title, opts.Description, opts.Decorated, !opts.HideClose, "dialog.close")
		if content != nil {
			ui.Scroll(c).Shrink(1).Children(content)
		}
		if opts.Actions != nil {
			dialogActions(c, opts.Actions)
		}
	})
	focus.shown(panel)
	switch {
	case closed:
		slog.Debug("mujica: dialog close", "by", "button")
		*open = false
	case !shown && backClicked && !opts.KeepOnBackdrop:
		slog.Debug("mujica: dialog close", "by", "backdrop")
		*open = false
	case !shown && !backClicked && !opts.KeepOnEscape:
		slog.Debug("mujica: dialog close", "by", "escape")
		*open = false
	}
	return panel
}

// dialogPanel styles a centered dialog panel at entrance progress p.
func dialogPanel(c *ui.Context, panel *ui.Element, width float32, decorated bool, p float32) {
	k := Tokens(c)
	u := Density(c).Unit()
	if width == 0 {
		width = 440
	}
	w, h := c.Size()
	overlayStyle(c, panel, p)
	panel.Width(width).MaxWidth(w-32).MaxHeight(h-32).Padding(u*6).Gap(u*3).Margin(8*(1-p), 0, 0, 0)
	if decorated {
		panel.Border(theme.BorderWidth, k.Ornament).Padding(u*6 + 6).DrawOver(func(pt *ui.Painter, r ui.Rect) {
			dialogOrnaments(pt, r, k.Ornament, k.Surface)
		})
	}
}

// dialogOrnaments draws the inner frame and the four corner diamonds.
func dialogOrnaments(p *ui.Painter, r ui.Rect, ink, ground ui.Color) {
	const in = 5
	f := ui.Rect{X: r.X + in, Y: r.Y + in, W: r.W - 2*in, H: r.H - 2*in}
	p.Stroke(f, ink, 1, 1)
	for _, pt := range [4][2]float32{{f.X, f.Y}, {f.X + f.W, f.Y}, {f.X, f.Y + f.H}, {f.X + f.W, f.Y + f.H}} {
		x, y := internal.Snap(p, pt[0]), internal.Snap(p, pt[1])
		p.FillPath(internal.Diamond(x, y, 11), ground)
		p.FillPath(internal.Diamond(x, y, 7), ink)
	}
}

// dialogHeader builds the title, description and close button; it reports a close click.
func dialogHeader(c *ui.Context, title, description string, decorated, closable bool, closeKey string) bool {
	k := Tokens(c)
	u := Density(c).Unit()
	clicked := false
	ui.Column(c).Gap(u).Shrink(0).Children(func() {
		ui.Row(c).Gap(u * 3).AlignItems(ui.Start).Children(func() {
			ui.Text(c, title).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).FontWeight(600).Grow(1)
			if closable {
				clicked = dialogCloseButton(c, msg(c, closeKey)).Clicked()
			}
		})
		if decorated {
			dialogRule(c, k.Ornament)
		}
		if description != "" {
			ui.Text(c, description).TextColor(k.TextMuted)
		}
	})
	return clicked
}

// dialogRule is the ornament line under a decorated title.
func dialogRule(c *ui.Context, ink ui.Color) {
	ui.Box(c).FillWidth().Height(9).Margin(2, 0, 2, 0).Draw(func(p *ui.Painter, r ui.Rect) {
		y, cx := r.Y+r.H/2, r.X+r.W/2
		p.Line(r.X, y, cx-10, y, 1, ink)
		p.Line(cx+10, y, r.X+r.W, y, 1, ink)
		p.FillPath(internal.Diamond(internal.Snap(p, cx), internal.Snap(p, y), 7), ink)
	})
}

// dialogCloseButton is the 28 DIP icon button closing a panel.
func dialogCloseButton(c *ui.Context, label string) *ui.Element {
	k := Tokens(c)
	b := ui.ButtonBase(c).Size(theme.MinHitSize, theme.MinHitSize).Radius(theme.ControlRadius).Label(label).Margin(-4, -4, 0, 0)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	b.TextColor(k.TextMuted)
	if b.Pressed() {
		b.Background(k.Selection).TextColor(k.Text)
	} else if b.Hovered() {
		b.Background(k.SurfaceHover).TextColor(k.Text)
	}
	b.Children(func() { ui.Icon(c, icons.Must("x")).Size(16, 16) })
	return b
}

// dialogActions lays out a dialog footer.
func dialogActions(c *ui.Context, actions func()) {
	u := Density(c).Unit()
	ui.Row(c).Gap(u*2).Justify(ui.End).Wrap().Shrink(0).Margin(u*2, 0, 0, 0).Children(actions)
}
