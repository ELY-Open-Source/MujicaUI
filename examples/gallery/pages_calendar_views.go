package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/datetime"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

// calSample is demo events around today in the local zone.
func calSample() []datetime.CalendarEvent {
	now := time.Now()
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	at := func(d, h, m, mins int) (time.Time, time.Time) {
		s := day.AddDate(0, 0, d).Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute)
		return s, s.Add(time.Duration(mins) * time.Minute)
	}
	mk := func(id, title string, color datetime.EventColor, d, h, m, mins int) datetime.CalendarEvent {
		s, e := at(d, h, m, mins)
		return datetime.CalendarEvent{ID: id, Title: title, Color: color, Start: s, End: e, Location: "Rehearsal hall"}
	}
	evs := []datetime.CalendarEvent{
		mk("1", "Rehearsal", datetime.EventAccent, 0, 9, 0, 120), mk("2", "Costume fitting", datetime.EventInfo, 0, 10, 0, 60),
		mk("3", "Lunch", datetime.EventSuccess, 0, 12, 30, 60), mk("4", "Lyrics review", datetime.EventWarning, 1, 14, 0, 90),
		mk("5", "Press", datetime.EventDanger, 2, 16, 0, 45), mk("6", "Sound check", datetime.EventNeutral, 3, 18, 0, 60),
		mk("7", "Encore notes", datetime.EventInfo, 0, 15, 0, 30), mk("8", "Stage walk", datetime.EventAccent, 0, 17, 0, 30),
	}
	holiday := mk("9", "Masquerade", datetime.EventAccent, 4, 0, 0, 24*60)
	holiday.AllDay, holiday.Start = true, day.AddDate(0, 0, 4)
	holiday.End = holiday.Start.AddDate(0, 0, 1)
	return append(evs, holiday)
}

// calApply applies and logs a view change.
func calApply(evs *[]datetime.CalendarEvent, ch datetime.CalendarChange, log func(string)) {
	switch ch.Kind {
	case datetime.CalendarCreate:
		id := fmt.Sprint(len(*evs) + 100)
		*evs = append(*evs, datetime.CalendarEvent{ID: id, Title: "New event", Start: ch.Start, End: ch.End})
		log("created " + ch.Start.Format("Jan 2 15:04") + " – " + ch.End.Format("15:04"))
		return
	}
	for i := range *evs {
		if (*evs)[i].ID == ch.ID {
			(*evs)[i].Start, (*evs)[i].End = ch.Start, ch.End
		}
	}
	verb := "moved"
	if ch.Kind == datetime.CalendarResize {
		verb = "resized"
	}
	log(fmt.Sprintf("%s %s to %s – %s", verb, ch.ID, ch.Start.Format("Jan 2 15:04"), ch.End.Format("15:04")))
}

// calViewResult logs what a view reported and applies its change.
func calViewResult(r datetime.CalendarViewResult, evs *[]datetime.CalendarEvent, version *int, log func(string)) {
	if id, ok := r.Opened(); ok {
		log("opened " + id)
	}
	if d, ok := r.DayChosen(); ok {
		log("day " + d.String())
	}
	if ch, ok := r.Changed(); ok {
		calApply(evs, ch, log)
		*version++
	}
}

var calWeekStarts = []string{"Locale", "Sunday", "Monday", "Saturday"}

func calWeekStart(s string) datetime.WeekStart {
	for i, n := range calWeekStarts {
		if n == s {
			return datetime.WeekStart(i)
		}
	}
	panic("gallery: unknown week start " + s)
}

func calMonthDemo(c *ui.Context, log func(string)) {
	type state struct {
		month             datetime.YearMonth
		events            []datetime.CalendarEvent
		version, maxRows  int
		weekStart, maxStr string
		disabled          bool
	}
	s := ui.Local(c.Root(), "421", func() state {
		now := time.Now()
		return state{month: datetime.YearMonth{Year: now.Year(), Month: now.Month()}, events: calSample(), weekStart: "Locale", maxStr: "3"}
	})
	dtSection(c, "CalendarMonthView", func() {
		dtControls(c, func() {
			input.Select(c, &s.weekStart, calOptions(calWeekStarts), input.SelectOptions{Label: "Week start"})
			input.Select(c, &s.maxStr, calOptions([]string{"1", "2", "3", "4"}), input.SelectOptions{Label: "Rows per day"})
			input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{})
		})
		s.maxRows = calAtoi(s.maxStr)
		r := datetime.CalendarMonthView(c, &s.month, s.events, datetime.CalendarMonthViewOptions{Location: time.Local, Version: s.version,
			WeekStart: calWeekStart(s.weekStart), MaxPerDay: s.maxRows, Disabled: s.disabled})
		calViewResult(r, &s.events, &s.version, log)
	})
}

func calWeekDemo(c *ui.Context, log func(string)) {
	type state struct {
		week     core.Date
		events   []datetime.CalendarEvent
		version  int
		slot     string
		hour12   bool
		disabled bool
	}
	s := ui.Local(c.Root(), "422", func() state {
		return state{week: core.DateOf(time.Now()), events: calSample(), slot: "30"}
	})
	dtSection(c, "CalendarWeekView", func() {
		dtControls(c, func() {
			input.Select(c, &s.slot, calOptions([]string{"15", "30", "60"}), input.SelectOptions{Label: "Slot minutes"})
			input.Checkbox(c, &s.hour12, "12-hour", input.CheckboxOptions{})
			input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{})
		})
		slot := calAtoi(s.slot)
		r := datetime.CalendarWeekView(c, &s.week, s.events, datetime.CalendarWeekViewOptions{Location: time.Local, Version: s.version,
			SlotMinutes: slot, Hour12: s.hour12, Disabled: s.disabled})
		calViewResult(r, &s.events, &s.version, log)
	})
}

func calDayDemo(c *ui.Context, log func(string)) {
	type state struct {
		day      core.Date
		events   []datetime.CalendarEvent
		version  int
		disabled bool
	}
	s := ui.Local(c.Root(), "423", func() state { return state{day: core.DateOf(time.Now()), events: calSample()} })
	dtSection(c, "CalendarDayView", func() {
		dtControls(c, func() { input.Checkbox(c, &s.disabled, "Disabled", input.CheckboxOptions{}) })
		r := datetime.CalendarDayView(c, &s.day, s.events, datetime.CalendarDayViewOptions{Location: time.Local, Version: s.version, Disabled: s.disabled})
		calViewResult(r, &s.events, &s.version, log)
	})
}

func calAgendaDemo(c *ui.Context, log func(string)) {
	type state struct {
		events []datetime.CalendarEvent
		days   string
	}
	s := ui.Local(c.Root(), "424", func() state { return state{events: calSample(), days: "7"} })
	dtSection(c, "AgendaView", func() {
		dtControls(c, func() {
			input.Select(c, &s.days, calOptions([]string{"1", "7", "30"}), input.SelectOptions{Label: "Days"})
		})
		days := calAtoi(s.days)
		ui.Column(c).Height(420).Children(func() {
			if id, ok := datetime.AgendaView(c, s.events, datetime.AgendaViewOptions{Location: time.Local, Days: days}).Opened(); ok {
				log("opened " + id)
			}
		})
	})
}

func calYearDemo(c *ui.Context, log func(string)) {
	type state struct {
		year   int
		events []datetime.CalendarEvent
	}
	s := ui.Local(c.Root(), "425", func() state { return state{year: time.Now().Year(), events: calSample()} })
	dtSection(c, "YearView", func() {
		if d, ok := datetime.YearView(c, &s.year, s.events, datetime.YearViewOptions{Location: time.Local}).DayChosen(); ok {
			log("day " + d.String())
		}
	})
}

func calAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func calOptions(values []string) []input.SelectOption[string] {
	out := make([]input.SelectOption[string], len(values))
	for i, v := range values {
		out[i] = input.SelectOption[string]{Value: v, Label: v}
	}
	return out
}
