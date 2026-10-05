package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSunburstReadoutAndDrill(t *testing.T) {
	var path []string
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return Sunburst(c, &path, chartTree(), SunburstOptions{Frame: ChartFrame{Label: "Rooms"}})
	}}
	tt := ui.NewTester(h.view, 560, 440)
	tabTo(t, tt, "Rooms")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "Library › Psalters · 30 · 34.1%" {
		t.Fatalf("readout %q", got)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeySpace)
	if fmt.Sprint(path) != "[Library]" || !tt.HasText("Abbey") {
		t.Fatalf("drill %v", path)
	}
	if err := tt.Click("Abbey"); err != nil || len(path) != 0 {
		t.Fatalf("breadcrumb back: %v %v", err, path)
	}
	b := h.res.Element.Bounds()
	plotTop := b.Y + b.H - 20 - 8 - 240
	cx, cy := b.X+b.W/2, plotTop+120
	tt.Key(0, ui.KeyEscape)
	tt.Move(cx+50, cy-30) // inner ring, upper right: Library
	if got := chartReadout(tt); got != "Library · 40 · 45.5% · Enter opens" {
		t.Fatalf("hover %q", got)
	}
	tt.Move(cx-48*0.866, cy-48*0.5) // inner ring at 300°: Garden spans 266°–348°
	if got := chartReadout(tt); got != "Garden · 20 · 22.7% · Enter opens" {
		t.Fatalf("hover at 300° %q", got)
	}
}

func TestSunburstBad(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for duplicate children")
		}
	}()
	var path []string
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		Sunburst(c, &path, ChartNode{Name: "r", Children: []ChartNode{{Name: "a", Value: 1}, {Name: "a", Value: 2}}}, SunburstOptions{Frame: ChartFrame{Label: "x"}})
	}, 300, 300)
}

func TestSunburstSnapshot(t *testing.T) {
	snapshot(t, "sunburst", 560, 360, func(c *ui.Context) {
		var path []string
		ui.Column(c).Padding(16).Children(func() {
			Sunburst(c, &path, chartTree(), SunburstOptions{Frame: ChartFrame{Label: "Rooms", Height: 260}})
		})
	})
}
