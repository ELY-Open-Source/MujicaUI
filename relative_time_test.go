package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestRelativeTimeWording(t *testing.T) {
	var en, zh []string
	cases := []time.Duration{10 * time.Second, -90 * time.Second, 3 * time.Minute, 5 * time.Hour, -25 * time.Hour, 40 * 24 * time.Hour, 800 * 24 * time.Hour}
	for _, loc := range []Locale{EnUS, ZhCN} {
		ui.NewTester(func(c *ui.Context) {
			Use(c, Settings{Locale: loc})
			var out []string
			for _, d := range cases {
				out = append(out, RelativeTimeText(c, d))
			}
			if loc == EnUS {
				en = out
			} else {
				zh = out
			}
		}, 10, 10)
	}
	wantEn := []string{"just now", "in 1 minute", "3 minutes ago", "5 hours ago", "in 1 day", "1 month ago", "2 years ago"}
	wantZh := []string{"刚刚", "1 分钟后", "3 分钟前", "5 小时前", "1 天后", "1 个月前", "2 年前"}
	for i := range cases {
		if en[i] != wantEn[i] || zh[i] != wantZh[i] {
			t.Fatalf("%v: %q %q", cases[i], en[i], zh[i])
		}
	}
}

func TestRelativeTimeRefreshFollowsDistance(t *testing.T) {
	for _, c := range []struct{ d, next time.Duration }{
		{10*time.Second + 300*time.Millisecond, 700 * time.Millisecond},
		{5*time.Minute + 20*time.Second, 40 * time.Second},
		{3*time.Hour + 15*time.Minute, 45 * time.Minute},
		{-(90 * time.Second), 1 * time.Second},
	} {
		if got := relativeTimeNext(c.d); got != c.next {
			t.Errorf("%v: next %v, want %v", c.d, got, c.next)
		}
	}
}

func TestRelativeTimeShowsExactTimeAndStopsWhenHidden(t *testing.T) {
	at := time.Date(2026, time.October, 5, 9, 0, 0, 0, time.UTC)
	shown := true
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if shown {
			RelativeTime(c, at, RelativeTimeOptions{Location: calTokyo, Now: at.Add(3 * time.Hour)})
		}
	}, 300, 100)
	if !tt.HasText("3 hours ago, 2026-10-05 18:00 JST") {
		t.Fatalf("texts %v", tt.Texts())
	}
	shown = false
	tt.Frame()
	if tt.HasText("ago") {
		t.Fatal("hidden text remains")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic without Location")
		}
	}()
	ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); RelativeTime(c, at, RelativeTimeOptions{}) }, 10, 10)
}

func TestRelativeTimeSnapshot(t *testing.T) {
	at := time.Date(2026, time.October, 5, 9, 0, 0, 0, calTokyo)
	snapshot(t, "relative_time", 260, 120, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			RelativeTime(c, at, RelativeTimeOptions{Location: calTokyo, Now: at.Add(20 * time.Second)})
			RelativeTime(c, at, RelativeTimeOptions{Location: calTokyo, Now: at.Add(-2 * time.Hour)})
			RelativeTime(c, at, RelativeTimeOptions{Location: calTokyo, Now: at.Add(72 * time.Hour)})
		})
	})
}
