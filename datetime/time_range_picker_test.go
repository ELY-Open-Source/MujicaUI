package datetime

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func timeRangeTester(value *TimeRange, opts TimeRangePickerOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Mode: core.Light})
		ui.Column(c).Padding(16).Children(func() {
			if TimeRangePicker(c, value, opts).Changed() {
				*changes++
			}
		})
	}, 520, 420)
}

func TestTimeRangePickerOrder(t *testing.T) {
	v, n := TimeRange{TimeOfDay{9, 0, 0}, TimeOfDay{17, 0, 0}}, 0
	tt := timeRangeTester(&v, TimeRangePickerOptions{}, &n)
	_ = tt.Click("End time: Hours")
	tt.Key(0, ui.Key0)
	tt.Key(0, ui.Key8)
	if v.End != (TimeOfDay{17, 0, 0}) || n != 0 {
		t.Fatalf("end before start accepted: %v", v)
	}
	if !tt.HasText("End must be after start; overnight ranges are off.") {
		t.Fatalf("missing order error: %v", tt.Texts())
	}
	_ = tt.Click("Start time: Hours")
	tt.Key(0, ui.Key0)
	tt.Key(0, ui.Key7)
	if v != (TimeRange{TimeOfDay{7, 0, 0}, TimeOfDay{8, 0, 0}}) || n != 2 || tt.HasText("End must be") {
		t.Fatalf("fixing the start should commit both: %v", v)
	}
	_ = tt.Click("End time: Hours")
	tt.Key(0, ui.KeyDown)
	if v.End != (TimeOfDay{8, 0, 0}) || !tt.HasText("Start and end must differ.") {
		t.Fatalf("equal ends accepted: %v", v)
	}
}

func TestTimeRangePickerOvernight(t *testing.T) {
	v, n := TimeRange{TimeOfDay{22, 0, 0}, TimeOfDay{23, 0, 0}}, 0
	tt := timeRangeTester(&v, TimeRangePickerOptions{AllowOvernight: true}, &n)
	_ = tt.Click("End time: Hours")
	tt.Key(0, ui.Key0)
	tt.Key(0, ui.Key6)
	if v.End != (TimeOfDay{6, 0, 0}) || n != 2 || !v.Overnight() || !tt.HasText("+1 day") {
		t.Fatalf("overnight not accepted: %v %v", v, tt.Texts())
	}
}

func TestTimeRangePickerRejectsOvernightValue(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	v := TimeRange{TimeOfDay{22, 0, 0}, TimeOfDay{6, 0, 0}}
	timeRangeTester(&v, TimeRangePickerOptions{}, new(int))
}

func TestTimeRangePickerSnapshot(t *testing.T) {
	a := TimeRange{TimeOfDay{9, 0, 0}, TimeOfDay{17, 30, 0}}
	b := TimeRange{TimeOfDay{22, 0, 0}, TimeOfDay{6, 0, 0}}
	snapshot(t, "time_range_picker", 420, 120, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			TimeRangePicker(c, &a, TimeRangePickerOptions{})
			TimeRangePicker(c, &b, TimeRangePickerOptions{AllowOvernight: true})
		})
	})
}
