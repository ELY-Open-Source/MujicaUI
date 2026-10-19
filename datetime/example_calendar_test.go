package datetime_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/datetime"
)

func ExampleCalendarMonthView() {
	loc := time.Local
	events := exampleEvents(loc)
	month := datetime.YearMonth{Year: 2026, Month: time.October}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := datetime.CalendarMonthView(c, &month, events, datetime.CalendarMonthViewOptions{Location: loc})
		if ch, ok := r.Changed(); ok {
			fmt.Println("move", ch.ID, "to", ch.Start)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCalendarWeekView() {
	loc := time.Local
	events := exampleEvents(loc)
	week := core.Date{Year: 2026, Month: time.October, Day: 14}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := datetime.CalendarWeekView(c, &week, events, datetime.CalendarWeekViewOptions{Location: loc, SlotMinutes: 15})
		if ch, ok := r.Changed(); ok && ch.Kind == datetime.CalendarCreate {
			fmt.Println("new event", ch.Start, ch.End)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCalendarDayView() {
	loc := time.Local
	events := exampleEvents(loc)
	day := core.Date{Year: 2026, Month: time.October, Day: 14}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if id, ok := datetime.CalendarDayView(c, &day, events, datetime.CalendarDayViewOptions{Location: loc}).Opened(); ok {
			fmt.Println("open", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAgendaView() {
	loc := time.Local
	events := exampleEvents(loc)
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		opts := datetime.AgendaViewOptions{Location: loc, From: core.Date{Year: 2026, Month: time.October, Day: 12}, Days: 14}
		if id, ok := datetime.AgendaView(c, events, opts).Opened(); ok {
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
		core.Use(c, core.Settings{})
		if day, ok := datetime.YearView(c, &year, events, datetime.YearViewOptions{Location: loc}).DayChosen(); ok {
			fmt.Println("go to", day)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleEventChip() {
	loc := time.Local
	ev := exampleEvents(loc)[0]
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if datetime.EventChip(c, ev, datetime.EventChipOptions{Location: loc}).Clicked() {
			fmt.Println("open", ev.ID)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleEventEditor() {
	loc := time.Local
	ev := exampleEvents(loc)[0]
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := datetime.EventEditor(c, &ev, datetime.EventEditorOptions{Location: loc})
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
		core.Use(c, core.Settings{})
		chip := datetime.EventChip(c, ev, datetime.EventChipOptions{Location: loc})
		if chip.Clicked() {
			open = true
		}
		r := datetime.EventPopover(c, chip, &open, ev, datetime.EventPopoverOptions{Location: loc})
		if r.Delete() {
			fmt.Println("delete", ev.ID)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCurrentTimeIndicator() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Box(c).Size(200, 240).Children(func() {
			datetime.CurrentTimeIndicator(c, datetime.CurrentTimeIndicatorOptions{Location: time.Local})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRecurrenceEditor() {
	rule := datetime.Recurrence{Unit: datetime.RepeatWeekly, Interval: 1, Weekdays: []time.Weekday{time.Monday}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := datetime.RecurrenceEditor(c, &rule, datetime.RecurrenceEditorOptions{Start: core.Date{Year: 2026, Month: time.October, Day: 5}})
		if r.Changed() && r.Problem() == "" {
			fmt.Println(datetime.RecurrenceSummary(c, rule))
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAvailabilityPicker() {
	days := []core.Date{{Year: 2026, Month: time.October, Day: 12}, {Year: 2026, Month: time.October, Day: 13}}
	mine := [][]bool{make([]bool, 18), make([]bool, 18)}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if datetime.AvailabilityPicker(c, &mine, datetime.AvailabilityPickerOptions{Days: days}).Changed() {
			fmt.Println("availability changed")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTimezoneSelect() {
	zone := "Asia/Shanghai"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if datetime.TimezoneSelect(c, &zone, datetime.TimezoneSelectOptions{}).Changed() {
			fmt.Println("zone", zone)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleReminderPicker() {
	reminders := []time.Duration{15 * time.Minute}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if datetime.ReminderPicker(c, &reminders, datetime.ReminderPickerOptions{Max: 3}).Changed() {
			fmt.Println("reminders", reminders)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAttendeeList() {
	people := []datetime.Attendee{{Name: "Sakiko", Status: datetime.AttendeeAccepted, Organizer: true}, {Name: "Mutsumi"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if i, ok := datetime.AttendeeList(c, people, datetime.AttendeeListOptions{Removable: true}).Removed(); ok {
			people = append(people[:i], people[i+1:]...)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleWeekPicker() {
	var week core.Date
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if datetime.WeekPicker(c, &week, datetime.WeekPickerOptions{WeekStart: datetime.WeekStartMonday, Clearable: true}).Changed() {
			fmt.Println("week of", week)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDurationPicker() {
	length := 90 * time.Minute
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		opts := datetime.DurationPickerOptions{Min: 15 * time.Minute, Max: 8 * time.Hour}
		if datetime.DurationPicker(c, &length, opts).Changed() {
			fmt.Println("length", length)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCronEditor() {
	expr := "0 9 * * MON-FRI"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if datetime.CronEditor(c, &expr, datetime.CronEditorOptions{Location: time.Local}).Changed() {
			if _, err := datetime.ParseCron(expr); err == nil {
				fmt.Println("schedule", expr)
			}
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleRelativeTime() {
	posted := time.Date(2026, time.October, 5, 9, 0, 0, 0, time.Local)
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		datetime.RelativeTime(c, posted, datetime.RelativeTimeOptions{Location: time.Local})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCountdown() {
	deadline := time.Date(2026, time.December, 31, 23, 59, 59, 0, time.Local)
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if datetime.Countdown(c, deadline, datetime.CountdownOptions{}).Finished() {
			fmt.Println("happy new year")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStopwatch() {
	var sw datetime.StopwatchState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if datetime.Stopwatch(c, &sw, datetime.StopwatchOptions{}).Changed() {
			fmt.Println("laps", len(sw.Laps))
		}
	}
	ui.Render(view, 400, 300, 1)
}
