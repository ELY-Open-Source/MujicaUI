package datetime

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DateRange is an inclusive span of days. The zero value means no range.
type DateRange struct {
	Start, End core.Date
}

// IsZero reports whether r is the empty range.
func (r DateRange) IsZero() bool { return r == DateRange{} }

// Days returns how many days r covers, counting both ends.
func (r DateRange) Days() int {
	if r.Start.IsZero() || r.End.IsZero() {
		return 0
	}
	return int(r.End.Civil().Sub(r.Start.Civil()).Hours()/24) + 1
}

// DateRangePickerOptions configure a DateRangePicker.
type DateRangePickerOptions struct {
	// Layout is the Go time layout shown; empty uses 2006-01-02.
	Layout string
	// Min and Max bound both ends; zero leaves a side open.
	Min, Max core.Date
	// DisableDate marks days that neither end nor the span may include.
	DisableDate func(core.Date) bool
	WeekStart   WeekStart
	// Today is marked in the calendar; zero uses the system clock.
	Today    core.Date
	Disabled bool
	// ReadOnly keeps the field focusable and copyable but blocks the popup.
	ReadOnly bool
	// Label names the field for assistive technology.
	Label string
}

// DateRangePicker picks start and end; Confirm writes *value.
func DateRangePicker(c *ui.Context, value *DateRange, opts DateRangePickerOptions) DateTimeResult {
	core.DatetimeCheckDate("DateRangePicker Start", value.Start)
	core.DatetimeCheckDate("DateRangePicker End", value.End)
	if value.Start.IsZero() != value.End.IsZero() || !value.IsZero() && value.End.Compare(value.Start) < 0 {
		panic(fmt.Sprintf("mujica: DateRangePicker value %v..%v is not a range", value.Start, value.End))
	}
	core.DatetimeCheckBounds("DateRangePicker", opts.Min, opts.Max)
	layout := opts.Layout
	if layout == "" {
		layout = "2006-01-02"
	}
	d := core.Density(c)
	k := core.Tokens(c)
	today := core.DatetimeToday(c, opts.Today)
	enabled := calendarEnabled(opts.Min, opts.Max, opts.DisableDate)
	col := ui.Column(c).Width(d.ControlHeight() * 9).Gap(d.Unit()).AlignItems(ui.Stretch).Disabled(opts.Disabled)
	open := ui.Local(col, "date-range-open", func() bool { return false })
	wasOpen := ui.Local(col, "date-range-was-open", func() bool { return false })
	draft := ui.Local(col, "date-range-draft", func() DateRange { return DateRange{} })
	cursor := ui.Local(col, "date-range-cursor", func() core.Date { return core.Date{} })
	hover := ui.Local(col, "date-range-hover", func() core.Date { return core.Date{} })
	problem := ui.Local(col, "date-range-error", func() string { return "" })
	changed, canceled, done := false, false, false
	finish := func(to DateRange, how string) {
		if to != *value {
			slog.Debug("mujica: date_range_picker change", "from", *value, "to", to, "by", how)
			*value = to
			changed = true
		}
		*open, done = false, true
	}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "date_range_picker.label")
	}
	text := ""
	if !value.IsZero() {
		text = value.Start.Format(layout) + " – " + value.End.Format(layout)
	}
	var field ui.Element
	col.Children(func() {
		field = core.DatetimeField(c, ui.Row(c), *problem != "", opts.ReadOnly)
		field.Children(func() {
			t := core.DatetimeTrigger(c, text, core.Msg(c, "date_range_picker.placeholder"), "calendar").Label(label).Description(text).Error(*problem)
			if !opts.ReadOnly && (t.Clicked() || !*open && t.Shortcut(0, ui.KeyDown)) {
				*open = !*open
				if *open {
					*draft, *problem, *hover = *value, "", core.Date{}
					*cursor = value.Start
					if cursor.IsZero() {
						*cursor = today
					}
					slog.Debug("mujica: date_range_picker open")
				}
			}
		})
		core.DatetimePopup(c, field, open, opts.ReadOnly, func() {
			ui.Text(c, dateRangeSummary(c, *draft, layout)).TextColor(k.TextMuted).
				FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum").Role(ui.RoleStatus)
			end := draft.End
			if end.IsZero() && !draft.Start.IsZero() {
				end = *hover
				if end.IsZero() {
					end = *cursor
				}
			}
			picked, hovered := calendarMonth(c, cursor, calendarConfig{
				chosen:    func(x core.Date) bool { return x == draft.Start || x == draft.End },
				inRange:   func(x core.Date) bool { return !draft.Start.IsZero() && x.Within(draft.Start, end) },
				enabled:   enabled,
				today:     today,
				weekStart: opts.WeekStart,
				autoFocus: true,
			})
			if hovered != *hover {
				*hover = hovered
				c.Invalidate()
			}
			if !picked.IsZero() {
				*problem = dateRangePick(c, draft, picked, enabled, layout)
			}
			core.DatetimeError(c, *problem)
			core.DatetimeFooter(c, func() {
				if input.Button(c, core.Msg(c, "calendar.clear"), input.ButtonOptions{Variant: input.Ghost, Disabled: value.IsZero() && draft.IsZero()}).Clicked() {
					finish(DateRange{}, "clear")
				}
				if input.Button(c, core.Msg(c, "calendar.cancel"), input.ButtonOptions{Variant: input.Secondary}).Clicked() {
					*open = false
				}
				ready := !draft.Start.IsZero() && !draft.End.IsZero()
				if input.Button(c, core.Msg(c, "calendar.confirm"), input.ButtonOptions{Disabled: !ready}).Clicked() {
					if p := dateRangeProblem(c, *draft, opts.Min, opts.Max, enabled, layout); p != "" {
						slog.Debug("mujica: date_range_picker invalid confirm", "draft", *draft)
						*problem = p
					} else {
						finish(*draft, "confirm")
					}
				}
			})
		})
	})
	if *wasOpen && !*open && !done {
		slog.Debug("mujica: date_range_picker cancel")
		canceled = true
	}
	*wasOpen = *open
	if !*open {
		*problem = ""
	}
	return DateTimeResult{Element: col, changed: changed, canceled: canceled}
}

// dateRangeProblem checks a whole draft against the current bounds and filter.
func dateRangeProblem(c *ui.Context, r DateRange, lo, hi core.Date, enabled func(core.Date) bool, layout string) string {
	format := func(d core.Date) string { return d.Format(layout) }
	for _, d := range []core.Date{r.Start, r.End} {
		if p := datePickerProblem(c, d, lo, hi, enabled, format); p != "" {
			return p
		}
	}
	for x := r.Start; x.Compare(r.End) <= 0; x = x.AddDays(1) {
		if !enabled(x) {
			return fmt.Sprintf(core.Msg(c, "date_range_picker.blocked"), format(x))
		}
	}
	return ""
}

// dateRangePick applies a picked day to the draft and returns any problem.
func dateRangePick(c *ui.Context, draft *DateRange, day core.Date, enabled func(core.Date) bool, layout string) string {
	if draft.Start.IsZero() || !draft.End.IsZero() || day.Compare(draft.Start) < 0 {
		*draft = DateRange{Start: day}
		return ""
	}
	for x := draft.Start; x.Compare(day) <= 0; x = x.AddDays(1) {
		if !enabled(x) {
			return fmt.Sprintf(core.Msg(c, "date_range_picker.blocked"), x.Format(layout))
		}
	}
	draft.End = day
	return ""
}

// dateRangeSummary describes the draft: what to pick next, or the span chosen.
func dateRangeSummary(c *ui.Context, r DateRange, layout string) string {
	switch {
	case r.Start.IsZero():
		return core.Msg(c, "date_range_picker.pick_start")
	case r.End.IsZero():
		return fmt.Sprintf(core.Msg(c, "date_range_picker.pick_end"), r.Start.Format(layout))
	}
	return fmt.Sprintf(core.Msg(c, "date_range_picker.summary"), r.Start.Format(layout), r.End.Format(layout), r.Days())
}
