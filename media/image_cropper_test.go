package media

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func cropperTester(src image.Image, crop *ImageCrop, opts ImageCropperOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Width(440).Children(func() {
			if ImageCropper(c, src, crop, opts).Changed() {
				*n++
			}
		})
	}, 500, 420)
}

func TestImageCropperDefaultAndKeys(t *testing.T) {
	crop, n := ImageCrop{}, 0
	tt := cropperTester(mediaTestImage(400, 200), &crop, ImageCropperOptions{Aspect: 1, Height: 200}, &n)
	if crop.Rect != image.Rect(100, 0, 300, 200) {
		t.Fatalf("default %v", crop.Rect)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if crop.Rect != image.Rect(104, 0, 304, 200) {
		t.Fatalf("Right: %v", crop.Rect)
	}
	tt.Key(ui.Shift, ui.KeyLeft)
	if crop.Rect.Dx() != crop.Rect.Dy() || crop.Rect.Dx() != 196 {
		t.Fatalf("Shift+Left keeps 1:1: %v", crop.Rect)
	}
	if n != 2 || !tt.HasText("104, 0 · 196 × 196 · 0°") {
		t.Fatalf("changes %d, texts %v", n, tt.Texts())
	}
}

func TestImageCropperRotate(t *testing.T) {
	crop, n := ImageCrop{Rect: image.Rect(10, 20, 110, 70)}, 0
	tt := cropperTester(mediaTestImage(400, 200), &crop, ImageCropperOptions{Height: 200}, &n)
	if err := tt.Click("Rotate right"); err != nil {
		t.Fatal(err)
	}
	if crop.Rotation != 90 || crop.Rect != image.Rect(130, 10, 180, 110) {
		t.Fatalf("rotated %+v", crop)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Shift, ui.KeyR)
	if crop.Rotation != 0 || crop.Rect != image.Rect(10, 20, 110, 70) {
		t.Fatalf("rotated back %+v", crop)
	}
}

func TestImageCropperDragCorner(t *testing.T) {
	crop, n := ImageCrop{Rect: image.Rect(0, 0, 200, 100)}, 0
	tt := cropperTester(mediaTestImage(400, 200), &crop, ImageCropperOptions{Aspect: 2, Height: 200}, &n)
	r, _ := tt.Find("Crop area")
	scale := min(r.W/400, r.H/200)
	ox, oy := r.X+(r.W-400*scale)/2, r.Y+(r.H-200*scale)/2
	cx, cy := ox+200*scale, oy+100*scale
	tt.Press(cx, cy)
	tt.Move(cx+50*scale, cy+5*scale)
	tt.Move(cx+100*scale, cy+10*scale)
	tt.Release(cx+100*scale, cy+10*scale)
	if crop.Rect.Min != (image.Point{}) || crop.Rect.Dx() < 290 || crop.Rect.Dx() > 310 || crop.Rect.Dx() != 2*crop.Rect.Dy() {
		t.Fatalf("corner drag %v", crop.Rect)
	}
}

func TestImageCropperDisabledAndBad(t *testing.T) {
	crop, n := ImageCrop{}, 0
	tt := cropperTester(mediaTestImage(100, 100), &crop, ImageCropperOptions{Disabled: true}, &n)
	start := crop
	_ = tt.Click("Rotate right")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if crop != start || n != 0 {
		t.Fatalf("disabled cropper changed %+v", crop)
	}
	img := mediaTestImage(100, 100)
	mediaPanics(t, "outside", func(c *ui.Context) {
		ImageCropper(c, img, &ImageCrop{Rect: image.Rect(0, 0, 200, 10)}, ImageCropperOptions{})
	})
	mediaPanics(t, "rotation", func(c *ui.Context) { ImageCropper(c, img, &ImageCrop{Rotation: 45}, ImageCropperOptions{}) })
	mediaPanics(t, "aspect", func(c *ui.Context) {
		ImageCropper(c, img, &ImageCrop{Rect: image.Rect(0, 0, 50, 10)}, ImageCropperOptions{Aspect: 1})
	})
}

func TestImageCropperSnapshot(t *testing.T) {
	img := mediaTestImage(320, 180)
	crop := ImageCrop{Rect: image.Rect(60, 30, 220, 150)}
	snapshot(t, "image_cropper", 440, 300, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(400).Children(func() {
			ImageCropper(c, img, &crop, ImageCropperOptions{Height: 220})
		})
	})
}
