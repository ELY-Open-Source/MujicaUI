package datetime_test

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/datetime"
)

func ExampleCalendar() {
	day := core.Date{Year: 2026, Month: time.October, Day: 5}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if datetime.Calendar(c, &day, datetime.CalendarOptions{WeekStart: datetime.WeekStartMonday}).Changed() {
			fmt.Println("chosen", day)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDatePicker() {
	due := core.Date{Year: 2026, Month: time.October, Day: 5}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		opts := datetime.DatePickerOptions{Label: "Due date", Clearable: true,
			Min: core.Date{Year: 2026, Month: time.January, Day: 1}}
		if datetime.DatePicker(c, &due, opts).Changed() {
			fmt.Println("due", due)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDateRangePicker() {
	var stay datetime.DateRange
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := datetime.DateRangePicker(c, &stay, datetime.DateRangePickerOptions{Label: "Stay"})
		if r.Changed() {
			fmt.Println("stay", stay.Start, stay.End, stay.Days())
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTimePicker() {
	alarm := datetime.TimeOfDay{Hour: 7, Minute: 30}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		opts := datetime.TimePickerOptions{Label: "Alarm", Hour12: true, MinuteStep: 5}
		if datetime.TimePicker(c, &alarm, opts).Changed() {
			fmt.Println("alarm", alarm)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTimeRangePicker() {
	shift := datetime.TimeRange{Start: datetime.TimeOfDay{Hour: 22}, End: datetime.TimeOfDay{Hour: 6}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		opts := datetime.TimeRangePickerOptions{Label: "Night shift", AllowOvernight: true}
		if datetime.TimeRangePicker(c, &shift, opts).Changed() {
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
		core.Use(c, core.Settings{})
		opts := datetime.DateTimePickerOptions{Label: "Meeting", Location: tokyo, MinuteStep: 15}
		if datetime.DateTimePicker(c, &meeting, opts).Changed() {
			fmt.Println("meeting", meeting)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMonthPicker() {
	var month datetime.YearMonth
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		opts := datetime.MonthPickerOptions{Label: "Billing month", Clearable: true}
		if datetime.MonthPicker(c, &month, opts).Changed() {
			fmt.Println("month", month)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleYearPicker() {
	year := 2026
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		opts := datetime.YearPickerOptions{Label: "Vintage", Min: 1990, Max: 2030}
		if datetime.YearPicker(c, &year, opts).Changed() {
			fmt.Println("year", year)
		}
	}
	ui.Render(view, 400, 300, 1)
}
