package input

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// MyGo rounds its 1% move to 1e-9, which can flip the apparent direction; the key decides it.
func sliderDirectionRun(t *testing.T, key ui.Key, want float64) {
	t.Helper()
	r := &sliderRecord{v: 6e-10, low: 6e-10, high: 8e-10}
	tt := sliderGridTester(r, 0, 1e-8, 1e-10)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, key)
	if r.v != want || len(r.single) != 1 {
		t.Errorf("Slider: %v, reported %v; want %v", r.v, r.single, want)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, key)
	if r.low != want || len(r.lows) != 1 {
		t.Errorf("RangeSlider low: %v, reported %v; want %v", r.low, r.lows, want)
	}
}

func TestSliderLeftUnderMyGoRounding(t *testing.T) { sliderDirectionRun(t, ui.KeyLeft, 5e-10) }

func TestSliderRightUnderMyGoRounding(t *testing.T) { sliderDirectionRun(t, ui.KeyRight, 7e-10) }
