package chart

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestIndicatorOverlayValuesInReadoutAndMemo(t *testing.T) {
	candles := sampleCandles(200)
	view := MarketView{Count: 60}
	overlays := []IndicatorOverlay{{Kind: IndicatorMA, Params: []float64{5}}, {Kind: IndicatorBOLL}}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "BTC"}, Overlays: overlays})
	}, 800, 400)
	if !tt.HasText("MA(5)") || !tt.HasText("BOLL(20,2) ↑") {
		t.Fatalf("legend %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	ma := SMA(closesOf(candles), 5)
	want := "MA(5) " + core.PriceUngrouped(ma[199], 2)
	found := false
	for _, s := range tt.Texts() {
		found = found || strings.Contains(s, want)
	}
	if !found {
		t.Fatalf("readout lacks %q: %q", want, tt.Texts())
	}
	// Caller values replace the computation.
	own := make([]float64, 200)
	for i := range own {
		own[i] = 42
	}
	overlays = []IndicatorOverlay{{Kind: IndicatorEMA, Values: own}}
	tt.Frame()
	tt.Frame()
	if !tt.HasText("EMA(20) 42.00") {
		t.Fatalf("own values %q", tt.Texts())
	}
	for name, o := range map[string]IndicatorOverlay{
		"pane kind":   {Kind: IndicatorMACD},
		"params":      {Kind: IndicatorMA, Params: []float64{0}},
		"values size": {Kind: IndicatorMA, Values: []float64{1}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "x"}, Overlays: []IndicatorOverlay{o}})
			}, 300, 200)
		}()
	}
}

func closesOf(cs []Candle) []float64 {
	out := make([]float64, len(cs))
	for i, c := range cs {
		out[i] = c.Close
	}
	return out
}

func TestIndicatorParamsField(t *testing.T) {
	var params []float64
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		IndicatorParamsField(c, IndicatorMACD, &params)
	}, 500, 100)
	if len(params) != 3 || params[0] != 12 {
		t.Fatalf("defaults %v", params)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	if params[0] != 13 {
		t.Fatalf("Up on fast %v", params)
	}
}

func TestIndicatorOverlaySnapshot(t *testing.T) {
	candles := sampleCandles(120)
	view := MarketView{Count: 90}
	snapshot(t, "indicator_overlay", 640, 340, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "BTC", Height: 250},
				Overlays: []IndicatorOverlay{{Kind: IndicatorMA, Params: []float64{10}}, {Kind: IndicatorBOLL}}})
		})
	})
}
