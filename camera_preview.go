package mujica

import (
	"image"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// CameraFrame is the caller's latest camera picture.
type CameraFrame struct {
	// Image is the frame; nil shows "No signal".
	Image image.Image
	// Seq numbers frames; a new number redraws a reused Image.
	Seq uint64
}

// CameraPreviewOptions configure a CameraPreview.
type CameraPreviewOptions struct {
	// Label names the camera for assistive technology.
	Label    string
	Disabled bool
}

type cameraCache struct {
	img    image.Image
	seq    uint64
	mirror bool
	bmp    *ui.Bitmap
}

type cameraKey struct{}

// CameraPreview shows frame at 16:9, mirrored while *mirror.
func CameraPreview(c *ui.Context, frame CameraFrame, mirror *bool, opts CameraPreviewOptions) ChoiceResult {
	k := Tokens(c)
	old := *mirror
	label := opts.Label
	if label == "" {
		label = msg(c, "camera_preview.label")
	}
	box := ui.Box(c).FillWidth().AspectRatio(16.0 / 9).Background(videoBlack).Radius(theme.CardRadius).Clip().
		Center().Focusable().Role(ui.RoleImage).Label(label).Disabled(opts.Disabled)
	internal.Ring(box, k.Focus, theme.CardRadius)
	cache := ui.Local(box, cameraKey{}, func() cameraCache { return cameraCache{} })
	box.Children(func() {
		if frame.Image == nil {
			box.Description(msg(c, "camera_preview.noSignal"))
			ui.Column(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, mediaCamera).Size(28, 28).TextColor(ui.Hex("#BCAEAA"))
				ui.Text(c, msg(c, "camera_preview.noSignal")).TextColor(ui.Hex("#BCAEAA"))
			})
		} else {
			// Copies are made once per frame number, not per paint.
			if cache.img != frame.Image || cache.seq != frame.Seq || cache.mirror != *mirror || cache.bmp == nil {
				*cache = cameraCache{frame.Image, frame.Seq, *mirror, cameraBitmap(frame.Image, *mirror)}
			}
			ui.Image(c, cache.bmp).Fit(ui.Contain).Fill()
		}
		ui.Box(c).Absolute().Top(8).Right(8).Background(ui.RGBA(0, 0, 0, 0.5)).Radius(theme.ControlRadius).Children(func() {
			b := mediaToggleAction(c, mediaMirror, msg(c, "camera_preview.mirror"), *mirror)
			if b.Clicked() && !b.IsDisabled() {
				*mirror = !*mirror
			}
		})
	})
	if !box.IsDisabled() && box.Shortcut(0, ui.KeyM) {
		*mirror = !*mirror
	}
	changed := *mirror != old
	if changed {
		slog.Debug("mujica: camera_preview mirror", "on", *mirror)
	}
	return ChoiceResult{Element: box, changed: changed}
}

// cameraBitmap copies img into a bitmap, flipped when mirror.
func cameraBitmap(img image.Image, mirror bool) *ui.Bitmap {
	if !mirror {
		return ui.NewBitmap(img)
	}
	b := img.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			out.Set(b.Max.X-1-x, y-b.Min.Y, img.At(x, y))
		}
	}
	return ui.NewBitmap(out)
}
