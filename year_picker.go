package mujica

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"
)

// YearPickerOptions configure a YearPicker.
type YearPickerOptions struct {
	// Min and Max bound the years chosen; zero leaves a side open.
	Min, Max    int
	DisableYear func(int) bool
	Clearable   bool
	// Today marks the current year; zero uses the system clock.
	Today    Date
	Disabled bool
	// Label names the field for assistive technology.
	Label string
}

// YearPicker chooses a year into *value, zero meaning none, from a decade page.
// Arrows move by year and row, PageUp and PageDown by decade.
func YearPicker(c *ui.Context, value *int, opts YearPickerOptions) DateTimeResult {
	if *value < 0 || opts.Min < 0 || opts.Max < 0 || opts.Min != 0 && opts.Max != 0 && opts.Min > opts.Max {
		panic(fmt.Sprintf("mujica: YearPicker value %d or range %d..%d is not valid", *value, opts.Min, opts.Max))
	}
	now := datetimeToday(c, opts.Today).Year
	home := *value
	if home == 0 {
		home = now
	}
	decade := func(o int) int { return o / 10 * 10 }
	text := ""
	if *value != 0 {
		text = strconv.Itoa(*value)
	}
	label := opts.Label
	if label == "" {
		label = msg(c, "year_picker.placeholder")
	}
	col, picked, ok, cleared := datetimeGridPicker(c, datetimeGridField{
		text: text, placeholder: msg(c, "year_picker.placeholder"), label: label,
		clearable: opts.Clearable, hasValue: *value != 0, disabled: opts.Disabled, home: home,
	}, datetimeGridConfig{
		pageOf:  func(o int) (int, int) { return decade(o) - 1, 12 },
		cols:    3,
		page:    10,
		title:   func(start int) string { return fmt.Sprintf("%d – %d", start+1, start+10) },
		prev:    msg(c, "year_picker.prev"),
		next:    msg(c, "year_picker.next"),
		text:    strconv.Itoa,
		label:   strconv.Itoa,
		chosen:  func(o int) bool { return o == *value },
		current: func(o int) bool { return o == now },
		enabled: func(o int) bool {
			return o > 0 && (opts.Min == 0 || o >= opts.Min) && (opts.Max == 0 || o <= opts.Max) &&
				(opts.DisableYear == nil || !opts.DisableYear(o))
		},
		muted: func(o, start int) bool { return o == start || o == start+11 },
	})
	to := *value
	switch {
	case cleared:
		to = 0
	case ok:
		to = picked
	}
	changed := to != *value
	if changed {
		slog.Debug("mujica: year_picker change", "from", *value, "to", to)
		*value = to
	}
	return DateTimeResult{Element: col, changed: changed}
}
