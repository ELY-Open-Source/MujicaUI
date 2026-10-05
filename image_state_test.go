package mujica

import (
	"errors"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func TestImageErrorTextNeedsRoom(t *testing.T) {
	for _, tc := range []struct {
		w    float32
		text bool
	}{{200, true}, {80, false}} {
		ready := false
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			if ready {
				Image(c, nil, ImageOptions{Width: tc.w, Height: 60, Err: errors.New("404")})
			}
		}, 300, 200)
		ready = true
		tt.Frame()
		n := 0
		for _, s := range tt.Texts() {
			if s == "Image failed to load" {
				n++
			}
		}
		// One is the box's name; a second is the visible text.
		if got := n == 2; got != tc.text {
			t.Errorf("width %v: error text shown %v, want %v (%v)", tc.w, got, tc.text, tt.Texts())
		}
	}
}

// imageStateBoxes returns the box of loaded, loading and error states.
func imageStateBoxes(t *testing.T, opts ImageOptions, parent func(c *ui.Context, fn func())) []ui.Rect {
	t.Helper()
	bmp := imageTestBitmap()
	var got []ui.Rect
	for _, state := range []ImageOptions{{}, {Loading: true}, {Err: errors.New("x")}} {
		o := opts
		o.Loading, o.Err = state.Loading, state.Err
		var r ui.Rect
		tt := ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{})
			parent(c, func() { r = Image(c, bmp, o).Element.Bounds() })
		}, 400, 400)
		tt.Frame()
		got = append(got, r)
	}
	return got
}

func TestImageStatesShareBoxInColumn(t *testing.T) {
	got := imageStateBoxes(t, ImageOptions{}, func(c *ui.Context, fn func()) {
		ui.Column(c).Width(200).Children(fn)
	})
	if got[0].W != 200 || got[0].H != 100 || got[1] != got[0] || got[2] != got[0] {
		t.Fatalf("loaded, loading, error boxes %v", got)
	}
}

func TestImageStatesShareBoxInSizedParent(t *testing.T) {
	got := imageStateBoxes(t, ImageOptions{}, func(c *ui.Context, fn func()) {
		ui.Column(c).Size(120, 60).Children(fn)
	})
	if got[0].W != 120 || got[1] != got[0] || got[2] != got[0] {
		t.Fatalf("loaded, loading, error boxes %v", got)
	}
}

func TestImageStatesShareExplicitSize(t *testing.T) {
	got := imageStateBoxes(t, ImageOptions{Width: 120, Height: 60}, func(c *ui.Context, fn func()) {
		ui.Column(c).Width(300).Children(fn)
	})
	if got[0].W != 120 || got[0].H != 60 || got[1] != got[0] || got[2] != got[0] {
		t.Fatalf("loaded, loading, error boxes %v", got)
	}
}

func TestImageStatesShareBoxInNarrowRow(t *testing.T) {
	got := imageStateBoxes(t, ImageOptions{}, func(c *ui.Context, fn func()) {
		ui.Row(c).Width(40).AlignItems(ui.Start).Children(fn)
	})
	if got[0].W != 40 || got[0].H != 20 || got[1] != got[0] || got[2] != got[0] {
		t.Fatalf("loaded, loading, error boxes %v", got)
	}
}

func TestImageRetryUnclippedInSmallBox(t *testing.T) {
	retries := 0
	var box ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		ui.Row(c).Padding(40).AlignItems(ui.Start).Children(func() {
			ui.Row(c).Width(40).AlignItems(ui.Start).Children(func() {
				v := Image(c, imageTestBitmap(), ImageOptions{Err: errors.New("timeout")})
				box = v.Element.Bounds()
				if v.Retry {
					retries++
				}
			})
		})
	}, 300, 200)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if box.W != 40 || box.H != 20 {
		t.Fatalf("box %v", box)
	}
	r, ok := tt.Find("Retry")
	if !ok || r.W != 28 || r.H != 28 {
		t.Fatalf("retry bounds %v", r)
	}
	tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	if retries != 1 {
		t.Fatalf("retries %d, want 1", retries)
	}
	cx := int(r.X + r.W/2)
	ring := func() (above, below ui.Color) {
		img := tt.Image()
		a := img.RGBAAt(cx, int(r.Y-theme.FocusGap-theme.FocusWidth/2))
		b := img.RGBAAt(cx, int(r.Y+r.H+theme.FocusGap+theme.FocusWidth/2))
		return ui.Color{R: a.R, G: a.G, B: a.B, A: a.A}, ui.Color{R: b.R, G: b.G, B: b.B, A: b.A}
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Retry") {
		t.Fatal("retry not focused")
	}
	above, below := ring()
	want := theme.Light().Focus
	if above != want || below != want {
		t.Fatalf("ring above %v below %v, want %v", above, below, want)
	}
	if r.Y-theme.FocusGap >= box.Y || r.Y+r.H+theme.FocusGap <= box.Y+box.H {
		t.Fatalf("ring rows not outside box %v (retry %v)", box, r)
	}
}
