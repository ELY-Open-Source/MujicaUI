package datetime

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func weekPickerTester(v *core.Date, opts *WeekPickerOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		o := *opts
		o.Today, o.Label = calendarToday, "Sprint"
		ui.Column(c).Padding(16).Children(func() {
			if WeekPicker(c, v, o).Changed() {
				*n++
			}
		})
	}, 480, 480)
}

func TestWeekPickerPicksWeekStart(t *testing.T) {
	v, n := core.Date{}, 0
	opts := &WeekPickerOptions{}
	tt := weekPickerTester(&v, opts, &n)
	_ = tt.Click("Sprint")
	_ = tt.Click(datetimeTestLabel(oct(15)))
	if v != oct(11) || n != 1 || !tt.HasText("October 11 – October 17, 2026") {
		t.Fatalf("Sunday-start week of Oct 15: %v %v", v, tt.Texts())
	}
	v, opts.WeekStart = oct(12), WeekStartMonday
	tt.Frame()
	_ = tt.Click("Sprint")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if v != oct(19) || n != 2 {
		t.Fatalf("Monday-start keyboard pick %v", v)
	}
}

func TestWeekPickerBoundsClearDisabled(t *testing.T) {
	v, n := oct(11), 0
	opts := &WeekPickerOptions{Min: oct(4), Max: oct(24), Clearable: true}
	tt := weekPickerTester(&v, opts, &n)
	_ = tt.Click("Sprint")
	_ = tt.Click(datetimeTestLabel(oct(26)))
	_ = tt.Click(datetimeTestLabel(oct(2)))
	if v != oct(11) || n != 0 {
		t.Fatalf("out-of-range week chosen: %v", v)
	}
	tt.Key(0, ui.KeyEscape)
	_ = tt.Click("Clear")
	if !v.IsZero() || n != 1 {
		t.Fatalf("clear left %v", v)
	}
	opts.Disabled = true
	tt.Frame()
	_ = tt.Click("Sprint")
	if tt.HasText("October 2026") {
		t.Fatal("disabled picker opened")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for a mid-week value")
		}
	}()
	mid := oct(14)
	weekPickerTester(&mid, &WeekPickerOptions{}, &n)
}

func TestWeekPickerSnapshot(t *testing.T) {
	v := oct(11)
	snapshot(t, "week_picker", 400, 80, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { WeekPicker(c, &v, WeekPickerOptions{Today: calendarToday, Clearable: true}) })
	})
}
