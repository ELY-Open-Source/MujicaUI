package chart

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var sankeyNodes = []SankeyNode{{"tithe", "Tithes"}, {"alms", "Alms"}, {"fund", "Fund"}, {"roof", "Roof"}, {"bread", "Bread"}, {"bells", "Bells"}}
var sankeyLinks = []SankeyLink{
	{"tithe", "fund", 60}, {"alms", "fund", 30}, {"alms", "bread", 10},
	{"fund", "roof", 50}, {"fund", "bells", 40},
}

func TestSankeyPlaceStable(t *testing.T) {
	l := sankeyPlace(sankeyNodes, sankeyLinks)
	if fmt.Sprint(l.col) != "[0 0 1 2 1 2]" || l.cols != 3 || l.in[2] != 90 || l.out[1] != 40 {
		t.Fatalf("layout %+v", l)
	}
	lit := sankeyLit(l, len(sankeyNodes), 0) // Tithes: its fund, roof and bells, not alms
	if !lit[2] || !lit[3] || !lit[5] || lit[1] || lit[4] || !lit[6] || lit[7] {
		t.Fatalf("lit from Tithes %v", lit)
	}
	lit = sankeyLit(l, len(sankeyNodes), 6+3) // the fund → roof link
	if !lit[0] || !lit[1] || !lit[3] || lit[5] || lit[6+4] {
		t.Fatalf("lit from fund → roof %v", lit)
	}
}

func TestSankeyChartReadout(t *testing.T) {
	h := &chartHarness{build: func(c *ui.Context) ChartResult {
		return SankeyChart(c, sankeyNodes, sankeyLinks, SankeyChartOptions{Frame: ChartFrame{Label: "Money"}})
	}}
	tt := ui.NewTester(h.view, 560, 440)
	tabTo(t, tt, "Money")
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	if got := chartReadout(tt); got != "Fund · in 90 · out 90" {
		t.Fatalf("node %q", got)
	}
	tt.Key(0, ui.KeyEnd)
	if got := chartReadout(tt); got != "Fund → Bells · 40" {
		t.Fatalf("link %q", got)
	}
	tt.Key(0, ui.KeyEscape)
	b := h.res.Element.Bounds()
	tt.Move(b.X+8+4, b.Y+10) // the first node, at the top left
	if got := chartReadout(tt); got != "Tithes · in 0 · out 60" {
		t.Fatalf("hover %q", got)
	}
}

func TestSankeyChartBad(t *testing.T) {
	for name, links := range map[string][]SankeyLink{
		"cycle":   {{"tithe", "fund", 1}, {"fund", "tithe", 1}},
		"unknown": {{"tithe", "nowhere", 1}},
		"zero":    {{"tithe", "fund", 0}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				SankeyChart(c, sankeyNodes, links, SankeyChartOptions{Frame: ChartFrame{Label: "x"}})
			}, 300, 300)
		}()
	}
}

func TestSankeyChartSnapshot(t *testing.T) {
	snapshot(t, "sankey_chart", 560, 320, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			SankeyChart(c, sankeyNodes, sankeyLinks, SankeyChartOptions{Frame: ChartFrame{Label: "Money"}})
		})
	})
}
