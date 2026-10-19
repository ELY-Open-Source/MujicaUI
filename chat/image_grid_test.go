package chat

import (
	"errors"
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

// imageGridPicture is a plain colored picture standing in for a photo.
func imageGridPicture(color string) ui.ImageSource {
	return ui.MustParseSVG([]byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="200" height="150"><rect width="200" height="150" fill="%s"/></svg>`, color)))
}

func imageGridSet(n int) []GridImage {
	colors := []string{"#7D2034", "#AD9164", "#355E80", "#326247", "#805714", "#6A5B55", "#9C3C53", "#D2C3AF", "#211C1D"}
	out := make([]GridImage, n)
	for i := range out {
		out[i] = GridImage{Source: imageGridPicture(colors[i]), Alt: fmt.Sprintf("Photo %d", i+1)}
	}
	return out
}

func TestImageGridLayoutsByCount(t *testing.T) {
	for n, cols := range map[int]int{1: 1, 2: 2, 3: 3, 4: 2, 5: 3, 9: 3} {
		var xs map[float32]bool
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			ImageGrid(c, imageGridSet(n))
		}, 900, 900)
		tt.Frame()
		xs = map[float32]bool{}
		for i := range n {
			r, ok := tt.Find(fmt.Sprintf("Photo %d", i+1))
			if !ok {
				t.Fatalf("%d images: Photo %d missing", n, i+1)
			}
			xs[r.X] = true
		}
		if len(xs) != cols {
			t.Errorf("%d images in %d columns, want %d", n, len(xs), cols)
		}
	}
}

func TestImageGridOpensViewerAndRetries(t *testing.T) {
	imgs := imageGridSet(3)
	imgs[1] = GridImage{Alt: "Lost", Err: errors.New("404")}
	retried := -1
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if i, ok := ImageGrid(c, imgs).Retried(); ok {
			retried = i
		}
	}, 800, 600)
	if err := tt.Click("Photo 3"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Image 2 of 2") {
		t.Fatalf("viewer not at the third picture (second viewable): %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	if err := tt.Click("Retry"); err != nil {
		t.Fatal(err)
	}
	if retried != 1 {
		t.Fatalf("retried %d, want 1", retried)
	}
}

func TestImageGridCountLimits(t *testing.T) {
	for _, n := range []int{0, 10} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%d images: no panic", n)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ImageGrid(c, make([]GridImage, n))
			}, 100, 100)
		}()
	}
}

func TestImageGridSnapshot(t *testing.T) {
	imgs := imageGridSet(5)
	imgs[3] = GridImage{Alt: "Loading", Loading: true}
	imgs[4] = GridImage{Alt: "Lost", Err: errors.New("404")}
	snapshot(t, "image_grid", 420, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { ImageGrid(c, imgs) })
	})
}
