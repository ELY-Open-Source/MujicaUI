package datetime

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TimeRange is a span of day; End before Start means next day.
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
	// ReadOnly keeps both ends focusable and copyable but blocks edits.
	ReadOnly bool
	// Label names the pair for assistive technology.
	Label string
}

// TimeRangePicker writes *value only when the range is valid.
func TimeRangePicker(c *ui.Context, value *TimeRange, opts TimeRangePickerOptions) DateTimeResult {
	cfg := timePickerConfig("TimeRangePicker", TimePickerOptions{Seconds: opts.Seconds, Hour12: opts.Hour12,
		MinuteStep: opts.MinuteStep, SecondStep: opts.SecondStep, Min: opts.Min, Max: opts.Max, ReadOnly: opts.ReadOnly})
	if !value.Start.Valid() || !value.End.Valid() {
		panic(fmt.Sprintf("mujica: TimeRangePicker value %v is not a pair of times", *value))
	}
	if value.Overnight() && !opts.AllowOvernight {
		panic(fmt.Sprintf("mujica: TimeRangePicker value %v crosses midnight without AllowOvernight", *value))
	}
	k := core.Tokens(c)
	d := core.Density(c)
	col := ui.Column(c).Gap(d.Unit()).AlignItems(ui.Start).Disabled(opts.Disabled).Role(ui.RoleGroup)
	if opts.Label != "" {
		col.Label(opts.Label)
	}
	draft := ui.Local(col, "time-range-draft", func() TimeRange { return *value })
	seen := ui.Local(col, "time-range-seen", func() TimeRange { return *value })
	problem := ui.Local(col, "time-range-error", func() string { return "" })
	pending := ui.Local(col, "time-range-pending", func() bool { return false })
	if *seen != *value {
		*draft, *seen, *pending = *value, *value, false
	}
	changed := false
	col.Children(func() {
		edited, bad := false, ""
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			start, end := cfg, cfg
			start.label, end.label = core.Msg(c, "time_range_picker.start"), core.Msg(c, "time_range_picker.end")
			start.problem, end.problem = *problem, *problem
			_, e1, b1 := timeEditor(c, &draft.Start, start)
			ui.Text(c, "–").TextColor(k.TextMuted).Role(ui.RoleNone)
			_, e2, b2 := timeEditor(c, &draft.End, end)
			edited, bad = e1 || e2, b1
			if bad == "" {
				bad = b2
			}
			if draft.Overnight() && opts.AllowOvernight {
				ui.Text(c, core.Msg(c, "time_range_picker.next_day")).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.AccentText)
			}
		})
		p := bad
		if p == "" {
			p = timeRangeProblem(c, *draft, cfg, opts.AllowOvernight)
		}
		if p != *problem {
			*problem = p
			c.Invalidate()
		}
		*pending = *pending || edited
		if *pending && p == "" {
			if *draft != *value {
				slog.Debug("mujica: time_range_picker change", "from", *value, "to", *draft)
				*value, *seen, changed = *draft, *draft, true
			}
			*pending = false
		}
		core.DatetimeError(c, *problem)
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
		return core.Msg(c, "time_range_picker.same")
	case r.Overnight() && !overnight:
		return core.Msg(c, "time_range_picker.order")
	}
	return ""
}
