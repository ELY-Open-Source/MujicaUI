package media

import (
	"image"
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal/snaptest"
)

// mediaPanics fails t unless fn panics.
func mediaPanics(t *testing.T, name string, fn func(c *ui.Context)) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s: no panic", name)
		}
	}()
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		fn(c)
	}, 400, 300)
	tt.Frame()
}

// snapshot compares view in both themes with testdata/snapshots/<os>/<name>-<mode>.png.
func snapshot(t *testing.T, name string, w, h int, view func(c *ui.Context)) {
	t.Helper()
	snaptest.Run(t, name, w, h, snapshotTheme, view, nil)
}

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

func snapshotTheme(c *ui.Context, mode string) {
	m := core.Light
	if mode == "dark" {
		m = core.Dark
	}
	core.Use(c, core.Settings{Mode: m})
}
