package chart

import (
	"bytes"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestOHLCChartSharesDataAndDiffersInLook(t *testing.T) {
	candles := sampleCandles(60)
	view := MarketView{Count: 40}
	ohlc := true
	var res ChartResult
	var st *mktState
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ohlc {
			res = OHLCChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "Bars"}})
		} else {
			res = CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "Bars"}})
		}
		st = mktState0(res)
	}, 600, 320)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	if res.Point != 59 || st.drawn != 40 {
		t.Fatalf("point %d drawn %d", res.Point, st.drawn)
	}
	a := encode(t, tt.Image())
	ohlc = false
	tt.Frame()
	if bytes.Equal(a, encode(t, tt.Image())) {
		t.Fatal("OHLC bars look the same as candles")
	}
}

func TestOHLCChartSnapshot(t *testing.T) {
	candles := sampleCandles(60)
	view := MarketView{Count: 50}
	snapshot(t, "ohlc_chart", 640, 300, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			OHLCChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "Bars", Height: 240}})
		})
	})
}
