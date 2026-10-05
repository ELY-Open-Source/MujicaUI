package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func timePickerTester(value *TimeOfDay, opts TimePickerOptions, changes *int) *ui.Tester {
	opts.Label = "Alarm"
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		ui.Column(c).Padding(16).Children(func() {
			if TimePicker(c, value, opts).Changed() {
				*changes++
			}
		})
	}, 420, 420)
}

func TestTimePickerArrowsAndTyping(t *testing.T) {
	v, n := TimeOfDay{9, 30, 0}, 0
	tt := timePickerTester(&v, TimePickerOptions{}, &n)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Alarm: Hours") {
		t.Fatal("Tab should focus the hours")
	}
	tt.Key(0, ui.KeyUp)
	if v != (TimeOfDay{10, 30, 0}) || n != 1 {
		t.Fatalf("Up gave %v", v)
	}
	tt.Key(0, ui.Key1)
	tt.Key(0, ui.Key7)
	if v != (TimeOfDay{17, 30, 0}) || !tt.Focused("Alarm: Minutes") {
		t.Fatalf("typed hours %v, focus moved %v", v, tt.Focused("Alarm: Minutes"))
	}
	tt.Key(0, ui.Key0)
	tt.Key(0, ui.Key5)
	if v != (TimeOfDay{17, 5, 0}) {
		t.Fatalf("typed minutes %v", v)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyDown)
	if v != (TimeOfDay{16, 5, 0}) {
		t.Fatalf("Left then Down gave %v", v)
	}
	tt.Key(0, ui.Key0)
	tt.Key(0, ui.KeyDown)
	if v != (TimeOfDay{23, 5, 0}) {
		t.Fatalf("Down should wrap round the clock: %v", v)
	}
}

func TestTimePickerHour12AndSeconds(t *testing.T) {
	v, n := TimeOfDay{9, 0, 0}, 0
	tt := timePickerTester(&v, TimePickerOptions{Hour12: true, Seconds: true}, &n)
	if !tt.HasText("AM") || !tt.HasText(" 9") {
		t.Fatalf("12-hour display: %v", tt.Texts())
	}
	_ = tt.Click("Alarm: AM/PM")
	tt.Key(0, ui.KeyP)
	if v != (TimeOfDay{21, 0, 0}) {
		t.Fatalf("P gave %v", v)
	}
	_ = tt.Click("Alarm: Hours")
	tt.Key(0, ui.Key1)
	tt.Key(0, ui.Key2)
	if v != (TimeOfDay{12, 0, 0}) {
		t.Fatalf("12 PM is noon, got %v", v)
	}
	_ = tt.Click("Alarm: Seconds")
	tt.Key(0, ui.Key4)
	tt.Key(0, ui.Key5)
	if v != (TimeOfDay{12, 0, 45}) {
		t.Fatalf("seconds %v", v)
	}
}

func TestTimePickerRangeAndStep(t *testing.T) {
	v, n := TimeOfDay{9, 0, 0}, 0
	opts := TimePickerOptions{Min: TimeOfDay{9, 0, 0}, Max: TimeOfDay{17, 0, 0}, MinuteStep: 15}
	tt := timePickerTester(&v, opts, &n)
	_ = tt.Click("Alarm: Hours")
	tt.Key(0, ui.KeyDown)
	if v != (TimeOfDay{9, 0, 0}) || n != 0 || !tt.HasText("Choose a time from 09:00 to 17:00.") {
		t.Fatalf("out of range accepted: %v %v", v, tt.Texts())
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyUp)
	if v != (TimeOfDay{10, 0, 0}) || tt.HasText("Choose a time") {
		t.Fatalf("back in range: %v", v)
	}
	_ = tt.Click("Alarm: Minutes")
	tt.Key(0, ui.KeyUp)
	if v != (TimeOfDay{10, 15, 0}) {
		t.Fatalf("step 15 gave %v", v)
	}
	tt.Key(0, ui.Key0)
	tt.Key(0, ui.Key7)
	if v != (TimeOfDay{10, 0, 0}) || !tt.HasText("Minutes must be a multiple of 15.") {
		t.Fatalf("off-step minute accepted: %v", v)
	}
}

func TestTimePickerPopup(t *testing.T) {
	v, n := TimeOfDay{9, 30, 0}, 0
	opts := TimePickerOptions{Min: TimeOfDay{9, 0, 0}, MinuteStep: 5}
	tt := timePickerTester(&v, opts, &n)
	if err := tt.Click("Choose time"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("25"); err != nil {
		t.Fatal(err)
	}
	if v != (TimeOfDay{9, 25, 0}) || n != 1 {
		t.Fatalf("minute pick %v", v)
	}
	if err := tt.Click("08"); err != nil {
		t.Fatal(err)
	}
	if v.Hour != 9 || n != 1 {
		t.Fatalf("hour before Min was chosen: %v", v)
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("05") {
		t.Fatal("Escape should close the popup")
	}
}

func TestTimePickerRejectsBadOptions(t *testing.T) {
	for _, opts := range []TimePickerOptions{{MinuteStep: 7}, {Min: TimeOfDay{10, 0, 0}, Max: TimeOfDay{9, 0, 0}}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: expected panic", opts)
				}
			}()
			v := TimeOfDay{}
			timePickerTester(&v, opts, new(int))
		}()
	}
}

func TestTimePickerSnapshot(t *testing.T) {
	a, b := TimeOfDay{9, 30, 0}, TimeOfDay{21, 5, 40}
	snapshot(t, "time_picker", 520, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			TimePicker(c, &a, TimePickerOptions{})
			TimePicker(c, &b, TimePickerOptions{Hour12: true, Seconds: true})
			TimePicker(c, &a, TimePickerOptions{Disabled: true})
		})
	})
}
