package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func paletteSwatches() []ColorSwatch {
	var s []ColorSwatch
	for i, c := range colorPickerPresets {
		s = append(s, ColorSwatch{Color: c, Name: "Tone " + string(rune('A'+i))})
	}
	return s
}

func TestColorPaletteChooseAndKeyboard(t *testing.T) {
	v := colorPickerPresets[0]
	dis := false
	n := 0
	sw := paletteSwatches()
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if ColorPalette(c, &v, sw, ColorPaletteOptions{Columns: 4, Custom: true, Disabled: dis, Label: "Tag color"}).Changed() {
			n++
		}
	}, 500, 400)
	if err := tt.Click("Tone C"); err != nil || v != sw[2].Color || n != 1 {
		t.Fatalf("click: %v %v %d", err, v, n)
	}
	tt.Key(0, ui.KeyDown) // focus stays on C: down a row is G
	tt.Key(0, ui.KeyEnter)
	if v != sw[6].Color {
		t.Fatalf("Down then Enter chose %v, want %v", v, sw[6].Color)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeySpace)
	if v != sw[3].Color {
		t.Fatalf("Right, Up, Space chose %v", v)
	}
	if _, ok := tt.Find("Custom color"); !ok {
		t.Fatal("no custom color entry")
	}
	dis = true
	tt.Frame()
	_ = tt.Click("Tone A")
	if v != sw[3].Color || n != 3 {
		t.Fatalf("disabled palette changed: %v %d", v, n)
	}
}

func TestColorPaletteRejects(t *testing.T) {
	v := ui.Hex("#000000")
	for name, f := range map[string]func(c *ui.Context){
		"no name":  func(c *ui.Context) { ColorPalette(c, &v, []ColorSwatch{{Color: v}}, ColorPaletteOptions{}) },
		"columns":  func(c *ui.Context) { ColorPalette(c, &v, paletteSwatches(), ColorPaletteOptions{Columns: -1}) },
		"no items": func(c *ui.Context) { ColorPalette(c, &v, nil, ColorPaletteOptions{}) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); f(c) }, 300, 200)
		}()
	}
}

func TestColorPaletteSnapshot(t *testing.T) {
	v := colorPickerPresets[2]
	snapshot(t, "color_palette", 300, 150, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ColorPalette(c, &v, paletteSwatches(), ColorPaletteOptions{Columns: 6, Custom: true})
		})
	})
}
