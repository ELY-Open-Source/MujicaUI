package mujica

import (
	"errors"
	"image"
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"
)

func imageTestBitmap() *ui.Bitmap {
	img := image.NewRGBA(image.Rect(0, 0, 80, 40))
	for y := range 40 {
		for x := range 80 {
			img.Set(x, y, color.RGBA{uint8(100 + x), 32, 52, 255})
		}
	}
	return ui.NewBitmap(img)
}

func TestImageLoadedSizes(t *testing.T) {
	bmp := imageTestBitmap()
	var natural, sized ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			natural = Image(c, bmp, ImageOptions{Alt: "Gallery"}).Element.Bounds()
			sized = Image(c, bmp, ImageOptions{Width: 50, Height: 50, Crop: true}).Element.Bounds()
		})
	}, 300, 300)
	tt.Frame()
	if natural.W != 80 || natural.H != 40 {
		t.Fatalf("natural %v", natural)
	}
	if sized.W != 50 || sized.H != 50 {
		t.Fatalf("sized %v", sized)
	}
	if _, ok := tt.Find("Gallery"); !ok {
		t.Fatal("alt text not exposed")
	}
}

func TestImageStates(t *testing.T) {
	for _, tc := range []struct {
		name  string
		src   ui.ImageSource
		opts  ImageOptions
		label string
	}{
		{"placeholder", nil, ImageOptions{}, "No image"},
		{"loading", imageTestBitmap(), ImageOptions{Loading: true}, "Loading image"},
		{"error", imageTestBitmap(), ImageOptions{Loading: true, Err: errors.New("404")}, "Image failed to load"},
	} {
		var r ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			ui.Column(c).AlignItems(ui.Start).Children(func() {
				r = Image(c, tc.src, tc.opts).Element.Bounds()
			})
		}, 300, 300)
		tt.Frame()
		if r.W != 160 || r.H != 120 {
			t.Errorf("%s: size %v", tc.name, r)
		}
		if _, ok := tt.Find(tc.label); !ok {
			t.Errorf("%s: no %q", tc.name, tc.label)
		}
	}
}

func TestImageRetry(t *testing.T) {
	retries := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if Image(c, nil, ImageOptions{Width: 200, Height: 140, Err: errors.New("timeout")}).Retry {
			retries++
		}
	}, 300, 300)
	if err := tt.Click("Retry"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if retries != 2 {
		t.Fatalf("retries %d, want 2", retries)
	}
}

func TestImageNegativeSizePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Image(c, nil, ImageOptions{Width: -1})
	}, 100, 100)
	tt.Frame()
}

func TestImageSnapshot(t *testing.T) {
	bmp := imageTestBitmap()
	snapshot(t, "image", 600, 180, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(12).AlignItems(ui.Start).Children(func() {
			Image(c, bmp, ImageOptions{Width: 120, Height: 120, Crop: true, Alt: "Cropped"})
			Image(c, nil, ImageOptions{Width: 120, Height: 120})
			Image(c, nil, ImageOptions{Width: 120, Height: 120, Loading: true})
			Image(c, nil, ImageOptions{Width: 160, Height: 140, Err: errors.New("404")})
		})
	})
}
