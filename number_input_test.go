package mujica

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func numberInputTester(v *float64, changes *int, opts NumberInputOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if NumberInput(c, v, opts).Changed() {
			*changes++
		}
		Button(c, "Elsewhere", ButtonOptions{})
	}, 400, 300)
}

func TestNumberInputStepsAndBounds(t *testing.T) {
	v, n := 9.0, 0
	tt := numberInputTester(&v, &n, NumberInputOptions{Min: 0, Max: 10, Label: "Copies"})
	_ = tt.Click("Copies")
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyUp)
	if v != 10 || n != 1 {
		t.Fatalf("value %v changes %d, want 10 and 1", v, n)
	}
	_ = tt.Click("Decrease")
	_ = tt.Click("Decrease")
	if v != 8 || n != 3 {
		t.Fatalf("after buttons %v changes %d", v, n)
	}
	_ = tt.Click("Increase")
	_ = tt.Click("Increase")
	_ = tt.Click("Increase")
	if v != 10 {
		t.Fatalf("increase passed Max: %v", v)
	}
}

func TestNumberInputTypingPrecisionAndErrors(t *testing.T) {
	v, n := 1.0, 0
	tt := numberInputTester(&v, &n, NumberInputOptions{Min: 0, Max: 2, Step: 0.25, Label: "Ratio"})
	_ = tt.Click("Ratio")
	tt.Command("selectAll")
	tt.Type("1.5")
	if v != 1.5 {
		t.Fatalf("typed value %v", v)
	}
	tt.Key(0, ui.KeyDown)
	if v != 1.25 {
		t.Fatalf("Down from 1.5 by 0.25 gave %v", v)
	}
	tt.Command("selectAll")
	tt.Type("abc")
	if v != 1.25 || !tt.HasText("Enter a number") {
		t.Fatalf("invalid text changed value to %v or showed no error: %q", v, tt.Texts())
	}
	tt.Command("selectAll")
	tt.Type("7")
	if v != 1.25 || !tt.HasText("Enter a number from 0.00 to 2.00") {
		t.Fatalf("out of range: value %v texts %q", v, tt.Texts())
	}
	_ = tt.Click("Elsewhere")
	if tt.HasText("Enter a number") {
		t.Fatal("error stays after leaving the field, which shows the value again")
	}
}

func TestNumberInputEmpty(t *testing.T) {
	v, n := 3.0, 0
	tt := numberInputTester(&v, &n, NumberInputOptions{AllowEmpty: true, Label: "Limit"})
	_ = tt.Click("Limit")
	tt.Command("selectAll")
	tt.Key(0, ui.KeyBackspace)
	if !math.IsNaN(v) {
		t.Fatalf("emptied value %v, want NaN", v)
	}
	tt.Key(0, ui.KeyUp)
	if v != 1 {
		t.Fatalf("Up from empty gave %v, want 1", v)
	}
}

func TestNumberInputStepperOnly(t *testing.T) {
	v, n := 2.0, 0
	tt := numberInputTester(&v, &n, NumberInputOptions{Min: 1, Max: 3, StepperOnly: true, Label: "Guests"})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyUp)
	if v != 3 || n != 1 || !tt.HasText("3") {
		t.Fatalf("value %v changes %d texts %q", v, n, tt.Texts())
	}
	tt.Key(0, ui.KeyHome)
	if v != 1 {
		t.Fatalf("Home gave %v", v)
	}
}

func TestNumberInputDisabled(t *testing.T) {
	for _, stepper := range []bool{false, true} {
		v, n := 2.0, 0
		tt := numberInputTester(&v, &n, NumberInputOptions{Min: 1, Max: 3, StepperOnly: stepper, Disabled: true, Label: "Off"})
		_ = tt.Click("Increase")
		_ = tt.Click("Off")
		tt.Key(0, ui.KeyUp)
		if v != 2 || n != 0 {
			t.Fatalf("stepper %v: disabled changed to %v (%d changes)", stepper, v, n)
		}
	}
}

func TestNumberInputBadOptionsPanic(t *testing.T) {
	for _, opts := range []NumberInputOptions{{Min: 3, Max: 1}, {StepperOnly: true}, {Step: -1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: no panic", opts)
				}
			}()
			v := 0.0
			ui.NewTester(func(c *ui.Context) {
				Use(c, Settings{})
				NumberInput(c, &v, opts)
			}, 200, 100)
		}()
	}
}

func TestNumberInputSnapshot(t *testing.T) {
	snapshot(t, "number_input", 420, 260, func(c *ui.Context) {
		a, b, e, s := 12.0, math.NaN(), 7.5, 4.0
		ui.Column(c).Padding(24).Gap(12).Width(260).Children(func() {
			NumberInput(c, &a, NumberInputOptions{Min: 0, Max: 99, Label: "A"})
			NumberInput(c, &b, NumberInputOptions{AllowEmpty: true, Placeholder: "No limit", Label: "B"})
			NumberInput(c, &e, NumberInputOptions{Step: 0.5, Label: "E", Error: "Must be whole"})
			NumberInput(c, &s, NumberInputOptions{Min: 1, Max: 8, StepperOnly: true, Label: "S"})
		})
	})
}
