package mujica

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func agendaTester(events []CalendarEvent, opts AgendaViewOptions, opened *[]string) *ui.Tester {
	opts.Location, opts.Today = calTokyo, calendarToday
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Height(500).Children(func() {
			if id, ok := AgendaView(c, events, opts).Opened(); ok {
				*opened = append(*opened, id)
			}
		})
	}, 600, 520)
}

func TestAgendaViewGroupsByDay(t *testing.T) {
	evs := []CalendarEvent{calEvent("b", 7, 14, 0, 60), calEvent("a", 7, 9, 0, 30), calEvent("c", 9, 10, 0, 60), calEvent("old", 1, 9, 0, 60)}
	evs[0].Location = "Rehearsal hall"
	var opened []string
	tt := agendaTester(evs, AgendaViewOptions{}, &opened)
	texts := tt.Texts()
	idx := func(s string) int {
		for i, x := range texts {
			if x == s {
				return i
			}
		}
		return -1
	}
	d7, a, b, d9, c := idx("Wednesday, October 7, 2026"), idx("Event a, 09:00 – 09:30"), idx("Event b, 14:00 – 15:00"), idx("Friday, October 9, 2026"), idx("Event c, 10:00 – 11:00")
	if d7 < 0 || !(d7 < a && a < b && b < d9 && d9 < c) || tt.HasText("Event old") || !tt.HasText("Rehearsal hall") {
		t.Fatalf("order %d %d %d %d %d in %v", d7, a, b, d9, c, texts)
	}
	_ = tt.Click("Event b, 14:00 – 15:00")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(opened) != 2 || opened[0] != "b" || opened[1] != "c" {
		t.Fatalf("opened %v", opened)
	}
}

func TestAgendaViewMultiDayEmptyAndDisabled(t *testing.T) {
	trip := calEvent("t", 6, 20, 0, 26*60)
	var opened []string
	tt := agendaTester([]CalendarEvent{trip}, AgendaViewOptions{Days: 3}, &opened)
	if !tt.HasText("Tuesday, October 6, 2026") || !tt.HasText("Wednesday, October 7, 2026") {
		t.Fatalf("a multi-day event lists on each day: %v", tt.Texts())
	}
	tt = agendaTester(nil, AgendaViewOptions{Days: 7}, &opened)
	if !tt.HasText("No events in the next 7 days.") {
		t.Fatal("no empty state")
	}
	tt = agendaTester([]CalendarEvent{calEvent("a", 6, 9, 0, 30)}, AgendaViewOptions{Disabled: true}, &opened)
	_ = tt.Click("Event a, 09:00 – 09:30")
	if len(opened) != 0 {
		t.Fatal("disabled agenda opened")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	agendaTester(nil, AgendaViewOptions{Days: -1}, &opened)
}

func TestAgendaViewLargeData(t *testing.T) {
	var evs []CalendarEvent
	base := time.Date(2026, time.October, 5, 0, 0, 0, 0, calTokyo)
	for i := range 100000 {
		s := base.Add(time.Duration(i%(365*24)) * time.Hour)
		evs = append(evs, CalendarEvent{ID: fmt.Sprint(i), Title: "E", Start: s, End: s.Add(30 * time.Minute)})
	}
	var opened []string
	tt := agendaTester(evs, AgendaViewOptions{Days: 365}, &opened)
	if n := len(tt.Texts()); n > 200 {
		t.Fatalf("built %d labels", n)
	}
}

func TestAgendaViewSnapshot(t *testing.T) {
	evs := []CalendarEvent{calEvent("a", 5, 9, 0, 30), calEvent("b", 5, 14, 0, 60), calEvent("c", 7, 10, 0, 60), calEvent("d", 8, 18, 0, 120)}
	evs[1].Location, evs[1].Color, evs[3].Color = "Rehearsal hall", EventInfo, EventSuccess
	snapshot(t, "agenda_view", 520, 420, func(c *ui.Context) {
		ui.Column(c).Padding(16).Height(420).Children(func() {
			AgendaView(c, evs, AgendaViewOptions{Location: calTokyo, Today: calendarToday})
		})
	})
}
