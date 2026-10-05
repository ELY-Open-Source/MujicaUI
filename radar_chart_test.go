package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

var radarAxes = []string{"Voice", "Latin", "Chant", "Organ", "Bells"}

func radarSeries() []ChartSeries {
	return []ChartSeries{
		{Name: "Agnès", Values: []float64{8, 6, 9, 4, 7}},
		{Name: "Clare", Values: []float64{5, 9, 6, 8, 3}},
	}
}

func TestRadarChartReadout(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return RadarChart(c, radarAxes, radarSeries(), RadarChartOptions{Frame: ChartFrame{Label: "Skills"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Skills")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "Latin · Agnès 6 · Clare 9" {
		t.Fatalf("readout %q", got)
	}
	b := h.res.Element.Bounds()
	cx, cy := b.X+b.W/2, b.Y+38+120+4
	tt.Key(0, ui.KeyEscape)
	tt.Move(cx, cy-60) // straight up: the first axis
	if got := chartReadout(tt); got != "Voice · Agnès 8 · Clare 5" {
		t.Fatalf("hover %q", got)
	}
	if err := tt.Click("Clare"); err != nil {
		t.Fatal(err)
	}
	tt.Move(cx, cy-60)
	if got := chartReadout(tt); got != "Voice · Agnès 8" {
		t.Fatalf("hidden series read: %q", got)
	}
}

func TestRadarChartBad(t *testing.T) {
	for name, build := range map[string]func(c *ui.Context){
		"two axes": func(c *ui.Context) {
			RadarChart(c, radarAxes[:2], nil, RadarChartOptions{Frame: ChartFrame{Label: "x"}})
		},
		"negative": func(c *ui.Context) {
			RadarChart(c, radarAxes[:3], []ChartSeries{{Name: "a", Values: []float64{1, -1, 2}}}, RadarChartOptions{Frame: ChartFrame{Label: "x"}})
		},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); build(c) }, 300, 300)
		}()
	}
}

func TestRadarChartSnapshot(t *testing.T) {
	snapshot(t, "radar_chart", 520, 340, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			RadarChart(c, radarAxes, radarSeries(), RadarChartOptions{Frame: ChartFrame{Label: "Skills"}})
		})
	})
}
