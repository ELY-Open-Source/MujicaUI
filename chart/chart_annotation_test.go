package chart

import (
	"math"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
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
	plot := ChartInset(ui.Rect{X: b.X, Y: b.Y, W: b.W, H: 240}, 52, 26)
	yAt := func(v float32) float32 { return plot.Y + plot.H*(1-v/100) }
	xAt := func(v float32) float32 { return plot.X + plot.W*v/10 }
	// Drag Goal from 50 to 70.
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
	plot := ChartInset(ui.Rect{X: b.X, Y: b.Y, W: b.W, H: 240}, 52, 26)
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
		core.Use(c, core.Settings{})
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

func TestChartAnnotationKeyboardEditsEveryKind(t *testing.T) {
	notes := []ChartNote{
		{Kind: NoteHLine, Value: 50, Draggable: true},
		{Kind: NoteVLine, Value: 3, Draggable: true},
		{Kind: NoteBand, Value: 6, End: 8, Draggable: true},
		{Kind: NoteText, Value: 1, End: 90, Text: "Feast", Draggable: true},
	}
	h := &chartHarness{build: func(c *ui.Context) ChartResult { return partsChart(c, notes, false) }}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Parts")
	tt.Key(ui.Alt, ui.KeyUp)
	if notes[0].Value != 51 {
		t.Fatalf("horizontal line %v", notes[0].Value)
	}
	if said := tt.Announcements(); len(said) == 0 || !strings.Contains(said[len(said)-1], "y 51") {
		t.Fatalf("horizontal readout %q", said)
	}
	tt.Key(ui.Alt, ui.KeyPageDown)
	tt.Key(ui.Alt, ui.KeyRight)
	if notes[1].Value != 3.1 {
		t.Fatalf("vertical line %v", notes[1].Value)
	}
	tt.Key(ui.Alt, ui.KeyPageDown)
	tt.Key(ui.Alt, ui.KeyRight)
	if math.Abs(notes[2].Value-6.1) > 1e-9 || notes[2].End-notes[2].Value != 2 {
		t.Fatalf("band %+v", notes[2])
	}
	tt.Key(ui.Alt, ui.KeyPageDown)
	tt.Key(ui.Alt, ui.KeyUp)
	if notes[3].End != 91 {
		t.Fatalf("text %+v", notes[3])
	}
	h.disabled = true
	tt.Frame()
	tt.Key(ui.Alt, ui.KeyUp)
	if notes[3].End != 91 {
		t.Fatalf("disabled text %+v", notes[3])
	}
}
