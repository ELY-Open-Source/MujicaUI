package chart

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var heatRows = []string{"Mon", "Tue", "Wed"}
var heatCols = []string{"Matins", "Lauds", "Vespers", "Compline"}

func heatValues() [][]float64 {
	return [][]float64{{1, 4, 9, 2}, {3, math.NaN(), 7, 5}, {8, 6, 0, 3}}
}

func TestHeatmapChartKeyboard2D(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return HeatmapChart(c, heatRows, heatCols, heatValues(), HeatmapChartOptions{Frame: ChartFrame{Label: "Hours"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Hours")
	tt.Key(0, ui.KeyRight) // enters at the first cell
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	if got := chartReadout(tt); got != "Tue · Lauds · no value" {
		t.Fatalf("readout %q", got)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnd)
	if got := chartReadout(tt); got != "Wed · Compline · 3" {
		t.Fatalf("End %q", got)
	}
}

func TestHeatmapChartHover(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return HeatmapChart(c, heatRows, heatCols, heatValues(), HeatmapChartOptions{Frame: ChartFrame{Label: "Hours"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	b := h.box(tt)
	plotW := b.W - 72 - 64
	tt.Move(b.X+72+plotW*0.6, b.Y+4+212*0.2) // third column, first row
	if got := chartReadout(tt); got != "Mon · Vespers · 9" {
		t.Fatalf("hover %q", got)
	}
}

func TestHeatmapChartBad(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a short row")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		HeatmapChart(c, heatRows[:1], heatCols, [][]float64{{1}}, HeatmapChartOptions{Frame: ChartFrame{Label: "x"}})
	}, 300, 300)
}

func TestHeatmapChartSnapshot(t *testing.T) {
	snapshot(t, "heatmap_chart", 520, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			HeatmapChart(c, heatRows, heatCols, heatValues(), HeatmapChartOptions{Frame: ChartFrame{Label: "Hours", Height: 220}})
		})
	})
}
