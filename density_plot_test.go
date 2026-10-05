package mujica

import (
	"math"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
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
		Use(c, Settings{})
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
