package chart

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestCandleCountdownIntervalsAndClock(t *testing.T) {
	now := time.Date(2026, 10, 5, 12, 34, 5, 0, time.UTC)
	iv := Interval("1m")
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CandleCountdown(c, iv, CandleCountdownOptions{Now: now})
	}, 300, 100)
	if !tt.HasText("00:55") {
		t.Fatalf("countdown %q", tt.Texts())
	}
	iv = "1h"
	tt.Frame()
	if !tt.HasText("25:55") {
		t.Fatalf("hour %q", tt.Texts())
	}
	for _, period := range Intervals {
		if !period.Close(now, time.UTC).After(now) {
			t.Fatalf("close %s", period)
		}
	}
	for _, tc := range []struct {
		d    time.Duration
		want string
	}{{time.Second, "00:01"}, {time.Hour + time.Second, "1:00:01"}, {48*time.Hour + time.Hour, "2d 01:00:00"}} {
		if got := candleCountdownText(tc.d); got != tc.want {
			t.Fatalf("text %q", got)
		}
	}
}
func TestCandleCountdownSnapshot(t *testing.T) {
	snapshot(t, "candle_countdown", 300, 80, func(c *ui.Context) {
		ui.Row(c).Padding(12).Children(func() {
			CandleCountdown(c, "1h", CandleCountdownOptions{Now: time.Date(2026, 10, 5, 12, 34, 5, 0, time.UTC)})
		})
	})
}
