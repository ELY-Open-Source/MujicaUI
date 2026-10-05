package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ImageOptions configure an Image.
type ImageOptions struct {
	// Width and Height size the box in DIPs; zero keeps the picture's size,
	// or 160×120 for the placeholder, loading and error states.
	Width, Height float32
	// Crop covers the box and crops the rest; false fits the picture inside.
	Crop bool
	// Alt describes the picture for assistive technology.
	Alt string
	// Loading shows the caller's load in progress.
	Loading bool
	// Err shows the caller's load failure with a retry button.
	Err error
}

// ImageView is a shown Image and what the user did to it this frame.
type ImageView struct {
	Element *ui.Element
	// Retry reports a click on the retry button of the error state.
	Retry bool
}

// Image shows a *ui.Bitmap or *ui.SVG. A nil src shows a placeholder.
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
	w, h := opts.Width, opts.Height
	if w == 0 {
		w = 160
	}
	if h == 0 {
		h = 120
	}
	var v ImageView
	frame := ui.Column(c).Size(w, h).Shrink(0).Center().Gap(u*2).Padding(u*2).Clip().
		Radius(theme.ControlRadius).Background(k.SurfaceHover).Border(theme.BorderWidth, k.Border)
	v.Element = frame
	frame.Children(func() {
		switch {
		case opts.Err != nil:
			frame.Role(ui.RoleGroup).Label(msg(c, "image.error")).Description(opts.Alt).Error(opts.Err.Error())
			Icon(c, "circle-alert", IconOptions{Size: 20, Tone: IconDanger})
			ui.Text(c, msg(c, "image.error")).FontSize(theme.CaptionSize).TextColor(k.Danger).TextAlign(ui.Center)
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
