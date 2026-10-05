package mujica

import (
	"github.com/egoist/mygo/ui"
	"testing"
)

func TestSymbolBadgeClickKeysAndDisabled(t *testing.T) {
	disabled := false
	count := 0
	tt := chatTester(500, 150, func(c *ui.Context) {
		if SymbolBadge(c, MarketSymbol{Code: "AAPL", Name: "Apple", Market: "NASDAQ", Type: "Stock"}, SymbolBadgeOptions{Clickable: true, Disabled: disabled}).Changed() {
			count++
		}
	})
	_ = tt.Click("AAPL")
	tt.Key(0, ui.KeyEnter)
	if count != 2 {
		t.Fatalf("clicks %d", count)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("AAPL")
	tt.Key(0, ui.KeyEnter)
	if count != 2 {
		t.Fatal("disabled activated")
	}
	marketsPanic(t, func(c *ui.Context) { SymbolBadge(c, MarketSymbol{}, SymbolBadgeOptions{}) })
}
func TestSymbolBadgeSnapshot(t *testing.T) {

	snapshot(t, "symbol_badge", 420, 100, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			ui.Row(c).Gap(8).Children(func() {
				SymbolBadge(c, MarketSymbol{Code: "AAPL", Market: "NASDAQ", Type: "Stock"}, SymbolBadgeOptions{Clickable: true})
				SymbolBadge(c, MarketSymbol{Code: "600519", Market: "SSE"}, SymbolBadgeOptions{})
			})

		})
	})
}
