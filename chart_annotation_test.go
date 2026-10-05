package mujica

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestChartAnnotationDrag(t *testing.T) {
	notes := []ChartNote{
		{Kind: NoteHLine, Value: 50, Text: "Goal", Draggable: true},
		{Kind: NoteVLine, Value: 3},
		{Kind: NoteBand, Value: 6, End: 8, Text: "Lent", Draggable: true},
		{Kind: NoteText, Value: 1, End: 90, Text: "Feast"},
	}
	h := &chartHarness{build: func(c *ui.Context) ChartResult { return partsChart(c, notes, false) }}
	tt := ui.NewTester(h.view, 520, 420)
	b := h.res.Element.Bounds()
	plot := chartInset(ui.Rect{X: b.X, Y: b.Y, W: b.W, H: 240}, 52, 26)
	yAt := func(v float32) float32 { return plot.Y + plot.H*(1-v/100) }
	xAt := func(v float32) float32 { return plot.X + plot.W*v/10 }
	// Drag the Goal line up by a fifth of the plot: 50 becomes 70.
	tt.Press(xAt(4), yAt(50))
	tt.Move(xAt(4), yAt(60))
	tt.Move(xAt(4), yAt(70))
	tt.Release(xAt(4), yAt(70))
	if math.Abs(notes[0].Value-70) > 0.5 {
		t.Fatalf("Goal at %v", notes[0].Value)
	}
	// The band moves as a whole and keeps its width.
	tt.Press(xAt(7), yAt(20))
	tt.Move(xAt(8), yAt(20))
	tt.Release(xAt(8), yAt(20))
	if math.Abs(notes[2].Value-7) > 0.1 || math.Abs(notes[2].End-notes[2].Value-2) > 1e-9 {
		t.Fatalf("band %v..%v", notes[2].Value, notes[2].End)
	}
	// The fixed line stays.
	tt.Press(xAt(3), yAt(20))
	tt.Move(xAt(5), yAt(20))
	tt.Release(xAt(5), yAt(20))
	if notes[1].Value != 3 {
		t.Fatalf("fixed line moved to %v", notes[1].Value)
	}
}

func TestChartAnnotationDisabledAndBad(t *testing.T) {
	notes := []ChartNote{{Kind: NoteHLine, Value: 50, Draggable: true}}
	h := &chartHarness{disabled: true, build: func(c *ui.Context) ChartResult { return partsChart(c, notes, false) }}
	tt := ui.NewTester(h.view, 520, 420)
	b := h.res.Element.Bounds()
	plot := chartInset(ui.Rect{X: b.X, Y: b.Y, W: b.W, H: 240}, 52, 26)
	tt.Press(plot.X+100, plot.Y+plot.H/2)
	tt.Move(plot.X+100, plot.Y+20)
	tt.Release(plot.X+100, plot.Y+20)
	if notes[0].Value != 50 {
		t.Fatalf("disabled chart dragged a note to %v", notes[0].Value)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for an empty band")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		partsChart(c, []ChartNote{{Kind: NoteBand, Value: 2, End: 2}}, false)
	}, 300, 300)
}

func TestChartAnnotationSnapshot(t *testing.T) {
	snapshot(t, "chart_annotation", 520, 340, func(c *ui.Context) {
		notes := []ChartNote{
			{Kind: NoteHLine, Value: 65, Text: "Goal", Draggable: true},
			{Kind: NoteVLine, Value: 3, Text: "Advent"},
			{Kind: NoteBand, Value: 6, End: 8, Text: "Lent"},
			{Kind: NoteText, Value: 1, End: 90, Text: "Feast"},
		}
		ui.Column(c).Padding(16).Children(func() { partsChart(c, notes, false) })
	})
}
