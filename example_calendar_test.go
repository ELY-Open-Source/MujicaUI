package mujica_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func exampleEvents(loc *time.Location) []mujica.CalendarEvent {
	start := time.Date(2026, time.October, 14, 9, 0, 0, 0, loc)
	return []mujica.CalendarEvent{
		{ID: "rehearsal", Title: "Rehearsal", Start: start, End: start.Add(2 * time.Hour), Location: "Hall"},
		{ID: "fitting", Title: "Fitting", Start: start.Add(time.Hour), End: start.Add(90 * time.Minute), Color: mujica.EventInfo},
	}
}

func ExampleCalendarMonthView() {
	loc := time.Local
	events := exampleEvents(loc)
	month := mujica.YearMonth{Year: 2026, Month: time.October}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.CalendarMonthView(c, &month, events, mujica.CalendarMonthViewOptions{Location: loc})
		if ch, ok := r.Changed(); ok {
			fmt.Println("move", ch.ID, "to", ch.Start)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCalendarWeekView() {
	loc := time.Local
	events := exampleEvents(loc)
	week := mujica.Date{Year: 2026, Month: time.October, Day: 14}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.CalendarWeekView(c, &week, events, mujica.CalendarWeekViewOptions{Location: loc, SlotMinutes: 15})
		if ch, ok := r.Changed(); ok && ch.Kind == mujica.CalendarCreate {
			fmt.Println("new event", ch.Start, ch.End)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCalendarDayView() {
	loc := time.Local
	events := exampleEvents(loc)
	day := mujica.Date{Year: 2026, Month: time.October, Day: 14}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if id, ok := mujica.CalendarDayView(c, &day, events, mujica.CalendarDayViewOptions{Location: loc}).Opened(); ok {
			fmt.Println("open", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAgendaView() {
	loc := time.Local
	events := exampleEvents(loc)
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		opts := mujica.AgendaViewOptions{Location: loc, From: mujica.Date{Year: 2026, Month: time.October, Day: 12}, Days: 14}
		if id, ok := mujica.AgendaView(c, events, opts).Opened(); ok {
			fmt.Println("open", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleYearView() {
	loc := time.Local
	events := exampleEvents(loc)
	year := 2026
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if day, ok := mujica.YearView(c, &year, events, mujica.YearViewOptions{Location: loc}).DayChosen(); ok {
			fmt.Println("go to", day)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleEventChip() {
	loc := time.Local
	ev := exampleEvents(loc)[0]
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.EventChip(c, ev, mujica.EventChipOptions{Location: loc}).Clicked() {
			fmt.Println("open", ev.ID)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleEventEditor() {
	loc := time.Local
	ev := exampleEvents(loc)[0]
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.EventEditor(c, &ev, mujica.EventEditorOptions{Location: loc})
		if r.Saved() {
			fmt.Println("saved", ev.Title, ev.Start)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleEventPopover() {
	loc := time.Local
	ev := exampleEvents(loc)[0]
	open := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		chip := mujica.EventChip(c, ev, mujica.EventChipOptions{Location: loc})
		if chip.Clicked() {
			open = true
		}
		r := mujica.EventPopover(c, chip, &open, ev, mujica.EventPopoverOptions{Location: loc})
		if r.Delete() {
			fmt.Println("delete", ev.ID)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCurrentTimeIndicator() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		ui.Box(c).Size(200, 240).Children(func() {
			mujica.CurrentTimeIndicator(c, mujica.CurrentTimeIndicatorOptions{Location: time.Local})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRecurrenceEditor() {
	rule := mujica.Recurrence{Unit: mujica.RepeatWeekly, Interval: 1, Weekdays: []time.Weekday{time.Monday}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.RecurrenceEditor(c, &rule, mujica.RecurrenceEditorOptions{Start: mujica.Date{Year: 2026, Month: time.October, Day: 5}})
		if r.Changed() && r.Problem() == "" {
			fmt.Println(mujica.RecurrenceSummary(c, rule))
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAvailabilityPicker() {
	days := []mujica.Date{{Year: 2026, Month: time.October, Day: 12}, {Year: 2026, Month: time.October, Day: 13}}
	mine := [][]bool{make([]bool, 18), make([]bool, 18)}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.AvailabilityPicker(c, &mine, mujica.AvailabilityPickerOptions{Days: days}).Changed() {
			fmt.Println("availability changed")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTimezoneSelect() {
	zone := "Asia/Shanghai"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.TimezoneSelect(c, &zone, mujica.TimezoneSelectOptions{}).Changed() {
			fmt.Println("zone", zone)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleReminderPicker() {
	reminders := []time.Duration{15 * time.Minute}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.ReminderPicker(c, &reminders, mujica.ReminderPickerOptions{Max: 3}).Changed() {
			fmt.Println("reminders", reminders)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAttendeeList() {
	people := []mujica.Attendee{{Name: "Sakiko", Status: mujica.AttendeeAccepted, Organizer: true}, {Name: "Mutsumi"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if i, ok := mujica.AttendeeList(c, people, mujica.AttendeeListOptions{Removable: true}).Removed(); ok {
			people = append(people[:i], people[i+1:]...)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleWeekPicker() {
	var week mujica.Date
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.WeekPicker(c, &week, mujica.WeekPickerOptions{WeekStart: mujica.WeekStartMonday, Clearable: true}).Changed() {
			fmt.Println("week of", week)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDurationPicker() {
	length := 90 * time.Minute
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		opts := mujica.DurationPickerOptions{Min: 15 * time.Minute, Max: 8 * time.Hour}
		if mujica.DurationPicker(c, &length, opts).Changed() {
			fmt.Println("length", length)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCronEditor() {
	expr := "0 9 * * MON-FRI"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.CronEditor(c, &expr, mujica.CronEditorOptions{Location: time.Local}).Changed() {
			if _, err := mujica.ParseCron(expr); err == nil {
				fmt.Println("schedule", expr)
			}
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRelativeTime() {
	posted := time.Date(2026, time.October, 5, 9, 0, 0, 0, time.Local)
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.RelativeTime(c, posted, mujica.RelativeTimeOptions{Location: time.Local})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCountdown() {
	deadline := time.Date(2026, time.December, 31, 23, 59, 59, 0, time.Local)
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.Countdown(c, deadline, mujica.CountdownOptions{}).Finished() {
			fmt.Println("happy new year")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStopwatch() {
	var sw mujica.StopwatchState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.Stopwatch(c, &sw, mujica.StopwatchOptions{}).Changed() {
			fmt.Println("laps", len(sw.Laps))
		}
	}
	ui.Render(view, 400, 300, 1)
}
