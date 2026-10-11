package datetime

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestStopwatchStartLapPauseReset(t *testing.T) {
	var s StopwatchState
	now := time.Date(2026, time.October, 5, 9, 0, 0, 0, time.UTC)
	n, disabled := 0, false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if Stopwatch(c, &s, StopwatchOptions{Now: now, Disabled: disabled}).Changed() {
			n++
		}
	}, 500, 500)
	_ = tt.Click("Start")
	now = now.Add(12*time.Second + 340*time.Millisecond)
	tt.Frame()
	if !tt.HasText("00:12.3") || !s.Running() {
		t.Fatalf("running %v", tt.Texts())
	}
	_ = tt.Click("Lap")
	now = now.Add(5 * time.Second)
	_ = tt.Click("Lap")
	if len(s.Laps) != 2 || !tt.HasText("Lap 2") || !tt.HasText("00:05.0") {
		t.Fatalf("laps %v %v", s.Laps, tt.Texts())
	}
	_ = tt.Click("Pause")
	now = now.Add(time.Minute)
	tt.Frame()
	if s.Running() || s.Elapsed(now) != 17340*time.Millisecond || !tt.HasText("Resume") {
		t.Fatalf("paused at %v", s.Elapsed(now))
	}
	tt.Key(0, ui.KeyEnter)
	now = now.Add(time.Hour)
	tt.Frame()
	if !tt.HasText("1:00:17.3") {
		t.Fatalf("resumed by keyboard: %v", tt.Texts())
	}
	_ = tt.Click("Pause")
	_ = tt.Click("Reset")
	if s.Elapsed(now) != 0 || len(s.Laps) != 0 || n != 7 {
		t.Fatalf("reset %v laps %v changes %d", s.Elapsed(now), s.Laps, n)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Start")
	if s.Running() {
		t.Fatal("disabled stopwatch started")
	}
}

func TestStopwatchSnapshot(t *testing.T) {
	now := time.Date(2026, time.October, 5, 9, 0, 0, 0, time.UTC)
	s := StopwatchState{base: 83*time.Second + 400*time.Millisecond, Laps: []time.Duration{31 * time.Second, 62*time.Second + 500*time.Millisecond}}
	snapshot(t, "stopwatch", 420, 260, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { Stopwatch(c, &s, StopwatchOptions{Now: now}) })
	})
}
