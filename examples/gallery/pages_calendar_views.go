package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

// calSample is demo events around today in the local zone.
func calSample() []mujica.CalendarEvent {
	now := time.Now()
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	at := func(d, h, m, mins int) (time.Time, time.Time) {
		s := day.AddDate(0, 0, d).Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute)
		return s, s.Add(time.Duration(mins) * time.Minute)
	}
	mk := func(id, title string, color mujica.EventColor, d, h, m, mins int) mujica.CalendarEvent {
		s, e := at(d, h, m, mins)
		return mujica.CalendarEvent{ID: id, Title: title, Color: color, Start: s, End: e, Location: "Rehearsal hall"}
	}
	evs := []mujica.CalendarEvent{
		mk("1", "Rehearsal", mujica.EventAccent, 0, 9, 0, 120), mk("2", "Costume fitting", mujica.EventInfo, 0, 10, 0, 60),
		mk("3", "Lunch", mujica.EventSuccess, 0, 12, 30, 60), mk("4", "Lyrics review", mujica.EventWarning, 1, 14, 0, 90),
		mk("5", "Press", mujica.EventDanger, 2, 16, 0, 45), mk("6", "Sound check", mujica.EventNeutral, 3, 18, 0, 60),
		mk("7", "Encore notes", mujica.EventInfo, 0, 15, 0, 30), mk("8", "Stage walk", mujica.EventAccent, 0, 17, 0, 30),
	}
	holiday := mk("9", "Masquerade", mujica.EventAccent, 4, 0, 0, 24*60)
	holiday.AllDay, holiday.Start = true, day.AddDate(0, 0, 4)
	holiday.End = holiday.Start.AddDate(0, 0, 1)
	return append(evs, holiday)
}

// calApply applies and logs a view change.
func calApply(evs *[]mujica.CalendarEvent, ch mujica.CalendarChange, log func(string)) {
	switch ch.Kind {
	case mujica.CalendarCreate:
		id := fmt.Sprint(len(*evs) + 100)
		*evs = append(*evs, mujica.CalendarEvent{ID: id, Title: "New event", Start: ch.Start, End: ch.End})
		log("created " + ch.Start.Format("Jan 2 15:04") + " – " + ch.End.Format("15:04"))
		return
	}
	for i := range *evs {
		if (*evs)[i].ID == ch.ID {
			(*evs)[i].Start, (*evs)[i].End = ch.Start, ch.End
		}
	}
	verb := "moved"
	if ch.Kind == mujica.CalendarResize {
		verb = "resized"
	}
	log(fmt.Sprintf("%s %s to %s – %s", verb, ch.ID, ch.Start.Format("Jan 2 15:04"), ch.End.Format("15:04")))
}

// calViewResult logs what a view reported and applies its change.
func calViewResult(r mujica.CalendarViewResult, evs *[]mujica.CalendarEvent, version *int, log func(string)) {
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

func calWeekStart(s string) mujica.WeekStart {
	for i, n := range calWeekStarts {
		if n == s {
			return mujica.WeekStart(i)
		}
	}
	panic("gallery: unknown week start " + s)
}

func calMonthDemo(c *ui.Context, log func(string)) {
	type state struct {
		month             mujica.YearMonth
		events            []mujica.CalendarEvent
		version, maxRows  int
		weekStart, maxStr string
		disabled          bool
	}
	s := ui.Local(c.Root(), "421", func() state {
		now := time.Now()
		return state{month: mujica.YearMonth{Year: now.Year(), Month: now.Month()}, events: calSample(), weekStart: "Locale", maxStr: "3"}
	})
	dtSection(c, "CalendarMonthView", func() {
		dtControls(c, func() {
			mujica.Select(c, &s.weekStart, calOptions(calWeekStarts), mujica.SelectOptions{Label: "Week start"})
			mujica.Select(c, &s.maxStr, calOptions([]string{"1", "2", "3", "4"}), mujica.SelectOptions{Label: "Rows per day"})
			mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{})
		})
		s.maxRows = calAtoi(s.maxStr)
		r := mujica.CalendarMonthView(c, &s.month, s.events, mujica.CalendarMonthViewOptions{Location: time.Local, Version: s.version,
			WeekStart: calWeekStart(s.weekStart), MaxPerDay: s.maxRows, Disabled: s.disabled})
		calViewResult(r, &s.events, &s.version, log)
	})
}

func calWeekDemo(c *ui.Context, log func(string)) {
	type state struct {
		week     mujica.Date
		events   []mujica.CalendarEvent
		version  int
		slot     string
		hour12   bool
		disabled bool
	}
	s := ui.Local(c.Root(), "422", func() state {
		return state{week: mujica.DateOf(time.Now()), events: calSample(), slot: "30"}
	})
	dtSection(c, "CalendarWeekView", func() {
		dtControls(c, func() {
			mujica.Select(c, &s.slot, calOptions([]string{"15", "30", "60"}), mujica.SelectOptions{Label: "Slot minutes"})
			mujica.Checkbox(c, &s.hour12, "12-hour", mujica.CheckboxOptions{})
			mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{})
		})
		slot := calAtoi(s.slot)
		r := mujica.CalendarWeekView(c, &s.week, s.events, mujica.CalendarWeekViewOptions{Location: time.Local, Version: s.version,
			SlotMinutes: slot, Hour12: s.hour12, Disabled: s.disabled})
		calViewResult(r, &s.events, &s.version, log)
	})
}

func calDayDemo(c *ui.Context, log func(string)) {
	type state struct {
		day      mujica.Date
		events   []mujica.CalendarEvent
		version  int
		disabled bool
	}
	s := ui.Local(c.Root(), "423", func() state { return state{day: mujica.DateOf(time.Now()), events: calSample()} })
	dtSection(c, "CalendarDayView", func() {
		dtControls(c, func() { mujica.Checkbox(c, &s.disabled, "Disabled", mujica.CheckboxOptions{}) })
		r := mujica.CalendarDayView(c, &s.day, s.events, mujica.CalendarDayViewOptions{Location: time.Local, Version: s.version, Disabled: s.disabled})
		calViewResult(r, &s.events, &s.version, log)
	})
}

func calAgendaDemo(c *ui.Context, log func(string)) {
	type state struct {
		events []mujica.CalendarEvent
		days   string
	}
	s := ui.Local(c.Root(), "424", func() state { return state{events: calSample(), days: "7"} })
	dtSection(c, "AgendaView", func() {
		dtControls(c, func() {
			mujica.Select(c, &s.days, calOptions([]string{"1", "7", "30"}), mujica.SelectOptions{Label: "Days"})
		})
		days := calAtoi(s.days)
		ui.Column(c).Height(420).Children(func() {
			if id, ok := mujica.AgendaView(c, s.events, mujica.AgendaViewOptions{Location: time.Local, Days: days}).Opened(); ok {
				log("opened " + id)
			}
		})
	})
}

func calYearDemo(c *ui.Context, log func(string)) {
	type state struct {
		year   int
		events []mujica.CalendarEvent
	}
	s := ui.Local(c.Root(), "425", func() state { return state{year: time.Now().Year(), events: calSample()} })
	dtSection(c, "YearView", func() {
		if d, ok := mujica.YearView(c, &s.year, s.events, mujica.YearViewOptions{Location: time.Local}).DayChosen(); ok {
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

func calOptions(values []string) []mujica.SelectOption[string] {
	out := make([]mujica.SelectOption[string], len(values))
	for i, v := range values {
		out[i] = mujica.SelectOption[string]{Value: v, Label: v}
	}
	return out
}
