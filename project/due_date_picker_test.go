package project

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func dueTester(v *core.Date, opts *DueDatePickerOptions, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		o := *opts
		o.Today = calendarToday
		ui.Column(c).Padding(16).Children(func() {
			if DueDatePicker(c, v, o).Changed() {
				*n++
			}
		})
	}, 480, 520)
}

func TestDueDatePickerShortcuts(t *testing.T) {
	v, n := core.Date{}, 0
	tt := dueTester(&v, &DueDatePickerOptions{}, &n)
	_ = tt.Click("Tomorrow")
	if v != oct(6) || n != 1 || !tt.HasText("Due tomorrow") {
		t.Fatalf("tomorrow %v", v)
	}
	_ = tt.Click("Next week")
	if v != oct(11) {
		t.Fatalf("next week (Sunday start) %v", v)
	}
	_ = tt.Click("Today")
	if v != oct(5) || !tt.HasText("Due today") {
		t.Fatalf("today %v", v)
	}
	_ = tt.Click("Open calendar")
	_ = tt.Click(datetimeTestLabel(oct(20)))
	if v != oct(20) || n != 4 {
		t.Fatalf("calendar pick %v (%d)", v, n)
	}
}

func TestDueDatePickerOverdueAndDone(t *testing.T) {
	v, n := oct(2), 0
	opts := &DueDatePickerOptions{}
	tt := dueTester(&v, opts, &n)
	if !tt.HasText("Overdue by 3 days") {
		t.Fatalf("overdue missing: %v", tt.Texts())
	}
	opts.Done = true
	tt.Frame()
	if tt.HasText("Overdue") {
		t.Fatal("a done task is not overdue")
	}
	opts.Disabled = true
	tt.Frame()
	_ = tt.Click("Today")
	if v != oct(2) || n != 0 {
		t.Fatal("disabled picker changed")
	}
}

func TestDueDatePickerSnapshot(t *testing.T) {
	v := oct(2)
	snapshot(t, "due_date_picker", 360, 140, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { DueDatePicker(c, &v, DueDatePickerOptions{Today: calendarToday}) })
	})
}

func TestDueDatePickerSelectedShortcutMark(t *testing.T) {
	v, n := oct(5), 0
	tt := dueTester(&v, &DueDatePickerOptions{}, &n)
	if !tt.HasText("✓ Today") {
		t.Fatal("selected shortcut needs a non-color mark")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if v != oct(6) || n != 1 {
		t.Fatalf("keyboard shortcut date %v, changes %d", v, n)
	}
}
