package chart

import (
	"math"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func densityGroups() []BoxGroup {
	a, b := make([]float64, 200), make([]float64, 200)
	for i := range a {
		u := float64(i)/200 - 0.5
		a[i] = 40 + 30*u*math.Abs(u)
		b[i] = 60 + 12*math.Sin(float64(i))
	}
	return []BoxGroup{{Name: "Matins", Values: a}, {Name: "Vespers", Values: b}}
}

func TestDensityKDEIntegratesToOne(t *testing.T) {
	v := []float64{1, 2, 2, 3, 7}
	xs := make([]float64, 2001)
	for i := range xs {
		xs[i] = -10 + 30*float64(i)/2000
	}
	d := densityKDE(v, densityBandwidth(v), xs)
	sum := 0.0
	for _, y := range d {
		sum += y * 30 / 2000
	}
	if math.Abs(sum-1) > 0.01 {
		t.Fatalf("density integrates to %v", sum)
	}
	if bw := densityBandwidth([]float64{4, 4}); bw != 1 {
		t.Fatalf("flat bandwidth %v", bw)
	}
}

func TestDensityPlotReadoutAndBandwidth(t *testing.T) {
	bw := 0.0
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return DensityPlot(c, densityGroups(), DensityPlotOptions{Frame: ChartFrame{Label: "Ages"}, Bandwidth: bw, Points: 64})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Ages")
	tt.Key(0, ui.KeyRight)
	for range 31 {
		tt.Key(0, ui.KeyRight)
	}
	first := chartReadout(tt)
	if !strings.Contains(first, "Matins ") || !strings.Contains(first, "Vespers ") {
		t.Fatalf("readout %q", first)
	}
	bw = 1
	tt.Frame()
	if chartReadout(tt) == first {
		t.Fatal("bandwidth did not change the curve")
	}
	if err := tt.Click("Vespers"); err != nil || strings.Contains(chartReadout(tt), "Vespers") {
		t.Fatalf("legend did not hide a group: %q", chartReadout(tt))
	}
}

func TestDensityPlotBadAndEmpty(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return DensityPlot(c, []BoxGroup{{Name: "a", Values: []float64{math.NaN()}}}, DensityPlotOptions{Frame: ChartFrame{Label: "Ages"}})
	}}
	if tt := ui.NewTester(h.view, 520, 420); !tt.HasText("No data") {
		t.Fatal("NaN sample is not empty")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for Points 1")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		DensityPlot(c, densityGroups(), DensityPlotOptions{Frame: ChartFrame{Label: "x"}, Points: 1})
	}, 300, 300)
}

func TestDensityPlotSnapshot(t *testing.T) {
	snapshot(t, "density_plot", 520, 340, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			DensityPlot(c, densityGroups(), DensityPlotOptions{Frame: ChartFrame{Label: "Ages"}})
		})
	})
}

func TestDensityPlotFiniteBandwidthAndNames(t *testing.T) {
	charts2MustPanic(t, func(c *ui.Context) {
		DensityPlot(c, densityGroups(), DensityPlotOptions{Frame: ChartFrame{Label: "Density"}, Bandwidth: math.Inf(1)})
	})
	charts2MustPanic(t, func(c *ui.Context) {
		DensityPlot(c, []BoxGroup{{Name: "A"}, {Name: "A"}}, DensityPlotOptions{Frame: ChartFrame{Label: "Density"}})
	})
}

func TestDensityPlotLargeDataBoundedAndDisabled(t *testing.T) {
	values := make([]float64, 100_000)
	for i := range values {
		values[i] = float64(i % 100)
	}
	groups := []BoxGroup{{Name: "Sample", Values: values}}
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return DensityPlot(c, groups, DensityPlotOptions{Frame: ChartFrame{Label: "Large density"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Large density")
	tt.Key(0, ui.KeyEnd)
	if h.res.Point != 127 {
		t.Fatalf("built sample index %d", h.res.Point)
	}
	h.disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyLeft)
	if h.res.Point != -1 {
		t.Fatalf("disabled density %d", h.res.Point)
	}
}
