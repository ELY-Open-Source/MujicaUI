package core

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// TestBasesBlockedByDisabledWrapper flips a focused MyGo base's wrapper to Disabled in the frame Space arrives.
func TestBasesBlockedByDisabledWrapper(t *testing.T) {
	for name, base := range map[string]func(c *ui.Context, on *bool) ui.Element{
		"CheckboxBase": func(c *ui.Context, on *bool) ui.Element { return ui.CheckboxBase(c, on) },
		"RadioBase":    func(c *ui.Context, on *bool) ui.Element { return ui.RadioBase(c, on, true) },
		"SwitchBase":   func(c *ui.Context, on *bool) ui.Element { return ui.SwitchBase(c, on) },
		"ToggleBase":   func(c *ui.Context, on *bool) ui.Element { return ui.ToggleBase(c, on) },
	} {
		for _, disable := range []bool{false, true} {
			on, off := false, false
			tt := ui.NewTester(func(c *ui.Context) {
				ui.Row(c).Disabled(off).Children(func() {
					base(c, &on).Label("Target").Size(40, 20)
				})
			}, 200, 100)
			tt.Key(0, ui.KeyTab)
			if !tt.Focused("Target") {
				t.Fatalf("%s: not focused", name)
			}
			off = disable
			tt.Key(0, ui.KeySpace)
			tt.Frame()
			if on == disable {
				t.Errorf("%s: disabled=%v, Space left the value %v", name, disable, on)
			}
		}
	}
}

// TestComboboxBaseBlockedByDisabledWrapper guards Combobox, which binds its base directly: a focused one takes no typing.
func TestComboboxBaseBlockedByDisabledWrapper(t *testing.T) {
	q, off := "", false
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Disabled(off).Children(func() {
			ui.ComboboxBase(c, &q).Input.Label("Target").Width(120)
		})
	}, 200, 100)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Target") {
		t.Fatal("not focused")
	}
	off = true
	tt.Frame()
	tt.Type("a")
	tt.Frame()
	if q != "" {
		t.Fatalf("typing under a disabled wrapper set %q; Combobox must edit a copy again", q)
	}
}
