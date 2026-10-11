package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

var contextSegmentsSample = []ContextSegment{{Label: "System", Tokens: 1200}, {Label: "History", Tokens: 58000}, {Label: "Attachments", Tokens: 24000}}

func TestContextWindowMeterLegendAndWarning(t *testing.T) {
	segs := contextSegmentsSample
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ContextWindowMeter(c, segs, 128000, ContextWindowMeterOptions{})
	}, 600, 200)
	for _, s := range []string{"System 1,200", "History 58,000", "Attachments 24,000", "Free 44,800", "83,200 of 128,000 tokens (65%)"} {
		if !tt.HasText(s) {
			t.Errorf("missing %q in %q", s, tt.Texts())
		}
	}
	if tt.HasText("Context is") {
		t.Fatal("warning below the threshold")
	}
	segs = append(segs, ContextSegment{Label: "Draft", Tokens: 30000})
	tt.Frame()
	if !tt.HasText("Context is 88% full") {
		t.Fatalf("warning missing: %q", tt.Texts())
	}
}

func TestContextWindowMeterPaletteContrast(t *testing.T) {
	for _, k := range []theme.Tokens{theme.Light(), theme.Dark()} {
		for i, col := range internal.Palette(k) {
			if r := theme.Contrast(col, k.Surface); r < 3 {
				t.Errorf("dark %v series %d contrast %.2f", k.Dark, i, r)
			}
		}
	}
}

func TestContextWindowMeterBadOptions(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a zero limit")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ContextWindowMeter(c, nil, 0, ContextWindowMeterOptions{})
	}, 100, 100)
}

func TestContextWindowMeterSnapshot(t *testing.T) {
	snapshot(t, "context_window_meter", 520, 140, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ContextWindowMeter(c, append(contextSegmentsSample, ContextSegment{Label: "Draft", Tokens: 30000}), 128000, ContextWindowMeterOptions{})
		})
	})
}
