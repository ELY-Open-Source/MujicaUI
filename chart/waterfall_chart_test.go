package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var waterfallSteps = []WaterfallStep{
	{Label: "Tithes", Value: 12},
	{Label: "Repairs", Value: -5},
	{Label: "Q1", Kind: WaterfallSubtotal},
	{Label: "Alms", Value: 4},
	{Label: "Year", Kind: WaterfallTotal},
}

func TestWaterfallChartReadout(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return WaterfallChart(c, waterfallSteps, WaterfallChartOptions{Frame: ChartFrame{Label: "Ledger"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Ledger")
	want := []string{"Tithes · +12 · running total 12", "Repairs · −5 · running total 7", "Q1 · 7", "Alms · +4 · running total 11", "Year · 11"}
	for i, w := range want {
		tt.Key(0, ui.KeyRight)
		if got := chartReadout(tt); got != w {
			t.Fatalf("step %d: %q, want %q", i, got, w)
		}
	}
}

func TestWaterfallChartBad(t *testing.T) {
	for _, s := range []WaterfallStep{{Label: "x", Value: 3, Kind: WaterfallTotal}, {Value: 1}, {Label: "x", Kind: 9}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: no panic", s)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				WaterfallChart(c, []WaterfallStep{s}, WaterfallChartOptions{Frame: ChartFrame{Label: "x"}})
			}, 300, 300)
		}()
	}
}

func TestWaterfallChartSnapshot(t *testing.T) {
	snapshot(t, "waterfall_chart", 520, 320, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			WaterfallChart(c, waterfallSteps, WaterfallChartOptions{Frame: ChartFrame{Label: "Ledger"}})
		})
	})
}
