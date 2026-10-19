package media

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
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
func ImageThumbnail(c *ui.Context, src ui.ImageSource, opts ImageThumbnailOptions) core.ChoiceResult {
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
	k := core.Tokens(c)
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
	b.Border(width, edge).Transition(core.ColorFade)
	b.Children(func() {
		switch {
		case opts.Err != nil:
			b.Description(opts.Err.Error()).Tooltip(opts.Err.Error())
			ui.Column(c).Gap(4).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(20, 20).TextColor(k.Danger)
				if h >= 64 {
					ui.Text(c, core.Msg(c, "image_thumbnail.failed")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Danger).SingleLine()
				}
			})
		case opts.Loading:
			core.Spinner(c, core.SpinnerOptions{Size: min(20, h/2), Label: core.Msg(c, "image_thumbnail.loading")})
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
	return core.NewChoiceResult(b, clicked)
}
