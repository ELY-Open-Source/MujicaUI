package media

import (
	"fmt"
	"image"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ImageCrop is Rect in pixels of the turned image.
type ImageCrop struct {
	Rect image.Rectangle
	// Rotation is 0, 90, 180 or 270.
	Rotation int
}

// ImageCropperOptions configure an ImageCropper.
type ImageCropperOptions struct {
	// Aspect locks width over height; zero is free.
	Aspect float64
	// Height is the stage's height; zero means 320.
	Height   float32
	Disabled bool
	Label    string
}

type cropperCache struct {
	src image.Image
	rot int
	bmp *ui.Bitmap
}

type cropperDrag struct {
	mode   int // 0 none, 1 move, 2-5 corners
	start  image.Rectangle
	dx, dy float32
}

type cropperKey struct{}

// ImageCropper edits *crop over src; empty Rect starts fitted.
func ImageCropper(c *ui.Context, src image.Image, crop *ImageCrop, opts ImageCropperOptions) core.ChoiceResult {
	if src == nil || src.Bounds().Empty() || opts.Aspect < 0 || opts.Height < 0 {
		panic("mujica: ImageCropper needs an image, Aspect >= 0 and Height >= 0")
	}
	if crop.Rotation%90 != 0 || crop.Rotation < 0 || crop.Rotation >= 360 {
		panic(fmt.Sprintf("mujica: ImageCropper Rotation %d is not 0, 90, 180 or 270", crop.Rotation))
	}
	size := cropperSize(src, crop.Rotation)
	if crop.Rect.Empty() {
		crop.Rect = cropperFit(size, opts.Aspect)
	}
	if !crop.Rect.In(image.Rectangle{Max: size}) {
		panic(fmt.Sprintf("mujica: ImageCropper Rect %v outside the %v image", crop.Rect, size))
	}
	if opts.Aspect > 0 && math.Abs(float64(crop.Rect.Dx())-float64(crop.Rect.Dy())*opts.Aspect) > 1+opts.Aspect {
		panic(fmt.Sprintf("mujica: ImageCropper Rect %v does not keep Aspect %v", crop.Rect, opts.Aspect))
	}
	h := opts.Height
	if h == 0 {
		h = 320
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	old := *crop
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "image_cropper.label")
	}
	col := ui.Column(c).Gap(u * 2).Disabled(opts.Disabled)
	col.Children(func() {
		stage := ui.Box(c).FillWidth().Height(h).Focusable().Clip().Radius(theme.ControlRadius).
			Background(VideoBlack).Label(label).Cursor(ui.CursorMove)
		internal.Ring(stage, k.Focus, theme.ControlRadius)
		r := stage.Bounds()
		scale, ox, oy := cropperPlace(r, size)
		live := !stage.IsDisabled()
		if live && scale > 0 {
			cropperPointer(stage, crop, r, size, scale, ox, oy, opts.Aspect)
			cropperKeys(stage, crop, size, opts.Aspect)
		}
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			left := core.IconAction(c, core.MediaRotateLeft, core.Msg(c, "image_cropper.rotateLeft"))
			right := core.IconAction(c, core.MediaRotate, core.Msg(c, "image_cropper.rotateRight"))
			switch {
			case live && (right.Clicked() || stage.Shortcut(0, ui.KeyR)):
				cropperRotate(crop, size, 90, opts.Aspect)
			case live && (left.Clicked() || stage.Shortcut(ui.Shift, ui.KeyR)):
				cropperRotate(crop, size, 270, opts.Aspect)
			}
			ui.Text(c, cropperReadout(*crop)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).FontFeatures("tnum")
		})
		stage.Description(cropperReadout(*crop))
		// Turned pixels are copied only when image or rotation change.
		cache := ui.Local(stage, cropperKey{}, func() cropperCache { return cropperCache{} })
		if cache.src != src || cache.rot != crop.Rotation || cache.bmp == nil {
			*cache = cropperCache{src, crop.Rotation, ui.NewBitmap(cropperTurn(src, crop.Rotation))}
		}
		size := cropperSize(src, crop.Rotation)
		bmp, rect := cache.bmp, crop.Rect
		stage.Draw(func(p *ui.Painter, r ui.Rect) {
			scale, ox, oy := cropperPlace(r, size)
			img := ui.Rect{X: ox, Y: oy, W: float32(size.X) * scale, H: float32(size.Y) * scale}
			p.Image(bmp, img, ui.FillBox)
			f := ui.Rect{X: internal.Snap(p, ox+float32(rect.Min.X)*scale), Y: internal.Snap(p, oy+float32(rect.Min.Y)*scale),
				W: internal.Snap(p, float32(rect.Dx())*scale), H: internal.Snap(p, float32(rect.Dy())*scale)}
			shade := ui.RGBA(0, 0, 0, 0.55)
			p.Fill(ui.Rect{X: img.X, Y: img.Y, W: img.W, H: f.Y - img.Y}, shade, 0)
			p.Fill(ui.Rect{X: img.X, Y: f.Y + f.H, W: img.W, H: img.Y + img.H - f.Y - f.H}, shade, 0)
			p.Fill(ui.Rect{X: img.X, Y: f.Y, W: f.X - img.X, H: f.H}, shade, 0)
			p.Fill(ui.Rect{X: f.X + f.W, Y: f.Y, W: img.X + img.W - f.X - f.W, H: f.H}, shade, 0)
			light := ui.Hex("#F2EBDD")
			p.Stroke(f, light, 0, 1)
			for i := 1; i < 3; i++ {
				p.Fill(ui.Rect{X: f.X + f.W*float32(i)/3, Y: f.Y, W: 1, H: f.H}, light.Alpha(0.4), 0)
				p.Fill(ui.Rect{X: f.X, Y: f.Y + f.H*float32(i)/3, W: f.W, H: 1}, light.Alpha(0.4), 0)
			}
			for _, pt := range [4][2]float32{{f.X, f.Y}, {f.X + f.W, f.Y}, {f.X, f.Y + f.H}, {f.X + f.W, f.Y + f.H}} {
				p.FillPath(internal.Diamond(internal.Snap(p, pt[0]), internal.Snap(p, pt[1]), 12), light)
			}
		})
	})
	changed := *crop != old
	if changed {
		slog.Debug("mujica: image_cropper changed", "rect", crop.Rect, "rotation", crop.Rotation)
	}
	return core.NewChoiceResult(col, changed)
}

func cropperReadout(cr ImageCrop) string {
	r := cr.Rect
	return fmt.Sprintf("%d, %d · %d × %d · %d°", r.Min.X, r.Min.Y, r.Dx(), r.Dy(), cr.Rotation)
}

// cropperSize is src's size after turning it rot degrees.
func cropperSize(src image.Image, rot int) image.Point {
	s := src.Bounds().Size()
	if rot == 90 || rot == 270 {
		return image.Pt(s.Y, s.X)
	}
	return s
}

// cropperFit is the largest centered aspect rectangle inside size.
func cropperFit(size image.Point, aspect float64) image.Rectangle {
	w, h := size.X, size.Y
	if aspect > 0 {
		if float64(w) > float64(h)*aspect {
			w = int(math.Round(float64(h) * aspect))
		} else {
			h = int(math.Round(float64(w) / aspect))
		}
	}
	x, y := (size.X-w)/2, (size.Y-h)/2
	return image.Rect(x, y, x+w, y+h)
}

// cropperPlace fits size in r: scale and top-left corner.
func cropperPlace(r ui.Rect, size image.Point) (scale, x, y float32) {
	if r.W <= 0 {
		return 0, 0, 0
	}
	scale = min(r.W/float32(size.X), r.H/float32(size.Y))
	return scale, r.X + (r.W-float32(size.X)*scale)/2, r.Y + (r.H-float32(size.Y)*scale)/2
}

// cropperTurn copies src turned rot degrees clockwise.
func cropperTurn(src image.Image, rot int) image.Image {
	if rot == 0 {
		return src
	}
	b := src.Bounds()
	size := cropperSize(src, rot)
	out := image.NewRGBA(image.Rectangle{Max: size})
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			var nx, ny int
			switch rot {
			case 90:
				nx, ny = b.Dy()-1-y, x
			case 180:
				nx, ny = b.Dx()-1-x, b.Dy()-1-y
			case 270:
				nx, ny = y, b.Dx()-1-x
			}
			out.Set(nx, ny, src.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return out
}

// cropperRotate turns the crop d degrees; locked aspects refit.
func cropperRotate(cr *ImageCrop, size image.Point, d int, aspect float64) {
	r := cr.Rect
	for range d / 90 {
		r = image.Rect(size.Y-r.Max.Y, r.Min.X, size.Y-r.Min.Y, r.Max.X)
		size = image.Pt(size.Y, size.X)
	}
	cr.Rotation = (cr.Rotation + d) % 360
	cr.Rect = r
	if aspect > 0 {
		cr.Rect = cropperFit(size, aspect)
	}
}

// cropperPointer moves the frame, or resizes it by a corner.
func cropperPointer(stage *ui.Element, cr *ImageCrop, r ui.Rect, size image.Point, scale, ox, oy float32, aspect float64) {
	d := ui.Local(stage, cropperDrag{}, func() cropperDrag { return cropperDrag{} })
	dx, dy, ok := stage.Dragged()
	if !stage.Pressed() {
		d.mode = 0
		return
	}
	if d.mode == 0 {
		px, py, _ := stage.PointerPosition()
		px, py = (px+r.X-ox)/scale, (py+r.Y-oy)/scale
		d.start, d.dx, d.dy, d.mode = cr.Rect, 0, 0, cropperHit(cr.Rect, px, py, 14/scale)
	}
	if !ok || d.mode == 0 {
		return
	}
	d.dx, d.dy = d.dx+dx/scale, d.dy+dy/scale
	cr.Rect = cropperDragRect(d.start, d.mode, int(math.Round(float64(d.dx))), int(math.Round(float64(d.dy))), size, aspect)
}

// cropperHit says what a press grabs: corner, inside or nothing.
func cropperHit(r image.Rectangle, x, y, slop float32) int {
	corners := [4][2]int{{r.Min.X, r.Min.Y}, {r.Max.X, r.Min.Y}, {r.Min.X, r.Max.Y}, {r.Max.X, r.Max.Y}}
	for i, c := range corners {
		if core.Abs32(x-float32(c[0])) <= slop && core.Abs32(y-float32(c[1])) <= slop {
			return 2 + i
		}
	}
	if x >= float32(r.Min.X) && x <= float32(r.Max.X) && y >= float32(r.Min.Y) && y <= float32(r.Max.Y) {
		return 1
	}
	return 0
}
