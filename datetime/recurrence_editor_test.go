package datetime

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

type recurrenceProbe struct {
	rule     Recurrence
	opts     RecurrenceEditorOptions
	changes  int
	problem  string
	summary  string
	disabled bool
}

func (p *recurrenceProbe) tester() *ui.Tester {
	p.opts.Start = oct(5)
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		o := p.opts
		o.Disabled = p.disabled
		r := RecurrenceEditor(c, &p.rule, o)
		if r.Changed() {
			p.changes++
		}
		p.problem = r.Problem()
		p.summary = RecurrenceSummary(c, p.rule)
	}, 720, 520)
}

func recurrencePick(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	if err := tt.Click("Repeat"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click(label); err != nil {
		t.Fatal(err)
	}
}

func TestRecurrenceEditorPresets(t *testing.T) {
	p := &recurrenceProbe{}
	tt := p.tester()
	recurrencePick(t, tt, "Weekly")
	if p.rule.Unit != RepeatWeekly || p.rule.Interval != 1 || !slices.Equal(p.rule.Weekdays, []time.Weekday{time.Monday}) || p.changes != 1 {
		t.Fatalf("weekly preset %+v", p.rule)
	}
	_ = tt.Click("Wednesday")
	if p.summary != "Weekly on Monday, Wednesday" {
		t.Fatalf("summary %q", p.summary)
	}
	for _, c := range []struct {
		label string
		unit  RepeatUnit
	}{{"Daily", RepeatDaily}, {"Monthly", RepeatMonthly}, {"Yearly", RepeatYearly}} {
		recurrencePick(t, tt, c.label)
		if p.rule.Unit != c.unit || p.rule.Interval != 1 {
			t.Fatalf("%s gave %+v", c.label, p.rule)
		}
	}
	recurrencePick(t, tt, "Does not repeat")
	if p.rule.Unit != RepeatNone || p.summary != "Does not repeat" {
		t.Fatalf("none gave %+v", p.rule)
	}
}

func TestRecurrenceEditorCustomAndEnds(t *testing.T) {
	p := &recurrenceProbe{rule: Recurrence{Unit: RepeatDaily, Interval: 1}}
	tt := p.tester()
	recurrencePick(t, tt, "Custom")
	if p.rule.Interval != 2 || !tt.HasText("Every") {
		t.Fatalf("custom %+v", p.rule)
	}
	_ = tt.Click("Interval")
	tt.Key(0, ui.KeyUp)
	_ = tt.Click("After")
	if p.rule.Interval != 3 || p.rule.End != RepeatCount || p.rule.Count != 10 || p.summary != "Every 3 days, 10 times" {
		t.Fatalf("custom count %+v %q", p.rule, p.summary)
	}
	_ = tt.Click("On date")
	if p.rule.End != RepeatUntil || p.rule.Until != (core.Date{Year: 2026, Month: time.November, Day: 5}) || p.problem != "" {
		t.Fatalf("until %+v %q", p.rule, p.problem)
	}
}

func TestRecurrenceEditorProblems(t *testing.T) {
	p := &recurrenceProbe{rule: Recurrence{Unit: RepeatWeekly, Interval: 1, Weekdays: []time.Weekday{time.Monday}}}
	tt := p.tester()
	_ = tt.Click("Monday")
	if len(p.rule.Weekdays) != 0 || p.problem != "Choose at least one day." || !tt.HasText("Choose at least one day.") {
		t.Fatalf("no-day problem %+v %q", p.rule, p.problem)
	}
	p.rule = Recurrence{Unit: RepeatDaily, Interval: 1, End: RepeatUntil, Until: core.Date{Year: 2026, Month: time.September, Day: 1}}
	tt.Frame()
	if p.problem != "The end date must not precede the start." {
		t.Fatalf("until problem %q", p.problem)
	}
	for _, bad := range []Recurrence{{Unit: 9}, {Unit: RepeatDaily}, {Interval: -1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v: expected panic", bad)
				}
			}()
			q := &recurrenceProbe{rule: bad}
			q.tester()
		}()
	}
}

func TestRecurrenceEditorDisabled(t *testing.T) {
	p := &recurrenceProbe{rule: Recurrence{Unit: RepeatWeekly, Interval: 1, Weekdays: []time.Weekday{time.Monday}}, disabled: true}
	tt := p.tester()
	_ = tt.Click("Repeat")
	_ = tt.Click("Daily")
	_ = tt.Click("Friday")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if p.rule.Unit != RepeatWeekly || len(p.rule.Weekdays) != 1 || p.changes != 0 {
		t.Fatalf("disabled editor changed: %+v", p.rule)
	}
}

func TestRecurrenceEditorSnapshot(t *testing.T) {
	r := Recurrence{Unit: RepeatWeekly, Interval: 2, Weekdays: []time.Weekday{time.Monday, time.Thursday}, End: RepeatCount, Count: 8}
	snapshot(t, "recurrence_editor", 640, 280, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { RecurrenceEditor(c, &r, RecurrenceEditorOptions{Start: oct(5)}) })
	})
}
