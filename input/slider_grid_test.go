package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

// sliderRecord builds a Slider and a RangeSlider and records the values each Changed reported.
type sliderRecord struct {
	v, low, high      float64
	single, lows, his []float64
}

func sliderGridTester(r *sliderRecord, min, max, step float64) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(216).Gap(10).Children(func() {
			if Slider(c, &r.v, SliderOptions{Min: min, Max: max, Step: step, Label: "One"}).Changed() {
				r.single = append(r.single, r.v)
			}
			if RangeSlider(c, &r.low, &r.high, RangeSliderOptions{Min: min, Max: max, Step: step, Label: "Two"}).Changed() {
				r.lows, r.his = append(r.lows, r.low), append(r.his, r.high)
			}
		})
	}, 300, 300)
}

func TestSliderEndSnapsToStepGrid(t *testing.T) {
	r := &sliderRecord{v: 3, low: 0, high: 3}
	tt := sliderGridTester(r, 0, 10, 3)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if len(r.single) != 1 || r.single[0] != 9 || r.v != 9 {
		t.Fatalf("Slider End reported %v, value %v", r.single, r.v)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if len(r.his) != 1 || r.his[0] != 9 || r.high != 9 {
		t.Fatalf("RangeSlider End reported %v, high %v", r.his, r.high)
	}
}

func TestSliderTinyStep(t *testing.T) {
	r := &sliderRecord{v: 0, low: 0, high: 1e-9}
	tt := sliderGridTester(r, 0, 1e-9, 1e-10)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if r.v != 1e-10 {
		t.Fatalf("Slider Right: %v", r.v)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if r.low != 1e-10 {
		t.Fatalf("RangeSlider Right: %v", r.low)
	}
}

func TestSliderKeepsPreciseMin(t *testing.T) {
	const lo = 0.1234567891
	r := &sliderRecord{v: 0.5, low: 0.5, high: 0.6}
	tt := sliderGridTester(r, lo, 1, 0.1)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	if r.v < lo || r.v != lo {
		t.Fatalf("Slider Home: %v, Min %v", r.v, lo)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	if r.low < lo || r.low != lo {
		t.Fatalf("RangeSlider Home: %v, Min %v", r.low, lo)
	}
}
