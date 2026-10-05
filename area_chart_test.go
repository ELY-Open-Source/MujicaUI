package mujica

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestAreaBands(t *testing.T) {
	series := []ChartSeries{{Name: "a", Values: []float64{1, math.NaN(), 3}}, {Name: "b", Values: []float64{2, 2, 2}}}
	none := func(string) bool { return false }
	st := areaBands(series, none, true, 10)
	if st[0][0][0] != 11 || st[1][1][0] != 11 || st[0][1][0] != 13 || !math.IsNaN(st[0][0][1]) || st[0][1][1] != 12 {
		t.Fatalf("stacked %v", st)
	}
	hideA := func(s string) bool { return s == "a" }
	st = areaBands(series, hideA, true, 0)
	if st[0][1][2] != 2 || st[1][1][2] != 0 {
		t.Fatalf("hidden series still stacked: %v", st)
	}
	flat := areaBands(series, none, false, -1)
	if flat[0][0][2] != 3 || flat[1][0][2] != -1 || flat[1][1][2] != -1 {
		t.Fatalf("baseline %v", flat)
	}
}

func TestAreaChartStackedReadout(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return AreaChart(c, lineSeries(), AreaChartOptions{Frame: ChartFrame{Label: "Nave"}, Stacked: true, Baseline: 1})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Nave")
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "0 · Pews 3 · Candles 5" {
		t.Fatalf("readout shows the series values, not the stack: %q", got)
	}
	b := h.res.Element.Bounds()
	tt.Move(b.X+52+(b.W-64)*0.5, b.Y+b.H/2)
	if got := chartReadout(tt); got != "2 · Pews no value · Candles 6" || h.res.Point != 2 {
		t.Fatalf("hover %q", got)
	}
}

func TestAreaChartSnapshot(t *testing.T) {
	snapshot(t, "area_chart", 520, 340, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			AreaChart(c, []ChartSeries{
				{Name: "Matins", Values: []float64{3, 4, 6, 5, 7, 8}},
				{Name: "Vespers", Values: []float64{2, 3, 2, 4, 3, 5}},
			}, AreaChartOptions{Frame: ChartFrame{Label: "Attendance"}, Stacked: true})
		})
	})
}
