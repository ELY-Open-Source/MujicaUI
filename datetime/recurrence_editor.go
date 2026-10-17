package datetime

import (
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// RepeatUnit is how often a rule repeats; zero never does.
type RepeatUnit int

const (
	RepeatNone RepeatUnit = iota
	RepeatDaily
	RepeatWeekly
	RepeatMonthly
	RepeatYearly
)

// RepeatEnd says when a recurrence stops.
type RepeatEnd int

const (
	RepeatForever RepeatEnd = iota
	RepeatCount
	RepeatUntil
)

// Recurrence repeats every Interval units, on Weekdays when weekly.
type Recurrence struct {
	Unit     RepeatUnit
	Interval int
	Weekdays []time.Weekday
	End      RepeatEnd
	Count    int
	Until    core.Date
}

// RecurrenceEditorOptions configure a RecurrenceEditor.
type RecurrenceEditorOptions struct {
	// Start is the first occurrence; Until must not precede it.
	Start     core.Date
	WeekStart WeekStart
	Disabled  bool
	Label     string
}

// RecurrenceResult is the editor, its change and its problem.
type RecurrenceResult struct {
	Element ui.Element
	changed bool
	problem string
}

// Changed reports whether the user changed the rule this frame.
func (r RecurrenceResult) Changed() bool { return r.changed }

// Problem explains why the rule cannot be saved.
func (r RecurrenceResult) Problem() string { return r.problem }

const recurrenceCustom = 5

// RecurrenceEditor edits a repeat rule and when it ends.
func RecurrenceEditor(c *ui.Context, rule *Recurrence, opts RecurrenceEditorOptions) RecurrenceResult {
	recurrenceCheck(*rule)
	d := core.Density(c)
	k := core.Tokens(c)
	before := recurrenceKey(*rule)
	r := RecurrenceResult{}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "recurrence_editor.label")
	}
	r.Element = ui.Column(c).Gap(d.Unit() * 3).Role(ui.RoleGroup).Label(label).Disabled(opts.Disabled)
	custom := ui.Local(r.Element, "recurrence-custom", func() bool { return rule.Interval > 1 })
	r.Element.Children(func() {
		mode := int(rule.Unit)
		if *custom && rule.Unit != RepeatNone {
			mode = recurrenceCustom
		}
		modes := []input.SelectOption[int]{}
		for i, key := range []string{"none", "daily", "weekly", "monthly", "yearly", "custom"} {
			modes = append(modes, input.SelectOption[int]{Value: i, Label: core.Msg(c, "recurrence_editor."+key)})
		}
		if input.Select(c, &mode, modes, input.SelectOptions{Label: core.Msg(c, "recurrence_editor.repeat"), Disabled: opts.Disabled}).Changed() {
			*custom = mode == recurrenceCustom
			switch {
			case *custom:
				rule.Unit = max(rule.Unit, RepeatDaily)
				rule.Interval = max(rule.Interval, 2)
			default:
				rule.Unit, rule.Interval = RepeatUnit(mode), 1
				if mode == int(RepeatNone) {
					*rule = Recurrence{}
				}
			}
			if rule.Unit == RepeatWeekly && len(rule.Weekdays) == 0 && !opts.Start.IsZero() {
				rule.Weekdays = []time.Weekday{opts.Start.Weekday()}
			}
		}
		if rule.Unit == RepeatNone {
			return
		}
		if *custom {
			ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
				ui.Text(c, core.Msg(c, "recurrence_editor.every"))
				ui.Box(c).Width(d.ControlHeight() * 4).Children(func() {
					n := float64(rule.Interval)
					if input.NumberInput(c, &n, input.NumberInputOptions{Min: 1, Max: 999, Label: core.Msg(c, "recurrence_editor.interval"), Disabled: opts.Disabled}).Changed() && !math.IsNaN(n) {
						rule.Interval = int(n)
					}
				})
				units := []input.SelectOption[RepeatUnit]{}
				for u := RepeatDaily; u <= RepeatYearly; u++ {
					units = append(units, input.SelectOption[RepeatUnit]{Value: u, Label: core.Msg(c, recurrenceUnitKeys[u])})
				}
				input.Select(c, &rule.Unit, units, input.SelectOptions{Label: core.Msg(c, "recurrence_editor.unit"), Disabled: opts.Disabled})
			})
		}
		if rule.Unit == RepeatWeekly {
			first := CalendarFirstWeekday(c, opts.WeekStart)
			items := []input.ToggleItem[time.Weekday]{}
			for i := range 7 {
				wd := time.Weekday((int(first) + i) % 7)
				items = append(items, input.ToggleItem[time.Weekday]{Value: wd, Label: core.Msg(c, fmt.Sprintf("calendar.weekday_%d", wd)),
					Name: core.Msg(c, fmt.Sprintf("calendar.weekday_long_%d", wd))})
			}
			input.ToggleGroup(c, &rule.Weekdays, items, input.ToggleGroupOptions{Multiple: true, AllowEmpty: true, Label: core.Msg(c, "recurrence_editor.on_days"), Disabled: opts.Disabled})
		}
		ends := []input.RadioItem[RepeatEnd]{{Value: RepeatForever, Label: core.Msg(c, "recurrence_editor.forever")},
			{Value: RepeatCount, Label: core.Msg(c, "recurrence_editor.after")}, {Value: RepeatUntil, Label: core.Msg(c, "recurrence_editor.until")}}
		ui.Row(c).Gap(d.Unit() * 3).AlignItems(ui.Center).Wrap().Children(func() {
			if input.RadioGroup(c, &rule.End, ends, input.RadioGroupOptions{Horizontal: true, Label: core.Msg(c, "recurrence_editor.ends"), Disabled: opts.Disabled}).Changed() {
				if rule.End == RepeatCount && rule.Count == 0 {
					rule.Count = 10
				}
				if rule.End == RepeatUntil && rule.Until.IsZero() && !opts.Start.IsZero() {
					rule.Until = opts.Start.AddMonths(1)
				}
			}
			switch rule.End {
			case RepeatCount:
				ui.Box(c).Width(d.ControlHeight() * 4).Children(func() {
					n := float64(rule.Count)
					if input.NumberInput(c, &n, input.NumberInputOptions{Min: 1, Max: 999, Label: core.Msg(c, "recurrence_editor.count"), Disabled: opts.Disabled}).Changed() && !math.IsNaN(n) {
						rule.Count = int(n)
					}
				})
			case RepeatUntil:
				DatePicker(c, &rule.Until, DatePickerOptions{Label: core.Msg(c, "recurrence_editor.until_date"), Min: opts.Start, Disabled: opts.Disabled})
			}
		})
		r.problem = recurrenceProblem(c, *rule, opts.Start)
		ui.Text(c, RecurrenceSummary(c, *rule)).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).Role(ui.RoleStatus)
		core.DatetimeError(c, r.problem)
	})
	if after := recurrenceKey(*rule); after != before && !r.Element.IsDisabled() {
		slog.Debug("mujica: recurrence_editor change", "rule", after)
		r.changed = true
	}
	return r
}

var recurrenceUnitKeys = [...]string{"", "recurrence_editor.days", "recurrence_editor.weeks", "recurrence_editor.months", "recurrence_editor.years"}

func recurrenceKey(r Recurrence) string { return fmt.Sprintf("%v", r) }

// recurrenceCheck panics on rules no editor can show.
func recurrenceCheck(r Recurrence) {
	if r.Unit < RepeatNone || r.Unit > RepeatYearly || r.End < RepeatForever || r.End > RepeatUntil || r.Interval < 0 || r.Count < 0 {
		panic(fmt.Sprintf("mujica: Recurrence %+v is out of range", r))
	}
	if r.Unit != RepeatNone && r.Interval == 0 {
		panic("mujica: Recurrence repeats with Interval 0")
	}
	core.DatetimeCheckDate("Recurrence Until", r.Until)
}

// recurrenceProblem explains why a rule cannot be saved.
func recurrenceProblem(c *ui.Context, r Recurrence, start core.Date) string {
	switch {
	case r.Unit == RepeatNone:
		return ""
	case r.Unit == RepeatWeekly && len(r.Weekdays) == 0:
		return core.Msg(c, "recurrence_editor.no_days")
	case r.End == RepeatCount && r.Count < 1:
		return core.Msg(c, "recurrence_editor.bad_count")
	case r.End == RepeatUntil && (r.Until.IsZero() || !start.IsZero() && r.Until.Compare(start) < 0):
		return core.Msg(c, "recurrence_editor.bad_until")
	}
	return ""
}

// RecurrenceSummary describes a rule in the window's language.
func RecurrenceSummary(c *ui.Context, r Recurrence) string {
	if r.Unit == RepeatNone {
		return core.Msg(c, "recurrence_editor.never")
	}
	s := fmt.Sprintf(core.Msg(c, "recurrence_editor.every_n"), r.Interval, core.Msg(c, recurrenceUnitKeys[r.Unit]))
	if r.Interval == 1 {
		s = core.Msg(c, []string{"", "recurrence_editor.daily", "recurrence_editor.weekly", "recurrence_editor.monthly", "recurrence_editor.yearly"}[r.Unit])
	}
	if r.Unit == RepeatWeekly && len(r.Weekdays) > 0 {
		var days []string
		for _, wd := range r.Weekdays {
			days = append(days, core.Msg(c, fmt.Sprintf("calendar.weekday_long_%d", wd)))
		}
		s += fmt.Sprintf(core.Msg(c, "recurrence_editor.on"), strings.Join(days, core.Msg(c, "recurrence_editor.list_sep")))
	}
	switch r.End {
	case RepeatCount:
		s += core.CalendarCount(c, "recurrence_editor.times", r.Count)
	case RepeatUntil:
		s += fmt.Sprintf(core.Msg(c, "recurrence_editor.until_text"), r.Until)
	}
	return s
}
