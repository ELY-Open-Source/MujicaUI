package finance

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestMarketStatusPhasesAndNext(t *testing.T) {
	for phase, word := range []string{"Closed", "Pre-market", "Open", "After hours", "Holiday"} {
		tt := chatTester(760, 150, func(c *ui.Context) {
			MarketStatus(c, MarketPhase(phase), MarketStatusOptions{Market: "NYSE", Now: mktNow, Next: mktNow.Add(90 * time.Minute), NextPhase: MarketAfterHours})
		})
		if !tt.HasText(word) || !tt.HasText("After hours at 16:00 (in 1h 30m)") {
			t.Fatalf("phase %d texts %q", phase, tt.Texts())
		}
	}
	marketsPanic(t, func(c *ui.Context) { MarketStatus(c, MarketPhase(99), MarketStatusOptions{}) })
}
func TestMarketStatusSnapshot(t *testing.T) {

	snapshot(t, "market_status", 520, 160, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			MarketStatus(c, MarketOpen, MarketStatusOptions{Market: "NYSE", Next: mktNow.Add(5*3600e9 + 30*60e9), NextPhase: MarketAfterHours, Now: mktNow})
			MarketStatus(c, MarketHoliday, MarketStatusOptions{})

		})
	})
}
