package mujica

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// CalendarDayViewOptions configure a CalendarDayView.
type CalendarDayViewOptions struct {
	// Location is the zone the day is cut in; required.
	Location *time.Location
	// Version must change when events change in place.
	Version int
	// SlotMinutes is the grid step; zero means 30.
	SlotMinutes int
	// Today is marked; zero uses the system clock.
	Today Date
	// Now places the current-time line; zero uses the system clock.
	Now      time.Time
	Hour12   bool
	Disabled bool
	Label    string
}

// CalendarDayView shows one day's time grid.
func CalendarDayView(c *ui.Context, day *Date, events []CalendarEvent, opts CalendarDayViewOptions) CalendarViewResult {
	datetimeCheckDate("CalendarDayView day", *day)
	if day.IsZero() {
		panic("mujica: CalendarDayView needs a day")
	}
	cfg := timeGridConfig{name: "calendar_day_view", loc: eventLocation("CalendarDayView", opts.Location), version: opts.Version,
		slot: timeGridSlot("CalendarDayView", opts.SlotMinutes), today: datetimeToday(c, opts.Today), now: opts.Now, hour12: opts.Hour12}
	eventsCheck("CalendarDayView", events)
	title := datetimeDayLabel(c, *day)
	label := opts.Label
	if label == "" {
		label = title
	}
	r := CalendarViewResult{}
	root := ui.Column(c).Gap(Density(c).Unit()).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	r.Element = root
	root.Children(func() {
		calendarViewHeader(c, title, msg(c, "calendar_day_view.prev"), msg(c, "calendar_day_view.next"), func(dir int) {
			*day = day.AddDays(dir)
		})
		timeGrid(c, root, []Date{*day}, events, cfg, &r)
	})
	return r
}
