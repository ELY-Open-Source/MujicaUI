package input

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// baselineAt returns the center of the element whose text is s.
func baselineAt(t *testing.T, tt *ui.Tester, s string) (float32, float32) {
	t.Helper()
	r, ok := tt.Find(s)
	if !ok {
		t.Fatalf("%q not found in %q", s, tt.Texts())
	}
	return r.X + r.W/2, r.Y + r.H/2
}

func TestBaselineButtonStates(t *testing.T) {
	view := func(c *ui.Context) {
		ui.Row(c).Padding(24).Gap(12).Children(func() {
			Button(c, "Save", ButtonOptions{})
			Button(c, "Cancel", ButtonOptions{Variant: Secondary})
		})
	}
	snapshotAfter(t, "button-hover", 300, 90, view, func(tt *ui.Tester) { tt.Move(baselineAt(t, tt, "Save")) })
	snapshotAfter(t, "button-pressed", 300, 90, view, func(tt *ui.Tester) { tt.Press(baselineAt(t, tt, "Save")) })
	snapshotAfter(t, "button-focus", 300, 90, view, func(tt *ui.Tester) { tt.Key(0, ui.KeyTab) })
}

func TestBaselineTextInputStates(t *testing.T) {
	name, date := "Uika", "1789-13-01"
	view := func(c *ui.Context) {
		ui.Column(c).Padding(24).Gap(12).Width(360).Children(func() {
			TextInput(c, &name, TextInputOptions{Label: "Name"})
			TextInput(c, &date, TextInputOptions{Label: "Date", Error: "No such date"})
		})
	}
	snapshotAfter(t, "text_input-focus", 400, 170, view, func(tt *ui.Tester) { tt.Key(0, ui.KeyTab) })
	snapshotAfter(t, "text_input-error-focus", 400, 170, view, func(tt *ui.Tester) {
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyTab)
	})
}

func TestBaselineCheckboxStates(t *testing.T) {
	on, off := true, false
	view := func(c *ui.Context) {
		ui.Column(c).Padding(24).Gap(12).Children(func() {
			Checkbox(c, &on, "Sealed", CheckboxOptions{})
			Checkbox(c, &off, "Open", CheckboxOptions{})
		})
	}
	snapshotAfter(t, "checkbox-hover", 240, 110, view, func(tt *ui.Tester) { tt.Move(baselineAt(t, tt, "Open")) })
	snapshotAfter(t, "checkbox-focus", 240, 110, view, func(tt *ui.Tester) { tt.Key(0, ui.KeyTab) })
	snapshotAfter(t, "checkbox-pressed", 240, 110, view, func(tt *ui.Tester) { tt.Press(baselineAt(t, tt, "Sealed")) })
}
