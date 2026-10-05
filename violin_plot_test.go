package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
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
		Use(c, Settings{})
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
