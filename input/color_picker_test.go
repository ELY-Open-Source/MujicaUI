package input

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

func colorPickerTester(col *ui.Color, opts ColorPickerOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).AlignItems(ui.Start).Children(func() {
			if ColorPicker(c, col, opts).Changed() {
				*changes++
			}
		})
	}, 400, 560)
}

func TestColorPickerPresetAndHex(t *testing.T) {
	col, n := ui.Hex("#7D2034"), 0
	tt := colorPickerTester(&col, ColorPickerOptions{Label: "Accent"}, &n)
	if err := tt.Click("Accent"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("#355E80"); err != nil {
		t.Fatal(err)
	}
	if col != ui.Hex("#355E80") || n != 1 {
		t.Fatalf("preset: %v changes %d", col, n)
	}
	if err := tt.Click("HEX"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("ff000080")
	if col != (ui.Color{R: 255, A: 128}) {
		t.Fatalf("hex: %v", col)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("zz")
	if col != (ui.Color{R: 255, A: 128}) || !tt.HasText(colorPickerInvalidEN) {
		t.Fatalf("invalid hex changed the color or showed no error: %v", col)
	}
}

const colorPickerInvalidEN = "Use HEX RGB, RRGGBB or RRGGBBAA; channels 0–255; opacity 0–100."

func TestColorPickerRGBAndAlphaFields(t *testing.T) {
	col, n := ui.RGB(0, 0, 0), 0
	tt := colorPickerTester(&col, ColorPickerOptions{Label: "Ink"}, &n)
	_ = tt.Click("Ink")
	_ = tt.Click("Green")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("200")
	_ = tt.Click("Opacity percent")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("50")
	if col != (ui.Color{G: 200, A: 128}) {
		t.Fatalf("rgb: %v", col)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("300")
	if col.A != 128 || !tt.HasText(colorPickerInvalidEN) {
		t.Fatalf("opacity 300 accepted: %v", col)
	}
}

func TestColorPickerKeyboardAndEscape(t *testing.T) {
	col, n := ui.Hex("#808080"), 0
	tt := colorPickerTester(&col, ColorPickerOptions{Label: "Shade"}, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !tt.Focused("Saturation and brightness") {
		t.Fatal("square not focused on open")
	}
	tt.Key(0, ui.KeyUp)
	if col == ui.Hex("#808080") || col.R != col.G {
		t.Fatalf("up arrow: %v", col)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Hue") {
		t.Fatal("hue slider not next")
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("HEX") || !tt.Focused("Shade") {
		t.Fatal("escape did not close and return focus")
	}
}

func TestColorPickerDisabledAndParse(t *testing.T) {
	col, n := ui.Hex("#7D2034"), 0
	tt := colorPickerTester(&col, ColorPickerOptions{Label: "Accent", Disabled: true}, &n)
	_ = tt.Click("Accent")
	if tt.HasText("HEX") {
		t.Fatal("disabled picker opened")
	}
	for in, want := range map[string]ui.Color{"#abc": ui.Hex("#aabbcc"), "7D2034": ui.Hex("#7D2034"), "#00000000": {}} {
		got, err := colorPickerParseHex(in)
		if err != nil || got != want {
			t.Errorf("%q: %v %v", in, got, err)
		}
	}
	for _, in := range []string{"", "#12", "#12345g", "#1234567"} {
		if _, err := colorPickerParseHex(in); err == nil {
			t.Errorf("%q parsed", in)
		}
	}
}

func TestColorPickerSnapshot(t *testing.T) {
	a, b := ui.Hex("#7D2034"), ui.Hex("#355E8099")
	snapshot(t, "color_picker", 340, 560, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(12).AlignItems(ui.Start).Children(func() {
			ColorPicker(c, &b, ColorPickerOptions{Disabled: true})
			p := ColorPicker(c, &a, ColorPickerOptions{})
			st := ui.Local(p.Trigger, "snap", func() bool { return false })
			if !*st {
				*st = true
				ui.Local(p.Trigger, "picker", func() colorPickerState { return colorPickerState{} }).open = true
				c.Invalidate()
			}
		})
	})
}

// Disabling an open picker closes it before keys can change the color.
func TestColorPickerDisabledWhileOpen(t *testing.T) {
	col, n := ui.Hex("#808080"), 0
	opts := ColorPickerOptions{Label: "Shade"}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ColorPicker(c, &col, opts).Changed() {
			n++
		}
	}, 400, 560)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !tt.Focused("Saturation and brightness") {
		t.Fatal("not open")
	}
	opts.Disabled = true
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyRight)
	if col != ui.Hex("#808080") || n != 0 || tt.HasText("HEX") {
		t.Fatalf("disabled picker changed to %v", col)
	}
}

func TestColorPickerSwatchHitArea(t *testing.T) {
	col, n := ui.Hex("#7D2034"), 0
	tt := colorPickerTester(&col, ColorPickerOptions{Label: "Accent"}, &n)
	_ = tt.Click("Accent")
	r, ok := tt.Find("#355E80")
	if !ok || r.W < 28 || r.H < 28 {
		t.Fatalf("swatch %v", r)
	}
}

func TestColorPickerSliderHitAreaCompact(t *testing.T) {
	col := ui.Hex("#7D2034")
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Density: theme.Compact})
		ColorPicker(c, &col, ColorPickerOptions{Label: "Accent"})
	}, 400, 560)
	_ = tt.Click("Accent")
	for _, name := range []string{"Hue", "Opacity"} {
		if r, ok := tt.Find(name); !ok || r.H < 28 {
			t.Fatalf("%s slider %v", name, r)
		}
	}
}
