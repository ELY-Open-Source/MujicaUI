package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func rangeTester(low, high *float64, opts RangeSliderOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(216).Children(func() {
			if RangeSlider(c, low, high, opts).Changed() {
				*changes++
			}
		})
	}, 300, 120)
}

func TestRangeSliderKeysPerThumb(t *testing.T) {
	low, high, changes := 20.0, 80.0, 0
	tt := rangeTester(&low, &high, RangeSliderOptions{Step: 10, Label: "Price"}, &changes)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Price minimum") {
		t.Fatal("first Tab did not focus the low thumb")
	}
	tt.Key(0, ui.KeyRight)
	if low != 30 {
		t.Fatalf("low %v", low)
	}
	tt.Key(0, ui.KeyEnd)
	if low != 80 || high != 80 {
		t.Fatalf("End took low to %v, high %v", low, high)
	}
	tt.Key(0, ui.KeyRight)
	if low != 80 {
		t.Fatalf("low crossed high: %v", low)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Price maximum") {
		t.Fatal("second Tab did not focus the high thumb")
	}
	tt.Key(0, ui.KeyLeft)
	if high != 80 {
		t.Fatalf("high crossed low: %v", high)
	}
	tt.Key(0, ui.KeyEnd)
	if high != 100 || changes != 3 {
		t.Fatalf("high %v changes %d", high, changes)
	}
}

func TestRangeSliderDragMovesNearerThumb(t *testing.T) {
	low, high, changes := 0.0, 100.0, 0
	tt := rangeTester(&low, &high, RangeSliderOptions{Step: 1}, &changes)
	// The track runs from x=8 to x=208.
	tt.Press(8+200*0.9, 18)
	tt.Move(8+200*0.7, 18)
	tt.Release(8+200*0.7, 18)
	if low != 0 || high != 70 {
		t.Fatalf("low %v high %v", low, high)
	}
	tt.Press(8+200*0.1, 18)
	tt.Move(8+200*0.95, 18)
	tt.Release(8+200*0.95, 18)
	if low != 70 || high != 70 {
		t.Fatalf("dragged low past high: low %v high %v", low, high)
	}
}

func TestRangeSliderOrdersSwappedEnds(t *testing.T) {
	low, high, changes := 60.0, 10.0, 0
	tt := rangeTester(&low, &high, RangeSliderOptions{}, &changes)
	tt.Frame()
	if low != 10 || high != 60 {
		t.Fatalf("low %v high %v", low, high)
	}
}

func TestRangeSliderSnapshot(t *testing.T) {
	snapshot(t, "range_slider", 420, 170, func(c *ui.Context) {
		a, b := 20.0, 70.0
		d, e := 30.0, 60.0
		marks := []SliderMark{{0, "0"}, {50, "50"}, {100, "100"}}
		ui.Column(c).Padding(20).Gap(8).Width(380).Children(func() {
			RangeSlider(c, &a, &b, RangeSliderOptions{Step: 5, Marks: marks, ShowValue: true, Label: "Price"})
			RangeSlider(c, &d, &e, RangeSliderOptions{Disabled: true})
		})
	})
}
