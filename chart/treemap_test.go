package chart

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func chartTree() ChartNode {
	return ChartNode{Name: "Abbey", Children: []ChartNode{
		{Name: "Library", Children: []ChartNode{{Name: "Psalters", Value: 30}, {Name: "Maps", Value: 10}}},
		{Name: "Kitchen", Value: 25},
		{Name: "Garden", Children: []ChartNode{{Name: "Herbs", Value: 12}, {Name: "Vines", Value: 8}}},
		{Name: "Cell of a very long-winded brother", Value: 3},
	}}
}

func TestSquarifyFills(t *testing.T) {
	rects := squarify([]float64{6, 6, 4, 3, 2, 2, 1}, ui.Rect{W: 600, H: 400})
	area := float32(0)
	for _, r := range rects {
		area += r.W * r.H
		if r.W <= 0 || r.H <= 0 {
			t.Fatalf("empty tile %v", r)
		}
	}
	if area < 600*400-1 || area > 600*400+1 {
		t.Fatalf("tiles cover %v of %v", area, 600*400)
	}
}

func TestTreemapDrillAndBack(t *testing.T) {
	var path []string
	values := 0
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		r := Treemap(c, &path, chartTree(), TreemapOptions{Frame: ChartFrame{Label: "Rooms"}})
		if r.ValueChanged() {
			values++
		}
		return r
	}}
	tt := ui.NewTester(h.view, 560, 440)
	tabTo(t, tt, "Rooms")
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "Library · 40 · 45.5% · Enter opens" {
		t.Fatalf("readout %q", got)
	}
	tt.Key(0, ui.KeyEnter)
	if fmt.Sprint(path) != "[Library]" || values != 1 || !tt.HasText("Library") {
		t.Fatalf("drill path %v events %d", path, values)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter) // a leaf stays
	if fmt.Sprint(path) != "[Library]" {
		t.Fatalf("leaf drilled: %v", path)
	}
	tt.Key(0, ui.KeyBackspace)
	if len(path) != 0 {
		t.Fatalf("Backspace kept %v", path)
	}
	path = []string{"Garden"}
	tt.Frame()
	if err := tt.Click("Back"); err != nil || len(path) != 0 {
		t.Fatalf("Back: %v %v", err, path)
	}
}

func TestTreemapBadPath(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for an unknown path")
		}
	}()
	path := []string{"Nowhere"}
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Treemap(c, &path, chartTree(), TreemapOptions{Frame: ChartFrame{Label: "x"}})
	}, 300, 300)
}

func TestTreemapSnapshot(t *testing.T) {
	snapshot(t, "treemap", 560, 340, func(c *ui.Context) {
		var path []string
		ui.Column(c).Padding(16).Children(func() {
			Treemap(c, &path, chartTree(), TreemapOptions{Frame: ChartFrame{Label: "Rooms"}})
		})
	})
}
