package datetime

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func dayViewTester(day *core.Date, events []CalendarEvent, opts CalendarDayViewOptions, changes *[]CalendarChange) *ui.Tester {
	opts.Location, opts.Today = calTokyo, calendarToday
	opts.Now = time.Date(2026, time.October, 14, 10, 30, 0, 0, calTokyo)
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ch, ok := CalendarDayView(c, day, events, opts).Changed(); ok {
			*changes = append(*changes, ch)
		}
	}, 600, 900)
}

func TestCalendarDayViewOverlapSideBySide(t *testing.T) {
	day := oct(14)
	evs := []CalendarEvent{calEvent("a", 14, 9, 0, 120), calEvent("b", 14, 9, 30, 60), calEvent("c", 14, 10, 0, 60), calEvent("d", 14, 12, 0, 30)}
	var ch []CalendarChange
	tt := dayViewTester(&day, evs, CalendarDayViewOptions{}, &ch)
	a, _ := tt.Find("Event a, 09:00 – 11:00")
	b, _ := tt.Find("Event b, 09:30 – 10:30")
	c, _ := tt.Find("Event c, 10:00 – 11:00")
	d, _ := tt.Find("Event d, 12:00 – 12:30")
	if !(a.X < b.X && b.X < c.X) || a.W > d.W/2 || d.X != a.X {
		t.Fatalf("three overlapping events need three columns, the lone one the full width: %v %v %v %v", a, b, c, d)
	}
}

func TestCalendarDayViewNavigationSlotAndDisabled(t *testing.T) {
	day := oct(14)
	var ch []CalendarChange
	tt := dayViewTester(&day, nil, CalendarDayViewOptions{SlotMinutes: 15}, &ch)
	if !tt.HasText("Wednesday, October 14, 2026") {
		t.Fatalf("title %v", tt.Texts()[:4])
	}
	_ = tt.Click("Schedule, Wednesday, October 14, 2026")
	tt.Key(0, ui.KeyEnter)
	if len(ch) != 1 || ch[0].End.Sub(ch[0].Start) != 15*time.Minute {
		t.Fatalf("15-minute slot create %+v", ch)
	}
	_ = tt.Click("Next day")
	if day != oct(15) {
		t.Fatalf("day %v", day)
	}
	d := oct(14)
	ch = nil
	tt = dayViewTester(&d, nil, CalendarDayViewOptions{Disabled: true}, &ch)
	_ = tt.Click("Schedule, Wednesday, October 14, 2026")
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("Previous day")
	if len(ch) != 0 || d != oct(14) {
		t.Fatal("disabled day view reacted")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic without a day")
		}
	}()
	zero := core.Date{}
	dayViewTester(&zero, nil, CalendarDayViewOptions{}, &ch)
}

func TestCalendarDayViewSnapshot(t *testing.T) {
	day := oct(14)
	evs := []CalendarEvent{calEvent("a", 14, 9, 0, 120), calEvent("b", 14, 9, 30, 60), calEvent("c", 14, 10, 0, 60), calEvent("d", 14, 12, 0, 30)}
	evs[1].Color, evs[2].Color = EventInfo, EventSuccess
	snapshot(t, "calendar_day_view", 560, 860, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			CalendarDayView(c, &day, evs, CalendarDayViewOptions{Location: calTokyo, Today: calendarToday,
				Now: time.Date(2026, time.October, 14, 10, 30, 0, 0, calTokyo)})
		})
	})
}
