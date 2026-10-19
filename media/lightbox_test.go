package media

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func lightboxImages() []data.ViewerImage {
	return []data.ViewerImage{
		{Source: mediaTestBitmap(160, 90), Alt: "The grand hall"},
		{Source: mediaTestBitmap(90, 160), Alt: "The east tower"},
		{Source: mediaTestBitmap(120, 120), Alt: "The garden"},
	}
}

func TestLightboxBrowseAndClose(t *testing.T) {
	open, index, n := false, 0, 0
	images := lightboxImages()
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Button(c, "Open gallery", input.ButtonOptions{}).Clicked() {
			open = true
		}
		if Lightbox(c, &open, &index, images).Changed() {
			n++
		}
	}, 640, 480)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !open || !tt.HasText("The grand hall") || !tt.HasText("1 of 3") {
		t.Fatalf("not open: %v", tt.Texts())
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if index != 2 || n != 2 || !tt.HasText("The garden") {
		t.Fatalf("index %d changes %d", index, n)
	}
	if err := tt.Click("Previous image"); err != nil {
		t.Fatal(err)
	}
	if index != 1 {
		t.Fatalf("previous: %d", index)
	}
	tt.Key(0, ui.KeyEscape)
	if open {
		t.Fatal("Escape did not close")
	}
	if !tt.Focused("Open gallery") {
		t.Fatal("focus not restored to the opener")
	}
}

func TestLightboxBadIndex(t *testing.T) {
	open, index := true, 5
	mediaPanics(t, "index", func(c *ui.Context) { Lightbox(c, &open, &index, lightboxImages()) })
}

func TestLightboxSnapshot(t *testing.T) {
	open, index := true, 0
	images := lightboxImages()
	snapshot(t, "lightbox", 520, 360, func(c *ui.Context) {
		Lightbox(c, &open, &index, images)
	})
}
