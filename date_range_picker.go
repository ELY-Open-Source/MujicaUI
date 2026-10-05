package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// DateRange is an inclusive span of days. The zero value means no range.
type DateRange struct {
	Start, End Date
}

// IsZero reports whether r is the empty range.
func (r DateRange) IsZero() bool { return r == DateRange{} }

// Days returns how many days r covers, counting both ends.
func (r DateRange) Days() int {
	if r.Start.IsZero() || r.End.IsZero() {
		return 0
	}
	return int(r.End.civil().Sub(r.Start.civil()).Hours()/24) + 1
}

// DateRangePickerOptions configure a DateRangePicker.
type DateRangePickerOptions struct {
	// Layout is the Go time layout shown; empty uses 2006-01-02.
	Layout string
	// Min and Max bound both ends; zero leaves a side open.
	Min, Max Date
	// DisableDate marks days that neither end nor the span may include.
	DisableDate func(Date) bool
	WeekStart   WeekStart
	// Today is marked in the calendar; zero uses the system clock.
	Today    Date
	Disabled bool
	// Label names the field for assistive technology.
	Label string
}

// DateRangePicker edits *value in a popup: pick a start, preview the span, pick an end, then confirm.
// Clear empties the value; Escape, Cancel or a click outside discards the draft.
func DateRangePicker(c *ui.Context, value *DateRange, opts DateRangePickerOptions) DateTimeResult {
	datetimeCheckDate("DateRangePicker Start", value.Start)
	datetimeCheckDate("DateRangePicker End", value.End)
	if value.Start.IsZero() != value.End.IsZero() || !value.IsZero() && value.End.Compare(value.Start) < 0 {
		panic(fmt.Sprintf("mujica: DateRangePicker value %v..%v is not a range", value.Start, value.End))
	}
	datetimeCheckBounds("DateRangePicker", opts.Min, opts.Max)
	layout := opts.Layout
	if layout == "" {
		layout = "2006-01-02"
	}
	d := Density(c)
	k := Tokens(c)
	today := datetimeToday(c, opts.Today)
	enabled := calendarEnabled(opts.Min, opts.Max, opts.DisableDate)
	col := ui.Column(c).Width(d.ControlHeight() * 9).Gap(d.Unit()).AlignItems(ui.Stretch).Disabled(opts.Disabled)
	open := ui.Local(col, "date-range-open", func() bool { return false })
	wasOpen := ui.Local(col, "date-range-was-open", func() bool { return false })
	draft := ui.Local(col, "date-range-draft", func() DateRange { return DateRange{} })
	cursor := ui.Local(col, "date-range-cursor", func() Date { return Date{} })
	hover := ui.Local(col, "date-range-hover", func() Date { return Date{} })
	problem := ui.Local(col, "date-range-error", func() string { return "" })
	changed, cancelled, done := false, false, false
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
		label = msg(c, "date_range_picker.label")
	}
	text := ""
	if !value.IsZero() {
		text = value.Start.Format(layout) + " – " + value.End.Format(layout)
	}
	var field *ui.Element
	col.Children(func() {
		field = datetimeField(c, ui.Row(c), false)
		field.Children(func() {
			t := datetimeTrigger(c, text, msg(c, "date_range_picker.placeholder"), "calendar").Label(label).Description(text)
			if t.Clicked() || !*open && t.Shortcut(0, ui.KeyDown) {
				*open = !*open
				if *open {
					*draft, *problem, *hover = *value, "", Date{}
					*cursor = value.Start
					if cursor.IsZero() {
						*cursor = today
					}
					slog.Debug("mujica: date_range_picker open")
				}
			}
		})
		datetimePopup(c, field, open, func() {
			ui.Text(c, dateRangeSummary(c, *draft, layout)).TextColor(k.TextMuted).
				FontSize(theme.CaptionSize).FontFeatures("tnum").Role(ui.RoleStatus)
			end := draft.End
			if end.IsZero() && !draft.Start.IsZero() {
				end = *hover
				if end.IsZero() {
					end = *cursor
				}
			}
			picked, hovered := calendarMonth(c, cursor, calendarConfig{
				chosen:    func(x Date) bool { return x == draft.Start || x == draft.End },
				inRange:   func(x Date) bool { return !draft.Start.IsZero() && x.within(draft.Start, end) },
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
			datetimeError(c, *problem)
			datetimeFooter(c, func() {
				if Button(c, msg(c, "calendar.clear"), ButtonOptions{Variant: Ghost, Disabled: value.IsZero() && draft.IsZero()}).Clicked() {
					finish(DateRange{}, "clear")
				}
				if Button(c, msg(c, "calendar.cancel"), ButtonOptions{Variant: Secondary}).Clicked() {
					*open = false
				}
				ready := !draft.Start.IsZero() && !draft.End.IsZero()
				if Button(c, msg(c, "calendar.confirm"), ButtonOptions{Disabled: !ready}).Clicked() {
					finish(*draft, "confirm")
				}
			})
		})
	})
	if *wasOpen && !*open && !done {
		slog.Debug("mujica: date_range_picker cancel")
		cancelled = true
	}
	*wasOpen = *open
	return DateTimeResult{Element: col, changed: changed, cancelled: cancelled}
}

// dateRangePick applies a picked day to the draft and returns any problem.
func dateRangePick(c *ui.Context, draft *DateRange, day Date, enabled func(Date) bool, layout string) string {
	if draft.Start.IsZero() || !draft.End.IsZero() || day.Compare(draft.Start) < 0 {
		*draft = DateRange{Start: day}
		return ""
	}
	for x := draft.Start; x.Compare(day) <= 0; x = x.AddDays(1) {
		if !enabled(x) {
			return fmt.Sprintf(msg(c, "date_range_picker.blocked"), x.Format(layout))
		}
	}
	draft.End = day
	return ""
}

// dateRangeSummary describes the draft: what to pick next, or the span chosen.
func dateRangeSummary(c *ui.Context, r DateRange, layout string) string {
	switch {
	case r.Start.IsZero():
		return msg(c, "date_range_picker.pick_start")
	case r.End.IsZero():
		return fmt.Sprintf(msg(c, "date_range_picker.pick_end"), r.Start.Format(layout))
	}
	return fmt.Sprintf(msg(c, "date_range_picker.summary"), r.Start.Format(layout), r.End.Format(layout), r.Days())
}
