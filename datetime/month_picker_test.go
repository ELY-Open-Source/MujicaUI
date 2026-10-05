package datetime

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func monthPickerTester(value *YearMonth, opts MonthPickerOptions, changes *int) *ui.Tester {
	opts.Today = calendarToday
	opts.Label = "Billing month"
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		ui.Column(c).Padding(16).Children(func() {
			if MonthPicker(c, value, opts).Changed() {
				*changes++
			}
		})
	}, 420, 420)
}

func TestMonthPickerClickAndYear(t *testing.T) {
	v, n := YearMonth{}, 0
	tt := monthPickerTester(&v, MonthPickerOptions{}, &n)
	_ = tt.Click("Billing month")
	if !tt.HasText("2026") {
		t.Fatalf("popup should open on the current year: %v", tt.Texts())
	}
	_ = tt.Click("Next year")
	if err := tt.Click("March 2027"); err != nil {
		t.Fatal(err)
	}
	if v != (YearMonth{2027, time.March}) || n != 1 || tt.HasText("April 2027") {
		t.Fatalf("picked %v, popup closed %v", v, !tt.HasText("April 2027"))
	}
}

func TestMonthPickerKeyboard(t *testing.T) {
	v, n := YearMonth{2026, time.May}, 0
	tt := monthPickerTester(&v, MonthPickerOptions{}, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("May 2026") {
		t.Fatal("popup should focus the chosen month")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyPageUp)
	if !tt.Focused("September 2025") {
		t.Fatal("arrows and PageUp should move the cursor")
	}
	tt.Key(0, ui.KeyEnter)
	if v != (YearMonth{2025, time.September}) || n != 1 {
		t.Fatalf("Enter chose %v", v)
	}
}

func TestMonthPickerLimitsAndClear(t *testing.T) {
	v, n := YearMonth{2026, time.June}, 0
	opts := MonthPickerOptions{Min: YearMonth{2026, time.March}, Max: YearMonth{2026, time.October}, Clearable: true,
		DisableMonth: func(m YearMonth) bool { return m.Month == time.August }}
	tt := monthPickerTester(&v, opts, &n)
	_ = tt.Click("Billing month")
	for _, m := range []string{"February 2026", "August 2026", "November 2026"} {
		_ = tt.Click(m)
	}
	if v != (YearMonth{2026, time.June}) || n != 0 {
		t.Fatalf("disabled month chosen: %v", v)
	}
	tt.Key(0, ui.KeyEscape)
	if err := tt.Click("Clear"); err != nil {
		t.Fatal(err)
	}
	if !v.IsZero() || n != 1 {
		t.Fatalf("clear left %v", v)
	}
}

func TestYearMonthOrdinal(t *testing.T) {
	for _, m := range []YearMonth{{2026, time.January}, {2026, time.December}, {1, time.March}} {
		if got := yearMonthOf(m.ordinal()); got != m {
			t.Fatalf("%v round-tripped to %v", m, got)
		}
	}
}

func TestMonthPickerSnapshot(t *testing.T) {
	a, b := YearMonth{2026, time.October}, YearMonth{}
	snapshot(t, "month_picker", 520, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			MonthPicker(c, &b, MonthPickerOptions{Today: calendarToday})
			MonthPicker(c, &a, MonthPickerOptions{Today: calendarToday, Clearable: true})
		})
	})
}
