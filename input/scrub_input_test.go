package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func scrubTester(v *float64, opts *ScrubInputOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Padding(20).Children(func() {
			if ScrubInput(c, "Opacity", v, *opts).Changed() {
				*n++
			}
		})
	}, 500, 200)
}

func TestScrubInputDragAndRange(t *testing.T) {
	v := 0.5
	opts := ScrubInputOptions{Min: 0, Max: 1, Step: 0.1}
	n := 0
	tt := scrubTester(&v, &opts, &n)
	r, _ := tt.Find("Opacity")
	x, y := r.X+r.W/2, r.Y+r.H/2
	tt.Press(x, y)
	tt.Move(x+6, y) // one step, 2 DIPs carried
	tt.Move(x+10, y)
	tt.Release(x+10, y)
	if v != 0.7 {
		t.Fatalf("10 DIPs right: %v", v)
	}
	tt.Press(x, y)
	tt.Move(x+400, y)
	tt.Release(x+400, y)
	if v != 1 {
		t.Fatalf("past the max: %v", v)
	}
	tt.Press(x, y)
	tt.Move(x-8, y)
	tt.Release(x-8, y)
	if v != 0.8 || n != 4 {
		t.Fatalf("8 DIPs left: %v (%d)", v, n)
	}
}

func TestScrubInputTypingAndDisabled(t *testing.T) {
	v := 2.0
	opts := ScrubInputOptions{Step: 1}
	n := 0
	tt := scrubTester(&v, &opts, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	if v != 3 {
		t.Fatalf("Up in the field: %v", v)
	}
	tt.Command("selectAll")
	tt.Type("12")
	if v != 12 {
		t.Fatalf("typed: %v", v)
	}
	opts.Disabled = true
	tt.Frame()
	r, _ := tt.Find("Opacity")
	tt.Press(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(r.X+r.W/2+40, r.Y+r.H/2)
	tt.Release(r.X+r.W/2+40, r.Y+r.H/2)
	if v != 12 {
		t.Fatalf("disabled scrub: %v", v)
	}
}

func TestScrubInputSnapshot(t *testing.T) {
	a, b := 0.6, 24.0
	snapshot(t, "scrub_input", 520, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			ScrubInput(c, "Opacity", &a, ScrubInputOptions{Max: 1, Step: 0.1})
			ScrubInput(c, "Size", &b, ScrubInputOptions{Disabled: true})
		})
	})
}
