package chart

import (
	"math"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var barCats = []string{"Nave", "Choir", "Crypt"}

func barSeries() []ChartSeries {
	return []ChartSeries{
		{Name: "Gold", Values: []float64{4, -2, 3}},
		{Name: "Silver", Values: []float64{2, 5, math.NaN()}},
	}
}

func TestBarChartKeyboardPerBar(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return BarChart(c, barCats, barSeries(), BarChartOptions{Frame: ChartFrame{Label: "Tithes"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	tabTo(t, tt, "Tithes")
	want := []string{"Nave · Gold 4", "Nave · Silver 2", "Choir · Gold -2", "Choir · Silver 5", "Crypt · Gold 3"}
	for i, w := range want {
		tt.Key(0, ui.KeyRight)
		if got := chartReadout(tt); got != w {
			t.Fatalf("bar %d readout %q, want %q", i, got, w)
		}
	}
	tt.Key(0, ui.KeyRight)
	if h.res.Point != 4 {
		t.Fatalf("stepped past the missing bar: %d", h.res.Point)
	}
}

func TestBarChartHoverStackedHorizontal(t *testing.T) {
	for _, o := range []BarChartOptions{{Stacked: true}, {Horizontal: true, Labels: true}, {Horizontal: true, Stacked: true}} {
		o.Frame = ChartFrame{Label: "Tithes"}
		h := &chartHarness{build: func(c *ui.Context) ChartResult { return BarChart(c, barCats, barSeries(), o) }}
		tt := ui.NewTester(h.view, 520, 420)
		tabTo(t, tt, "Tithes")
		found := map[string]bool{}
		b := h.res.Element.Bounds()
		for y := b.Y + 30; y < b.Y+b.H-40; y += 4 {
			for x := b.X + 60; x < b.X+b.W-16; x += 8 {
				tt.Move(x, y)
				if h.res.Point >= 0 {
					found[chartReadout(tt)] = true
				}
			}
		}
		if len(found) != 5 || !found["Choir · Gold -2"] {
			t.Fatalf("%+v: hover found %v", o, found)
		}
	}
}

func TestBarChartBadAndEmpty(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for mismatched categories")
		}
	}()
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return BarChart(c, barCats, nil, BarChartOptions{Frame: ChartFrame{Label: "Tithes"}})
	}}
	tt := ui.NewTester(h.view, 520, 420)
	if !tt.HasText("No data") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		BarChart(c, barCats[:2], barSeries(), BarChartOptions{Frame: ChartFrame{Label: "x"}})
	}, 300, 300)
}

func TestBarChartSnapshot(t *testing.T) {
	snapshot(t, "bar_chart", 640, 320, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			ui.Column(c).Width(296).Children(func() {
				BarChart(c, barCats, barSeries(), BarChartOptions{Frame: ChartFrame{Label: "Grouped", Height: 220}, Labels: true})
			})
			ui.Column(c).Width(296).Children(func() {
				BarChart(c, barCats, barSeries(), BarChartOptions{Frame: ChartFrame{Label: "Stacked", Height: 220}, Horizontal: true, Stacked: true})
			})
		})
	})
}
