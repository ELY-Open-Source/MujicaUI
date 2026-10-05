package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
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

// imageBlank sizes the box of a nil src: 160×120.
var imageBlank = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="160" height="120"/>`))

type imageRoomKey struct{}

// imageErrorTextWidth is the room the error text needs.
const imageErrorTextWidth = 120

// Image shows a bitmap or SVG; nil shows a placeholder.
func Image(c *ui.Context, src ui.ImageSource, opts ImageOptions) ImageView {
	if opts.Width < 0 || opts.Height < 0 {
		panic("mujica: Image Width and Height must not be negative")
	}
	loaded := opts.Err == nil && !opts.Loading && src != nil
	var v ImageView
	v.Element = imageSized(c, src, opts, func(box, pic *ui.Element) {
		if loaded {
			if opts.Alt != "" {
				pic.Label(opts.Alt)
			}
			return
		}
		pic.Invisible()
		v.Retry = imageLayers(c, box, opts)
	})
	return v
}

// imageLayers builds the state layers over the hidden picture.
func imageLayers(c *ui.Context, box *ui.Element, opts ImageOptions) (retry bool) {
	k := Tokens(c)
	u := Density(c).Unit()
	room := box.Bounds()
	// Build again when the box changed, settling the error text.
	if last := ui.Local(box, imageRoomKey{}, func() ui.Rect { return ui.Rect{W: -1} }); *last != room {
		*last = room
		c.AnimationFrame()
	}
	{
		// The clipped backdrop; Retry sits above it, unclipped.
		back := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Center().Clip().
			Radius(theme.ControlRadius).Background(k.SurfaceHover).Border(theme.BorderWidth, k.Border)
		switch {
		case opts.Err != nil:
			box.Role(ui.RoleGroup).Label(msg(c, "image.error")).Description(opts.Alt).Error(opts.Err.Error())
			ui.Row(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Center().Gap(u * 2).Children(func() {
				if imageRetry(c, k, opts.Err.Error()).Clicked() {
					retry = true
					slog.Debug("mujica: image retry", "alt", opts.Alt)
				}
				if room.W >= theme.MinHitSize+u*4+imageErrorTextWidth && room.H >= theme.MinHitSize {
					ui.Text(c, msg(c, "image.error")).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.Danger).SingleLine()
				}
			})
		case opts.Loading:
			box.Role(ui.RoleImage).Label(msg(c, "image.loading")).Description(opts.Alt)
			back.Children(func() {
				Spinner(c, SpinnerOptions{Size: max(8, min(20, room.H-4)), Label: msg(c, "image.loading")})
			})
		default:
			box.Role(ui.RoleImage).Label(msg(c, "image.placeholder")).Description(opts.Alt)
			back.Children(func() { Icon(c, "image", IconOptions{Size: 24, Tone: IconMuted}) })
		}
	}
	return retry
}

// imageSized builds the outer box and its picture, sized alike in every state.
func imageSized(c *ui.Context, src ui.ImageSource, opts ImageOptions, layers func(box, pic *ui.Element)) *ui.Element {
	if src == nil {
		src = imageBlank
	}
	box := ui.Box(c)
	box.Children(func() {
		pic := ui.Image(c, src).Radius(theme.ControlRadius)
		if opts.Crop {
			pic.Fit(ui.Cover)
		}
		if opts.Width > 0 {
			box.Width(opts.Width)
			pic.Width(opts.Width)
		}
		if opts.Height > 0 {
			box.Height(opts.Height)
			pic.Height(opts.Height)
		}
		layers(box, pic)
	})
	return box
}

// imageRetry is the 28 DIP icon-only retry button.
func imageRetry(c *ui.Context, k theme.Tokens, problem string) *ui.Element {
	b := ui.ButtonBase(c).Size(theme.MinHitSize, theme.MinHitSize).Shrink(0).Radius(theme.ControlRadius).
		Label(msg(c, "image.retry")).Description(problem).Tooltip(problem)
	bg := k.Surface
	if b.Pressed() {
		bg = k.Selection
	} else if b.Hovered() {
		bg = k.SurfaceHover
	}
	b.Background(bg).Border(theme.BorderWidth, k.ControlBorder).TextColor(k.Danger)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	b.Transition(colorFade)
	b.Children(func() { Icon(c, "refresh-cw", IconOptions{Size: 16}) })
	return b
}
