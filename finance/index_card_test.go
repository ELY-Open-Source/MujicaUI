package finance

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestIndexCardQuoteAndKeyboard(t *testing.T) {
	count := 0
	disabled := false
	q := IndexQuote{Name: "Index", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{100, 110, 105}}
	tt := chatTester(500, 300, func(c *ui.Context) {
		if IndexCard(c, q, IndexCardOptions{Clickable: true, Disabled: disabled}).Changed() {
			count++
		}
	})
	if !tt.HasText("5,751.07") || !tt.HasText("+51.13 (+0.90%)") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if count != 1 {
		t.Fatalf("count %d", count)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Index")
	if count != 1 {
		t.Fatal("disabled click")
	}
	marketsPanic(t, func(c *ui.Context) { IndexCard(c, IndexQuote{}, IndexCardOptions{}) })
}
func TestIndexCardSnapshot(t *testing.T) {
	snapshot(t, "index_card", 260, 220, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			IndexCard(c, IndexQuote{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: mktSeries(30, 0)}, IndexCardOptions{})
		})
	})
}
