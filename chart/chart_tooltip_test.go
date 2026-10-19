package chart

import (
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

func TestChartTooltipFlipsAtEdges(t *testing.T) {
	var at [2]float32
	var tip ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		tip = ChartTooltip(c, at[0], at[1], []ChartTooltipLine{{Text: "Lauds"}, {Text: "Pews 5", Color: theme.Light().Info, Shape: 1}}).Bounds()
	}, 400, 300)
	at = [2]float32{50, 50}
	tt.Frame()
	tt.Frame()
	if tip.X != 62 || tip.Y != 62 || tip.W < 40 {
		t.Fatalf("tooltip at %v", tip)
	}
	at = [2]float32{395, 295}
	tt.Frame()
	tt.Frame()
	if tip.X+tip.W > 395-12+0.5 || tip.Y+tip.H > 295-12+0.5 {
		t.Fatalf("tooltip did not flip: %v", tip)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without lines")
		}
	}()
	ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); ChartTooltip(c, 0, 0, nil) }, 100, 100)
}

func TestChartTooltipFollowsFocusAndPointer(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return LineChart(c, lineSeries(), LineChartOptions{Frame: ChartFrame{Label: "Chapel"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	tabTo(t, tt, "Chapel")
	tt.Key(0, ui.KeyEnd)
	if !tooltipShown(tt) {
		t.Fatal("no tooltip for the keyboard point")
	}
	tt.Key(0, ui.KeyEscape)
	if tooltipShown(tt) {
		t.Fatal("tooltip stayed after Escape")
	}
	b := h.box(tt)
	tt.Move(b.X+b.W/2, b.Y+b.H/2)
	if !tooltipShown(tt) {
		t.Fatal("no tooltip under the pointer")
	}
}

// tooltipShown finds Surface runs longer than marker rings.
func tooltipShown(tt *ui.Tester) bool {
	img := tt.Image()
	want := toRGBA(theme.Light().Surface)
	r := img.Bounds()
	for y := r.Min.Y; y < r.Max.Y; y += 2 {
		run := 0
		for x := r.Min.X; x < r.Max.X; x++ {
			if color.RGBAModel.Convert(img.At(x, y)).(color.RGBA) == want {
				if run++; run >= 30 {
					return true
				}
			} else {
				run = 0
			}
		}
	}
	return false
}

func TestChartTooltipSnapshot(t *testing.T) {
	snapshot(t, "chart_tooltip", 320, 140, func(c *ui.Context) {
		pal := internal.Palette(core.Tokens(c))
		ChartTooltip(c, 20, 20, []ChartTooltipLine{{Text: "Lauds"}, {Text: "Pews 5", Color: pal[0]}, {Text: "Candles 2", Color: pal[1], Shape: 1}})
	})
}

func TestChartTooltipFitsLongContent(t *testing.T) {
	var b ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		lines := make([]ChartTooltipLine, 40)
		for i := range lines {
			lines[i].Text = "A very long series name that wraps across the tooltip panel"
		}
		b = ChartTooltip(c, 159, 99, lines).Bounds()
	}, 160, 100)
	tt.Frame()
	if b.X < 4 || b.Y < 4 || b.X+b.W > 156 || b.Y+b.H > 96 {
		t.Fatalf("tooltip outside window: %v", b)
	}
}
