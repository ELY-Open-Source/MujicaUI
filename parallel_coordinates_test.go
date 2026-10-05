package mujica

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

var pcAxes = []string{"Age", "Hymns", "Vigils"}
var pcRows = [][]float64{{20, 3, 9}, {35, 8, 2}, {50, 5, 5}, {65, 1, math.NaN()}}
var pcNames = []string{"Agnès", "Clare", "Hild", "Bede"}

func TestParallelCoordinatesReadoutAndBrush(t *testing.T) {
	filters := map[string][2]float64{}
	values := 0
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		r := ParallelCoordinates(c, &filters, pcAxes, pcRows, ParallelCoordinatesOptions{Frame: ChartFrame{Label: "Monks"}, Names: pcNames})
		if r.ValueChanged() {
			values++
		}
		return r
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Monks")
	tt.Key(0, ui.KeyEnd)
	if got := chartReadout(tt); got != "Bede · Age 65 · Hymns 1 · Vigils no value" {
		t.Fatalf("readout %q", got)
	}
	b := h.res.Element.Bounds()
	plot := ui.Rect{X: b.X + 36, Y: b.Y + 24, W: b.W - 72, H: 240 - 32}
	// Age runs 20..80 on the first axis: brush 30..55.
	y := func(v float32) float32 { return plot.Y + plot.H*(1-(v-20)/60) }
	tt.Press(plot.X, y(55))
	tt.Move(plot.X, y(40))
	tt.Move(plot.X, y(30))
	tt.Release(plot.X, y(30))
	f, ok := filters["Age"]
	if !ok || math.Abs(f[0]-30) > 0.5 || math.Abs(f[1]-55) > 0.5 || values != 1 {
		t.Fatalf("brush %v %v (%d)", f, ok, values)
	}
	tt.Key(0, ui.KeyHome)
	if got := chartReadout(tt); got != "Agnès · Age 20 · Hymns 3 · Vigils 9 · filtered out" {
		t.Fatalf("filtered readout %q", got)
	}
	tt.Press(plot.X, y(45)) // a click on the axis clears it
	tt.Release(plot.X, y(45))
	if _, ok := filters["Age"]; ok {
		t.Fatal("click did not clear the brush")
	}
}

func TestParallelCoordinatesHoverAndBad(t *testing.T) {
	filters := map[string][2]float64{}
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return ParallelCoordinates(c, &filters, pcAxes, pcRows, ParallelCoordinatesOptions{Frame: ChartFrame{Label: "Monks"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	b := h.res.Element.Bounds()
	plot := ui.Rect{X: b.X + 36, Y: b.Y + 24, W: b.W - 72, H: 240 - 32}
	// A quarter into the first gap, row 2 (Age 50 of 20..80, Hymns 5 of 0..8) is 0.53 up.
	tt.Move(plot.X+plot.W/8, plot.Y+plot.H*(1-0.53125))
	if h.res.Point != 2 {
		t.Fatalf("hover picked %d", h.res.Point)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a short row")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ParallelCoordinates(c, &filters, pcAxes, [][]float64{{1}}, ParallelCoordinatesOptions{Frame: ChartFrame{Label: "x"}})
	}, 300, 300)
}

func TestParallelCoordinatesSnapshot(t *testing.T) {
	snapshot(t, "parallel_coordinates", 520, 320, func(c *ui.Context) {
		filters := map[string][2]float64{"Hymns": {2, 9}}
		ui.Column(c).Padding(16).Children(func() {
			ParallelCoordinates(c, &filters, pcAxes, pcRows, ParallelCoordinatesOptions{Frame: ChartFrame{Label: "Monks"}, Names: pcNames})
		})
	})
}
