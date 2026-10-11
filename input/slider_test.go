package input

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func sliderTester(value *float64, opts SliderOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(0, 0).Width(216).Children(func() {
			if Slider(c, value, opts).Changed() {
				*changes++
			}
		})
	}, 300, 120)
}

func TestSliderKeysMoveOneStep(t *testing.T) {
	v, changes := 40.0, 0
	tt := sliderTester(&v, SliderOptions{Step: 10, Label: "Volume"}, &changes)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if v != 50 {
		t.Fatalf("Right: %v, want 50", v)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	if v != 30 {
		t.Fatalf("Down twice: %v, want 30", v)
	}
	tt.Key(0, ui.KeyEnd)
	if v != 100 {
		t.Fatalf("End: %v", v)
	}
	tt.Key(0, ui.KeyRight)
	if v != 100 || changes != 4 {
		t.Fatalf("past the end: %v, changes %d", v, changes)
	}
	tt.Key(0, ui.KeyHome)
	if v != 0 {
		t.Fatalf("Home: %v", v)
	}
}

func TestSliderDragSnaps(t *testing.T) {
	v, changes := 0.0, 0
	tt := sliderTester(&v, SliderOptions{Min: 0, Max: 10, Step: 2}, &changes)
	// The track runs from x=8 to x=208: 200 DIPs for 0 to 10.
	tt.Press(8, 18)
	tt.Move(8+200*0.55, 18)
	if v != 6 {
		t.Fatalf("drag to 5.5: %v, want 6", v)
	}
	tt.Move(400, 18)
	tt.Release(400, 18)
	if v != 10 {
		t.Fatalf("drag past the end: %v", v)
	}
}

func TestSliderClampsOutOfRangeValue(t *testing.T) {
	v, changes := 250.0, 0
	tt := sliderTester(&v, SliderOptions{}, &changes)
	tt.Frame()
	if v != 100 || changes != 0 {
		t.Fatalf("value %v changes %d", v, changes)
	}
}

func TestSliderBadRange(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	v := 0.0
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Slider(c, &v, SliderOptions{Min: 5, Max: 1})
	}, 100, 100, 1)
}

func TestSliderSnapshot(t *testing.T) {
	snapshot(t, "slider", 420, 230, func(c *ui.Context) {
		a, b, d := 30.0, 60.0, 50.0
		marks := []SliderMark{{0, "0"}, {25, "25"}, {50, "50"}, {75, "75"}, {100, "100"}}
		ui.Column(c).Padding(20).Gap(8).Width(380).Children(func() {
			Slider(c, &a, SliderOptions{Step: 5})
			s := Slider(c, &b, SliderOptions{Step: 25, Marks: marks, ShowValue: true, Format: func(v float64) string { return fmt.Sprintf("%.0f%%", v) }})
			s.Element.Focus()
			Slider(c, &d, SliderOptions{Disabled: true})
		})
	})
}
