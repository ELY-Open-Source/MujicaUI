package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestTimePickerDisabledIgnoresSegmentKeys(t *testing.T) {
	v, n := TimeOfDay{9, 30, 0}, 0
	opts := TimePickerOptions{Label: "Alarm"}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		if TimePicker(c, &v, opts).Changed() {
			n++
		}
	}, 420, 300)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Alarm: Hours") {
		t.Fatal("setup: hours not focused")
	}
	opts.Disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.Key1)
	tt.Key(0, ui.Key5)
	if v != (TimeOfDay{9, 30, 0}) || n != 0 {
		t.Fatalf("disabled picker took keys: %v (%d)", v, n)
	}
}

func TestTimeRangePickerDisabledIgnoresSegmentKeys(t *testing.T) {
	v, n := TimeRange{TimeOfDay{9, 0, 0}, TimeOfDay{17, 0, 0}}, 0
	opts := TimeRangePickerOptions{}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		if TimeRangePicker(c, &v, opts).Changed() {
			n++
		}
	}, 520, 300)
	_ = tt.Click("End time: Hours")
	opts.Disabled = true
	tt.Frame()
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	if v != (TimeRange{TimeOfDay{9, 0, 0}, TimeOfDay{17, 0, 0}}) || n != 0 {
		t.Fatalf("disabled range took keys: %v (%d)", v, n)
	}
}

func TestDatePickerDisabledBlurDoesNotCommit(t *testing.T) {
	v, n := oct(5), 0
	opts := DatePickerOptions{}
	tt := datePickerLive(&v, &opts, &n)
	datePickerType(tt, "Due", "2026-10-20")
	opts.Disabled = true
	tt.Frame()
	_ = tt.Click("Open calendar")
	tt.Key(0, ui.KeyTab)
	tt.Frame()
	if v != oct(5) || n != 0 {
		t.Fatalf("disabled picker committed its draft on blur: %v (%d)", v, n)
	}
}

func TestTimePickerSegmentErrorSurvivesOtherSegments(t *testing.T) {
	v, n := TimeOfDay{9, 30, 0}, 0
	tt := timePickerTester(&v, TimePickerOptions{}, &n)
	_ = tt.Click("Alarm: Hours")
	tt.Key(0, ui.Key2)
	tt.Key(0, ui.Key9)
	_ = tt.Click("Alarm: Minutes")
	tt.Key(0, ui.KeyUp)
	if !tt.HasText("29") || !tt.HasText("Hours must be from 0 to 23.") {
		t.Fatalf("editing minutes cleared the hour error: %v", tt.Texts())
	}
	if v != (TimeOfDay{2, 30, 0}) {
		t.Fatalf("value written while a segment is invalid: %v", v)
	}
	_ = tt.Click("Alarm: Hours")
	tt.Key(0, ui.Key0)
	tt.Key(0, ui.Key8)
	if v != (TimeOfDay{8, 31, 0}) || tt.HasText("Hours must be") {
		t.Fatalf("fixing the hour should commit both: %v", v)
	}
}

func TestDateTimePickerBlocksConfirmWithInvalidSegment(t *testing.T) {
	v := time.Date(2026, time.October, 5, 9, 30, 0, 0, time.UTC)
	n, x := 0, 0
	tt := dateTimeTester(t, &v, DateTimePickerOptions{Location: time.UTC}, &n, &x)
	_ = tt.Click("Meeting")
	_ = tt.Click("Time: Hours")
	tt.Key(0, ui.Key2)
	tt.Key(0, ui.Key9)
	_ = tt.Click("Time: Minutes")
	tt.Key(0, ui.KeyUp)
	_ = tt.Click("Confirm")
	if n != 0 || !v.Equal(time.Date(2026, time.October, 5, 9, 30, 0, 0, time.UTC)) {
		t.Fatalf("confirmed with an invalid hour: %v", v)
	}
	if !tt.HasText("Hours must be from 0 to 23.") || !tt.HasText("Confirm") {
		t.Fatalf("expected the popup to stay with the error: %v", tt.Texts())
	}
}

func TestDateRangePickerConfirmChecksCurrentBounds(t *testing.T) {
	v, n := DateRange{}, 0
	opts := DateRangePickerOptions{}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		o := opts
		o.Today, o.Label = calendarToday, "Stay"
		if DateRangePicker(c, &v, o).Changed() {
			n++
		}
	}, 420, 520)
	_ = tt.Click("Stay")
	_ = tt.Click(datetimeTestLabel(oct(10)))
	_ = tt.Click(datetimeTestLabel(oct(15)))
	opts.Max = oct(12)
	tt.Frame()
	_ = tt.Click("Confirm")
	if n != 0 || !v.IsZero() || !tt.HasText("Choose 2026-10-12 or earlier.") {
		t.Fatalf("out-of-range draft committed: %v %v", v, tt.Texts())
	}
	opts.Max = Date{}
	opts.DisableDate = func(d Date) bool { return d == oct(11) }
	tt.Frame()
	_ = tt.Click("Confirm")
	if n != 0 || !tt.HasText("The range cannot include 2026-10-11.") {
		t.Fatalf("draft over a newly disabled day committed: %v", v)
	}
	opts.DisableDate = nil
	tt.Frame()
	_ = tt.Click("Confirm")
	if n != 1 || v != (DateRange{oct(10), oct(15)}) {
		t.Fatalf("valid draft not committed: %v", v)
	}
}
