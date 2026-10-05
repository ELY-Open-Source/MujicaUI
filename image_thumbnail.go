package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// ImageThumbnailOptions configure an ImageThumbnail.
type ImageThumbnailOptions struct {
	// Width is the thumbnail's width; zero means 96.
	Width float32
	// Aspect is width over height, cropped; zero is 1.
	Aspect   float32
	Selected bool
	// Loading and Err show the load state at one size.
	Loading bool
	Err     error
	// Alt names the picture; required.
	Alt      string
	Disabled bool
}

// ImageThumbnail is a cropped picture; Changed reports activation.
func ImageThumbnail(c *ui.Context, src ui.ImageSource, opts ImageThumbnailOptions) ChoiceResult {
	if opts.Width < 0 || opts.Aspect < 0 {
		panic("mujica: ImageThumbnail needs Width and Aspect >= 0")
	}
	if opts.Alt == "" {
		panic("mujica: ImageThumbnail needs Alt")
	}
	w, aspect := opts.Width, opts.Aspect
	if w == 0 {
		w = 96
	}
	if aspect == 0 {
		aspect = 1
	}
	k := Tokens(c)
	h := w / aspect
	b := ui.ButtonBase(c).Size(w, h).Radius(theme.ControlRadius).Clip().
		Disabled(opts.Disabled).Label(opts.Alt).Selected(opts.Selected).Background(k.SurfaceHover)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	edge, width := k.Border, theme.BorderWidth
	switch {
	case opts.Selected:
		edge, width = k.AccentText, 2
	case b.Hovered():
		edge = k.ControlBorder
	}
	b.Border(width, edge).Transition(colorFade)
	b.Children(func() {
		switch {
		case opts.Err != nil:
			b.Description(opts.Err.Error()).Tooltip(opts.Err.Error())
			ui.Column(c).Gap(4).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(20, 20).TextColor(k.Danger)
				if h >= 64 {
					ui.Text(c, msg(c, "image_thumbnail.failed")).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.Danger).SingleLine()
				}
			})
		case opts.Loading:
			Spinner(c, SpinnerOptions{Size: min(20, h/2), Label: msg(c, "image_thumbnail.loading")})
		case src == nil:
			ui.Icon(c, icons.Must("image")).Size(20, 20).TextColor(k.TextMuted)
		default:
			ui.Image(c, src).Fit(ui.Cover).Size(w, h)
		}
		if opts.Selected {
			// A diamond mark, so selection never rests on color alone.
			ui.Box(c).Absolute().Top(4).Right(4).Size(14, 14).PassThrough().Draw(func(p *ui.Painter, r ui.Rect) {
				cx, cy := internal.Snap(p, r.X+r.W/2), internal.Snap(p, r.Y+r.H/2)
				p.FillPath(internal.Diamond(cx, cy, 14), k.Surface)
				p.FillPath(internal.Diamond(cx, cy, 10), k.AccentText)
			})
		}
	})
	clicked := b.Clicked() && !b.IsDisabled()
	if clicked {
		slog.Debug("mujica: image_thumbnail click", "alt", opts.Alt)
	}
	return ChoiceResult{Element: b, changed: clicked}
}
