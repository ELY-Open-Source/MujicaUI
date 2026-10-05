package datetime

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// YearMonth is a month of a year. The zero value means no month.
type YearMonth struct {
	Year  int
	Month time.Month
}

// IsZero reports whether m is the empty month.
func (m YearMonth) IsZero() bool { return m == YearMonth{} }

// String formats m as 2006-01, or "" when empty.
func (m YearMonth) String() string {
	if m.IsZero() {
		return ""
	}
	return fmt.Sprintf("%04d-%02d", m.Year, m.Month)
}

func (m YearMonth) ordinal() int { return m.Year*12 + int(m.Month) - 1 }

func yearMonthOf(o int) YearMonth {
	y := o / 12
	if o < 0 && o%12 != 0 {
		y--
	}
	return YearMonth{y, time.Month(o-y*12) + 1}
}

// MonthPickerOptions configure a MonthPicker.
type MonthPickerOptions struct {
	// Min and Max bound the months chosen; zero leaves a side open.
	Min, Max     YearMonth
	DisableMonth func(YearMonth) bool
	Clearable    bool
	// Today marks the current month; zero uses the system clock.
	Today    core.Date
	Disabled bool
	// ReadOnly keeps the field focusable and copyable but blocks the popup and clearing.
	ReadOnly bool
	// Label names the field for assistive technology.
	Label string
}

// MonthPicker chooses a month into *value from a year page.
func MonthPicker(c *ui.Context, value *YearMonth, opts MonthPickerOptions) DateTimeResult {
	for _, m := range []YearMonth{*value, opts.Min, opts.Max} {
		if !m.IsZero() && (m.Month < time.January || m.Month > time.December) {
			panic(fmt.Sprintf("mujica: MonthPicker month %v is not valid", m))
		}
	}
	if !opts.Min.IsZero() && !opts.Max.IsZero() && opts.Min.ordinal() > opts.Max.ordinal() {
		panic(fmt.Sprintf("mujica: MonthPicker Min %v is after Max %v", opts.Min, opts.Max))
	}
	today := core.DatetimeToday(c, opts.Today)
	now := YearMonth{today.Year, today.Month}
	title := func(m YearMonth) string { return core.DatetimeMonthTitle(c, m.Year, m.Month) }
	home := *value
	if home.IsZero() {
		home = now
	}
	text := ""
	if !value.IsZero() {
		text = title(*value)
	}
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "month_picker.placeholder")
	}
	col, picked, ok, cleared := datetimeGridPicker(c, datetimeGridField{
		text: text, placeholder: core.Msg(c, "month_picker.placeholder"), label: label,
		clearable: opts.Clearable, hasValue: !value.IsZero(), disabled: opts.Disabled, readOnly: opts.ReadOnly, home: home.ordinal(),
	}, datetimeGridConfig{
		pageOf:  func(o int) (int, int) { return yearMonthOf(o).Year * 12, 12 },
		cols:    3,
		page:    12,
		title:   func(start int) string { return fmt.Sprintf(core.Msg(c, "month_picker.year"), yearMonthOf(start).Year) },
		prev:    core.Msg(c, "calendar.prev_year"),
		next:    core.Msg(c, "calendar.next_year"),
		text:    func(o int) string { return core.DatetimeMonthName(c, yearMonthOf(o).Month) },
		label:   func(o int) string { return title(yearMonthOf(o)) },
		chosen:  func(o int) bool { return !value.IsZero() && o == value.ordinal() },
		current: func(o int) bool { return o == now.ordinal() },
		enabled: func(o int) bool {
			m := yearMonthOf(o)
			return (opts.Min.IsZero() || o >= opts.Min.ordinal()) && (opts.Max.IsZero() || o <= opts.Max.ordinal()) &&
				(opts.DisableMonth == nil || !opts.DisableMonth(m))
		},
		muted: func(int, int) bool { return false },
	})
	to := *value
	switch {
	case cleared:
		to = YearMonth{}
	case ok:
		to = yearMonthOf(picked)
	}
	changed := to != *value
	if changed {
		slog.Debug("mujica: month_picker change", "from", *value, "to", to)
		*value = to
	}
	return DateTimeResult{Element: col, changed: changed}
}
