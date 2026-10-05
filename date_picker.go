package mujica

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// DatePickerOptions configure a DatePicker.
type DatePickerOptions struct {
	// Layout is the Go time layout typed and shown; empty uses 2006-01-02.
	Layout string
	// Min and Max bound the dates accepted; zero leaves a side open.
	Min, Max    Date
	DisableDate func(Date) bool
	WeekStart   WeekStart
	// Today is marked in the calendar; zero uses the system clock.
	Today Date
	// Clearable shows a clear button and accepts empty text as no date.
	Clearable bool
	// Placeholder shows in the empty field; empty uses the built-in copy.
	Placeholder string
	Disabled    bool
	// ReadOnly keeps the text focusable and copyable but blocks edits and the popup.
	ReadOnly bool
	// Label names the field for assistive technology.
	Label string
}

// DatePicker types or picks *value; invalid text is kept and reported.
func DatePicker(c *ui.Context, value *Date, opts DatePickerOptions) DateTimeResult {
	datetimeCheckDate("DatePicker value", *value)
	datetimeCheckBounds("DatePicker", opts.Min, opts.Max)
	layout := opts.Layout
	if layout == "" {
		layout = "2006-01-02"
	}
	format := func(d Date) string {
		if d.IsZero() {
			return ""
		}
		return d.Format(layout)
	}
	d := Density(c)
	today := datetimeToday(c, opts.Today)
	enabled := calendarEnabled(opts.Min, opts.Max, opts.DisableDate)
	col := ui.Column(c).Width(d.ControlHeight() * 7).Gap(d.Unit()).AlignItems(ui.Stretch).Disabled(opts.Disabled)
	text := ui.Local(col, "date-picker-text", func() string { return format(*value) })
	seen := ui.Local(col, "date-picker-seen", func() Date { return *value })
	problem := ui.Local(col, "date-picker-error", func() string { return "" })
	open := ui.Local(col, "date-picker-open", func() bool { return false })
	cursor := ui.Local(col, "date-picker-cursor", func() Date { return Date{} })
	wasFocused := ui.Local(col, "date-picker-focused", func() bool { return false })
	shownLayout := ui.Local(col, "date-picker-layout", func() string { return layout })
	if *shownLayout != layout {
		if *text == Date.Format(*value, *shownLayout) || value.IsZero() && *text == "" {
			*text = format(*value)
		}
		*shownLayout = layout
	}
	if *seen != *value {
		*seen, *text, *problem = *value, format(*value), ""
	}
	changed := false
	commit := func(to Date, how string) {
		if to != *value {
			slog.Debug("mujica: date_picker change", "from", *value, "to", to, "by", how)
			*value = to
			changed = true
		}
		*seen, *text, *problem = to, format(to), ""
	}
	submit := func() {
		s := strings.TrimSpace(*text)
		if s == "" && opts.Clearable {
			commit(Date{}, "text")
			return
		}
		t, err := time.Parse(layout, s)
		if err != nil {
			*problem = fmt.Sprintf(msg(c, "date_picker.invalid"), format(today))
			slog.Debug("mujica: date_picker invalid", "text", s)
			return
		}
		if p := datePickerProblem(c, DateOf(t), opts.Min, opts.Max, enabled, format); p != "" {
			*problem = p
			return
		}
		commit(DateOf(t), "text")
	}
	openAt := func() {
		*open = true
		*cursor = *value
		if cursor.IsZero() {
			*cursor = today
		}
	}
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = msg(c, "date_picker.placeholder")
	}
	var field *ui.Element
	col.Children(func() {
		field = datetimeField(c, ui.Row(c), *problem != "", opts.ReadOnly)
		field.Children(func() {
			in := textInputEditor(c, field, text, opts.ReadOnly, ui.TextInputBase).Grow(1).Shrink(1).Placeholder(placeholder).FontFeatures("tnum")
			label := opts.Label
			if label == "" {
				label = placeholder
			}
			in.Label(label).Error(*problem)
			focused := in.Focused()
			if !opts.ReadOnly && !in.IsDisabled() && (in.Submitted() || *wasFocused && !focused && *text != format(*value)) {
				submit()
			}
			*wasFocused = focused
			if opts.ReadOnly {
				return
			}
			if opts.Clearable && !value.IsZero() && datetimeIconButton(c, "x", msg(c, "calendar.clear")).Clicked() {
				commit(Date{}, "clear")
			}
			if datetimeIconButton(c, "calendar", msg(c, "date_picker.open")).Role(ui.RolePopUpButton).Clicked() {
				if *open {
					*open = false
				} else {
					openAt()
				}
			}
		})
		if !*open && !opts.ReadOnly && (field.Shortcut(0, ui.KeyDown) || field.Shortcut(ui.Alt, ui.KeyDown)) {
			openAt()
		}
		datetimePopup(c, field, open, opts.ReadOnly, func() {
			picked, _ := calendarMonth(c, cursor, calendarConfig{
				chosen:    func(x Date) bool { return x == *value },
				enabled:   enabled,
				today:     today,
				weekStart: opts.WeekStart,
				autoFocus: true,
			})
			if !picked.IsZero() {
				commit(picked, "calendar")
				*open = false
			}
		})
		datetimeError(c, *problem)
	})
	return DateTimeResult{Element: col, changed: changed}
}

// datePickerProblem explains why an existing day d cannot be chosen, or returns "".
func datePickerProblem(c *ui.Context, d Date, lo, hi Date, enabled func(Date) bool, format func(Date) string) string {
	switch {
	case !lo.IsZero() && d.Compare(lo) < 0:
		return fmt.Sprintf(msg(c, "calendar.before_min"), format(lo))
	case !hi.IsZero() && d.Compare(hi) > 0:
		return fmt.Sprintf(msg(c, "calendar.after_max"), format(hi))
	case !enabled(d):
		return fmt.Sprintf(msg(c, "calendar.unavailable"), format(d))
	}
	return ""
}
