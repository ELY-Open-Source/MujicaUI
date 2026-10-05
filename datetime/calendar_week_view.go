package datetime

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
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
	Today core.Date
	// Now places the current-time line; zero uses the system clock.
	Now      time.Time
	Hour12   bool
	Disabled bool
	Label    string
}

// CalendarWeekView shows *week as a time grid; drags edit events.
func CalendarWeekView(c *ui.Context, week *core.Date, events []CalendarEvent, opts CalendarWeekViewOptions) CalendarViewResult {
	core.DatetimeCheckDate("CalendarWeekView week", *week)
	if week.IsZero() {
		panic("mujica: CalendarWeekView needs a week")
	}
	cfg := timeGridConfig{name: "calendar_week_view", loc: eventLocation("CalendarWeekView", opts.Location), version: opts.Version,
		slot: timeGridSlot("CalendarWeekView", opts.SlotMinutes), today: core.DatetimeToday(c, opts.Today), now: opts.Now, hour12: opts.Hour12}
	eventsCheck("CalendarWeekView", events)
	first := CalendarFirstWeekday(c, opts.WeekStart)
	start := week.AddDays(-((int(week.Weekday()) - int(first) + 7) % 7))
	days := make([]core.Date, 7)
	for i := range days {
		days[i] = start.AddDays(i)
	}
	title := WeekViewTitle(c, start, start.AddDays(6))
	label := opts.Label
	if label == "" {
		label = title
	}
	r := CalendarViewResult{}
	root := ui.Column(c).Gap(core.Density(c).Unit()).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	r.Element = root
	root.Children(func() {
		calendarViewHeader(c, title, core.Msg(c, "calendar_week_view.prev"), core.Msg(c, "calendar_week_view.next"), func(dir int) {
			*week = week.AddDays(7 * dir)
		})
		timeGrid(c, root, days, events, cfg, &r)
	})
	return r
}

// calendarViewHeader is the title row with previous and next buttons.
func calendarViewHeader(c *ui.Context, title, prev, next string, step func(dir int)) {
	k := core.Tokens(c)
	ui.Row(c).AlignItems(ui.Center).Gap(core.Density(c).Unit()).Children(func() {
		if b := core.IconAction(c, icons.Must("chevron-left"), prev); b.Clicked() && !b.IsDisabled() {
			step(-1)
		}
		ui.Text(c, title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).TextColor(k.Text).Grow(1).TextAlign(ui.Center)
		if b := core.IconAction(c, icons.Must("chevron-right"), next); b.Clicked() && !b.IsDisabled() {
			step(1)
		}
	})
}

// WeekViewTitle names a span of days.
func WeekViewTitle(c *ui.Context, a, b core.Date) string {
	return fmt.Sprintf(core.Msg(c, "calendar_week_view.title"), core.DatetimeMonthName(c, a.Month), a.Day, core.DatetimeMonthName(c, b.Month), b.Day, b.Year)
}
