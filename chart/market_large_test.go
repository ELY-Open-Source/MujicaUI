package chart

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestDepthChartLargeDataBounded(t *testing.T) {
	bids, asks := make([]DepthLevel, 100000), make([]DepthLevel, 100000)
	for i := range bids {
		bids[i] = DepthLevel{100 - float64(i)*0.0001, 1}
		asks[i] = DepthLevel{101 + float64(i)*0.0001, 2}
	}
	var res ChartResult
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res = DepthChart(c, bids, asks, DepthChartOptions{Frame: ChartFrame{Label: "Depth"}})
	}, 640, 340)
	drawn := *ui.Local(res.Element, "mujica.depth.drawn", func() int { return 0 })
	if drawn <= 0 || drawn > 4*640+8 {
		t.Fatalf("painted %d depth vertices", drawn)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if res.Point != 199999 {
		t.Fatalf("last point %d", res.Point)
	}
}

func TestVolumeProfileLargeDataBounded(t *testing.T) {
	bins := make([]ProfileBin, 100000)
	for i := range bins {
		bins[i] = ProfileBin{float64(i), float64(i + 1), float64(i%19 + 1)}
	}
	var res ChartResult
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res = VolumeProfile(c, bins, VolumeProfileOptions{Frame: ChartFrame{Label: "Profile"}})
	}, 640, 340)
	drawn := *ui.Local(res.Element, "mujica.profile.drawn", func() int { return 0 })
	if drawn <= 0 || drawn > 4*340+4 {
		t.Fatalf("painted %d profile bins", drawn)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if res.Point != 99999 {
		t.Fatalf("last bin %d", res.Point)
	}
}
