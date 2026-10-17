package chart

import (
	"fmt"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

var boxGroups = []BoxGroup{
	{Name: "Choir", Values: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 40}},
	{Name: "Nave", Values: []float64{4, 5, 5, 6, 7, math.NaN()}},
}

func TestBoxSummary(t *testing.T) {
	s := boxSummary(boxGroups[0].Values)
	if s.q1 != 3.25 || s.median != 5.5 || s.q3 != 7.75 || s.lo != 1 || s.hi != 9 || fmt.Sprint(s.outliers) != "[40]" || s.n != 10 {
		t.Fatalf("summary %+v", s)
	}
	if s := boxSummary([]float64{math.NaN()}); s.n != 0 {
		t.Fatal("NaN counted")
	}
}

func TestBoxPlotReadout(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return BoxPlot(c, boxGroups, BoxPlotOptions{Frame: ChartFrame{Label: "Hymns"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Hymns")
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "Choir · median 5.5 · quartiles 3.25–7.75 · whiskers 1–9 · 1 outliers" {
		t.Fatalf("readout %q", got)
	}
	b := h.box(tt)
	tt.Key(0, ui.KeyEscape)
	tt.Move(b.X+52+(b.W-64)*0.75, b.Y+100)
	if got := chartReadout(tt); got != "Nave · median 5 · quartiles 5–6 · whiskers 4–7 · 0 outliers" {
		t.Fatalf("hover %q", got)
	}
}

func TestBoxPlotEmpty(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return BoxPlot(c, []BoxGroup{{Name: "a"}}, BoxPlotOptions{Frame: ChartFrame{Label: "Hymns"}})
	}}
	if tt := ui.NewTester(h.view, 520, 420); !tt.HasText("No data") {
		t.Fatalf("empty: %q", tt.Texts())
	}
}

func TestBoxPlotSnapshot(t *testing.T) {
	snapshot(t, "box_plot", 520, 320, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			BoxPlot(c, append(boxGroups, BoxGroup{Name: "Crypt", Values: []float64{12, 14, 15, 18, 22, 25}}), BoxPlotOptions{Frame: ChartFrame{Label: "Hymns"}})
		})
	})
}
