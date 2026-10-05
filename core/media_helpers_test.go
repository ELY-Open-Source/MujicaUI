package core

import (
	"image"
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"
)

// mediaTestImage is a w×h picture whose left half is wine and right half ivory.
func mediaTestImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			c := color.RGBA{125, 32, 52, 255}
			if x >= w/2 {
				c = color.RGBA{242, 235, 221, 255}
			}
			if y < h/8 {
				c = color.RGBA{173, 145, 100, 255}
			}
			img.Set(x, y, c)
		}
	}
	return img
}

func mediaTestBitmap(w, h int) *ui.Bitmap { return ui.NewBitmap(mediaTestImage(w, h)) }

// mediaPanics fails t unless fn panics.
func mediaPanics(t *testing.T, name string, fn func(c *ui.Context)) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s: no panic", name)
		}
	}()
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		fn(c)
	}, 400, 300)
	tt.Frame()
}
