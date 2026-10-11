package chart

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestCompareSymbolPercentAndRemove(t *testing.T) {
	candles := sampleCandles(100)
	view := MarketView{Count: 50}
	eth := make([]float64, 100)
	for i := range eth {
		eth[i] = 20 + float64(i)*0.1
	}
	cmp := []CompareSymbol{{Name: "ETH", Closes: eth}}
	removed := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if CompareLegend(c, &cmp).Changed() {
			removed++
		}
		CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "BTC"}, Compare: &cmp})
	}, 800, 420)
	tt.Key(0, ui.KeyTab) // remove ETH
	tt.Key(0, ui.KeyTab) // the chart
	tt.Key(0, ui.KeyEnd)
	// Normalized from the first shown candle, index 50: (29.9/25 - 1) = 19.6%.
	want := "ETH 19.60"
	found := false
	for _, s := range tt.Texts() {
		found = found || strings.Contains(s, want)
	}
	if !found {
		t.Fatalf("readout lacks %q: %q", want, tt.Texts())
	}
	if err := tt.Click("Remove ETH"); err != nil {
		t.Fatal(err)
	}
	if len(cmp) != 0 || removed != 1 {
		t.Fatalf("remove: %v %d", cmp, removed)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("short closes accepted")
		}
	}()
	bad := []CompareSymbol{{Name: "X", Closes: []float64{1}}}
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "x"}, Compare: &bad})
	}, 300, 200)
}

func TestCompareSymbolSnapshot(t *testing.T) {
	candles := sampleCandles(100)
	view := MarketView{Count: 80}
	eth, sol := make([]float64, 100), make([]float64, 100)
	for i := range eth {
		eth[i] = 20 + float64(i)*0.08 + float64(i%7)*0.2
		sol[i] = 5 - float64(i)*0.01 + float64(i%5)*0.05
	}
	cmp := []CompareSymbol{{Name: "ETH", Closes: eth}, {Name: "SOL", Closes: sol}}
	snapshot(t, "compare_symbol", 640, 360, func(c *ui.Context) {
		ui.Column(c).Padding(12).Gap(8).Children(func() {
			CompareLegend(c, &cmp)
			CandlestickChart(c, candles, &view, CandlestickOptions{Frame: ChartFrame{Label: "BTC", Height: 240}, Type: ChartLine, Compare: &cmp})
		})
	})
}
