package chart

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
)

type chartHarness struct {
	res      ChartResult
	bounds   ui.Rect
	changes  int
	disabled bool
	build    func(c *ui.Context) ChartResult
	built    func(r ChartResult) // reads the live element during the build
}

// box is the plot's committed bounds: Bounds reads the previous frame, so one more frame settles it.
func (h *chartHarness) box(tt *ui.Tester) ui.Rect {
	tt.Frame()
	return h.bounds
}

func (h *chartHarness) view(c *ui.Context) {
	core.Use(c, core.Settings{Mode: core.Light})
	ui.Column(c).Padding(16).Width(480).Disabled(h.disabled).Children(func() {
		h.res = h.build(c)
		h.bounds = h.res.Element.Bounds()
		if h.built != nil {
			h.built(h.res)
		}
		if h.res.Changed() {
			h.changes++
		}
	})
}

// chartReadout returns the readout line, the text after the plot.
func chartReadout(tt *ui.Tester) string {
	texts := tt.Texts()
	return texts[len(texts)-1]
}

func lineSeries() []ChartSeries {
	return []ChartSeries{
		{Name: "Pews", Values: []float64{3, 5, math.NaN(), 4, 8}},
		{Name: "Candles", Values: []float64{5, 2, 6, 7, 1}},
	}
}

func TestLineChartKeyboardReadout(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return LineChart(c, lineSeries(), LineChartOptions{Frame: ChartFrame{Label: "Chapel"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	if !strings.Contains(chartReadout(tt), "arrow keys") {
		t.Fatalf("hint %q", chartReadout(tt))
	}
	tabTo(t, tt, "Chapel")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "1 · Pews 5 · Candles 2" || h.res.Point != 1 || h.changes != 2 {
		t.Fatalf("readout %q point %d", got, h.res.Point)
	}
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "2 · Pews no value · Candles 6" {
		t.Fatalf("gap readout %q", got)
	}
	if a := tt.Announcements(); len(a) == 0 || a[len(a)-1] != "2 · Pews no value · Candles 6" {
		t.Fatalf("announcements %q", a)
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyRight)
	if h.res.Point != 4 {
		t.Fatalf("End then Right: %d", h.res.Point)
	}
	tt.Key(0, ui.KeyEscape)
	if h.res.Point != -1 {
		t.Fatalf("Escape kept point %d", h.res.Point)
	}
}

func TestLineChartHoverAndToggle(t *testing.T) {
	var toggled string
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		r := LineChart(c, lineSeries(), LineChartOptions{Frame: ChartFrame{Label: "Chapel"}, X: []float64{0, 10, 20, 30, 40}})
		if s, ok := r.Toggled(); ok {
			toggled = s
		}
		return r
	}}
	tt := ui.NewTester(h.view, 520, 420)
	b := h.box(tt)
	// The plot spans 52..W-12; x 30 is three quarters across.
	px := b.X + 52 + (b.W-64)*0.75
	tt.Move(px, b.Y+b.H/2)
	if got := chartReadout(tt); got != "30 · Pews 4 · Candles 7" {
		t.Fatalf("hover readout %q", got)
	}
	if err := tt.Click("Candles"); err != nil {
		t.Fatal(err)
	}
	tt.Move(px, b.Y+b.H/2)
	if got := chartReadout(tt); got != "30 · Pews 4" || toggled != "Candles" {
		t.Fatalf("hidden series still read: %q toggled %q", got, toggled)
	}
}

func TestLineChartTimeAxisAndStates(t *testing.T) {
	t0 := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	times := []time.Time{t0, t0.Add(time.Hour), t0.Add(2 * time.Hour), t0.Add(3 * time.Hour), t0.Add(4 * time.Hour)}
	status := data.DataReady
	var series []ChartSeries
	var heights []float32
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		r := LineChart(c, series, LineChartOptions{Frame: ChartFrame{Label: "Chapel", Status: status}, Times: times[:len(series[0].Values)]})
		return r
	}}
	series = lineSeries()
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Chapel")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); !strings.HasPrefix(got, "2026-10-05 10:00:00 · ") {
		t.Fatalf("time readout %q", got)
	}
	heights = append(heights, h.box(tt).H)
	for _, st := range []data.DataStatus{data.DataLoading, data.DataFailed} {
		status = st
		tt.Frame()
		heights = append(heights, h.box(tt).H)
	}
	if !tt.HasText("Retry") || heights[0] != heights[1] || heights[1] != heights[2] {
		t.Fatalf("states changed size: %v", heights)
	}
	status, series = data.DataReady, []ChartSeries{{Name: "A", Values: []float64{math.NaN()}}}
	tt.Frame()
	if !tt.HasText("No data") {
		t.Fatalf("all-missing chart: %q", tt.Texts())
	}
}

func TestLineChartDisabledIgnoresKeys(t *testing.T) {
	h := &chartHarness{disabled: true, build: func(c *ui.Context) ChartResult {
		return LineChart(c, lineSeries(), LineChartOptions{Frame: ChartFrame{Label: "Chapel"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	for range 4 {
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyRight)
	}
	if h.res.Point != -1 {
		t.Fatalf("disabled chart moved to %d", h.res.Point)
	}
}

func TestLineChartBadOptions(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context){
		"lengths": func(c *ui.Context) {
			LineChart(c, []ChartSeries{{Name: "a", Values: []float64{1}}, {Name: "b"}}, LineChartOptions{Frame: ChartFrame{Label: "x"}})
		},
		"order": func(c *ui.Context) {
			LineChart(c, lineSeries(), LineChartOptions{Frame: ChartFrame{Label: "x"}, X: []float64{0, 2, 1, 3, 4}})
		},
		"no label": func(c *ui.Context) { LineChart(c, lineSeries(), LineChartOptions{}) },
		"duplicate": func(c *ui.Context) {
			LineChart(c, []ChartSeries{{Name: "a"}, {Name: "a"}}, LineChartOptions{Frame: ChartFrame{Label: "x"}})
		},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); build(c) }, 300, 300)
		}()
	}
}

func TestLineChartLarge(t *testing.T) {
	vals := make([]float64, 100_000)
	for i := range vals {
		vals[i] = math.Sin(float64(i) / 500)
	}
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return LineChart(c, []ChartSeries{{Name: "Wave", Values: vals}}, LineChartOptions{Frame: ChartFrame{Label: "Tide"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Tide")
	tt.Key(0, ui.KeyEnd)
	if h.res.Point != 99_999 || !strings.HasPrefix(chartReadout(tt), "99999 · Wave ") {
		t.Fatalf("End: %d %q", h.res.Point, chartReadout(tt))
	}
}

func TestLineChartSnapshot(t *testing.T) {
	snapshot(t, "line_chart", 520, 340, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			LineChart(c, lineSeries(), LineChartOptions{Frame: ChartFrame{Label: "Chapel", Height: 240}})
		})
	})
}
