package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ImageOptions configure an Image.
type ImageOptions struct {
	// Width, Height in DIPs; zero uses src size or 160×120.
	Width, Height float32
	// Crop covers and crops the box; false fits inside.
	Crop bool
	// Alt describes the picture for assistive technology.
	Alt string
	// Loading shows the caller's load in progress.
	Loading bool
	// Err shows the caller's load failure with a retry button.
	Err error
}

// ImageView is a shown Image and this frame's events.
type ImageView struct {
	Element *ui.Element
	// Retry reports the error state's retry button was pressed.
	Retry bool
}

// Image shows a bitmap or SVG; nil shows a placeholder.
func Image(c *ui.Context, src ui.ImageSource, opts ImageOptions) ImageView {
	if opts.Width < 0 || opts.Height < 0 {
		panic("mujica: Image Width and Height must not be negative")
	}
	if opts.Err == nil && !opts.Loading && src != nil {
		e := ui.Image(c, src).Radius(theme.ControlRadius)
		if opts.Crop {
			e.Fit(ui.Cover)
		}
		if opts.Width > 0 {
			e.Width(opts.Width)
		}
		if opts.Height > 0 {
			e.Height(opts.Height)
		}
		if opts.Alt != "" {
			e.Label(opts.Alt)
		}
		return ImageView{Element: e}
	}
	k := Tokens(c)
	u := Density(c).Unit()
	w, h := imageBox(src, opts.Width, opts.Height)
	var v ImageView
	frame := ui.Column(c).Size(w, h).Shrink(0).Center().Gap(u*2).Padding(u*2).Clip().
		Radius(theme.ControlRadius).Background(k.SurfaceHover).Border(theme.BorderWidth, k.Border)
	v.Element = frame
	frame.Children(func() {
		switch {
		case opts.Err != nil:
			frame.Role(ui.RoleGroup).Label(msg(c, "image.error")).Description(opts.Alt).Error(opts.Err.Error())
			Icon(c, "circle-alert", IconOptions{Size: 20, Tone: IconDanger})
			ui.Text(c, msg(c, "image.error")).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.Danger).TextAlign(ui.Center)
			if Button(c, msg(c, "image.retry"), ButtonOptions{Variant: Secondary}).Clicked() {
				v.Retry = true
				slog.Debug("mujica: image retry", "alt", opts.Alt)
			}
		case opts.Loading:
			frame.Role(ui.RoleImage).Label(msg(c, "image.loading")).Description(opts.Alt)
			Spinner(c, SpinnerOptions{Size: 20, Label: msg(c, "image.loading")})
		default:
			frame.Role(ui.RoleImage).Label(msg(c, "image.placeholder")).Description(opts.Alt)
			Icon(c, "image", IconOptions{Size: 24, Tone: IconMuted})
		}
	})
	return v
}

// imageBox fills unset sides from src's aspect, else 160×120.
func imageBox(src ui.ImageSource, w, h float32) (float32, float32) {
	var sw, sh float32
	switch s := src.(type) {
	case *ui.Bitmap:
		iw, ih := s.Size()
		sw, sh = float32(iw), float32(ih)
	case *ui.SVG:
		sw, sh = s.Size()
	}
	if sw <= 0 || sh <= 0 {
		sw, sh = 160, 120
	}
	switch {
	case w == 0 && h == 0:
		return sw, sh
	case w == 0:
		return h * sw / sh, h
	case h == 0:
		return w, w * sh / sw
	}
	return w, h
}
