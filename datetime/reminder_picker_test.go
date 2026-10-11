package datetime

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func reminderTester(v *[]time.Duration, opts ReminderPickerOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if ReminderPicker(c, v, opts).Changed() {
			*n++
		}
	}, 640, 200)
}

func TestReminderPickerTogglesInChoiceOrder(t *testing.T) {
	v, n := []time.Duration{}, 0
	tt := reminderTester(&v, ReminderPickerOptions{}, &n)
	_ = tt.Click("1 day before")
	_ = tt.Click("15 minutes before")
	_ = tt.Click("At start")
	if !slices.Equal(v, []time.Duration{0, 15 * time.Minute, 24 * time.Hour}) || n != 3 {
		t.Fatalf("value %v (%d)", v, n)
	}
	_ = tt.Click("15 minutes before")
	if !slices.Equal(v, []time.Duration{0, 24 * time.Hour}) {
		t.Fatalf("untoggle left %v", v)
	}
}

func TestReminderPickerMaxKeyboardDisabled(t *testing.T) {
	v, n := []time.Duration{time.Hour}, 0
	tt := reminderTester(&v, ReminderPickerOptions{Max: 2}, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	_ = tt.Click("1 week before")
	_ = tt.Click("5 minutes before")
	if !slices.Equal(v, []time.Duration{0, time.Hour}) {
		t.Fatalf("Max ignored: %v", v)
	}
	d := []time.Duration{}
	tt = reminderTester(&d, ReminderPickerOptions{Disabled: true}, &n)
	_ = tt.Click("At start")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if len(d) != 0 {
		t.Fatal("disabled picker changed")
	}
}

func TestReminderTextAndInvalid(t *testing.T) {
	var got []string
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{Locale: core.ZhCN})
		got = []string{ReminderText(c, time.Minute), ReminderText(c, 2*time.Hour), ReminderText(c, 14*24*time.Hour)}
	}, 10, 10)
	if !slices.Equal(got, []string{"提前 1 分钟", "提前 2 小时", "提前 2 周"}) {
		t.Fatalf("texts %q", got)
	}
	for _, opts := range []ReminderPickerOptions{{Choices: []time.Duration{30 * time.Second}}, {Choices: []time.Duration{0, 0}}, {Max: -1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: expected panic", opts)
				}
			}()
			v := []time.Duration{}
			reminderTester(&v, opts, new(int))
		}()
	}
}

func TestReminderPickerSnapshot(t *testing.T) {
	v := []time.Duration{15 * time.Minute, 24 * time.Hour}
	snapshot(t, "reminder_picker", 300, 330, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { ReminderPicker(c, &v, ReminderPickerOptions{}) })
	})
}
