package datetime

import (
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func dateTimeTester(t *testing.T, value *time.Time, opts DateTimePickerOptions, changes, cancels *int) *ui.Tester {
	t.Helper()
	opts.Today = calendarToday
	opts.Label = "Meeting"
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		ui.Column(c).Padding(16).Children(func() {
			r := DateTimePicker(c, value, opts)
			if r.Changed() {
				*changes++
			}
			if r.Canceled() {
				*cancels++
			}
		})
	}, 480, 640)
}

func newYork(t *testing.T) *time.Location {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	return loc
}

func TestDateTimeResolve(t *testing.T) {
	ny := newYork(t)
	if _, n := dateTimeResolve(core.Date{Year: 2026, Month: time.March, Day: 8}, TimeOfDay{2, 30, 0}, ny); n != 0 {
		t.Fatalf("spring-forward gap gave %d instants", n)
	}
	if _, n := dateTimeResolve(core.Date{Year: 2026, Month: time.November, Day: 1}, TimeOfDay{1, 30, 0}, ny); n != 2 {
		t.Fatalf("fall-back overlap gave %d instants", n)
	}
	at, n := dateTimeResolve(core.Date{Year: 2026, Month: time.July, Day: 4}, TimeOfDay{12, 0, 0}, ny)
	if n != 1 || !at.Equal(time.Date(2026, time.July, 4, 16, 0, 0, 0, time.UTC)) {
		t.Fatalf("plain time %v (%d)", at, n)
	}
}

func TestDateTimePickerConfirmInLocation(t *testing.T) {
	tokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	v := time.Date(2026, time.October, 5, 23, 30, 0, 0, time.UTC) // 08:30 on Oct 6 in Tokyo
	n, x := 0, 0
	tt := dateTimeTester(t, &v, DateTimePickerOptions{Location: tokyo}, &n, &x)
	if !tt.HasText("2026-10-06 08:30 JST") {
		t.Fatalf("value should show in the caller's zone: %v", tt.Texts())
	}
	_ = tt.Click("Meeting")
	_ = tt.Click(datetimeTestLabel(oct(9)))
	_ = tt.Click("Time: Hours")
	tt.Key(0, ui.KeyUp)
	if n != 0 {
		t.Fatal("the draft should wait for Confirm")
	}
	_ = tt.Click("Confirm")
	want := time.Date(2026, time.October, 9, 9, 30, 0, 0, tokyo)
	if !v.Equal(want) || n != 1 || x != 0 || v.Location() != tokyo {
		t.Fatalf("confirmed %v, want %v", v, want)
	}
	_ = tt.Click("Meeting")
	_ = tt.Click(datetimeTestLabel(oct(20)))
	_ = tt.Click("Cancel")
	if !v.Equal(want) || x != 1 {
		t.Fatalf("cancel changed %v", v)
	}
}

func TestDateTimePickerReportsDSTGap(t *testing.T) {
	ny := newYork(t)
	v := time.Date(2026, time.March, 7, 2, 30, 0, 0, ny)
	n, x := 0, 0
	tt := dateTimeTester(t, &v, DateTimePickerOptions{Location: ny, Today: core.Date{Year: 2026, Month: time.March, Day: 1}}, &n, &x)
	_ = tt.Click("Meeting")
	_ = tt.Click("Sunday, March 8, 2026")
	_ = tt.Click("Confirm")
	if n != 0 || !v.Equal(time.Date(2026, time.March, 7, 2, 30, 0, 0, ny)) {
		t.Fatalf("gap time was shifted into %v", v)
	}
	if !tt.HasText("2026-03-08 02:30 does not exist in America/New_York: the clocks skip it.") {
		t.Fatalf("missing gap error: %v", tt.Texts())
	}
	_ = tt.Click("Time: Hours")
	tt.Key(0, ui.KeyUp)
	_ = tt.Click("Confirm")
	if n != 1 || !v.Equal(time.Date(2026, time.March, 8, 3, 30, 0, 0, ny)) {
		t.Fatalf("after fixing: %v", v)
	}
}

func TestDateTimePickerReportsAmbiguous(t *testing.T) {
	ny := newYork(t)
	v := time.Date(2026, time.October, 31, 1, 30, 0, 0, ny)
	n, x := 0, 0
	tt := dateTimeTester(t, &v, DateTimePickerOptions{Location: ny}, &n, &x)
	_ = tt.Click("Meeting")
	_ = tt.Click("Sunday, November 1, 2026")
	_ = tt.Click("Confirm")
	if n != 0 || !tt.HasText("2026-11-01 01:30 happens twice in America/New_York: the clocks repeat it.") {
		t.Fatalf("ambiguous time accepted: %v %v", v, tt.Texts())
	}
}

func TestDateTimePickerRange(t *testing.T) {
	v := time.Time{}
	n, x := 0, 0
	maxT := time.Date(2026, time.October, 10, 12, 0, 0, 0, time.UTC)
	tt := dateTimeTester(t, &v, DateTimePickerOptions{Location: time.UTC, Max: maxT}, &n, &x)
	_ = tt.Click("Meeting")
	_ = tt.Click(datetimeTestLabel(oct(11)))
	_ = tt.Click(datetimeTestLabel(oct(10)))
	_ = tt.Click("Time: Hours")
	tt.Key(0, ui.Key1)
	tt.Key(0, ui.Key3)
	_ = tt.Click("Confirm")
	if n != 0 || !v.IsZero() || !tt.HasText("Choose 2026-10-10 12:00 UTC or earlier.") {
		t.Fatalf("past Max accepted: %v %v", v, tt.Texts())
	}
}

func TestDateTimePickerNeedsLocation(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic without Location")
		}
	}()
	v := time.Time{}
	dateTimeTester(t, &v, DateTimePickerOptions{}, new(int), new(int))
}

func TestDateTimePickerSnapshot(t *testing.T) {
	a, b := time.Date(2026, time.October, 14, 18, 45, 0, 0, time.UTC), time.Time{}
	snapshot(t, "date_time_picker", 720, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			DateTimePicker(c, &b, DateTimePickerOptions{Location: time.UTC, Today: calendarToday})
			DateTimePicker(c, &a, DateTimePickerOptions{Location: time.UTC, Today: calendarToday, Hour12: true})
		})
	})
}
