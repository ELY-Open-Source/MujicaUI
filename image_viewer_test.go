package mujica

import (
	"image"
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"
)

// viewerPicture is red on the left half and blue on the right.
func viewerPicture() *ui.Bitmap {
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))
	for y := range 100 {
		for x := range 200 {
			c := color.RGBA{200, 30, 30, 255}
			if x >= 100 {
				c = color.RGBA{30, 30, 200, 255}
			}
			img.SetRGBA(x, y, c)
		}
	}
	return ui.NewBitmap(img)
}

type viewerHarness struct {
	open    bool
	index   int
	images  []ViewerImage
	changes int
}

func newViewerHarness() *viewerHarness {
	pic := viewerPicture()
	return &viewerHarness{open: true, images: []ViewerImage{{pic, "Rose window"}, {pic, "Choir stalls"}, {pic, "Crypt door"}}}
}

func (h *viewerHarness) view(c *ui.Context) {
	Use(c, Settings{Mode: Light})
	if ImageViewer(c, &h.open, &h.index, h.images).Changed() {
		h.changes++
	}
}

func centerPixel(tt *ui.Tester) color.RGBA {
	r, _ := tt.Find("Rose window") // the stage, by its Label
	return color.RGBAModel.Convert(tt.Image().At(int(r.X+r.W/2), int(r.Y+r.H/2))).(color.RGBA)
}

func TestImageViewerSwitchAndClose(t *testing.T) {
	h := newViewerHarness()
	tt := ui.NewTester(h.view, 640, 480)
	if !tt.HasText("Image 1 of 3") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Previous image")
	if h.index != 0 {
		t.Fatal("previous moved before the first image")
	}
	_ = tt.Click("Next image")
	if h.index != 1 || h.changes != 1 || !tt.HasText("Choir stalls") {
		t.Fatalf("next: index %d changes %d", h.index, h.changes)
	}
	tt.Key(0, ui.KeyRight)
	if h.index != 2 {
		t.Fatalf("Right: %d", h.index)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyEscape)
	if h.open || h.index != 1 {
		t.Fatalf("Escape: open %v index %d", h.open, h.index)
	}
	h.open = true
	tt.Frame()
	_ = tt.Click("Close")
	if h.open {
		t.Fatal("close button left it open")
	}
}

func TestImageViewerZoomAndPan(t *testing.T) {
	h := newViewerHarness()
	tt := ui.NewTester(h.view, 640, 480)
	tt.Key(0, ui.KeyEqual)
	if !tt.HasText("125%") {
		t.Fatalf("+ did not zoom: %q", tt.Texts())
	}
	_ = tt.Click("Zoom in")
	_ = tt.Click("Zoom in")
	if !tt.HasText("195%") {
		t.Fatalf("zoom in: %q", tt.Texts())
	}
	for range 10 {
		tt.Key(ui.Shift, ui.KeyLeft) // the picture moves right: its left half shows
	}
	if p := centerPixel(tt); p.R < 150 || p.B > 100 {
		t.Fatalf("pan right shows %v, want red", p)
	}
	for range 20 {
		tt.Key(ui.Shift, ui.KeyRight)
	}
	if p := centerPixel(tt); p.B < 150 || p.R > 100 {
		t.Fatalf("pan left shows %v, want blue", p)
	}
	tt.Key(0, ui.Key0)
	if !tt.HasText("100%") {
		t.Fatalf("0 did not fit: %q", tt.Texts())
	}
	for range 20 {
		tt.Key(0, ui.KeyMinus)
	}
	if !tt.HasText("25%") {
		t.Fatalf("zoom floor: %q", tt.Texts())
	}
}

func TestImageViewerClosedAndBad(t *testing.T) {
	h := newViewerHarness()
	h.open = false
	tt := ui.NewTester(h.view, 300, 200)
	if len(tt.Texts()) != 0 {
		t.Fatalf("closed viewer shows %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic on bad index")
		}
	}()
	h.open, h.index = true, 9
	tt.Frame()
}

// viewerArt is a wine diamond on old gold, for the snapshot.
func viewerArt() *ui.Bitmap {
	img := image.NewRGBA(image.Rect(0, 0, 240, 160))
	for y := range 160 {
		for x := range 240 {
			c := color.RGBA{173, 145, 100, 255}
			if viewerAbs(x-120)*2/3+viewerAbs(y-80) < 60 {
				c = color.RGBA{125, 32, 52, 255}
			}
			img.SetRGBA(x, y, c)
		}
	}
	return ui.NewBitmap(img)
}

func viewerAbs(v int) int { return max(v, -v) }

func TestImageViewerSnapshot(t *testing.T) {
	pic := viewerArt()
	snapshot(t, "image_viewer", 560, 360, func(c *ui.Context) {
		open, i := true, 0
		ImageViewer(c, &open, &i, []ViewerImage{{pic, "Rose window"}, {pic, "Choir stalls"}})
	})
}
