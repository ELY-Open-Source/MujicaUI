package chart

import (
	"fmt"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestHistogramBin(t *testing.T) {
	vals := []float64{1, 2, 2, 3, 4, 5, 9, 10, math.NaN()}
	if got := histogramBin(vals, 3, 0, false); fmt.Sprint(got.counts) != "[4 2 2]" || got.start != 1 || got.width != 3 {
		t.Fatalf("3 bins %+v", got)
	}
	if got := histogramBin(vals, 0, 5, false); fmt.Sprint(got.counts) != "[5 2 1]" || got.start != 0 {
		t.Fatalf("width 5 %+v", got)
	}
	if got := histogramBin(vals, 0, 5, true); fmt.Sprint(got.counts) != "[5 7 8]" {
		t.Fatalf("cumulative %+v", got)
	}
	if got := histogramBin(vals, 0, 0, false); len(got.counts) != 4 {
		t.Fatalf("Sturges for 8 values: %d bins", len(got.counts))
	}
	if got := histogramBin([]float64{math.NaN()}, 3, 0, false); got.counts != nil {
		t.Fatal("no finite values made bins")
	}
}

func TestHistogramReadoutAndModes(t *testing.T) {
	opts := HistogramOptions{Frame: ChartFrame{Label: "Ages"}, BinWidth: 5}
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return Histogram(c, []float64{1, 2, 2, 3, 4, 5, 9, 10}, opts)
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Ages")
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "[0, 5) · 5" {
		t.Fatalf("first bin %q", got)
	}
	tt.Key(0, ui.KeyEnd)
	if got := chartReadout(tt); got != "[10, 15) · 1" {
		t.Fatalf("last bin %q", got)
	}
	opts.Cumulative = true
	tt.Frame()
	if got := chartReadout(tt); got != "≤ 15 · 8" {
		t.Fatalf("cumulative %q", got)
	}
}

func TestHistogramBad(t *testing.T) {
	for _, o := range []HistogramOptions{{Bins: -1}, {BinWidth: 1e-9}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: no panic", o)
				}
			}()
			o.Frame.Label = "x"
			ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); Histogram(c, []float64{0, 1}, o) }, 300, 300)
		}()
	}
}

func TestHistogramSnapshot(t *testing.T) {
	vals := make([]float64, 400)
	for i := range vals {
		vals[i] = 50 + 15*math.Sin(float64(i)*1.7)*math.Cos(float64(i)*0.3)
	}
	snapshot(t, "histogram", 520, 320, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			Histogram(c, vals, HistogramOptions{Frame: ChartFrame{Label: "Ages"}, Bins: 12})
		})
	})
}
