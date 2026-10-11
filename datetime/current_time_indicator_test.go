package datetime

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

func TestCurrentTimeIndicatorPositionAndZone(t *testing.T) {
	now := time.Date(2026, time.October, 14, 3, 0, 0, 0, time.UTC) // 12:00 in Tokyo
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		ui.Box(c).Size(200, 240).Children(func() {
			CurrentTimeIndicator(c, CurrentTimeIndicatorOptions{Location: calTokyo, Day: oct(14), Now: now})
		})
	}, 200, 240)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if !tt.HasText("Now, 12:00") {
		t.Fatalf("label %v", tt.Texts())
	}
	img := tt.Image()
	if !datetimeNear(img.RGBAAt(100, 120), theme.Light().Danger) || datetimeNear(img.RGBAAt(100, 60), theme.Light().Danger) {
		t.Fatal("the line should sit halfway down at noon")
	}
}

func TestCurrentTimeIndicatorOtherDayHidden(t *testing.T) {
	now := time.Date(2026, time.October, 14, 12, 0, 0, 0, calTokyo)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Box(c).Size(200, 240).Children(func() {
			CurrentTimeIndicator(c, CurrentTimeIndicatorOptions{Location: calTokyo, Day: oct(15), Now: now})
		})
	}, 200, 240)
	if tt.HasText("Now,") {
		t.Fatal("no line on another day")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic without Location")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CurrentTimeIndicator(c, CurrentTimeIndicatorOptions{})
	}, 10, 10)
}

func TestCurrentTimeIndicatorUsesClockWhileShown(t *testing.T) {
	shown := true
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Box(c).Size(200, 240).Children(func() {
			if shown {
				CurrentTimeIndicator(c, CurrentTimeIndicatorOptions{Location: time.Local})
			}
		})
	}, 200, 240)
	if !tt.HasText("Now, ") {
		t.Fatalf("today's line missing: %v", tt.Texts())
	}
	shown = false
	tt.Frame()
	if tt.HasText("Now, ") {
		t.Fatal("hidden indicator still shows")
	}
}

func TestCurrentTimeIndicatorSnapshot(t *testing.T) {
	now := time.Date(2026, time.October, 14, 15, 30, 0, 0, calTokyo)
	snapshot(t, "current_time_indicator", 240, 160, func(c *ui.Context) {
		ui.Box(c).Margin(16).Size(208, 128).Border(1, core.Tokens(c).Border).Children(func() {
			CurrentTimeIndicator(c, CurrentTimeIndicatorOptions{Location: calTokyo, Now: now})
		})
	})
}
