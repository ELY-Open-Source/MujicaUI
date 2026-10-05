package chart

import (
	"image/color"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestChartCrosshairSnapsAndFollowsKeys(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult { return partsChart(c, nil, true) }}
	tt := ui.NewTester(h.view, 520, 420)
	b := h.res.Element.Bounds()
	plot := ChartInset(ui.Rect{X: b.X, Y: b.Y, W: b.W, H: 240}, 52, 26)
	// Point 1 is ten DIPs from the pointer.
	px, py := plot.X+plot.W/2, plot.Y+plot.H*0.2
	tt.Move(px+10, py-6)
	if h.res.Point != 1 || chartReadout(tt) != "x 5 · y 80" {
		t.Fatalf("snap: point %d %q", h.res.Point, chartReadout(tt))
	}
	if !crosshairAt(tt, int(px), int(plot.Y+plot.H-2)) {
		t.Fatal("vertical line does not pass through the snapped point")
	}
	tt.Move(px+60, py+60)
	if h.res.Point != -1 {
		t.Fatalf("snapped from too far: %d", h.res.Point)
	}
	tabTo(t, tt, "Parts")
	tt.Key(0, ui.KeyEnd)
	if !crosshairAt(tt, int(plot.X+plot.W*0.8), int(plot.Y+plot.H-2)) {
		t.Fatal("keyboard point has no crosshair")
	}
}

// crosshairAt checks a crosshair pixel.
func crosshairAt(tt *ui.Tester, x, y int) bool {
	want := toRGBA(theme.Light().TextMuted)
	for dx := -1; dx <= 1; dx++ {
		if color.RGBAModel.Convert(tt.Image().At(x+dx, y)).(color.RGBA) == want {
			return true
		}
	}
	return false
}

func TestChartCrosshairBad(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without axes")
		}
	}()
	ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); ChartCrosshair(c, ChartCrosshairOptions{}) }, 100, 100)
}

func TestChartCrosshairSnapshot(t *testing.T) {
	snapshot(t, "chart_crosshair", 420, 240, func(c *ui.Context) {
		x := ChartAxis(c, ChartAxisOptions{Min: 0, Max: 10})
		y := ChartAxis(c, ChartAxisOptions{Side: SideLeft, Min: 0, Max: 100})
		hair := ChartCrosshair(c, ChartCrosshairOptions{X: &x, Y: &y})
		ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
			plot := ui.Rect{X: r.X + 48, Y: r.Y + 16, W: r.W - 64, H: r.H - 48}
			x.Draw(p, plot)
			y.Draw(p, plot)
			hair.Draw(p, plot, x.Map(plot, 6.5), y.Map(plot, 42))
		})
	})
}

func TestChartCrosshairDisablesFocusedReadout(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult { return partsChart(c, nil, true) }}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Parts")
	tt.Key(0, ui.KeyRight)
	h.disabled = true
	tt.Frame()
	if h.res.Point != -1 {
		t.Fatalf("disabled readout %d", h.res.Point)
	}
}

func TestChartCrosshairFiniteReach(t *testing.T) {
	charts2MustPanic(t, func(c *ui.Context) {
		a := ChartAxis(c, ChartAxisOptions{Max: 1})
		ChartCrosshair(c, ChartCrosshairOptions{X: &a, Reach: float32(math.NaN())})
	})
}
