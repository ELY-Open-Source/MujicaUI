package media

import (
	"errors"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestImageThumbnailClickAndKeys(t *testing.T) {
	n := 0
	bmp := mediaTestBitmap(120, 80)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Padding(20).AlignItems(ui.Start).Children(func() {
			if ImageThumbnail(c, bmp, ImageThumbnailOptions{Alt: "Hall", Aspect: 4.0 / 3}).Changed() {
				n++
			}
		})
	}, 300, 200)
	r, ok := tt.Find("Hall")
	if !ok || r.W != 96 || r.H != 72 {
		t.Fatalf("bounds %v", r)
	}
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeySpace)
	if n != 3 {
		t.Fatalf("activations %d", n)
	}
}

func TestImageThumbnailStatesKeepSize(t *testing.T) {
	var got []ui.Rect
	for _, o := range []ImageThumbnailOptions{{}, {Loading: true}, {Err: errors.New("404")}, {Selected: true}} {
		o.Alt, o.Width = "Hall", 120
		var r ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			ui.Row(c).AlignItems(ui.Start).Children(func() {
				r = ImageThumbnail(c, mediaTestBitmap(40, 40), o).Element.Bounds()
			})
		}, 300, 200)
		tt.Frame()
		got = append(got, r)
	}
	for _, r := range got {
		if r != got[0] || r.W != 120 || r.H != 120 {
			t.Fatalf("sizes %v", got)
		}
	}
}

func TestImageThumbnailDisabledAndBad(t *testing.T) {
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ImageThumbnail(c, nil, ImageThumbnailOptions{Alt: "Hall", Disabled: true}).Changed() {
			n++
		}
	}, 300, 200)
	_ = tt.Click("Hall")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if n != 0 {
		t.Fatal("disabled thumbnail activated")
	}
	mediaPanics(t, "no alt", func(c *ui.Context) { ImageThumbnail(c, nil, ImageThumbnailOptions{}) })
}

func TestImageThumbnailSnapshot(t *testing.T) {
	bmp := mediaTestBitmap(160, 90)
	snapshot(t, "image_thumbnail", 520, 120, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(12).AlignItems(ui.Start).Children(func() {
			ImageThumbnail(c, bmp, ImageThumbnailOptions{Alt: "Hall", Width: 80})
			ImageThumbnail(c, bmp, ImageThumbnailOptions{Alt: "Hall", Width: 80, Selected: true})
			ImageThumbnail(c, bmp, ImageThumbnailOptions{Alt: "Hall", Width: 80, Loading: true})
			ImageThumbnail(c, bmp, ImageThumbnailOptions{Alt: "Hall", Width: 80, Err: errors.New("404")})
			ImageThumbnail(c, bmp, ImageThumbnailOptions{Alt: "Hall", Width: 80, Disabled: true})
		})
	})
}
