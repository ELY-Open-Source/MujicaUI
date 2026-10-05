package mujica

import (
	"github.com/egoist/mygo/ui"
	"testing"
	"time"
)

func TestTradingSessionClockZonesBoundariesAndHoliday(t *testing.T) {
	ny := SessionMarket{Name: "NYSE", Location: mktNewYork(), Sessions: []TradingSession{{MarketPreMarket, 240, 570}, {MarketOpen, 570, 960}, {MarketAfterHours, 960, 1200}}}
	for _, tc := range []struct {
		minute int
		phase  MarketPhase
	}{{239, MarketClosed}, {240, MarketPreMarket}, {569, MarketPreMarket}, {570, MarketOpen}, {960, MarketAfterHours}, {1200, MarketClosed}} {
		local := time.Date(2026, 10, 5, tc.minute/60, tc.minute%60, 0, 0, ny.Location)
		phase, next, _ := tradingSessionPhase(ny, local)
		if phase != tc.phase || !next.After(local) {
			t.Fatalf("at %d phase %d next %v", tc.minute, phase, next)
		}
	}
	tt := chatTester(960, 260, func(c *ui.Context) {
		TradingSessionClock(c, []SessionMarket{ny, {Name: "SSE", Location: time.FixedZone("CST", 8*3600), Holiday: true}}, TradingSessionClockOptions{Now: mktNow})
	})
	if !tt.HasText("10:30 EDT") || !tt.HasText("22:30 CST") || !tt.HasText("Holiday") {
		t.Fatalf("zones %q", tt.Texts())
	}
	marketsPanic(t, func(c *ui.Context) {
		TradingSessionClock(c, []SessionMarket{{Name: "Bad"}}, TradingSessionClockOptions{})
	})
}
func TestTradingSessionClockSnapshot(t *testing.T) {

	snapshot(t, "trading_session_clock", 720, 220, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Fill().Children(func() {
			TradingSessionClock(c, []SessionMarket{
				{Name: "NYSE", Location: mktNewYork(), Sessions: []TradingSession{{MarketPreMarket, 240, 570}, {MarketOpen, 570, 960}, {MarketAfterHours, 960, 1200}}},
				{Name: "SSE", Location: time.FixedZone("CST", 8*3600), Sessions: []TradingSession{{MarketOpen, 570, 690}, {MarketOpen, 780, 900}}},
			}, TradingSessionClockOptions{Now: mktNow})

		})
	})
}

func TestTradingSessionClockDST(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	m := SessionMarket{Name: "NYSE", Location: loc, Sessions: []TradingSession{{MarketOpen, 570, 960}}}
	local := time.Date(2026, 3, 8, 8, 0, 0, 0, loc)
	_, next, _ := tradingSessionPhase(m, local)
	if next.Hour() != 9 || next.Minute() != 30 {
		t.Fatalf("DST open %v", next)
	}
}
