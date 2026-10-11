package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var funnelStages = []ChartSlice{{"Visitors", 1000}, {"Novices", 400}, {"Vows", 100}, {"Abbots", 0}}

func TestFunnelChartReadout(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return FunnelChart(c, funnelStages, FunnelChartOptions{Frame: ChartFrame{Label: "Order"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Order")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "Novices · 400 · 40.0% of the stage before · 40.0% of the first" {
		t.Fatalf("readout %q", got)
	}
	tt.Key(0, ui.KeyEscape)
	b := h.res.Element.Bounds()
	tt.Move(b.X+b.W/3, b.Y+4+232*0.6) // third stage band
	if got := chartReadout(tt); got != "Vows · 100 · 25.0% of the stage before · 10.0% of the first" {
		t.Fatalf("hover %q", got)
	}
}

func TestFunnelChartBadAndEmpty(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return FunnelChart(c, []ChartSlice{{"a", 0}}, FunnelChartOptions{Frame: ChartFrame{Label: "Order"}})
	}}
	if tt := ui.NewTester(h.view, 520, 420); !tt.HasText("No data") {
		t.Fatal("all-zero funnel is not empty")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a negative stage")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		FunnelChart(c, []ChartSlice{{"a", -1}}, FunnelChartOptions{Frame: ChartFrame{Label: "x"}})
	}, 300, 300)
}

func TestFunnelChartSnapshot(t *testing.T) {
	snapshot(t, "funnel_chart", 520, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			FunnelChart(c, funnelStages[:3], FunnelChartOptions{Frame: ChartFrame{Label: "Order", Height: 220}})
		})
	})
}
