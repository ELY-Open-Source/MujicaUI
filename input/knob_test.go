package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func knobTester(v *float64, opts *KnobOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Padding(40).Children(func() {
			if Knob(c, v, *opts).Changed() {
				*n++
			}
		})
	}, 300, 300)
}

func TestKnobKeysAndStepGrid(t *testing.T) {
	v := 0.0
	opts := KnobOptions{Max: 1, Step: 0.1, Label: "Gain"}
	n := 0
	tt := knobTester(&v, &opts, &n)
	tt.Key(0, ui.KeyTab)
	for range 3 {
		tt.Key(0, ui.KeyUp)
	}
	if v != 0.3 || n != 3 || !tt.HasText("0.3") {
		t.Fatalf("three Ups: %v (%d) %q", v, n, tt.Texts())
	}
	tt.Key(0, ui.KeyLeft)
	if v != 0.2 {
		t.Fatalf("Left: %v", v)
	}
	tt.Key(0, ui.KeyEnd)
	if v != 1 {
		t.Fatalf("End: %v", v)
	}
	tt.Key(0, ui.KeyUp)
	if v != 1 || n != 5 {
		t.Fatalf("past the end: %v %d", v, n)
	}
	tt.Key(0, ui.KeyHome)
	if v != 0 {
		t.Fatalf("Home: %v", v)
	}
}

func TestKnobDragAndDisabled(t *testing.T) {
	v := 50.0
	opts := KnobOptions{Step: 1, Label: "Volume"}
	n := 0
	tt := knobTester(&v, &opts, &n)
	// The dial's center: padding 40, caption and gap above it.
	x, y := float32(68), float32(40+21+30)
	tt.Press(x, y)
	tt.Move(x, y-32) // a fifth of 160 DIPs: 20 of 100
	tt.Release(x, y-32)
	if v != 70 {
		t.Fatalf("drag up 32 DIPs: %v", v)
	}
	opts.Disabled = true
	tt.Frame()
	tt.Press(x, y)
	tt.Move(x, y-32)
	tt.Release(x, y-32)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	if v != 70 {
		t.Fatalf("disabled knob moved: %v", v)
	}
}

func TestKnobRejects(t *testing.T) {
	v := 0.0
	for name, o := range map[string]KnobOptions{"ticks": {Ticks: -1}, "size": {Size: -1}, "range": {Min: 5, Max: 1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); Knob(c, &v, o) }, 200, 200)
		}()
	}
}

func TestKnobSnapshot(t *testing.T) {
	a, b, d := 35.0, 0.8, 10.0
	snapshot(t, "knob", 320, 130, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(24).Children(func() {
			Knob(c, &a, KnobOptions{Label: "Gain", Ticks: 11})
			Knob(c, &b, KnobOptions{Max: 1, Step: 0.1, Label: "Mix", Size: 72})
			Knob(c, &d, KnobOptions{Label: "Off", Disabled: true})
		})
	})
}
