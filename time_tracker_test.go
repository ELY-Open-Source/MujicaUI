package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestTimeTrackerStartStopCorrect(t *testing.T) {
	var s TimeTrackerState
	now := time.Date(2026, time.October, 5, 9, 0, 0, 0, time.UTC)
	n, disabled := 0, false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if TimeTracker(c, &s, TimeTrackerOptions{Now: now, Disabled: disabled}).Changed() {
			n++
		}
	}, 900, 160)
	_ = tt.Click("Start")
	now = now.Add(25*time.Minute + 30*time.Second)
	tt.Frame()
	if !s.Running() || !tt.HasText("Time tracked, 25 min 30 s") || !tt.HasText("Tracking…") {
		t.Fatalf("running %v", tt.Texts())
	}
	_ = tt.Click("Correct the total")
	if tt.HasText("Apply") {
		t.Fatal("correction must wait until tracking stops")
	}
	_ = tt.Click("Stop")
	now = now.Add(time.Hour)
	_ = tt.Click("Start")
	now = now.Add(4*time.Minute + 30*time.Second)
	_ = tt.Click("Stop")
	if s.Total(now) != 30*time.Minute || len(s.Entries) != 2 || n != 4 {
		t.Fatalf("total %v entries %d changes %d", s.Total(now), len(s.Entries), n)
	}
	_ = tt.Click("Correct the total")
	_ = tt.Click("Minutes")
	tt.Command("selectAll")
	tt.Type("45")
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("Apply")
	if s.Total(now) != 45*time.Minute || s.Adjust != 15*time.Minute || !tt.HasText("Time tracked, 45 min") {
		t.Fatalf("corrected to %v (adjust %v)", s.Total(now), s.Adjust)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Start")
	if s.Running() {
		t.Fatal("disabled tracker started")
	}
}

func TestTimeTrackerInvalid(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for two open entries")
		}
	}()
	now := time.Now()
	s := TimeTrackerState{Entries: []TimeEntry{{Start: now}, {Start: now}}}
	ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); TimeTracker(c, &s, TimeTrackerOptions{Now: now}) }, 10, 10)
}

func TestTimeTrackerSnapshot(t *testing.T) {
	now := time.Date(2026, time.October, 5, 9, 0, 0, 0, time.UTC)
	s := TimeTrackerState{Entries: []TimeEntry{{Start: now.Add(-time.Hour), End: now.Add(-30 * time.Minute)}, {Start: now.Add(-12 * time.Minute)}}}
	snapshot(t, "time_tracker", 460, 80, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { TimeTracker(c, &s, TimeTrackerOptions{Now: now}) })
	})
}

func TestTimeTrackerRejectsInvalidSpans(t *testing.T) {
	now := time.Date(2026, time.October, 5, 9, 0, 0, 0, time.UTC)
	for _, state := range []TimeTrackerState{
		{Entries: []TimeEntry{{Start: now, End: now.Add(time.Hour)}, {Start: now.Add(time.Minute), End: now.Add(2 * time.Hour)}}},
		{Entries: []TimeEntry{{Start: now.Add(time.Hour)}}},
		{Adjust: -time.Second},
	} {
		t.Run("span", func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected invalid span panic")
				}
			}()
			ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); TimeTracker(c, &state, TimeTrackerOptions{Now: now}) }, 600, 200)
		})
	}
}

func TestTimeTrackerCorrectionRemovesFractionalRemainder(t *testing.T) {
	now := time.Date(2026, time.October, 5, 9, 0, 0, 0, time.UTC)
	s := TimeTrackerState{Entries: []TimeEntry{{Start: now.Add(-30500 * time.Millisecond), End: now}}}
	tt := ui.NewTester(func(c *ui.Context) { Use(c, Settings{}); TimeTracker(c, &s, TimeTrackerOptions{Now: now}) }, 900, 180)
	if err := tt.Click("Correct the total"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Apply"); err != nil {
		t.Fatal(err)
	}
	if s.Total(now) != 30*time.Second {
		t.Fatalf("total = %v", s.Total(now))
	}
}
