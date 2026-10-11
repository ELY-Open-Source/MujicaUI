package datetime

import (
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

var calTokyo = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	return loc
}()

// calEvent is an event on calendarToday's month in Tokyo.
func calEvent(id string, day, h, m, mins int) CalendarEvent {
	s := time.Date(2026, time.October, day, h, m, 0, 0, calTokyo)
	return CalendarEvent{ID: id, Title: "Event " + id, Start: s, End: s.Add(time.Duration(mins) * time.Minute)}
}

func TestEventChipShowsZoneTimeAndClicks(t *testing.T) {
	ev := calEvent("a", 6, 9, 30, 90)
	ev.Start, ev.End = ev.Start.UTC(), ev.End.UTC()
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if EventChip(c, ev, EventChipOptions{Location: calTokyo}).Clicked() {
			n++
		}
	}, 300, 200)
	if !tt.HasText("09:30 – 11:00") || !tt.HasText("Event a") {
		t.Fatalf("chip should show Tokyo wall time: %v", tt.Texts())
	}
	label := "Event a, 09:30 – 11:00"
	_ = tt.Click(label)
	tt.Key(0, ui.KeyEnter)
	if n != 2 {
		t.Fatalf("clicks %d", n)
	}
}

func TestEventChipTruncatesAndAllDay(t *testing.T) {
	ev := calEvent("b", 6, 0, 0, 24*60)
	ev.AllDay, ev.Title = true, strings.Repeat("Long title ", 20)
	var w float32
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(120).Children(func() {
			w = EventChip(c, ev, EventChipOptions{Location: calTokyo, Inline: true}).Bounds().W
		})
	}, 300, 200)
	tt.Frame()
	if w > 120 || !tt.HasText("All day") {
		t.Fatalf("width %v texts %v", w, tt.Texts())
	}
}

func TestEventChipDisabledAndInvalid(t *testing.T) {
	ev := calEvent("c", 6, 9, 0, 30)
	n := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if EventChip(c, ev, EventChipOptions{Location: calTokyo, Disabled: true}).Clicked() {
			n++
		}
	}, 300, 200)
	_ = tt.Click("Event c, 09:00 – 09:30")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if n != 0 {
		t.Fatal("disabled chip clicked")
	}
	for _, f := range []func(){
		func() { EventChip(nil, ev, EventChipOptions{}) },
		func() { ev.Color = 9; EventTone(theme.Light(), ev.Color) },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("expected panic")
				}
			}()
			f()
		}()
	}
}

func TestEventPaletteContrast(t *testing.T) {
	for _, k := range []theme.Tokens{theme.Light(), theme.Dark()} {
		for i, col := range eventPalette(k) {
			if r := theme.Contrast(col, k.Surface); r < 3 {
				t.Errorf("dark=%v color %d contrast %.2f", k.Dark, i, r)
			}
		}
	}
}

func TestEventColumns(t *testing.T) {
	col, cols := eventColumns([]eventSpan{{0, 2}, {1, 3}, {1.5, 2.5}, {2, 4}, {5, 6}})
	want := []int{0, 1, 2, 0, 0}
	wantN := []int{3, 3, 3, 3, 1}
	for i := range want {
		if col[i] != want[i] || cols[i] != wantN[i] {
			t.Fatalf("columns %v of %v", col, cols)
		}
	}
}

func TestEventChipSnapshot(t *testing.T) {
	evs := []CalendarEvent{calEvent("1", 6, 9, 0, 60), calEvent("2", 6, 13, 0, 30), calEvent("3", 6, 15, 0, 45)}
	evs[1].Color, evs[2].Color = EventSuccess, EventWarning
	snapshot(t, "event_chip", 420, 150, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(8).Children(func() {
			EventChip(c, evs[0], EventChipOptions{Location: calTokyo})
			EventChip(c, evs[1], EventChipOptions{Location: calTokyo, Selected: true})
			EventChip(c, evs[2], EventChipOptions{Location: calTokyo, Inline: true})
		})
	})
}
