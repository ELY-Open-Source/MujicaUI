package messaging

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestSnoozeChoicesPreserveWallHour(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	for _, now := range []time.Time{
		time.Date(2026, 3, 7, 14, 0, 0, 0, loc),
		time.Date(2026, 10, 31, 14, 0, 0, 0, loc),
	} {
		choices := SnoozeChoices(now, 9)
		if at := choices[1].At; at.Hour() != 9 || at.Day() != now.AddDate(0, 0, 1).Day() {
			t.Fatalf("tomorrow wall time: %v", at)
		}
	}
}

func TestSnoozeChoices(t *testing.T) {
	wed := time.Date(2026, 10, 7, 14, 20, 0, 0, time.UTC)
	got := SnoozeChoices(wed, 9)
	want := []SnoozeChoice{
		{"snooze_picker.later", time.Date(2026, 10, 7, 17, 0, 0, 0, time.UTC)},
		{"snooze_picker.tomorrow", time.Date(2026, 10, 8, 9, 0, 0, 0, time.UTC)},
		{"snooze_picker.weekend", time.Date(2026, 10, 10, 9, 0, 0, 0, time.UTC)},
		{"snooze_picker.nextWeek", time.Date(2026, 10, 12, 9, 0, 0, 0, time.UTC)},
	}
	if !slices.Equal(got, want) {
		t.Fatalf("choices %v", got)
	}
	late := SnoozeChoices(time.Date(2026, 10, 10, 22, 0, 0, 0, time.UTC), 9)
	if len(late) != 2 || late[0].Key != "snooze_picker.tomorrow" {
		t.Fatalf("Saturday night %v", late)
	}
}

func TestSnoozePickerChoose(t *testing.T) {
	var until time.Time
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if SnoozePicker(c, &until, SnoozePickerOptions{Location: time.UTC}).Changed() {
			n++
		}
	}, 500, 600)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	now := time.Now().UTC()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, time.UTC)
	later := now.Add(3 * time.Hour).Truncate(time.Hour)
	want := tomorrow
	if later.Day() != now.Day() {
		want = SnoozeChoices(now, 9)[1].At
	}
	if !until.Equal(want) || n != 1 {
		t.Fatalf("until %v, want %v", until, want)
	}
	if err := tt.Click("Pick date and time"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Snooze") {
		t.Fatal("custom picker not shown")
	}
	mediaPanics(t, "hour", func(c *ui.Context) { SnoozePicker(c, &until, SnoozePickerOptions{MorningHour: 30}) })
}

func TestSnoozePickerSnapshot(t *testing.T) {
	until := time.Time{}
	snapshot(t, "snooze_picker", 340, 260, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			SnoozePicker(c, &until, SnoozePickerOptions{Location: time.UTC, Now: time.Date(2026, 10, 7, 14, 20, 0, 0, time.UTC)})
		})
	})
}

func TestSnoozePickerCustomCommitAndDisabled(t *testing.T) {
	now := time.Date(2026, 10, 7, 14, 20, 0, 0, time.UTC)
	var until time.Time
	disabled, changes := false, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if SnoozePicker(c, &until, SnoozePickerOptions{Now: now, Location: time.UTC, Disabled: disabled}).Changed() {
			changes++
		}
	}, 500, 640)
	for _, label := range []string{"Pick date and time", "2026-10-07 17:00 UTC", "Time: Hours"} {
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	tt.Key(0, ui.KeyUp)
	if err := tt.Click("Confirm"); err != nil {
		t.Fatal(err)
	}
	if !until.IsZero() || changes != 0 {
		t.Fatal("custom draft committed before Snooze")
	}
	if err := tt.Click("Snooze"); err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 10, 7, 18, 0, 0, 0, time.UTC)
	if !until.Equal(want) || changes != 1 {
		t.Fatalf("until=%v changes=%d", until, changes)
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Tomorrow")
	tt.Key(0, ui.KeyEnter)
	if !until.Equal(want) || changes != 1 {
		t.Fatal("disabled snooze changed")
	}
}
