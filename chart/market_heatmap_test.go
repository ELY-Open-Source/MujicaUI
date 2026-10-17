package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestMarketHeatmapAreaDrillAndKeys(t *testing.T) {
	root := mktHeat()
	var path []string
	built := false
	tt := chatTester(640, 400, func(c *ui.Context) {
		built = MarketHeatmap(c, &path, root, MarketHeatmapOptions{Frame: ChartFrame{Label: "Market"}}).Element.Valid()
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	if !tt.HasText("Technology · up +1.18% · 8,000 · Enter opens") {
		t.Fatalf("readout %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if len(path) != 1 || path[0] != "Technology" {
		t.Fatalf("path %v", path)
	}
	tt.Key(0, ui.KeyBackspace)
	if len(path) != 0 {
		t.Fatalf("back %v", path)
	}
	if !built {
		t.Fatal("element missing")
	}
	if root.total() != 9830 {
		t.Fatalf("cap %v", root.total())
	}
	marketsPanic(t, func(c *ui.Context) {
		MarketHeatmap(c, &path, HeatNode{Name: "Bad", Cap: -1}, MarketHeatmapOptions{Frame: ChartFrame{Label: "x"}})
	})
}
func TestMarketHeatmapSnapshot(t *testing.T) {
	var path []string
	snapshot(t, "market_heatmap", 640, 400, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			MarketHeatmap(c, &path, mktHeat(), MarketHeatmapOptions{Frame: ChartFrame{Label: "S&P 500", Height: 280}})

		})
	})
}
