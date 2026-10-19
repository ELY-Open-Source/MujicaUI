package datetime

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func dateRangeTester(value *DateRange, opts DateRangePickerOptions, changes, cancels *int) *ui.Tester {
	opts.Today = calendarToday
	opts.Label = "Stay"
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		ui.Column(c).Padding(16).Children(func() {
			r := DateRangePicker(c, value, opts)
			if r.Changed() {
				*changes++
			}
			if r.Canceled() {
				*cancels++
			}
		})
	}, 420, 520)
}

func oct(day int) core.Date { return core.Date{Year: 2026, Month: time.October, Day: day} }

func TestDateRangePickerPickAndConfirm(t *testing.T) {
	v, n, x := DateRange{}, 0, 0
	tt := dateRangeTester(&v, DateRangePickerOptions{}, &n, &x)
	if err := tt.Click("Stay"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Choose the start date.") {
		t.Fatalf("popup missing: %v", tt.Texts())
	}
	_ = tt.Click(datetimeTestLabel(oct(12)))
	if !tt.HasText("From 2026-10-12, choose the end date.") {
		t.Fatal("start not taken")
	}
	_ = tt.Click(datetimeTestLabel(oct(8)))
	if !tt.HasText("From 2026-10-08, choose the end date.") {
		t.Fatal("a day before the start should restart the range")
	}
	_ = tt.Click(datetimeTestLabel(oct(14)))
	if !tt.HasText("2026-10-08 – 2026-10-14 · 7 days") || n != 0 {
		t.Fatalf("draft should wait for Confirm: %v %d", v, n)
	}
	if err := tt.Click("Confirm"); err != nil {
		t.Fatal(err)
	}
	if v != (DateRange{oct(8), oct(14)}) || n != 1 || x != 0 || tt.HasText("Confirm") {
		t.Fatalf("confirm gave %v, n=%d x=%d", v, n, x)
	}
}

func TestDateRangePickerCancelKeepsValue(t *testing.T) {
	v, n, x := DateRange{oct(1), oct(3)}, 0, 0
	tt := dateRangeTester(&v, DateRangePickerOptions{}, &n, &x)
	_ = tt.Click("Stay")
	_ = tt.Click(datetimeTestLabel(oct(20)))
	tt.Key(0, ui.KeyEscape)
	if v != (DateRange{oct(1), oct(3)}) || n != 0 || x != 1 {
		t.Fatalf("Escape: %v n=%d x=%d", v, n, x)
	}
	_ = tt.Click("Stay")
	_ = tt.Click("Cancel")
	if x != 2 || tt.HasText("Confirm") {
		t.Fatalf("Cancel: x=%d", x)
	}
	_ = tt.Click("Stay")
	_ = tt.Click("Clear")
	if !v.IsZero() || n != 1 || x != 2 {
		t.Fatalf("Clear: %v n=%d x=%d", v, n, x)
	}
}

func TestDateRangePickerLimits(t *testing.T) {
	v, n, x := DateRange{}, 0, 0
	opts := DateRangePickerOptions{Min: oct(5), DisableDate: func(d core.Date) bool { return d == oct(16) }}
	tt := dateRangeTester(&v, opts, &n, &x)
	_ = tt.Click("Stay")
	_ = tt.Click(datetimeTestLabel(oct(2)))
	if !tt.HasText("Choose the start date.") {
		t.Fatal("a day before Min was taken")
	}
	_ = tt.Click(datetimeTestLabel(oct(14)))
	_ = tt.Click(datetimeTestLabel(oct(18)))
	if !tt.HasText("The range cannot include 2026-10-16.") {
		t.Fatalf("span over a disabled day accepted: %v", tt.Texts())
	}
	_ = tt.Click("Confirm")
	if n != 0 || !v.IsZero() {
		t.Fatal("confirm without an end changed the value")
	}
	_ = tt.Click(datetimeTestLabel(oct(15)))
	_ = tt.Click("Confirm")
	if v != (DateRange{oct(14), oct(15)}) {
		t.Fatalf("got %v", v)
	}
}

func TestDateRangePickerKeyboard(t *testing.T) {
	v, n, x := DateRange{}, 0, 0
	tt := dateRangeTester(&v, DateRangePickerOptions{}, &n, &x)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused(datetimeTestLabel(calendarToday)) {
		t.Fatal("the popup should focus today")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyDown)
	if !tt.HasText("From 2026-10-05, choose the end date.") {
		t.Fatal("Enter should pick the start")
	}
	tt.Key(0, ui.KeySpace)
	if !tt.HasText("2026-10-05 – 2026-10-12 · 8 days") {
		t.Fatalf("Space should pick the end: %v", tt.Texts())
	}
}

func TestDateRangeDays(t *testing.T) {
	r := DateRange{core.Date{Year: 2026, Month: time.March, Day: 28}, core.Date{Year: 2026, Month: time.April, Day: 2}}
	if r.Days() != 6 || (DateRange{}).Days() != 0 {
		t.Fatalf("days %d", r.Days())
	}
}

func TestDateRangePickerSnapshot(t *testing.T) {
	a, b := DateRange{oct(8), oct(14)}, DateRange{}
	snapshot(t, "date_range_picker", 640, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			DateRangePicker(c, &b, DateRangePickerOptions{Today: calendarToday})
			DateRangePicker(c, &a, DateRangePickerOptions{Today: calendarToday})
		})
	})
}
