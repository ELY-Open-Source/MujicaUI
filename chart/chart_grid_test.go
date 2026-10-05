package chart

import (
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestChartGridPixelAlignedWithMinor(t *testing.T) {
	for _, scale := range []float32{1, 2} {
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{Mode: core.Light})
			y := ChartAxis(c, ChartAxisOptions{Side: SideLeft, Min: 0, Max: 100, Ticks: []float64{0, 50, 100}})
			grid := ChartGrid(c, ChartGridOptions{Y: &y, Minor: 2})
			ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) { grid.Draw(p, ui.Rect{X: 0, Y: 10.3, W: 200, H: 100}) })
		}, 200, 130)
		tt.SetScale(scale)
		img := tt.Image()
		major, minor := toRGBA(theme.Light().Border), toRGBA(theme.Light().Border.Mix(theme.Light().Surface, 0.55))
		count := map[color.RGBA]int{}
		for py := range img.Bounds().Dy() {
			c := color.RGBAModel.Convert(img.At(100, py)).(color.RGBA)
			count[c]++
		}
		// Three major and two minor lines stay pixel-aligned.
		if count[major] != 3*int(scale) || count[minor] != 2*int(scale) {
			t.Fatalf("scale %v: %d major and %d minor pixel rows", scale, count[major], count[minor])
		}
	}
}

func TestChartGridBad(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context){
		"no axis": func(c *ui.Context) { ChartGrid(c, ChartGridOptions{}) },
		"category minor": func(c *ui.Context) {
			a := ChartAxis(c, ChartAxisOptions{Kind: AxisCategory, Categories: []string{"a"}})
			ChartGrid(c, ChartGridOptions{X: &a, Minor: 2})
		},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); build(c) }, 100, 100)
		}()
	}
}

func TestChartGridSnapshot(t *testing.T) {
	snapshot(t, "chart_grid", 420, 240, func(c *ui.Context) {
		x := ChartAxis(c, ChartAxisOptions{Min: 0, Max: 10})
		y := ChartAxis(c, ChartAxisOptions{Side: SideLeft, Kind: AxisLog, Min: 1, Max: 1000})
		grid := ChartGrid(c, ChartGridOptions{X: &x, Y: &y, Minor: 4})
		ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
			plot := ui.Rect{X: r.X + 48, Y: r.Y + 16, W: r.W - 64, H: r.H - 48}
			grid.Draw(p, plot)
			x.Draw(p, plot)
			y.Draw(p, plot)
		})
	})
}
