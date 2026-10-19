package media

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func compareTester(split *float64, disabled bool, n *int) *ui.Tester {
	before, after := mediaTestBitmap(160, 90), mediaTestBitmap(160, 90)
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Width(320).Children(func() {
			if ImageCompare(c, before, after, split, ImageCompareOptions{Disabled: disabled}).Changed() {
				*n++
			}
		})
	}, 400, 300)
}

func TestImageCompareKeysAndDrag(t *testing.T) {
	split, n := 0.5, 0
	tt := compareTester(&split, false, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if split != 0.51 {
		t.Fatalf("Right: %v", split)
	}
	tt.Key(0, ui.KeyHome)
	if split != 0 {
		t.Fatalf("Home: %v", split)
	}
	r, ok := tt.Find("Before / After")
	if !ok {
		t.Fatal("not named")
	}
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(r.X+r.W*0.75, r.Y+r.H/2)
	tt.Release(r.X+r.W*0.75, r.Y+r.H/2)
	if split < 0.7 || split > 0.8 {
		t.Fatalf("drag: %v", split)
	}
	if n < 3 {
		t.Fatalf("changes %d", n)
	}
}

func TestImageCompareDisabled(t *testing.T) {
	split, n := 0.5, 0
	tt := compareTester(&split, true, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	r, _ := tt.Find("Before / After")
	tt.ClickAt(r.X+10, r.Y+10)
	if split != 0.5 || n != 0 {
		t.Fatalf("disabled moved to %v", split)
	}
	bad := 2.0
	mediaPanics(t, "split over 1", func(c *ui.Context) {
		ImageCompare(c, mediaTestBitmap(4, 4), mediaTestBitmap(4, 4), &bad, ImageCompareOptions{})
	})
}

func TestImageCompareSnapshot(t *testing.T) {
	split := 0.4
	before := mediaTestBitmap(160, 90)
	after := ui.NewBitmap(mediaTestImage(90, 160))
	snapshot(t, "image_compare", 360, 220, func(c *ui.Context) {
		ui.Column(c).Padding(20).Width(320).Children(func() {
			ImageCompare(c, before, after, &split, ImageCompareOptions{})
		})
	})
}
