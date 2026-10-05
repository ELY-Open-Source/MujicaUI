package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestBidAskBarShares(t *testing.T) {
	bid, ask := 620.0, 380.0
	tt := tradingTester(GreenUp, 400, 100, func(c *ui.Context) {
		ui.Box(c).Width(300).Children(func() { BidAskBar(c, bid, ask, BidAskBarOptions{}) })
	})
	if !tt.HasText("Buy 62.0% · 620") || !tt.HasText("380 · Sell 38.0%") {
		t.Fatalf("texts %q", tt.Texts())
	}
	bid, ask = 0, 0
	tt.Frame()
	if !tt.HasText("Buy 50.0% · 0") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("negative volume accepted")
		}
	}()
	tradingTester(GreenUp, 100, 100, func(c *ui.Context) { BidAskBar(c, -1, 2, BidAskBarOptions{}) })
}

func TestBidAskBarSnapshot(t *testing.T) {
	snapshot(t, "bid_ask_bar", 340, 80, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() { BidAskBar(c, 6234, 3820, BidAskBarOptions{}) })
	})
}
