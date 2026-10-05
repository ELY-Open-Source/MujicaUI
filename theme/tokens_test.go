package theme

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestContrastFloors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		tokens Tokens
	}{{"light", Light()}, {"dark", Dark()}} {
		k := tc.tokens
		text := map[string][2]ui.Color{
			"Text/Background":      {k.Text, k.Background},
			"Text/Surface":         {k.Text, k.Surface},
			"Text/SurfaceHover":    {k.Text, k.SurfaceHover},
			"Text/Selection":       {k.Text, k.Selection},
			"TextMuted/Background": {k.TextMuted, k.Background},
			"TextMuted/Surface":    {k.TextMuted, k.Surface},
			"OnAccent/Accent":      {k.OnAccent, k.Accent},
			"OnAccent/AccentHover": {k.OnAccent, k.AccentHover},
			"OnAccent/Pressed":     {k.OnAccent, k.AccentPressed},
			"OnDanger/DangerFill":  {k.OnDanger, k.DangerFill},
			"AccentText/Surface":   {k.AccentText, k.Surface},
			"Success/Surface":      {k.Success, k.Surface},
			"Warning/Surface":      {k.Warning, k.Surface},
			"Danger/Surface":       {k.Danger, k.Surface},
			"Info/Surface":         {k.Info, k.Surface},
		}
		for pair, c := range text {
			if r := Contrast(c[0], c[1]); r < 4.5 {
				t.Errorf("%s %s = %.2f, want >= 4.5", tc.name, pair, r)
			}
		}
		graphic := map[string][2]ui.Color{
			"ControlBorder/Background": {k.ControlBorder, k.Background},
			"ControlBorder/Surface":    {k.ControlBorder, k.Surface},
			"Focus/Background":         {k.Focus, k.Background},
		}
		for pair, c := range graphic {
			if r := Contrast(c[0], c[1]); r < 3 {
				t.Errorf("%s %s = %.2f, want >= 3", tc.name, pair, r)
			}
		}
	}
}

func TestDensity(t *testing.T) {
	for _, tc := range []struct {
		d            Density
		unit, height float32
	}{{Compact, 3, 28}, {Default, 4, 36}, {Comfortable, 5, 44}} {
		if tc.d.Unit() != tc.unit || tc.d.ControlHeight() != tc.height {
			t.Errorf("%v: unit %v height %v", tc.d, tc.d.Unit(), tc.d.ControlHeight())
		}
	}
}
