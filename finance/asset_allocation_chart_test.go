package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestAssetAllocationChartWeightsDriftAndKeys(t *testing.T) {
	items := []AllocationItem{{"Stocks", 60, 50}, {"Bonds", 30, 40}, {"Cash", 10, 10}}
	var last chart.ChartResult
	tt := chatTester(960, 420, func(c *ui.Context) {
		last = AssetAllocationChart(c, items, AssetAllocationChartOptions{Frame: chart.ChartFrame{Label: "Allocation"}})
	})
	if !tt.HasText("60.0%") || !tt.HasText("+10.0 pt") || !tt.HasText("−10.0 pt") || !tt.HasText("Over") || !tt.HasText("Under") {
		t.Fatalf("weights %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyHome)
	if last.Point != 0 {
		t.Fatalf("point %d", last.Point)
	}
	tt.Key(0, ui.KeyEnd)
	if last.Point != 2 {
		t.Fatalf("point %d", last.Point)
	}
	marketsPanic(t, func(c *ui.Context) {
		AssetAllocationChart(c, []AllocationItem{{"Bad", 1, 101}}, AssetAllocationChartOptions{Frame: chart.ChartFrame{Label: "x"}})
	})
}
func TestAssetAllocationChartSnapshot(t *testing.T) {

	snapshot(t, "asset_allocation_chart", 760, 380, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			AssetAllocationChart(c, []AllocationItem{{"Stocks", 62000, 55}, {"Bonds", 25000, 30}, {"Cash", 8000, 10}, {"Gold", 5000, 5}}, AssetAllocationChartOptions{Frame: chart.ChartFrame{Label: "Allocation"}})

		})
	})
}

func TestAssetAllocationChartMergedMarkers(t *testing.T) {
	items := []AllocationItem{{"Tiny A", 0.1, 0}, {"Tiny B", 0.1, 0}, {"Stocks", 99.8, 0}}
	tt := chatTester(960, 420, func(c *ui.Context) {
		AssetAllocationChart(c, items, AssetAllocationChartOptions{Frame: chart.ChartFrame{Label: "Allocation"}})
	})
	r, _ := tt.Find("99.8%")
	asset, _ := tt.Find("Asset")
	palette := internal.Palette(theme.Light())
	found := false
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(asset.X); x < int(r.X); x++ {
			pixel := tt.Image().RGBAAt(x, y)
			if pixel.R == palette[0].R && pixel.G == palette[0].G && pixel.B == palette[0].B {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("Stocks list marker differs from its donut marker")
	}
}
