package mujica

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

func TestMarketOverviewBreadthStampAndChoice(t *testing.T) {
	var got MarketOverviewView
	disabled := false
	tt := chatTester(760, 360, func(c *ui.Context) {
		v := MarketOverview(c, []IndexQuote{{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9}}, MarketOverviewOptions{Breadth: MarketBreadth{2100, 1400, 300}, AsOf: mktNow, Disabled: disabled})
		if v.Chosen >= 0 {
			got = v
		}
	})
	if !tt.HasText("2100 up") || !tt.HasText("1400 down") || !tt.HasText("As of 2026-10-05 14:30") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("S&P 500")
	if got.Chosen != 0 {
		t.Fatalf("chosen %d", got.Chosen)
	}
	got.Chosen = -1
	disabled = true
	tt.Frame()
	_ = tt.Click("S&P 500")
	if got.Chosen != -1 {
		t.Fatal("disabled chose")
	}
	marketsPanic(t, func(c *ui.Context) {
		MarketOverview(c, nil, MarketOverviewOptions{Breadth: MarketBreadth{Advancing: -1}})
	})
}
func TestMarketOverviewSnapshot(t *testing.T) {

	snapshot(t, "market_overview", 760, 360, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			MarketOverview(c, []IndexQuote{{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: mktSeries(30, 0)},
				{Name: "Nasdaq", Price: 18137.85, Change: -219.22, ChangePercent: -1.2, Series: mktSeries(30, 2)}},
				MarketOverviewOptions{Breadth: MarketBreadth{2100, 1400, 300}, AsOf: mktNow})

		})
	})
}
