package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

var radioGroupTestItems = []RadioItem[string]{
	{Value: "s", Label: "Small"},
	{Value: "m", Label: "Medium", Disabled: true},
	{Value: "l", Label: "Large"},
}

func radioGroupTester(v *string, changes *int, opts RadioGroupOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if RadioGroup(c, v, radioGroupTestItems, opts).Changed() {
			*changes++
		}
		Button(c, "After", ButtonOptions{})
	}, 400, 300)
}

func TestRadioGroupClickAndArrows(t *testing.T) {
	v, n := "s", 0
	tt := radioGroupTester(&v, &n, RadioGroupOptions{Label: "Size"})
	_ = tt.Click("Medium")
	if v != "s" || n != 0 {
		t.Fatalf("disabled item chosen: %q", v)
	}
	_ = tt.Click("Large")
	if v != "l" || n != 1 {
		t.Fatalf("click: %q changes %d", v, n)
	}
	tt.Key(0, ui.KeyDown)
	if v != "s" || n != 2 {
		t.Fatalf("Down wraps past the disabled item to Small, got %q changes %d", v, n)
	}
	tt.Key(0, ui.KeyUp)
	if v != "l" {
		t.Fatalf("Up skipped to %q", v)
	}
	tt.Key(0, ui.KeyHome)
	if v != "s" {
		t.Fatalf("Home chose %q", v)
	}
}

func TestRadioGroupTabStopsAtChosen(t *testing.T) {
	v, n := "l", 0
	tt := radioGroupTester(&v, &n, RadioGroupOptions{Horizontal: true})
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Large") {
		t.Fatal("Tab did not go to the chosen radio")
	}
	tt.Key(0, ui.KeyLeft)
	if v != "s" {
		t.Fatalf("Left in a row chose %q", v)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("Tab did not leave the group")
	}
}

func TestRadioGroupDisabledGroup(t *testing.T) {
	v, n := "s", 0
	tt := radioGroupTester(&v, &n, RadioGroupOptions{Disabled: true})
	_ = tt.Click("Large")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if v != "s" || n != 0 {
		t.Fatalf("disabled group chose %q", v)
	}
}

func TestRadioGroupDuplicateValuePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for duplicate values")
		}
	}()
	v := 1
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		RadioGroup(c, &v, []RadioItem[int]{{Value: 1, Label: "A"}, {Value: 1, Label: "B"}}, RadioGroupOptions{})
	}, 200, 100)
}

func TestRadioGroupSnapshot(t *testing.T) {
	snapshot(t, "radio_group", 460, 220, func(c *ui.Context) {
		v, h := "m", 2
		ui.Column(c).Padding(24).Gap(16).Children(func() {
			RadioGroup(c, &v, []RadioItem[string]{{Value: "s", Label: "Sonnet"}, {Value: "m", Label: "Madrigal"}, {Value: "o", Label: "Ode", Disabled: true}}, RadioGroupOptions{})
			RadioGroup(c, &h, []RadioItem[int]{{Value: 1, Label: "Matins"}, {Value: 2, Label: "Vespers"}, {Value: 3, Label: "Compline"}}, RadioGroupOptions{Horizontal: true})
		})
	})
}
