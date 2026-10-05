package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// choiceBetween builds control between two buttons, for Tab order checks.
func choiceBetween(control func(c *ui.Context)) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).Gap(10).Width(300).Children(func() {
			Button(c, "Before", ButtonOptions{})
			control(c)
			Button(c, "After", ButtonOptions{})
		})
	}, 400, 400)
}

// choiceTabSkips reports whether Tab from Before lands on After.
func choiceTabSkips(t *testing.T, tt *ui.Tester) bool {
	t.Helper()
	if err := tt.Click("Before"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	return tt.Focused("After")
}

func TestToggleDisabled(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		on := false
		tt := choiceBetween(func(c *ui.Context) {
			Toggle(c, &on, "Italic", ToggleOptions{Disabled: disabled})
		})
		_ = tt.Click("Italic")
		if on == disabled {
			t.Fatalf("disabled=%v: click left on=%v", disabled, on)
		}
		if choiceTabSkips(t, tt) != disabled {
			t.Fatalf("disabled=%v: Tab skip wrong", disabled)
		}
	}
}

func TestSliderDisabled(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		v := 30.0
		tt := choiceBetween(func(c *ui.Context) {
			Slider(c, &v, SliderOptions{Step: 1, Disabled: disabled, Label: "Volume"})
		})
		_ = tt.Click("Volume")
		if (v != 30) == disabled {
			t.Fatalf("disabled=%v: click left %v", disabled, v)
		}
		if choiceTabSkips(t, tt) != disabled {
			t.Fatalf("disabled=%v: Tab skip wrong", disabled)
		}
	}
}

func TestRangeSliderDisabled(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		low, high := 20.0, 40.0
		tt := choiceBetween(func(c *ui.Context) {
			RangeSlider(c, &low, &high, RangeSliderOptions{Step: 1, Disabled: disabled, Label: "Price"})
		})
		_ = tt.Click("Price")
		if (high != 40) == disabled || low != 20 {
			t.Fatalf("disabled=%v: click left %v..%v", disabled, low, high)
		}
		if choiceTabSkips(t, tt) != disabled {
			t.Fatalf("disabled=%v: Tab skip wrong", disabled)
		}
	}
}

// Disabling a focused control in the frame a key arrives must block the key.
func TestChoiceDisabledWhileFocused(t *testing.T) {
	disabled := false
	on, pressed, group := false, false, []string{"left"}
	items := []ToggleItem[string]{{Value: "left", Label: "Left"}, {Value: "right", Label: "Right"}}
	for name, build := range map[string]func(c *ui.Context){
		"switch": func(c *ui.Context) { Switch(c, &on, "Sync", SwitchOptions{Disabled: disabled}) },
		"toggle": func(c *ui.Context) { Toggle(c, &pressed, "Bold", ToggleOptions{Disabled: disabled}) },
		"group": func(c *ui.Context) {
			items[1].Disabled = disabled
			ToggleGroup(c, &group, items, ToggleGroupOptions{})
		},
	} {
		disabled, on, pressed, group = false, false, false, []string{"left"}
		tt := choiceBetween(build)
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyTab)
		if name == "group" {
			tt.Key(0, ui.KeyRight)
		}
		disabled = true
		tt.Key(0, ui.KeySpace)
		if on || pressed || group[0] != "left" {
			t.Errorf("%s: Space changed a disabled control: %v %v %v", name, on, pressed, group)
		}
	}
}
