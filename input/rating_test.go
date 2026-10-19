package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

func ratingTester(v *float64, opts RatingOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).AlignItems(ui.Start).Children(func() {
			if Rating(c, v, opts).Changed() {
				*changes++
			}
		})
	}, 300, 100)
}

func TestRatingKeyboardSteps(t *testing.T) {
	v, n := 2.0, 0
	tt := ratingTester(&v, RatingOptions{Step: 0.5, Label: "Taste"}, &n)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Taste") {
		t.Fatal("rating not focusable")
	}
	tt.Key(0, ui.KeyRight)
	if v != 2.5 {
		t.Fatalf("right: %v", v)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	if v != 1.5 {
		t.Fatalf("down twice: %v", v)
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyRight)
	if v != 5 {
		t.Fatalf("end then right: %v", v)
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyLeft)
	if v != 0 || n != 5 {
		t.Fatalf("home then left: %v, %d changes", v, n)
	}
}

func TestRatingPointer(t *testing.T) {
	v, n := 0.0, 0
	tt := ratingTester(&v, RatingOptions{Step: 0.5, Label: "Taste"}, &n)
	size := theme.MinHitSize
	// The left half of the third star.
	tt.ClickAt(10+2*size+size*0.25, 10+size/2)
	if v != 2.5 || n != 1 {
		t.Fatalf("click: %v, %d changes", v, n)
	}
	tt.ClickAt(10+2*size+size*0.25, 10+size/2)
	if n != 1 {
		t.Fatalf("same value reported a change")
	}
}

func TestRatingReadOnlyAndDisabled(t *testing.T) {
	for _, opts := range []RatingOptions{{ReadOnly: true, Label: "Score"}, {Disabled: true, Label: "Score"}} {
		v, n := 3.0, 0
		tt := ratingTester(&v, opts, &n)
		tt.ClickAt(12, 20)
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyRight)
		if v != 3 || n != 0 {
			t.Fatalf("%+v: value %v changes %d", opts, v, n)
		}
		if opts.ReadOnly != tt.Focused("Score") {
			t.Fatalf("%+v: focus %v", opts, tt.Focused("Score"))
		}
	}
}

func TestRatingShowsValueAndPanics(t *testing.T) {
	v, n := 3.5, 0
	tt := ratingTester(&v, RatingOptions{Step: 0.5, ShowValue: true}, &n)
	if !tt.HasText("3.5 / 5") {
		t.Fatalf("texts %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("Step 0.3 accepted")
		}
	}()
	ratingTester(&v, RatingOptions{Step: 0.3}, &n)
}

func TestRatingSnapshot(t *testing.T) {
	a, b, d := 3.5, 4.0, 2.0
	snapshot(t, "rating", 360, 170, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(14).AlignItems(ui.Start).Children(func() {
			Rating(c, &a, RatingOptions{Step: 0.5, ShowValue: true})
			Rating(c, &b, RatingOptions{ReadOnly: true, Max: 7, ShowValue: true})
			Rating(c, &d, RatingOptions{Disabled: true})
		})
	})
}

func TestRatingHitAreaCompact(t *testing.T) {
	v := 2.0
	var b ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Density: theme.Compact})
		ui.Column(c).AlignItems(ui.Start).Children(func() {
			b = Rating(c, &v, RatingOptions{}).Element.Bounds()
		})
	}, 300, 100)
	tt.Frame()
	if b.W < 5*28 || b.H < 28 {
		t.Fatalf("stars %v", b)
	}
}
