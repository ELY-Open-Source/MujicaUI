package chart

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

var pieSlices = []ChartSlice{{"Wax", 50}, {"Oil", 30}, {"Myrrh", 12}, {"Salt", 4}, {"Ash", 4}}

func TestPieMerge(t *testing.T) {
	got := PieMerge(pieSlices, 20, "Other")
	if fmt.Sprint(got) != "[{Wax 50} {Oil 30} {Myrrh 12} {Other 8}]" {
		t.Fatalf("merge 20°: %v", got)
	}
	if got := PieMerge(pieSlices[:4], 20, "Other"); len(got) != 4 {
		t.Fatalf("a single narrow slice merged: %v", got)
	}
}

func TestPieChartReadoutAndHover(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return PieChart(c, pieSlices, PieChartOptions{Frame: ChartFrame{Label: "Stores"}, Donut: true, MinAngle: 20})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	if !tt.HasText("Other") {
		t.Fatalf("legend lacks Other: %q", tt.Texts())
	}
	tabTo(t, tt, "Stores")
	tt.Key(0, ui.KeyEnd)
	if got := chartReadout(tt); got != "Other · 8 · 8.0%" {
		t.Fatalf("readout %q", got)
	}
	b := h.box(tt)
	plotTop := b.Y + 38 // below the one-line legend
	cy := plotTop + 120
	tt.Move(b.X+b.W/2+60, cy-40) // upper right: the first slice
	if got := chartReadout(tt); got != "Wax · 50 · 50.0%" {
		t.Fatalf("hover %q", got)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Move(b.X+b.W/2, cy) // the hole
	if h.res.Point != -1 {
		t.Fatalf("hole picked %d", h.res.Point)
	}
}

func TestPieChartBad(t *testing.T) {
	for _, s := range [][]ChartSlice{{{"a", -1}}, {{"", 1}}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%v: no panic", s)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				PieChart(c, s, PieChartOptions{Frame: ChartFrame{Label: "x"}})
			}, 300, 300)
		}()
	}
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return PieChart(c, []ChartSlice{{"a", 0}}, PieChartOptions{Frame: ChartFrame{Label: "x"}})
	}}
	if tt := ui.NewTester(h.view, 300, 300); !tt.HasText("No data") {
		t.Fatal("zero total is not empty")
	}
}

func TestPieChartSnapshot(t *testing.T) {
	snapshot(t, "pie_chart", 640, 360, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			ui.Column(c).Width(296).Children(func() {
				PieChart(c, pieSlices, PieChartOptions{Frame: ChartFrame{Label: "Pie", Height: 220}})
			})
			ui.Column(c).Width(296).Children(func() {
				PieChart(c, pieSlices, PieChartOptions{Frame: ChartFrame{Label: "Donut", Height: 220}, Donut: true, MinAngle: 20})
			})
		})
	})
}
