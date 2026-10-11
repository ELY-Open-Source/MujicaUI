package chart

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestViolinPlotReadout(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return ViolinPlot(c, append(densityGroups(), boxGroups[0]), ViolinPlotOptions{Frame: ChartFrame{Label: "Hymns"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Hymns")
	tt.Key(0, ui.KeyEnd)
	if got := chartReadout(tt); got != "Choir · median 5.5 · quartiles 3.25–7.75 · 10 values" {
		t.Fatalf("readout %q", got)
	}
	tt.Key(0, ui.KeyEscape)
	b := h.res.Element.Bounds()
	tt.Move(b.X+52+(b.W-64)/6, b.Y+120)
	if h.res.Point != 0 {
		t.Fatalf("hover %d", h.res.Point)
	}
}

func TestViolinPlotBad(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for an unnamed group")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ViolinPlot(c, []BoxGroup{{Values: []float64{1}}}, ViolinPlotOptions{Frame: ChartFrame{Label: "x"}})
	}, 300, 300)
}

func TestViolinPlotSnapshot(t *testing.T) {
	snapshot(t, "violin_plot", 520, 320, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ViolinPlot(c, densityGroups(), ViolinPlotOptions{Frame: ChartFrame{Label: "Hymns"}})
		})
	})
}

func TestViolinPlotFiniteBandwidthAndFlatShape(t *testing.T) {
	charts2MustPanic(t, func(c *ui.Context) {
		ViolinPlot(c, []BoxGroup{{Name: "Flat", Values: []float64{4, 4}}}, ViolinPlotOptions{Frame: ChartFrame{Label: "Flat"}, Bandwidth: math.NaN()})
	})
}

func TestViolinPlotLargeDataAndDisabled(t *testing.T) {
	values := make([]float64, 100_000)
	for i := range values {
		values[i] = float64(i % 100)
	}
	groups := []BoxGroup{{Name: "Sample", Values: values}}
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return ViolinPlot(c, groups, ViolinPlotOptions{Frame: ChartFrame{Label: "Large violin"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Large violin")
	tt.Key(0, ui.KeyRight)
	if chartReadout(tt) != "Sample · median 49.5 · quartiles 24.75–74.25 · 100000 values" {
		t.Fatalf("large readout %q", chartReadout(tt))
	}
	shown := ui.Local(h.res.Element, "mujica.violin.visible", func() []int { return nil })
	if len(*shown) != 1 {
		t.Fatalf("built %d violins", len(*shown))
	}
	h.disabled = true
	tt.Frame()
	if h.res.Point != -1 {
		t.Fatalf("disabled violin %d", h.res.Point)
	}
}
