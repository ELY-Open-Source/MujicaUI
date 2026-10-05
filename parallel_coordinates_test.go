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
	// Row 2 crosses this point between the first axes.
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

func TestParallelCoordinatesKeyboardBrushAndDisabled(t *testing.T) {
	filters := map[string][2]float64{}
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return ParallelCoordinates(c, &filters, pcAxes, pcRows, ParallelCoordinatesOptions{Frame: ChartFrame{Label: "Monks"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Monks")
	tt.Key(ui.Alt, ui.KeyUp)
	f, ok := filters["Age"]
	if !ok || math.Abs(f[0]-20.6) > 1e-9 || f[1] != 80 {
		t.Fatalf("keyboard brush %v", filters)
	}
	h.disabled = true
	tt.Frame()
	tt.Key(ui.Alt, ui.KeyUp)
	if filters["Age"] != f {
		t.Fatal("disabled keyboard brush")
	}
	h.disabled = false
	tt.Frame()
	tabTo(t, tt, "Monks")
	tt.Key(ui.Alt, ui.KeyDelete)
	if len(filters) != 0 {
		t.Fatalf("clear brush %v", filters)
	}
}

func TestParallelCoordinatesLargeDataBounded(t *testing.T) {
	rows := make([][]float64, 100_000)
	for i := range rows {
		rows[i] = []float64{float64(i % 70), float64(i % 11), float64(i % 17)}
	}
	filters := map[string][2]float64{}
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return ParallelCoordinates(c, &filters, pcAxes, rows, ParallelCoordinatesOptions{Frame: ChartFrame{Label: "Large parallel"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	shown := ui.Local(h.res.Element, "mujica.pc.visible", func() []int { return nil })
	if len(*shown) == 0 || len(*shown) > 480 {
		t.Fatalf("built %d rows", len(*shown))
	}
	tabTo(t, tt, "Large parallel")
	tt.Key(0, ui.KeyEnd)
	if h.res.Point != 99_999 || len(*shown) > 481 {
		t.Fatalf("last row %d, built %d", h.res.Point, len(*shown))
	}
}
