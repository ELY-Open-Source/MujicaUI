package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestPriceChangeBadgeSignsAndWords(t *testing.T) {
	tt := tradingTester(GreenUp, 400, 160, func(c *ui.Context) {
		PriceChangeBadge(c, 1.25, 0.85, PriceChangeBadgeOptions{})
		PriceChangeBadge(c, -3, -1.5, PriceChangeBadgeOptions{Decimals: 1})
		PriceChangeBadge(c, 0, 0, PriceChangeBadgeOptions{HidePercent: true})
	})
	for _, want := range []string{"+1.25 (+0.85%)", "−3.0 (−1.50%)", "0.00"} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	for _, name := range []string{"up +1.25 (+0.85%)", "down −3.0 (−1.50%)", "unchanged 0.00"} {
		if _, ok := tt.Find(name); !ok {
			t.Fatalf("not spoken: %q", name)
		}
	}
}

func TestPriceChangeBadgeSnapshot(t *testing.T) {
	snapshot(t, "price_change_badge", 360, 80, func(c *ui.Context) {
		ui.Row(c).Padding(20).Gap(10).Children(func() {
			PriceChangeBadge(c, 1.25, 0.85, PriceChangeBadgeOptions{})
			PriceChangeBadge(c, -3.4, -1.52, PriceChangeBadgeOptions{})
			PriceChangeBadge(c, 0, 0, PriceChangeBadgeOptions{})
		})
	})
}
