package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

var scatterGroups = []ScatterGroup{
	{Name: "Monks", Points: []ScatterPoint{{1, 2, 4}, {3, 5, 9}, {6, 4, 1}}},
	{Name: "Nuns", Points: []ScatterPoint{{2, 7, 16}, {8, 9, 2}}},
}

func TestScatterChartKeyboardSelect(t *testing.T) {
	var sel []ScatterKey
	values := 0
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		r := ScatterChart(c, &sel, scatterGroups, ScatterChartOptions{Frame: ChartFrame{Label: "Choir"}, Bubbles: true})
		if r.ValueChanged() {
			values++
		}
		return r
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Choir")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	if fmt.Sprint(sel) != "[{0 1}]" || values != 1 {
		t.Fatalf("Space selected %v (%d events)", sel, values)
	}
	if got := chartReadout(tt); got != "Monks · x 3 · y 5 · size 9 · selected" {
		t.Fatalf("readout %q", got)
	}
	tt.Key(0, ui.KeyEnter)
	if len(sel) != 0 {
		t.Fatalf("Enter did not unselect: %v", sel)
	}
}

func TestScatterChartBoxSelect(t *testing.T) {
	var sel []ScatterKey
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return ScatterChart(c, &sel, scatterGroups, ScatterChartOptions{Frame: ChartFrame{Label: "Choir"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	b := h.res.Element.Bounds()
	top := b.Y + 38
	// Nice domains: x 0..10 over 52..W-12; y 0..10 over top+8..top+214.
	x := func(v float32) float32 { return b.X + 52 + (b.W-64)*v/10 }
	y := func(v float32) float32 { return top + 8 + 206*(1-v/10) }
	tt.Press(x(0.5), y(6))
	tt.Move(x(4), y(4))
	tt.Move(x(6.5), y(1.5))
	tt.Release(x(6.5), y(1.5))
	if fmt.Sprint(sel) != "[{0 0} {0 1} {0 2}]" {
		t.Fatalf("box chose %v", sel)
	}
}

func TestScatterChartBad(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for negative size")
		}
	}()
	var sel []ScatterKey
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ScatterChart(c, &sel, []ScatterGroup{{Name: "a", Points: []ScatterPoint{{1, 1, -1}}}}, ScatterChartOptions{Frame: ChartFrame{Label: "x"}})
	}, 300, 300)
}

func TestScatterChartSnapshot(t *testing.T) {
	snapshot(t, "scatter_chart", 520, 340, func(c *ui.Context) {
		sel := []ScatterKey{{1, 0}}
		ui.Column(c).Padding(16).Children(func() {
			ScatterChart(c, &sel, scatterGroups, ScatterChartOptions{Frame: ChartFrame{Label: "Choir"}, Bubbles: true})
		})
	})
}
