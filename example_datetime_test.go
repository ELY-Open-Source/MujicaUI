package mujica_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleCalendar() {
	day := mujica.Date{Year: 2026, Month: time.October, Day: 5}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.Calendar(c, &day, mujica.CalendarOptions{WeekStart: mujica.WeekStartMonday}).Changed() {
			fmt.Println("chosen", day)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDatePicker() {
	due := mujica.Date{Year: 2026, Month: time.October, Day: 5}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		opts := mujica.DatePickerOptions{Label: "Due date", Clearable: true,
			Min: mujica.Date{Year: 2026, Month: time.January, Day: 1}}
		if mujica.DatePicker(c, &due, opts).Changed() {
			fmt.Println("due", due)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDateRangePicker() {
	var stay mujica.DateRange
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		r := mujica.DateRangePicker(c, &stay, mujica.DateRangePickerOptions{Label: "Stay"})
		if r.Changed() {
			fmt.Println("stay", stay.Start, stay.End, stay.Days())
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTimePicker() {
	alarm := mujica.TimeOfDay{Hour: 7, Minute: 30}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		opts := mujica.TimePickerOptions{Label: "Alarm", Hour12: true, MinuteStep: 5}
		if mujica.TimePicker(c, &alarm, opts).Changed() {
			fmt.Println("alarm", alarm)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTimeRangePicker() {
	shift := mujica.TimeRange{Start: mujica.TimeOfDay{Hour: 22}, End: mujica.TimeOfDay{Hour: 6}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		opts := mujica.TimeRangePickerOptions{Label: "Night shift", AllowOvernight: true}
		if mujica.TimeRangePicker(c, &shift, opts).Changed() {
			fmt.Println("shift", shift.Start, shift.End, shift.Overnight())
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDateTimePicker() {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	var meeting time.Time
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		opts := mujica.DateTimePickerOptions{Label: "Meeting", Location: tokyo, MinuteStep: 15}
		if mujica.DateTimePicker(c, &meeting, opts).Changed() {
			fmt.Println("meeting", meeting)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMonthPicker() {
	var month mujica.YearMonth
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		opts := mujica.MonthPickerOptions{Label: "Billing month", Clearable: true}
		if mujica.MonthPicker(c, &month, opts).Changed() {
			fmt.Println("month", month)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleYearPicker() {
	year := 2026
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		opts := mujica.YearPickerOptions{Label: "Vintage", Min: 1990, Max: 2030}
		if mujica.YearPicker(c, &year, opts).Changed() {
			fmt.Println("year", year)
		}
	}
	ui.Render(view, 400, 300, 1)
}
