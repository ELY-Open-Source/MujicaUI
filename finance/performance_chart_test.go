package finance

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
)

func TestPerformanceChartBenchmarkDrawdownAndKeys(t *testing.T) {
	values := []float64{100, 110, 88, 120}
	times := mktTimes(4)
	var last chart.ChartResult
	tt := chatTester(760, 420, func(c *ui.Context) {
		last = PerformanceChart(c, values, PerformanceChartOptions{Frame: chart.ChartFrame{Label: "Performance"}, Times: times, Benchmark: []float64{100, 105, 108, 112}})
	})
	if dd := MaxDrawdown(values); dd.Peak != 1 || dd.Trough != 2 || math.Abs(dd.Percent+20) > 1e-9 {
		t.Fatalf("drawdown %+v", dd)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if last.Point != 3 || !tt.HasText("2026-01-04 · Portfolio +20.00% · Benchmark +12.00%") {
		t.Fatalf("point %d text %q", last.Point, tt.Texts())
	}
	_ = tt.Click("Benchmark")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if !tt.HasText("2026-01-04 · Portfolio +20.00%") {
		t.Fatalf("legend %q", tt.Texts())
	}
	marketsPanic(t, func(c *ui.Context) {
		PerformanceChart(c, []float64{-1}, PerformanceChartOptions{Frame: chart.ChartFrame{Label: "x"}})
	})
}
func TestPerformanceChartSnapshot(t *testing.T) {

	snapshot(t, "performance_chart", 720, 380, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			PerformanceChart(c, mktSeries(120, 0), PerformanceChartOptions{Frame: chart.ChartFrame{Label: "Performance"}, Times: mktTimes(120), Benchmark: mktSeries(120, 1)})

		})
	})
}

func TestPerformanceChartRequiresTimes(t *testing.T) {
	marketsPanic(t, func(c *ui.Context) {
		PerformanceChart(c, []float64{100, 110}, PerformanceChartOptions{Frame: chart.ChartFrame{Label: "x"}})
	})
}

func TestPerformanceChartDateVersion(t *testing.T) {
	values := []float64{100, 105, 110}
	times := mktTimes(3)
	version := 0
	var axis chart.ChartX
	tt := chatTester(760, 420, func(c *ui.Context) {
		r := PerformanceChart(c, values, PerformanceChartOptions{Frame: chart.ChartFrame{Label: "Performance"}, Times: times, Version: version})
		memo := ui.Local(r.Element, "mujica.memo.x", func() chart.ChartMemoEntry { return chart.ChartMemoEntry{} })
		axis = memo.Val.(chart.ChartX)
	})
	times[2] = times[0].AddDate(0, 0, 9)
	values[2] = 120
	version++
	tt.Frame()
	if got := axis.Vals[2]; got != float64(times[2].UnixNano())/1e9 {
		t.Fatalf("axis time %v", got)
	}
}
