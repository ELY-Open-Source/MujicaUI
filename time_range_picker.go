package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// TimeRange is a start and an end time of day. End before Start means it ends the next day.
type TimeRange struct {
	Start, End TimeOfDay
}

// Overnight reports whether r crosses midnight.
func (r TimeRange) Overnight() bool { return r.End.Compare(r.Start) < 0 }

// TimeRangePickerOptions configure a TimeRangePicker.
type TimeRangePickerOptions struct {
	Seconds bool
	Hour12  bool
	// MinuteStep and SecondStep must divide 60; zero means 1.
	MinuteStep, SecondStep int
	// Min and Max bound both ends; a zero Max means 23:59:59.
	Min, Max TimeOfDay
	// AllowOvernight accepts an end before the start, meaning the next day.
	AllowOvernight bool
	Disabled       bool
	// Label names the pair for assistive technology.
	Label string
}

// TimeRangePicker edits a start and an end time. It writes *value only when both ends are valid
// and the end follows the start, or crosses midnight with AllowOvernight.
func TimeRangePicker(c *ui.Context, value *TimeRange, opts TimeRangePickerOptions) DateTimeResult {
	cfg := timePickerConfig("TimeRangePicker", TimePickerOptions{Seconds: opts.Seconds, Hour12: opts.Hour12,
		MinuteStep: opts.MinuteStep, SecondStep: opts.SecondStep, Min: opts.Min, Max: opts.Max})
	if !value.Start.Valid() || !value.End.Valid() {
		panic(fmt.Sprintf("mujica: TimeRangePicker value %v is not a pair of times", *value))
	}
	if value.Overnight() && !opts.AllowOvernight {
		panic(fmt.Sprintf("mujica: TimeRangePicker value %v crosses midnight without AllowOvernight", *value))
	}
	k := Tokens(c)
	d := Density(c)
	col := ui.Column(c).Gap(d.Unit()).AlignItems(ui.Start).Disabled(opts.Disabled).Role(ui.RoleGroup)
	if opts.Label != "" {
		col.Label(opts.Label)
	}
	draft := ui.Local(col, "time-range-draft", func() TimeRange { return *value })
	seen := ui.Local(col, "time-range-seen", func() TimeRange { return *value })
	problem := ui.Local(col, "time-range-error", func() string { return "" })
	if *seen != *value {
		*draft, *seen, *problem = *value, *value, ""
	}
	changed := false
	col.Children(func() {
		edited := false
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			start, end := cfg, cfg
			start.label, end.label = msg(c, "time_range_picker.start"), msg(c, "time_range_picker.end")
			start.invalid, end.invalid = *problem != "", *problem != ""
			if _, e := timeEditor(c, &draft.Start, start); e {
				edited = true
			}
			ui.Text(c, "–").TextColor(k.TextMuted).Role(ui.RoleNone)
			if _, e := timeEditor(c, &draft.End, end); e {
				edited = true
			}
			if draft.Overnight() && opts.AllowOvernight {
				ui.Text(c, msg(c, "time_range_picker.next_day")).FontSize(theme.CaptionSize).TextColor(k.AccentText)
			}
		})
		if edited {
			*problem = timeRangeProblem(c, *draft, cfg, opts.AllowOvernight)
			if *problem == "" && *draft != *value {
				slog.Debug("mujica: time_range_picker change", "from", *value, "to", *draft)
				*value, *seen, changed = *draft, *draft, true
			}
		}
		datetimeError(c, *problem)
	})
	return DateTimeResult{Element: col, changed: changed}
}

// timeRangeProblem explains why r cannot be written, or returns "".
func timeRangeProblem(c *ui.Context, r TimeRange, cfg timeConfig, overnight bool) string {
	for _, t := range []TimeOfDay{r.Start, r.End} {
		if p := timePickerProblem(c, t, cfg); p != "" {
			return p
		}
	}
	switch {
	case r.Start == r.End:
		return msg(c, "time_range_picker.same")
	case r.Overnight() && !overnight:
		return msg(c, "time_range_picker.order")
	}
	return ""
}
