package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/theme"
)

const dtCategory = "Date & Time"

var datetimePages = []page{
	{ID: "051", Name: "Calendar", ZhName: "日历", Category: dtCategory, Demo: dtCalendarDemo,
		Code: dtCalendarCode, Keys: "Tab enters the days; arrows move by day and week; PageUp/PageDown by month, with Shift by year; Home/End go to the first and last day; Enter or Space chooses."},
	{ID: "052", Name: "DatePicker", ZhName: "日期选择器", Category: dtCategory, Demo: dtDatePickerDemo,
		Code: dtDatePickerCode, Keys: "Type a date and press Enter, or leave the field, to commit; Down opens the calendar; in it arrows, PageUp/PageDown, Enter choose and Escape closes."},
	{ID: "053", Name: "DateRangePicker", ZhName: "日期范围选择器", Category: dtCategory, Demo: dtDateRangeDemo,
		Code: dtDateRangeCode, Keys: "Enter, Space or Down opens; Enter picks the start, then the end; Tab reaches Clear, Cancel and Confirm; Escape discards the draft."},
	{ID: "054", Name: "TimePicker", ZhName: "时间选择器", Category: dtCategory, Demo: dtTimePickerDemo,
		Code: dtTimePickerCode, Keys: "Tab enters the segments; Left/Right move between them; Up/Down step; digits type; A/P set the period; Alt+Down or the clock button opens the lists."},
	{ID: "055", Name: "TimeRangePicker", ZhName: "时间范围选择器", Category: dtCategory, Demo: dtTimeRangeDemo,
		Code: dtTimeRangeCode, Keys: "Each end works as a TimePicker; Tab moves from the start to the end."},
	{ID: "056", Name: "DateTimePicker", ZhName: "日期时间选择器", Category: dtCategory, Demo: dtDateTimeDemo,
		Code: dtDateTimeCode, Keys: "Enter, Space or Down opens; the calendar and time segments work as in DatePicker and TimePicker; Tab reaches Cancel and Confirm; Escape cancels."},
	{ID: "057", Name: "MonthPicker", ZhName: "月份选择器", Category: dtCategory, Demo: dtMonthDemo,
		Code: dtMonthCode, Keys: "Enter, Space or Down opens; arrows move by month and row; PageUp/PageDown by year; Enter chooses; Escape closes."},
	{ID: "058", Name: "YearPicker", ZhName: "年份选择器", Category: dtCategory, Demo: dtYearDemo,
		Code: dtYearCode, Keys: "Enter, Space or Down opens; arrows move by year and row; PageUp/PageDown by decade; Enter chooses; Escape closes."},
}

// dtSection lays out a captioned block of a demo.
func dtSection(c *ui.Context, title string, fn func()) {
	k := mujica.Tokens(c)
	ui.Column(c).Gap(8).Children(func() {
		ui.Text(c, title).Font(theme.SerifFont).FontSize(theme.H3Size).TextColor(k.Text)
		fn()
	})
}

// dtStates lays labeled variants side by side.
func dtStates(c *ui.Context, fn func(state func(name string, body func()))) {
	k := mujica.Tokens(c)
	ui.Row(c).Gap(24).Wrap().AlignItems(ui.Start).Children(func() {
		fn(func(name string, body func()) {
			ui.Column(c).Gap(6).AlignItems(ui.Start).Children(func() {
				ui.Text(c, name).FontSize(theme.CaptionSize).TextColor(k.TextMuted)
				body()
			})
		})
	})
}

func dtControls(c *ui.Context, fn func()) {
	ui.Row(c).Gap(16).Wrap().AlignItems(ui.Center).Children(fn)
}

var dtWeekStarts = []string{"Locale", "Sunday", "Monday", "Saturday"}

func dtWeekStart(s string) mujica.WeekStart {
	for i, n := range dtWeekStarts {
		if n == s {
			return mujica.WeekStart(i)
		}
	}
	panic("gallery: unknown week start " + s)
}

func dtCalendarDemo(c *ui.Context, log func(string)) {
	type state struct {
		day                     mujica.Date
		weekStart               string
		weekends, bound, frozen bool
	}
	s := ui.Local(c.Root(), "051", func() state { return state{day: mujica.DateOf(time.Now()), weekStart: "Locale"} })
	dtSection(c, "Calendar", func() {
		dtControls(c, func() {
			ui.Select(c, &s.weekStart, dtWeekStarts)
			ui.Checkbox(c, &s.weekends, "Disable weekends")
			ui.Checkbox(c, &s.bound, "Only this month")
			ui.Checkbox(c, &s.frozen, "Disabled")
		})
		opts := mujica.CalendarOptions{WeekStart: dtWeekStart(s.weekStart), Disabled: s.frozen, Label: "Calendar demo"}
		if s.weekends {
			opts.DisableDate = func(d mujica.Date) bool { return d.Weekday() == time.Saturday || d.Weekday() == time.Sunday }
		}
		if s.bound {
			now := mujica.DateOf(time.Now())
			opts.Min = mujica.Date{Year: now.Year, Month: now.Month, Day: 1}
			opts.Max = opts.Min.AddMonths(1).AddDays(-1)
		}
		if mujica.Calendar(c, &s.day, opts).Changed() {
			log("Calendar changed: " + s.day.String())
		}
	})
}

func dtDatePickerDemo(c *ui.Context, log func(string)) {
	type state struct {
		a, b, d          mujica.Date
		clearable, slash bool
	}
	s := ui.Local(c.Root(), "052", func() state {
		t := mujica.DateOf(time.Now())
		return state{a: t, d: t, clearable: true}
	})
	dtSection(c, "DatePicker", func() {
		dtControls(c, func() {
			ui.Checkbox(c, &s.clearable, "Clearable")
			ui.Checkbox(c, &s.slash, "Layout 02/01/2006")
		})
		layout := ""
		if s.slash {
			layout = "02/01/2006"
		}
		now := mujica.DateOf(time.Now())
		dtStates(c, func(state func(string, func())) {
			state("Default", func() {
				if mujica.DatePicker(c, &s.a, mujica.DatePickerOptions{Label: "Due date", Layout: layout, Clearable: s.clearable}).Changed() {
					log("DatePicker changed: " + s.a.String())
				}
			})
			state("Next 30 days only", func() {
				opts := mujica.DatePickerOptions{Label: "Delivery", Layout: layout, Clearable: s.clearable, Min: now, Max: now.AddDays(30)}
				if mujica.DatePicker(c, &s.b, opts).Changed() {
					log("DatePicker (bounded) changed: " + s.b.String())
				}
			})
			state("Disabled", func() {
				mujica.DatePicker(c, &s.d, mujica.DatePickerOptions{Label: "Locked", Layout: layout, Disabled: true})
			})
		})
	})
}

func dtDateRangeDemo(c *ui.Context, log func(string)) {
	type state struct {
		a            mujica.DateRange
		noTuesdays   bool
		frozen       bool
		fromThisWeek bool
	}
	s := ui.Local(c.Root(), "053", func() state { return state{} })
	dtSection(c, "DateRangePicker", func() {
		dtControls(c, func() {
			ui.Checkbox(c, &s.noTuesdays, "Tuesdays unavailable")
			ui.Checkbox(c, &s.fromThisWeek, "Start from today")
			ui.Checkbox(c, &s.frozen, "Disabled")
		})
		opts := mujica.DateRangePickerOptions{Label: "Stay", Disabled: s.frozen}
		if s.noTuesdays {
			opts.DisableDate = func(d mujica.Date) bool { return d.Weekday() == time.Tuesday }
		}
		if s.fromThisWeek {
			opts.Min = mujica.DateOf(time.Now())
		}
		r := mujica.DateRangePicker(c, &s.a, opts)
		if r.Changed() {
			log(fmt.Sprintf("DateRangePicker changed: %v – %v (%d days)", s.a.Start, s.a.End, s.a.Days()))
		}
		if r.Cancelled() {
			log("DateRangePicker cancelled")
		}
	})
}

func dtTimePickerDemo(c *ui.Context, log func(string)) {
	type state struct {
		a, b            mujica.TimeOfDay
		seconds, hour12 bool
		step            string
	}
	s := ui.Local(c.Root(), "054", func() state {
		return state{a: mujica.TimeOfDay{Hour: 9, Minute: 30}, b: mujica.TimeOfDay{Hour: 10}, step: "1"}
	})
	dtSection(c, "TimePicker", func() {
		dtControls(c, func() {
			ui.Checkbox(c, &s.seconds, "Seconds")
			ui.Checkbox(c, &s.hour12, "12-hour")
			ui.Text(c, "Minute step")
			ui.Select(c, &s.step, []string{"1", "5", "15", "30"})
		})
		step := 0
		if _, err := fmt.Sscan(s.step, &step); err != nil {
			panic(err)
		}
		base := mujica.TimePickerOptions{Seconds: s.seconds, Hour12: s.hour12, MinuteStep: step}
		dtStates(c, func(state func(string, func())) {
			state("Any time", func() {
				o := base
				o.Label = "Alarm"
				if mujica.TimePicker(c, &s.a, o).Changed() {
					log("TimePicker changed: " + s.a.String())
				}
			})
			state("Office hours 09:00–18:00", func() {
				o := base
				o.Label, o.Min, o.Max = "Meeting", mujica.TimeOfDay{Hour: 9}, mujica.TimeOfDay{Hour: 18}
				if mujica.TimePicker(c, &s.b, o).Changed() {
					log("TimePicker (bounded) changed: " + s.b.String())
				}
			})
			state("Disabled", func() {
				o := base
				o.Label, o.Disabled = "Locked", true
				mujica.TimePicker(c, &s.a, o)
			})
		})
	})
}

func dtTimeRangeDemo(c *ui.Context, log func(string)) {
	type state struct {
		v         mujica.TimeRange
		overnight bool
	}
	s := ui.Local(c.Root(), "055", func() state {
		return state{v: mujica.TimeRange{Start: mujica.TimeOfDay{Hour: 9}, End: mujica.TimeOfDay{Hour: 17}}}
	})
	dtSection(c, "TimeRangePicker", func() {
		dtControls(c, func() {
			if ui.Checkbox(c, &s.overnight, "Allow overnight").Changed() && !s.overnight && s.v.Overnight() {
				s.v = mujica.TimeRange{Start: s.v.End, End: s.v.Start}
				log("TimeRangePicker swapped ends: overnight turned off")
			}
		})
		opts := mujica.TimeRangePickerOptions{Label: "Opening hours", AllowOvernight: s.overnight, MinuteStep: 15}
		if mujica.TimeRangePicker(c, &s.v, opts).Changed() {
			log(fmt.Sprintf("TimeRangePicker changed: %v – %v overnight=%v", s.v.Start, s.v.End, s.v.Overnight()))
		}
	})
}

var dtZones = []string{"Local", "UTC", "Asia/Shanghai", "America/New_York", "Europe/London"}

func dtDateTimeDemo(c *ui.Context, log func(string)) {
	type state struct {
		v      time.Time
		zone   string
		hour12 bool
	}
	s := ui.Local(c.Root(), "056", func() state { return state{zone: "Local"} })
	dtSection(c, "DateTimePicker", func() {
		dtControls(c, func() {
			ui.Text(c, "Time zone")
			ui.Select(c, &s.zone, dtZones)
			ui.Checkbox(c, &s.hour12, "12-hour")
		})
		loc, err := time.LoadLocation(s.zone)
		if err != nil {
			panic(err)
		}
		r := mujica.DateTimePicker(c, &s.v, mujica.DateTimePickerOptions{Label: "Meeting", Location: loc, Hour12: s.hour12})
		if r.Changed() {
			log("DateTimePicker changed: " + s.v.Format(time.RFC3339))
		}
		if r.Cancelled() {
			log("DateTimePicker cancelled")
		}
		k := mujica.Tokens(c)
		ui.Text(c, "Try 2026-03-08 02:30 in America/New_York: the picker reports the DST gap.").
			FontSize(theme.CaptionSize).TextColor(k.TextMuted)
	})
}

func dtMonthDemo(c *ui.Context, log func(string)) {
	type state struct {
		v               mujica.YearMonth
		clearable, past bool
	}
	s := ui.Local(c.Root(), "057", func() state { return state{clearable: true} })
	dtSection(c, "MonthPicker", func() {
		dtControls(c, func() {
			ui.Checkbox(c, &s.clearable, "Clearable")
			ui.Checkbox(c, &s.past, "Past months only")
		})
		opts := mujica.MonthPickerOptions{Label: "Billing month", Clearable: s.clearable,
			DisableMonth: func(m mujica.YearMonth) bool { return m.Month == time.February }}
		if s.past {
			now := time.Now()
			opts.Max = mujica.YearMonth{Year: now.Year(), Month: now.Month()}
		}
		dtStates(c, func(state func(string, func())) {
			state("February unavailable", func() {
				if mujica.MonthPicker(c, &s.v, opts).Changed() {
					log("MonthPicker changed: " + s.v.String())
				}
			})
			state("Disabled", func() {
				mujica.MonthPicker(c, &s.v, mujica.MonthPickerOptions{Label: "Locked", Disabled: true})
			})
		})
	})
}

func dtYearDemo(c *ui.Context, log func(string)) {
	type state struct {
		v         int
		clearable bool
	}
	s := ui.Local(c.Root(), "058", func() state { return state{v: time.Now().Year(), clearable: true} })
	dtSection(c, "YearPicker", func() {
		dtControls(c, func() { ui.Checkbox(c, &s.clearable, "Clearable") })
		opts := mujica.YearPickerOptions{Label: "Vintage", Min: 1990, Max: time.Now().Year(), Clearable: s.clearable,
			DisableYear: func(y int) bool { return y == 2020 }}
		dtStates(c, func(state func(string, func())) {
			state("1990 to now, 2020 unavailable", func() {
				if mujica.YearPicker(c, &s.v, opts).Changed() {
					log(fmt.Sprintf("YearPicker changed: %d", s.v))
				}
			})
			state("Disabled", func() {
				mujica.YearPicker(c, &s.v, mujica.YearPickerOptions{Label: "Locked", Disabled: true})
			})
		})
	})
}

// Code snippets mirror example_datetime_test.go.
const (
	dtCalendarCode = `day := mujica.Date{Year: 2026, Month: time.October, Day: 5}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.Calendar(c, &day, mujica.CalendarOptions{WeekStart: mujica.WeekStartMonday}).Changed() {
		fmt.Println("chosen", day)
	}
}
ui.Render(view, 400, 300, 1)`
	dtDatePickerCode = `due := mujica.Date{Year: 2026, Month: time.October, Day: 5}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	opts := mujica.DatePickerOptions{Label: "Due date", Clearable: true,
		Min: mujica.Date{Year: 2026, Month: time.January, Day: 1}}
	if mujica.DatePicker(c, &due, opts).Changed() {
		fmt.Println("due", due)
	}
}
ui.Render(view, 400, 300, 1)`
	dtDateRangeCode = `var stay mujica.DateRange
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.DateRangePicker(c, &stay, mujica.DateRangePickerOptions{Label: "Stay"})
	if r.Changed() {
		fmt.Println("stay", stay.Start, stay.End, stay.Days())
	}
}
ui.Render(view, 400, 300, 1)`
	dtTimePickerCode = `alarm := mujica.TimeOfDay{Hour: 7, Minute: 30}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	opts := mujica.TimePickerOptions{Label: "Alarm", Hour12: true, MinuteStep: 5}
	if mujica.TimePicker(c, &alarm, opts).Changed() {
		fmt.Println("alarm", alarm)
	}
}
ui.Render(view, 400, 300, 1)`
	dtTimeRangeCode = `shift := mujica.TimeRange{Start: mujica.TimeOfDay{Hour: 22}, End: mujica.TimeOfDay{Hour: 6}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	opts := mujica.TimeRangePickerOptions{Label: "Night shift", AllowOvernight: true}
	if mujica.TimeRangePicker(c, &shift, opts).Changed() {
		fmt.Println("shift", shift.Start, shift.End, shift.Overnight())
	}
}
ui.Render(view, 400, 300, 1)`
	dtDateTimeCode = `tokyo, err := time.LoadLocation("Asia/Tokyo")
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
ui.Render(view, 400, 300, 1)`
	dtMonthCode = `var month mujica.YearMonth
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	opts := mujica.MonthPickerOptions{Label: "Billing month", Clearable: true}
	if mujica.MonthPicker(c, &month, opts).Changed() {
		fmt.Println("month", month)
	}
}
ui.Render(view, 400, 300, 1)`
	dtYearCode = `year := 2026
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	opts := mujica.YearPickerOptions{Label: "Vintage", Min: 1990, Max: 2030}
	if mujica.YearPicker(c, &year, opts).Changed() {
		fmt.Println("year", year)
	}
}
ui.Render(view, 400, 300, 1)`
)
