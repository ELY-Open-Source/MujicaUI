package media

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func cameraTester(frame *CameraFrame, mirror *bool, disabled bool, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Width(320).Children(func() {
			if CameraPreview(c, *frame, mirror, CameraPreviewOptions{Disabled: disabled}).Changed() {
				*n++
			}
		})
	}, 400, 260)
}

func TestCameraPreviewNoSignalAndMirror(t *testing.T) {
	frame, mirror, n := CameraFrame{}, false, 0
	tt := cameraTester(&frame, &mirror, false, &n)
	if !tt.HasText("No signal") {
		t.Fatalf("texts %v", tt.Texts())
	}
	if err := tt.Click("Mirror"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Camera") {
		t.Fatal("preview not focusable")
	}
	tt.Key(0, ui.KeyM)
	if mirror || n != 2 {
		t.Fatalf("mirror %v changes %d", mirror, n)
	}
	frame = CameraFrame{Image: mediaTestImage(32, 18), Seq: 1}
	tt.Frame()
	if tt.HasText("No signal") {
		t.Fatal("frame did not replace No signal")
	}
}

func TestCameraPreviewDisabled(t *testing.T) {
	frame, mirror, n := CameraFrame{}, false, 0
	tt := cameraTester(&frame, &mirror, true, &n)
	_ = tt.Click("Mirror")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyM)
	if mirror || n != 0 {
		t.Fatal("disabled preview flipped")
	}
}

func TestCameraBitmapMirrors(t *testing.T) {
	img := mediaTestImage(8, 8)
	flipped := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := range 8 {
		for x := range 8 {
			flipped.Set(7-x, y, img.At(x, y))
		}
	}
	got := ui.Render(func(c *ui.Context) {
		ui.Image(c, cameraBitmap(img, true)).Size(8, 8)
	}, 8, 8, 1)
	want := ui.Render(func(c *ui.Context) {
		ui.Image(c, ui.NewBitmap(flipped)).Size(8, 8)
	}, 8, 8, 1)
	if got.RGBAAt(1, 6) != want.RGBAAt(1, 6) || got.RGBAAt(1, 6) == ui.Render(func(c *ui.Context) {
		ui.Image(c, ui.NewBitmap(img)).Size(8, 8)
	}, 8, 8, 1).RGBAAt(1, 6) {
		t.Fatal("mirrored bitmap is not flipped")
	}
}

func TestCameraPreviewSnapshot(t *testing.T) {
	img := mediaTestImage(160, 90)
	m1, m2 := true, false
	snapshot(t, "camera_preview", 560, 180, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(12).Children(func() {
			ui.Column(c).Grow(1).Children(func() { CameraPreview(c, CameraFrame{Image: img}, &m1, CameraPreviewOptions{}) })
			ui.Column(c).Grow(1).Children(func() { CameraPreview(c, CameraFrame{}, &m2, CameraPreviewOptions{}) })
		})
	})
}
