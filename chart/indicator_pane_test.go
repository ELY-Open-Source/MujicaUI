package chart

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestIndicatorPaneStackReadAndRemove(t *testing.T) {
	candles := sampleCandles(300)
	view := MarketView{Count: 80}
	panes := IndicatorPanes{Panes: []IndicatorPaneSpec{{Kind: IndicatorMACD}, {Kind: IndicatorRSI, Params: []float64{6}}}}
	removed := -1
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if r := IndicatorPane(c, candles, &view, &panes, IndicatorPaneOptions{Label: "BTC"}); r.Removed >= 0 {
			removed = r.Removed
		}
	}, 800, 500)
	tt.Frame()
	if !tt.HasText("MACD(12,26,9)") || !tt.HasText("RSI(6)") || len(panes.Sizes) != 2 {
		t.Fatalf("texts %q sizes %v", tt.Texts(), panes.Sizes)
	}
	tt.Key(0, ui.KeyTab) // remove MACD
	tt.Key(0, ui.KeyTab) // MACD plot
	tt.Key(0, ui.KeyEnd)
	found := false
	for _, s := range tt.Texts() {
		found = found || strings.Contains(s, "DIF ") && strings.Contains(s, "DEA ")
	}
	if !found {
		t.Fatalf("MACD readout missing: %q", tt.Texts())
	}
	if err := tt.Click("Remove MACD(12,26,9)"); err != nil {
		t.Fatal(err)
	}
	if len(panes.Panes) != 1 || panes.Panes[0].Kind != IndicatorRSI || removed != 0 {
		t.Fatalf("after remove %+v (%d)", panes.Panes, removed)
	}
	for name, p := range map[string]IndicatorPanes{
		"overlay kind": {Panes: []IndicatorPaneSpec{{Kind: IndicatorMA}}},
		"bad params":   {Panes: []IndicatorPaneSpec{{Kind: IndicatorRSI, Params: []float64{1, 2}}}},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				IndicatorPane(c, candles, &view, &p, IndicatorPaneOptions{})
			}, 300, 200)
		}()
	}
}

func TestIndicatorPaneSnapshot(t *testing.T) {
	candles := sampleCandles(120)
	view := MarketView{Count: 80}
	panes := IndicatorPanes{Panes: []IndicatorPaneSpec{{Kind: IndicatorMACD}, {Kind: IndicatorKDJ}}}
	snapshot(t, "indicator_pane", 640, 340, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() {
			IndicatorPane(c, candles, &view, &panes, IndicatorPaneOptions{Height: 316})
		})
	})
}
