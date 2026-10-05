package datetime

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DateTimePickerOptions configure a DateTimePicker.
type DateTimePickerOptions struct {
	// Location is the time zone the user edits in; required.
	Location *time.Location
	// Layout is the Go layout of the date part; empty uses 2006-01-02.
	Layout string
	// Min and Max bound the instant; zero leaves a side open.
	Min, Max   time.Time
	Seconds    bool
	Hour12     bool
	MinuteStep int
	WeekStart  WeekStart
	// Today is marked in the calendar; zero uses the system clock.
	Today    core.Date
	Disabled bool
	// ReadOnly keeps the field focusable and copyable but blocks the popup.
	ReadOnly bool
	// Label names the field for assistive technology.
	Label string
}

// DateTimePicker edits *value as wall time in opts.Location.
func DateTimePicker(c *ui.Context, value *time.Time, opts DateTimePickerOptions) DateTimeResult {
	loc := opts.Location
	if loc == nil {
		panic("mujica: DateTimePicker needs a Location")
	}
	if !opts.Min.IsZero() && !opts.Max.IsZero() && opts.Min.After(opts.Max) {
		panic(fmt.Sprintf("mujica: DateTimePicker Min %v is after Max %v", opts.Min, opts.Max))
	}
	cfg := timePickerConfig("DateTimePicker", TimePickerOptions{Seconds: opts.Seconds, Hour12: opts.Hour12, MinuteStep: opts.MinuteStep})
	cfg.label = core.Msg(c, "date_time_picker.time")
	k := core.Tokens(c)
	d := core.Density(c)
	today := core.DatetimeToday(c, opts.Today)
	var lo, hi core.Date
	if !opts.Min.IsZero() {
		lo = core.DateOf(opts.Min.In(loc))
	}
	if !opts.Max.IsZero() {
		hi = core.DateOf(opts.Max.In(loc))
	}
	col := ui.Column(c).Width(d.ControlHeight() * 9).Gap(d.Unit()).AlignItems(ui.Stretch).Disabled(opts.Disabled)
	open := ui.Local(col, "date-time-open", func() bool { return false })
	wasOpen := ui.Local(col, "date-time-was-open", func() bool { return false })
	day := ui.Local(col, "date-time-day", func() core.Date { return core.Date{} })
	clock := ui.Local(col, "date-time-clock", func() TimeOfDay { return TimeOfDay{} })
	cursor := ui.Local(col, "date-time-cursor", func() core.Date { return core.Date{} })
	problem := ui.Local(col, "date-time-error", func() string { return "" })
	typedBad := ui.Local(col, "date-time-typed-bad", func() bool { return false })
	if opts.Layout == "" {
		opts.Layout = "2006-01-02"
	}
	format := func(t time.Time) string { return dateTimeText(c, t.In(loc), opts.Layout, cfg) }
	changed, canceled, done := false, false, false
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "date_time_picker.placeholder")
	}
	text := ""
	if !value.IsZero() {
		text = format(*value)
	}
	col.Children(func() {
		field := core.DatetimeField(c, ui.Row(c), *problem != "", opts.ReadOnly)
		field.Children(func() {
			t := core.DatetimeTrigger(c, text, core.Msg(c, "date_time_picker.placeholder"), "calendar").Label(label).Description(text).Error(*problem)
			if !opts.ReadOnly && (t.Clicked() || !*open && t.Shortcut(0, ui.KeyDown)) {
				*open = !*open
				if *open {
					*day, *clock, *problem = today, TimeOfDay{}, ""
					if !value.IsZero() {
						w := value.In(loc)
						*day, *clock = core.DateOf(w), TimeOfDay{w.Hour(), w.Minute(), w.Second()}
					}
					*cursor = *day
					slog.Debug("mujica: date_time_picker open")
				}
			}
		})
		core.DatetimePopup(c, field, open, opts.ReadOnly, func() {
			picked, _ := calendarMonth(c, cursor, calendarConfig{
				chosen:    func(x core.Date) bool { return x == *day },
				enabled:   calendarEnabled(lo, hi, nil),
				today:     today,
				weekStart: opts.WeekStart,
				autoFocus: true,
			})
			if !picked.IsZero() {
				*day, *problem = picked, ""
			}
			bad := ""
			ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
				cfg.problem = *problem
				_, edited, b := timeEditor(c, clock, cfg)
				switch {
				case b != "":
					*problem, *typedBad, bad = b, true, b
				case edited || *typedBad:
					*problem, *typedBad = "", false
				}
				ui.Text(c, loc.String()).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
			})
			core.DatetimeError(c, *problem)
			core.DatetimeFooter(c, func() {
				if input.Button(c, core.Msg(c, "calendar.cancel"), input.ButtonOptions{Variant: input.Secondary}).Clicked() {
					*open = false
				}
				if input.Button(c, core.Msg(c, "calendar.confirm"), input.ButtonOptions{}).Clicked() {
					at, p := dateTimeConfirm(c, *day, *clock, loc, opts, cfg)
					if bad != "" {
						p = bad
					}
					*problem = p
					if p == "" {
						if !at.Equal(*value) {
							slog.Debug("mujica: date_time_picker change", "from", *value, "to", at)
							*value, changed = at, true
						}
						*open, done = false, true
					}
				}
			})
		})
	})
	if *wasOpen && !*open && !done {
		slog.Debug("mujica: date_time_picker cancel")
		canceled = true
	}
	*wasOpen = *open
	if !*open {
		*problem = ""
	}
	return DateTimeResult{Element: col, changed: changed, canceled: canceled}
}

// dateTimeConfirm turns the draft into an instant, or explains why it cannot.
func dateTimeConfirm(c *ui.Context, d core.Date, t TimeOfDay, loc *time.Location, opts DateTimePickerOptions, cfg timeConfig) (time.Time, string) {
	if p := timePickerProblem(c, t, cfg); p != "" {
		return time.Time{}, p
	}
	wall := d.Format(opts.Layout) + " " + timeText(c, t, cfg)
	at, n := dateTimeResolve(d, t, loc)
	switch {
	case n == 0:
		slog.Debug("mujica: date_time_picker gap", "wall", wall, "zone", loc.String())
		return at, fmt.Sprintf(core.Msg(c, "date_time_picker.gap"), wall, loc.String())
	case n > 1:
		slog.Debug("mujica: date_time_picker ambiguous", "wall", wall, "zone", loc.String())
		return at, fmt.Sprintf(core.Msg(c, "date_time_picker.ambiguous"), wall, loc.String())
	case !opts.Min.IsZero() && at.Before(opts.Min):
		return at, fmt.Sprintf(core.Msg(c, "calendar.before_min"), dateTimeText(c, opts.Min.In(loc), opts.Layout, cfg))
	case !opts.Max.IsZero() && at.After(opts.Max):
		return at, fmt.Sprintf(core.Msg(c, "calendar.after_max"), dateTimeText(c, opts.Max.In(loc), opts.Layout, cfg))
	}
	return at, ""
}

// dateTimeResolve counts instants for a wall time: 0 gap, 2 overlap.
func dateTimeResolve(d core.Date, t TimeOfDay, loc *time.Location) (time.Time, int) {
	wall := time.Date(d.Year, d.Month, d.Day, t.Hour, t.Minute, t.Second, 0, time.UTC)
	guess := time.Date(d.Year, d.Month, d.Day, t.Hour, t.Minute, t.Second, 0, loc)
	var hits []time.Time
	for _, probe := range []time.Time{guess.Add(-24 * time.Hour), guess, guess.Add(24 * time.Hour)} {
		_, off := probe.Zone()
		u := wall.Add(-time.Duration(off) * time.Second).In(loc)
		y, m, dd := u.Date()
		if (core.Date{Year: y, Month: m, Day: dd}) != d || u.Hour() != t.Hour || u.Minute() != t.Minute || u.Second() != t.Second {
			continue
		}
		seen := false
		for _, h := range hits {
			seen = seen || h.Equal(u)
		}
		if !seen {
			hits = append(hits, u)
		}
	}
	if len(hits) == 0 {
		return time.Time{}, 0
	}
	return hits[0], len(hits)
}

// dateTimeText formats a wall time with the picker's date layout and clock style.
func dateTimeText(c *ui.Context, w time.Time, layout string, cfg timeConfig) string {
	return w.Format(layout) + " " + timeText(c, TimeOfDay{w.Hour(), w.Minute(), w.Second()}, cfg) + " " + w.Format("MST")
}
