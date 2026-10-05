package mujica

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// CalendarWeekViewOptions configure a CalendarWeekView.
type CalendarWeekViewOptions struct {
	// Location is the zone days are cut in; required.
	Location *time.Location
	// Version must change when events change in place.
	Version   int
	WeekStart WeekStart
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

// CalendarWeekView shows *week as a time grid; drags edit events.
func CalendarWeekView(c *ui.Context, week *Date, events []CalendarEvent, opts CalendarWeekViewOptions) CalendarViewResult {
	datetimeCheckDate("CalendarWeekView week", *week)
	if week.IsZero() {
		panic("mujica: CalendarWeekView needs a week")
	}
	cfg := timeGridConfig{name: "calendar_week_view", loc: eventLocation("CalendarWeekView", opts.Location), version: opts.Version,
		slot: timeGridSlot("CalendarWeekView", opts.SlotMinutes), today: datetimeToday(c, opts.Today), now: opts.Now, hour12: opts.Hour12}
	eventsCheck("CalendarWeekView", events)
	first := calendarFirstWeekday(c, opts.WeekStart)
	start := week.AddDays(-((int(week.Weekday()) - int(first) + 7) % 7))
	days := make([]Date, 7)
	for i := range days {
		days[i] = start.AddDays(i)
	}
	title := weekViewTitle(c, start, start.AddDays(6))
	label := opts.Label
	if label == "" {
		label = title
	}
	r := CalendarViewResult{}
	root := ui.Column(c).Gap(Density(c).Unit()).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	r.Element = root
	root.Children(func() {
		calendarViewHeader(c, title, msg(c, "calendar_week_view.prev"), msg(c, "calendar_week_view.next"), func(dir int) {
			*week = week.AddDays(7 * dir)
		})
		timeGrid(c, root, days, events, cfg, &r)
	})
	return r
}

// calendarViewHeader is the title row with previous and next buttons.
func calendarViewHeader(c *ui.Context, title, prev, next string, step func(dir int)) {
	k := Tokens(c)
	ui.Row(c).AlignItems(ui.Center).Gap(Density(c).Unit()).Children(func() {
		if b := iconAction(c, icons.Must("chevron-left"), prev); b.Clicked() && !b.IsDisabled() {
			step(-1)
		}
		ui.Text(c, title).Font(theme.SerifFont).FontSize(FontSize(c, theme.H3Size)).TextColor(k.Text).Grow(1).TextAlign(ui.Center)
		if b := iconAction(c, icons.Must("chevron-right"), next); b.Clicked() && !b.IsDisabled() {
			step(1)
		}
	})
}

// weekViewTitle names a span of days.
func weekViewTitle(c *ui.Context, a, b Date) string {
	return fmt.Sprintf(msg(c, "calendar_week_view.title"), datetimeMonthName(c, a.Month), a.Day, datetimeMonthName(c, b.Month), b.Day, b.Year)
}
