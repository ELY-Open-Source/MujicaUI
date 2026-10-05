package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestQuoteCardShowsQuoteAndChartSlot(t *testing.T) {
	charted := false
	tt := tradingTester(RedUp, 400, 240, func(c *ui.Context) {
		QuoteCard(c, Quote{Symbol: "600519", Name: "Kweichow Moutai", Price: 1688.5, Change: 12.3, ChangePercent: 0.73},
			QuoteCardOptions{Chart: func() { charted = true; ui.Text(c, "chart") }})
	})
	for _, want := range []string{"600519", "Kweichow Moutai", "1,688.50", "+12.30 (+0.73%)", "chart"} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	if !charted {
		t.Fatal("chart slot not built")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("empty symbol accepted")
		}
	}()
	tradingTester(GreenUp, 100, 100, func(c *ui.Context) { QuoteCard(c, Quote{}, QuoteCardOptions{}) })
}

func TestQuoteCardSnapshot(t *testing.T) {
	snapshot(t, "quote_card", 300, 190, func(c *ui.Context) {
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
			QuoteCard(c, Quote{Symbol: "MUJC", Name: "Ave Mujica Holdings", Price: 412.08, Change: -5.62, ChangePercent: -1.35},
				QuoteCardOptions{Chart: func() {
					k := Tokens(c)
					ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
						p.Line(r.X, r.Y+r.H*0.3, r.X+r.W, r.Y+r.H*0.7, 1.5, k.Info)
					})
				}})
		})
	})
}
