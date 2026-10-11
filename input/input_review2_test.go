package input

import (
	"image/color"
	"runtime"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestInputsObeyDisabledAncestor(t *testing.T) {
	text, area, num, find, cur := "", "", 5.0, "a", 0
	off, events := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Gap(8).Disabled(off).Children(func() {
			for _, r := range []InputResult{
				TextInput(c, &text, TextInputOptions{Label: "Text"}),
				TextArea(c, &area, TextAreaOptions{Label: "Area", SubmitKey: ui.KeyEnter}),
				NumberInput(c, &num, NumberInputOptions{Label: "Num"}),
			} {
				if r.Changed() || r.Submitted() {
					events++
				}
			}
			if SearchInput(c, &find, SearchInputOptions{Label: "Find", Find: true, Matches: 3, Current: &cur}).Moved() {
				events++
			}
		})
	}, 400, 500)
	next := func() {
		if runtime.GOOS == "darwin" {
			tt.Key(ui.Cmd, ui.KeyG)
		} else {
			tt.Key(0, ui.KeyF3)
		}
	}
	for _, name := range []string{"Text", "Area", "Num"} {
		_ = tt.Click(name)
		events = 0
		off = true
		tt.Type("x")
		tt.Key(0, ui.KeyEnter)
		tt.Key(0, ui.KeyUp)
		next()
		if events != 0 {
			t.Errorf("%s: %d events under a disabled ancestor", name, events)
		}
		off = false
		tt.Frame()
	}
	if text != "" || area != "" || num != 5 || cur != 0 {
		t.Fatalf("values changed: %q %q %v current %d", text, area, num, cur)
	}
}

// buttonLabelInk counts pixels in the label's box that differ from its first pixel.
func buttonLabelInk(tt *ui.Tester, label string) int {
	r, _ := tt.Find(label)
	img := tt.Image()
	bg := img.At(int(r.X), int(r.Y))
	n := 0
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if !sameColor(img.At(x, y), bg) {
				n++
			}
		}
	}
	return n
}

func sameColor(a, b color.Color) bool {
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

func TestButtonLoadingKeepsLabelAndWidth(t *testing.T) {
	for _, d := range []theme.Density{theme.Compact, theme.Default, theme.Comfortable} {
		for _, icon := range []*ui.SVG{nil, icons.Must("check")} {
			loading := false
			var w float32
			tt := ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{Density: d})
				ui.Row(c).Children(func() {
					w = Button(c, "Save", ButtonOptions{Icon: icon, Loading: loading}).Bounds().W
				})
			}, 300, 200)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			tt.Frame()
			before := w
			loading = true
			tt.Frame()
			tt.Frame()
			if w != before {
				t.Errorf("%v icon %v: width %v -> %v", d, icon != nil, before, w)
			}
			// The spinner sits beside the label, never over it, and the label has ink.
			sp, _ := tt.Find("Loading")
			lb, _ := tt.Find("Save")
			if sp.X+sp.W > lb.X || buttonLabelInk(tt, "Save") == 0 {
				t.Errorf("%v icon %v: label hidden while loading (spinner %v, label %v)", d, icon != nil, sp, lb)
			}
		}
	}
}

func numberInputStepTester(v *float64, opts NumberInputOptions) *ui.Tester {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		NumberInput(c, v, opts)
	}, 400, 200)
	tt.Key(0, ui.KeyTab)
	return tt
}

func TestNumberInputStepGrid(t *testing.T) {
	for _, stepper := range []bool{false, true} {
		v := 0.2
		tt := numberInputStepTester(&v, NumberInputOptions{Min: 0, Max: 0.29999999, Step: 0.1, StepperOnly: stepper})
		tt.Key(0, ui.KeyUp)
		if v != 0.2 {
			t.Errorf("stepper %v: Up past Max 0.29999999 gave %v", stepper, v)
		}
		v = 0
		tt = numberInputStepTester(&v, NumberInputOptions{Min: -1, Max: 1, Step: 1e-10, StepperOnly: stepper})
		tt.Key(0, ui.KeyUp)
		if v != 1e-10 {
			t.Errorf("stepper %v: Up by 1e-10 from 0 gave %v", stepper, v)
		}
	}
	v := 0.0
	tt := numberInputStepTester(&v, NumberInputOptions{Min: 0, Max: 10, Step: 3, StepperOnly: true})
	var got []float64
	for _, k := range []ui.Key{ui.KeyUp, ui.KeyEnd, ui.KeyDown, ui.KeyHome, ui.KeyDown} {
		tt.Key(0, k)
		got = append(got, v)
	}
	want := []float64{3, 10, 9, 0, 0}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("stepper [0,10] by 3: got %v, want %v", got, want)
		}
	}
}

// radioSameName prints the same for every value, as fmt keys would see it.
type radioSameName struct{ id int }

func (radioSameName) String() string { return "same" }

func TestRadioGroupKeysByEquality(t *testing.T) {
	v := radioSameName{0}
	items := []RadioItem[radioSameName]{{Value: radioSameName{1}, Label: "One"}, {Value: radioSameName{2}, Label: "Two"}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		RadioGroup(c, &v, items, RadioGroupOptions{})
	}, 300, 200)
	_ = tt.Click("Two")
	if v.id != 2 {
		t.Fatalf("click chose %v", v)
	}
	tt.Key(0, ui.KeyUp)
	if v.id != 1 {
		t.Fatalf("Up chose %v", v)
	}
	items[0], items[1] = items[1], items[0]
	tt.Frame()
	tt.Key(0, ui.KeySpace)
	if v.id != 1 || !tt.Focused("One") {
		t.Fatalf("after reorder chose %v, One focused %v", v, tt.Focused("One"))
	}
}

func TestNumberInputStepExactBounds(t *testing.T) {
	for _, stepper := range []bool{false, true} {
		v := 1e12
		tt := numberInputStepTester(&v, NumberInputOptions{Min: 0, Max: 2e12, Step: 1, StepperOnly: stepper})
		tt.Key(0, ui.KeyUp)
		if v != 1e12+1 {
			t.Errorf("stepper %v: Up from 1e12 by 1 gave %.0f", stepper, v)
		}
		tt.Key(0, ui.KeyDown)
		tt.Key(0, ui.KeyDown)
		if v != 1e12-1 {
			t.Errorf("stepper %v: two Downs gave %.0f", stepper, v)
		}
		v = 0.2
		tt = numberInputStepTester(&v, NumberInputOptions{Min: 0, Max: 0.29999999995, Step: 0.1, StepperOnly: stepper})
		tt.Key(0, ui.KeyUp)
		if v != 0.2 {
			t.Errorf("stepper %v: Up past Max 0.29999999995 gave %v", stepper, v)
		}
		tt.Key(0, ui.KeyEnd)
		if stepper && v != 0.2 {
			t.Errorf("End went to %v, past the largest value at precision 1", v)
		}
	}
}
