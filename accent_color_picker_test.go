package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

func TestAccentColorPickerMarksLowContrast(t *testing.T) {
	col, n := ui.Hex("#7D2034"), 0
	cands := []ui.Color{ui.Hex("#7D2034"), ui.Hex("#F0E68C")}
	tt := accountTester(400, 200, func(c *ui.Context) {
		if AccentColorPicker(c, &col, AccentColorPickerOptions{Candidates: cands}).Changed() {
			n++
		}
	})
	k := theme.Light()
	if !accentPickerReadable(k, cands[0]) || accentPickerReadable(k, cands[1]) {
		t.Fatal("contrast check wrong")
	}
	if !tt.HasText("Marked colors are too faint for text on them.") {
		t.Fatalf("no low-contrast note: %q", tt.Texts())
	}
	if err := tt.Click("#F0E68C"); err != nil {
		t.Fatal(err)
	}
	if col != cands[1] || n != 1 || !tt.HasText("#F0E68C has low contrast; text on it may be hard to read.") {
		t.Fatalf("pick: %v %d %q", col, n, tt.Texts())
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeySpace)
	if col != cands[0] || n != 2 {
		t.Fatalf("keys: %v", col)
	}
}

// Built-in accents meet their promised light-theme contrast.
func TestAccentColorPickerDefaults(t *testing.T) {
	low := 0
	for _, col := range accentPickerColors {
		if !accentPickerReadable(theme.Light(), col) {
			low++
		}
	}
	if low != 2 {
		t.Fatalf("%d defaults fail; the palette shows two faint ones on purpose", low)
	}
}

func TestAccentColorPickerSnapshot(t *testing.T) {
	snapshot(t, "accent_color_picker", 420, 120, func(c *ui.Context) {
		col := ui.Local(c.Root(), "col", func() ui.Color { return ui.Hex("#355E80") })
		ui.Column(c).Padding(16).Children(func() { AccentColorPicker(c, col, AccentColorPickerOptions{}) })
	})
}

func TestAccentColorPickerDisabledAndEmptyCandidates(t *testing.T) {
	col := ui.Hex("#7D2034")
	tt := accountTester(600, 300, func(c *ui.Context) {
		ui.Column(c).Disabled(true).Children(func() { AccentColorPicker(c, &col, AccentColorPickerOptions{}) })
	})
	if err := tt.Click("#355E80"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	if col != ui.Hex("#7D2034") {
		t.Fatal("disabled swatch changes")
	}
	accountPanics(t, func(c *ui.Context) { AccentColorPicker(c, &col, AccentColorPickerOptions{Candidates: []ui.Color{}}) })
}
