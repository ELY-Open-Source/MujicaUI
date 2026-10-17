package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestBulletChartReadout(t *testing.T) {
	target := 250.0
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return BulletChart(c, 270, BulletChartOptions{Frame: ChartFrame{Label: "Alms"}, Target: &target, Bands: []float64{150, 225, 300}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Alms")
	want := []string{"Actual 270", "Target 250", "Range 1: 0–150", "Range 2: 150–225", "Range 3: 225–300"}
	for i, w := range want {
		tt.Key(0, ui.KeyRight)
		if got := chartReadout(tt); got != w {
			t.Fatalf("mark %d %q, want %q", i, got, w)
		}
	}
	tt.Key(0, ui.KeyEscape)
	b := h.box(tt)
	plotW := b.W - 24
	tt.Move(b.X+8+plotW*0.1, b.Y+6+21) // band 1, inside the value bar
	if got := chartReadout(tt); got != "Actual 270" {
		t.Fatalf("hover bar %q", got)
	}
	tt.Move(b.X+8+plotW*0.1, b.Y+8) // band 1, above the bar
	if got := chartReadout(tt); got != "Range 1: 0–150" {
		t.Fatalf("hover band %q", got)
	}
}

func TestBulletChartVerticalAndBad(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return BulletChart(c, 40, BulletChartOptions{Frame: ChartFrame{Label: "Alms", Height: 200}, Vertical: true, Max: 100})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	b := h.box(tt)
	tt.Move(b.X+b.W/2, b.Y+8+184*0.8) // low on the scale: inside the bar
	if got := chartReadout(tt); got != "Actual 40" {
		t.Fatalf("vertical hover %q", got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for unsorted bands")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		BulletChart(c, 1, BulletChartOptions{Frame: ChartFrame{Label: "x"}, Bands: []float64{5, 2}})
	}, 300, 300)
}

func TestBulletChartSnapshot(t *testing.T) {
	target := 250.0
	snapshot(t, "bullet_chart", 520, 300, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			ui.Column(c).Width(360).Children(func() {
				BulletChart(c, 270, BulletChartOptions{Frame: ChartFrame{Label: "Alms"}, Target: &target, Bands: []float64{150, 225, 300}})
			})
			ui.Column(c).Width(100).Children(func() {
				BulletChart(c, 40, BulletChartOptions{Frame: ChartFrame{Label: "Oil", Height: 200}, Vertical: true, Target: &[]float64{70}[0], Bands: []float64{50, 80, 100}})
			})
		})
	})
}
