package chart

import (
	"fmt"
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

type brushHarness struct {
	value    [2]float64
	changes  int
	disabled bool
	bounds   ui.Rect
}

func (h *brushHarness) view(c *ui.Context) {
	core.Use(c, core.Settings{Mode: core.Light})
	ui.Column(c).Padding(10).Width(420).Children(func() {
		r := ChartBrush(c, &h.value, ChartBrushOptions{Label: "Months", Min: 0, Max: 100, Step: 1, Disabled: h.disabled})
		h.bounds = r.Element.Bounds()
		if r.Changed() {
			h.changes++
		}
	})
}

func TestChartBrushKeyboard(t *testing.T) {
	h := &brushHarness{value: [2]float64{20, 40}}
	tt := ui.NewTester(h.view, 460, 200)
	tabTo(t, tt, "Months: 20 to 40")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if h.value != [2]float64{22, 42} || h.changes != 2 {
		t.Fatalf("window moved to %v (%d)", h.value, h.changes)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Range end 42") {
		t.Fatal("Tab did not reach the end handle")
	}
	tt.Key(0, ui.KeyPageUp)
	tt.Key(0, ui.KeyEnd)
	if h.value != [2]float64{22, 100} {
		t.Fatalf("end handle: %v", h.value)
	}
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(0, ui.KeyEnd) // the start cannot pass the end
	if h.value != [2]float64{99, 100} || !tt.HasText("99 – 100") {
		t.Fatalf("start to End: %v %q", h.value, tt.Texts())
	}
}

func TestChartBrushDrag(t *testing.T) {
	h := &brushHarness{value: [2]float64{20, 40}}
	tt := ui.NewTester(h.view, 460, 200)
	tt.Frame() // Bounds reads the previous frame
	b := h.bounds
	x := func(v float32) float32 { return b.X + 400*v/100 }
	y := b.Y + 24
	tt.Press(x(30), y) // inside the window
	tt.Move(x(40), y)
	tt.Move(x(50), y)
	tt.Release(x(50), y)
	if h.value != [2]float64{40, 60} {
		t.Fatalf("window drag: %v", h.value)
	}
	tt.Press(x(60), y) // the end handle
	tt.Move(x(70), y)
	tt.Move(x(75), y)
	tt.Release(x(75), y)
	if h.value != [2]float64{40, 75} {
		t.Fatalf("end drag: %v", h.value)
	}
	tt.Press(x(85), y) // outside: a new range
	tt.Move(x(90), y)
	tt.Move(x(95), y)
	tt.Release(x(95), y)
	if h.value != [2]float64{85, 95} {
		t.Fatalf("new range: %v", h.value)
	}
}

func TestChartBrushDisabledAndBad(t *testing.T) {
	h := &brushHarness{value: [2]float64{20, 40}, disabled: true}
	tt := ui.NewTester(h.view, 460, 200)
	tt.Frame() // Bounds reads the previous frame
	b := h.bounds
	tt.Press(b.X+120, b.Y+24)
	tt.Move(b.X+200, b.Y+24)
	tt.Release(b.X+200, b.Y+24)
	for range 3 {
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyRight)
	}
	if h.value != [2]float64{20, 40} || h.changes != 0 {
		t.Fatalf("disabled brush moved: %v", h.value)
	}
	for _, v := range [][2]float64{{50, 40}, {-1, 10}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%v: no panic", v)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ChartBrush(c, &v, ChartBrushOptions{Label: "x", Max: 100})
			}, 300, 100)
		}()
	}
}

func TestChartBrushSnapshot(t *testing.T) {
	vals := make([]float64, 60)
	for i := range vals {
		vals[i] = float64((i*7)%13) + float64(i)/6
	}
	snapshot(t, "chart_brush", 460, 110, func(c *ui.Context) {
		v := [2]float64{18, 42}
		ui.Column(c).Padding(16).Children(func() {
			ChartBrush(c, &v, ChartBrushOptions{Label: "Months", Min: 0, Max: 59, Step: 1, Values: vals, Format: func(f float64) string { return fmt.Sprintf("week %.0f", f) }})
		})
	})
}

func TestChartBrushWindowDragPreservesSpanAtEdges(t *testing.T) {
	h := &brushHarness{value: [2]float64{20, 60}}
	tt := ui.NewTester(h.view, 460, 200)
	tt.Frame() // Bounds reads the previous frame
	b := h.bounds
	tt.Press(b.X+160, b.Y+24)
	tt.Move(b.X+380, b.Y+24)
	tt.Release(b.X+380, b.Y+24)
	if h.value != [2]float64{60, 100} {
		t.Fatalf("right edge shrank window: %v", h.value)
	}
	tt.Press(b.X+320, b.Y+24)
	tt.Move(b.X-50, b.Y+24)
	tt.Release(b.X-50, b.Y+24)
	if h.value != [2]float64{0, 40} {
		t.Fatalf("left edge shrank window: %v", h.value)
	}
}

func TestChartBrushInvalidMetrics(t *testing.T) {
	for _, o := range []ChartBrushOptions{
		{Label: "Range", Max: 100, Step: math.Inf(1)},
		{Label: "Range", Max: math.Inf(1)},
		{Label: "Range", Max: 100, Height: 10},
	} {
		charts2MustPanic(t, func(c *ui.Context) { v := [2]float64{20, 40}; ChartBrush(c, &v, o) })
	}
}

func TestChartBrushNarrowRangeHitAndAnnouncement(t *testing.T) {
	h := &brushHarness{value: [2]float64{40, 41}}
	tt := ui.NewTester(h.view, 460, 200)
	b, ok := tt.Find("Months: 40 to 41")
	if !ok || b.W < 28 || b.H < 28 {
		t.Fatalf("window hit area %v", b)
	}
	tabTo(t, tt, "Months: 40 to 41")
	tt.Key(0, ui.KeyRight)
	if said := tt.Announcements(); len(said) == 0 || said[len(said)-1] != "Months: 41 to 42" {
		t.Fatalf("window announcement %q", said)
	}
}
