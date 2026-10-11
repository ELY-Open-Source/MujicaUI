package datetime

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var calendarToday = core.Date{Year: 2026, Month: time.October, Day: 5}

// datetimeTestLabel is the en-US spoken name of a day cell.
func datetimeTestLabel(d core.Date) string {
	return fmt.Sprintf("%v, %v %d, %d", d.Weekday(), d.Month, d.Day, d.Year)
}

func calendarTester(value *core.Date, opts CalendarOptions, changes *int) *ui.Tester {
	opts.Today = calendarToday
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		if Calendar(c, value, opts).Changed() {
			*changes++
		}
	}, 400, 400)
}

func TestCalendarClickChoosesDay(t *testing.T) {
	v, n := core.Date{}, 0
	tt := calendarTester(&v, CalendarOptions{}, &n)
	if !tt.HasText("October 2026") {
		t.Fatalf("empty value should show today's month, got %v", tt.Texts())
	}
	if err := tt.Click("Monday, October 12, 2026"); err != nil {
		t.Fatal(err)
	}
	if v != (core.Date{Year: 2026, Month: time.October, Day: 12}) || n != 1 {
		t.Fatalf("value %v changes %d", v, n)
	}
	if err := tt.Click("Monday, October 12, 2026"); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("same day reported change again: %d", n)
	}
}

func TestCalendarKeyboard(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 15}, 0
	tt := calendarTester(&v, CalendarOptions{}, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Thursday, October 15, 2026") {
		t.Fatal("Tab should land on the chosen day")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Friday, October 23, 2026") {
		t.Fatal("arrows should move the focus by day and week")
	}
	tt.Key(0, ui.KeyPageDown)
	if !tt.HasText("November 2026") || !tt.Focused("Monday, November 23, 2026") {
		t.Fatalf("PageDown should move a month: %v", tt.Texts()[:3])
	}
	tt.Key(ui.Shift, ui.KeyPageUp)
	if !tt.HasText("November 2025") {
		t.Fatal("Shift+PageUp should move a year")
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyEnter)
	if v != (core.Date{Year: 2025, Month: time.November, Day: 30}) || n != 1 {
		t.Fatalf("Enter chose %v (%d changes)", v, n)
	}
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeySpace)
	if v != (core.Date{Year: 2025, Month: time.November, Day: 1}) || n != 2 {
		t.Fatalf("Space chose %v", v)
	}
}

func TestCalendarLimits(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 10}, 0
	opts := CalendarOptions{
		Min:         core.Date{Year: 2026, Month: time.October, Day: 5},
		Max:         core.Date{Year: 2026, Month: time.October, Day: 20},
		DisableDate: func(d core.Date) bool { return d.Weekday() == time.Sunday },
	}
	tt := calendarTester(&v, opts, &n)
	for _, label := range []string{"Thursday, October 1, 2026", "Sunday, October 11, 2026", "Wednesday, October 21, 2026"} {
		if err := tt.Click(label); err != nil {
			t.Fatal(err)
		}
	}
	if v != (core.Date{Year: 2026, Month: time.October, Day: 10}) || n != 0 {
		t.Fatalf("disabled days changed the value to %v", v)
	}
	if err := tt.Click("Saturday, October 10, 2026"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Monday, October 12, 2026") {
		t.Fatal("arrows should skip the disabled Sunday")
	}
}

func TestCalendarWeekStart(t *testing.T) {
	for _, tc := range []struct {
		ws    WeekStart
		loc   core.Locale
		first string
	}{{WeekStartLocale, core.EnUS, "Sunday"}, {WeekStartLocale, core.ZhCN, "星期一"}, {WeekStartMonday, core.EnUS, "Monday"}, {WeekStartSaturday, core.EnUS, "Saturday"}} {
		v := core.Date{Year: 2026, Month: time.October, Day: 5}
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{Locale: tc.loc})
			Calendar(c, &v, CalendarOptions{WeekStart: tc.ws, Today: calendarToday})
		}, 400, 400)
		texts := tt.Texts()
		var first string
		for i, s := range texts {
			if s == "2026年10月" || s == "October 2026" {
				first = texts[i+3]
				break
			}
		}
		if first != tc.first {
			t.Errorf("%v %v: first column %q, want %q", tc.ws, tc.loc, first, tc.first)
		}
	}
}

func TestCalendarDisabled(t *testing.T) {
	v, n := core.Date{Year: 2026, Month: time.October, Day: 10}, 0
	tt := calendarTester(&v, CalendarOptions{Disabled: true}, &n)
	_ = tt.Click("Monday, October 12, 2026")
	_ = tt.Click("Next month")
	if v != (core.Date{Year: 2026, Month: time.October, Day: 10}) || n != 0 || !tt.HasText("October 2026") {
		t.Fatalf("disabled calendar changed: %v", v)
	}
}

func TestCalendarRejectsInvalidValue(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for February 30")
		}
	}()
	v := core.Date{Year: 2026, Month: time.February, Day: 30}
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Calendar(c, &v, CalendarOptions{})
	}, 400, 400)
}

func TestCalendarSnapshot(t *testing.T) {
	v := core.Date{Year: 2026, Month: time.October, Day: 14}
	snapshot(t, "calendar", 320, 340, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			Calendar(c, &v, CalendarOptions{Today: calendarToday, Max: core.Date{Year: 2026, Month: time.October, Day: 28},
				DisableDate: func(d core.Date) bool { return d.Day == 20 }})
		})
	})
}
