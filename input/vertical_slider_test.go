package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func verticalTester(v *float64, opts *VerticalSliderOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Padding(20).Children(func() {
			if VerticalSlider(c, v, *opts).Changed() {
				*n++
			}
		})
	}, 300, 300)
}

func TestVerticalSliderKeysAndPointer(t *testing.T) {
	v := 0.0
	opts := VerticalSliderOptions{Step: 10, Label: "Level", Height: 216}
	n := 0
	tt := verticalTester(&v, &opts, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyRight)
	if v != 20 || n != 2 {
		t.Fatalf("Up, Right: %v %d", v, n)
	}
	tt.Key(0, ui.KeyDown)
	if v != 10 {
		t.Fatalf("Down: %v", v)
	}
	tt.Key(0, ui.KeyEnd)
	if v != 100 {
		t.Fatalf("End: %v", v)
	}
	r, _ := tt.Find("Level")
	// Track runs from y+8 (100) to y+208 (0): a quarter up is 25, snapped to 30 or 20.
	tt.Press(r.X+r.W/2, r.Y+8+150)
	tt.Release(r.X+r.W/2, r.Y+8+150)
	if v != 30 && v != 20 {
		t.Fatalf("press at a quarter height: %v", v)
	}
	tt.Press(r.X+r.W/2, r.Y+8+150)
	tt.Move(r.X+r.W/2, r.Y)
	tt.Release(r.X+r.W/2, r.Y)
	if v != 100 {
		t.Fatalf("drag to the top: %v", v)
	}
	opts.Disabled = true
	tt.Frame()
	tt.Press(r.X+r.W/2, r.Y+200)
	tt.Release(r.X+r.W/2, r.Y+200)
	if v != 100 {
		t.Fatalf("disabled slider moved: %v", v)
	}
}

func TestVerticalSliderRejects(t *testing.T) {
	v := 0.0
	for name, o := range map[string]VerticalSliderOptions{"height": {Height: -1}, "step": {Step: -1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); VerticalSlider(c, &v, o) }, 200, 200)
		}()
	}
}

func TestVerticalSliderSnapshot(t *testing.T) {
	a, b := 40.0, 70.0
	snapshot(t, "vertical_slider", 260, 220, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(24).Children(func() {
			VerticalSlider(c, &a, VerticalSliderOptions{Marks: []SliderMark{{Value: 0, Label: "0"}, {Value: 50, Label: "50"}, {Value: 100, Label: "100"}}})
			VerticalSlider(c, &b, VerticalSliderOptions{Disabled: true})
		})
	})
}
