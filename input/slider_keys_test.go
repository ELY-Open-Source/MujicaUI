package input

import (
	"bytes"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestSliderFineStepNearEnds(t *testing.T) {
	v, changes := 99.2, 0
	tt := sliderTester(&v, SliderOptions{Step: 0.1, Label: "Gain"}, &changes)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if v != 99.3 {
		t.Fatalf("Right from 99.2: %v", v)
	}
	tt.Key(0, ui.KeyEnd)
	if v != 100 {
		t.Fatalf("End: %v", v)
	}
	v = 0.5
	tt.Frame()
	tt.Key(0, ui.KeyLeft)
	if v != 0.4 {
		t.Fatalf("Left from 0.5: %v", v)
	}
	tt.Key(0, ui.KeyHome)
	if v != 0 {
		t.Fatalf("Home: %v", v)
	}
}

func TestRangeSliderFineStepNearEnds(t *testing.T) {
	low, high, changes := 0.3, 99.8, 0
	tt := rangeTester(&low, &high, RangeSliderOptions{Step: 0.1, Label: "Band"}, &changes)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyLeft)
	if low != 0.2 {
		t.Fatalf("Left from 0.3: %v", low)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if high != 99.9 {
		t.Fatalf("Right from 99.8: %v", high)
	}
}

func TestSliderFormatKeepsStepDecimals(t *testing.T) {
	for _, tc := range []struct {
		g    SliderGeom
		v    float64
		want string
	}{
		{SliderGeom{lo: 0, hi: 1, step: 0.25}, 0.75, "0.75"},
		{SliderGeom{lo: 0.5, hi: 10, step: 1}, 2.5, "2.5"},
		{SliderGeom{lo: 0, hi: 100, step: 5}, 35, "35"},
	} {
		if got := sliderFormat(nil, tc.g)(tc.v); got != tc.want {
			t.Errorf("%+v: %q, want %q", tc.g, got, tc.want)
		}
	}
}

func TestSliderMarksFollowTextScale(t *testing.T) {
	render := func(scale float32) []byte {
		v, low, high := 50.0, 20.0, 80.0
		marks := []SliderMark{{Value: 0, Label: "Low"}, {Value: 100, Label: "High"}}
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			ui.Column(c).Width(300).Children(func() {
				Slider(c, &v, SliderOptions{Marks: marks})
				RangeSlider(c, &low, &high, RangeSliderOptions{Marks: marks})
			})
		}, 320, 200)
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: scale})
		tt.Frame()
		return encode(t, tt.Image())
	}
	if bytes.Equal(render(1), render(2)) {
		t.Fatal("mark captions ignore the text scale")
	}
}
