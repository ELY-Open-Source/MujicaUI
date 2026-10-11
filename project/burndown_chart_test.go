package project

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/chart"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
)

var burndownRemaining = []float64{40, 38, 35, 35, 30, 26, math.NaN(), 20, 19}

func burndownTester(remaining []float64, opts BurndownChartOptions, points *[]int) *ui.Tester {
	opts.Frame.Label = "Sprint 12 burndown"
	opts.Start, opts.End = oct(1), oct(14)
	if opts.Total == 0 {
		opts.Total = 40
	}
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(16).Width(600).Children(func() {
			if r := BurndownChart(c, remaining, opts); r.Changed() {
				*points = append(*points, r.Point)
			}
		})
	}, 640, 400)
}

func TestBurndownChartReadoutAndKeys(t *testing.T) {
	var points []int
	tt := burndownTester(burndownRemaining, BurndownChartOptions{}, &points)
	if !tt.HasText("Remaining") || !tt.HasText("Ideal") {
		t.Fatalf("legend %v", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("October 3 · Remaining 35 · Ideal 33.85") {
		t.Fatalf("readout %v", tt.Texts())
	}
	if a := tt.Announcements(); len(a) == 0 || a[len(a)-1] != "October 3 · Remaining 35 · Ideal 33.85" {
		t.Fatalf("announce %q", a)
	}
	tt.Key(0, ui.KeyEnd)
	if !tt.HasText("October 14 · Ideal 0") {
		t.Fatalf("future day shows only the ideal: %v", tt.Texts())
	}
	_ = tt.Click("Ideal")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	if !tt.HasText("October 1 · Remaining 40") || tt.HasText("Ideal 40") {
		t.Fatal("hidden series in readout")
	}
}

func TestBurndownChartEmptyAndInvalid(t *testing.T) {
	var points []int
	tt := burndownTester(nil, BurndownChartOptions{Frame: chart.ChartFrame{Status: data.DataLoading}}, &points)
	if !tt.HasText("Loading") || tt.HasText("October") {
		t.Fatalf("loading state %v", tt.Texts())
	}
	for _, bad := range []BurndownChartOptions{{Total: -1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("expected panic")
				}
			}()
			burndownTester(nil, bad, &points)
		}()
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for too many days")
		}
	}()
	burndownTester(make([]float64, 20), BurndownChartOptions{}, &points)
}

func TestBurndownChartSnapshot(t *testing.T) {
	snapshot(t, "burndown_chart", 620, 340, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			BurndownChart(c, burndownRemaining, BurndownChartOptions{Frame: chart.ChartFrame{Label: "Burndown"}, Start: oct(1), End: oct(14), Total: 40})
		})
	})
}

func TestBurndownChartInvalidAmounts(t *testing.T) {
	for _, values := range [][]float64{{-1}, {math.Inf(1)}} {
		t.Run("amount", func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected invalid amount panic")
				}
			}()
			var points []int
			burndownTester(values, BurndownChartOptions{}, &points)
		})
	}
}

func TestBurndownChartEmptyFailedDisabled(t *testing.T) {
	var points []int
	tt := burndownTester(nil, BurndownChartOptions{}, &points)
	if !tt.HasText("No data") {
		t.Fatal("missing empty state")
	}
	tt = burndownTester(nil, BurndownChartOptions{Frame: chart.ChartFrame{Status: data.DataFailed, Error: "Unavailable"}}, &points)
	if !tt.HasText("Unavailable") {
		t.Fatal("missing error state")
	}
	tt = ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(true).Children(func() {
			r := BurndownChart(c, burndownRemaining, BurndownChartOptions{Frame: chart.ChartFrame{Label: "Burndown"}, Start: oct(1), End: oct(14), Total: 40})
			if _, ok := r.Toggled(); ok {
				t.Fatal("disabled legend toggled")
			}
		})
	}, 600, 340)
	_ = tt.Click("Ideal")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if len(tt.Announcements()) != 0 {
		t.Fatal("disabled chart announced a point")
	}
}

func TestBurndownChartLargeData(t *testing.T) {
	values := make([]float64, 100000)
	for i := range values {
		values[i] = float64(100000 - i)
	}
	opts := BurndownChartOptions{Frame: chart.ChartFrame{Label: "Large burndown"}, Start: core.Date{Year: 1900, Month: 1, Day: 1}, End: core.Date{Year: 1900, Month: 1, Day: 1}.AddDays(99999), Total: 100000}
	ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); BurndownChart(c, values, opts) }, 640, 360)
	points := burndownSamples(values, 100000, 600)
	if len(points) > 2404 {
		t.Fatalf("built %d points", len(points))
	}
	if points[0].index != 0 || points[len(points)-1].index != 99999 {
		t.Fatal("downsampling lost endpoints")
	}
}
