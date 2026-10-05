package mujica

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestVolumeChartSharesViewAndReads(t *testing.T) {
	candles := sampleCandles(100000)
	view := MarketView{Count: 60}
	var price, vol ChartResult
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Width(800).Children(func() {
			price = CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "Price"}})
			vol = VolumeChart(c, candles, &view, VolumeChartOptions{Frame: ChartFrame{Label: "Volume"}})
		})
	}, 820, 600)
	tt.Frame()
	if d := mktState0(vol).drawn; d != 60 {
		t.Fatalf("volume drew %d bars", d)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	found := false
	for _, s := range tt.Texts() {
		found = found || strings.HasPrefix(s, candles[99999].Time.Format("2006-01-02 15:04")+" · Vol ")
	}
	if !found || vol.Point != 99999 {
		t.Fatalf("readout %q point %d", tt.Texts(), vol.Point)
	}
	tt.Key(ui.Shift, ui.KeyLeft)
	if view.End == 0 {
		t.Fatal("Shift+Left on volume did not pan the shared view")
	}
	_ = price
	defer func() {
		if recover() == nil {
			t.Fatal("nil view accepted")
		}
	}()
	ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); VolumeChart(c, candles, nil, VolumeChartOptions{}) }, 200, 100)
}

func TestVolumeChartSnapshot(t *testing.T) {
	candles := sampleCandles(80)
	view := MarketView{Count: 60}
	snapshot(t, "volume_chart", 640, 180, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			VolumeChart(c, candles, &view, VolumeChartOptions{Frame: ChartFrame{Label: "Volume"}})
		})
	})
}
