package mujica

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func durationTester(v *time.Duration, opts *DurationPickerOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if DurationPicker(c, v, *opts).Changed() {
			*n++
		}
	}, 640, 200)
}

func TestDurationPickerEditsFields(t *testing.T) {
	v, n := 90*time.Minute, 0
	opts := &DurationPickerOptions{Seconds: true}
	tt := durationTester(&v, opts, &n)
	_ = tt.Click("Hours")
	tt.Key(0, ui.KeyUp)
	_ = tt.Click("Seconds")
	tt.Command("selectAll")
	tt.Type("45")
	tt.Key(0, ui.KeyEnter)
	if v != 2*time.Hour+30*time.Minute+45*time.Second || n != 2 {
		t.Fatalf("value %v (%d)", v, n)
	}
}

func TestDurationPickerRangeReported(t *testing.T) {
	v, n := 30*time.Minute, 0
	opts := &DurationPickerOptions{Min: 15 * time.Minute, Max: 2 * time.Hour}
	tt := durationTester(&v, opts, &n)
	_ = tt.Click("Hours")
	tt.Key(0, ui.KeyUp)
	if v != 90*time.Minute || n != 1 {
		t.Fatalf("1 h 30 min should fit the range: %v", v)
	}
	tt.Key(0, ui.KeyUp)
	if v != 90*time.Minute || !tt.HasText("Choose a duration from 15 min to 2 h.") {
		t.Fatalf("over Max accepted: %v %v", v, tt.Texts())
	}
	_ = tt.Click("Minutes")
	tt.Command("selectAll")
	tt.Type("0")
	tt.Key(0, ui.KeyEnter)
	if v != 2*time.Hour || n != 2 || tt.HasText("Choose a duration") {
		t.Fatalf("fixing the minutes: %v", v)
	}
	if tt.HasText("Seconds") {
		t.Fatal("seconds field without Seconds")
	}
}

func TestDurationPickerDisabledAndInvalid(t *testing.T) {
	v, n := time.Hour, 0
	tt := durationTester(&v, &DurationPickerOptions{Disabled: true}, &n)
	_ = tt.Click("Hours")
	tt.Key(0, ui.KeyUp)
	if v != time.Hour || n != 0 {
		t.Fatal("disabled picker changed")
	}
	for _, c := range []struct {
		v    time.Duration
		opts DurationPickerOptions
	}{{30 * time.Second, DurationPickerOptions{}}, {0, DurationPickerOptions{Min: time.Hour, Max: time.Minute}}, {-time.Minute, DurationPickerOptions{}}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%v %+v: expected panic", c.v, c.opts)
				}
			}()
			durationTester(&c.v, &c.opts, &n)
		}()
	}
}

func TestDurationPickerSnapshot(t *testing.T) {
	v := 2*time.Hour + 15*time.Minute + 30*time.Second
	snapshot(t, "duration_picker", 640, 80, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { DurationPicker(c, &v, DurationPickerOptions{Seconds: true}) })
	})
}
