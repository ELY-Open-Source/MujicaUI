package datetime

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestCountdownFormatsAndFinishesOnce(t *testing.T) {
	deadline := time.Date(2026, time.October, 5, 12, 0, 0, 0, calTokyo)
	now, n := deadline.Add(-(26*time.Hour + 3*time.Minute + 4*time.Second + 500*time.Millisecond)), 0
	format := CountdownClock
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if Countdown(c, deadline, CountdownOptions{Now: now, Format: format}).Finished() {
			n++
		}
	}, 400, 120)
	if !tt.HasText("Time left, 1 d 2 h 3 min 5 s") {
		t.Fatalf("label %v", tt.Texts())
	}
	var clock, words string
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		clock, words = CountdownText(c, 26*time.Hour+3*time.Minute+5*time.Second, CountdownClock), CountdownText(c, 0, CountdownWords)
	}, 10, 10)
	if clock != "1d 02:03:05" || words != "0 s" {
		t.Fatalf("formats %q %q", clock, words)
	}
	format = CountdownWords
	now = deadline.Add(time.Second)
	tt.Frame()
	tt.Frame()
	if n != 1 || !tt.HasText("Time left, 0 s") {
		t.Fatalf("finished %d times, %v", n, tt.Texts())
	}
	if a := tt.Announcements(); len(a) != 1 || a[0] != "Time left: time is up." {
		t.Fatalf("announcements %q", a)
	}
	deadline = deadline.Add(time.Minute)
	now = deadline
	tt.Frame()
	if n != 2 {
		t.Fatal("a new deadline should finish again")
	}
}

func TestCountdownInvalidFormat(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Countdown(c, time.Now(), CountdownOptions{Format: 7})
	}, 10, 10)
}

func TestCountdownSnapshot(t *testing.T) {
	deadline := time.Date(2026, time.October, 5, 12, 0, 0, 0, calTokyo)
	snapshot(t, "countdown", 320, 140, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			Countdown(c, deadline, CountdownOptions{Now: deadline.Add(-(3*time.Hour + 25*time.Second))})
			Countdown(c, deadline, CountdownOptions{Now: deadline.Add(-90 * time.Second), Format: CountdownWords})
		})
	})
}
