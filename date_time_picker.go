package mujica

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// DateTimePickerOptions configure a DateTimePicker.
type DateTimePickerOptions struct {
	// Location is the time zone the user edits in. Required.
	Location *time.Location
	// Min and Max bound the instant; zero leaves a side open.
	Min, Max   time.Time
	Seconds    bool
	Hour12     bool
	MinuteStep int
	WeekStart  WeekStart
	// Today is marked in the calendar; zero uses the system clock.
	Today    Date
	Disabled bool
	// Label names the field for assistive technology.
	Label string
}

// DateTimePicker edits *value as a date and a wall-clock time in opts.Location, confirmed together.
// Times a DST change skips or repeats are reported, never shifted.
func DateTimePicker(c *ui.Context, value *time.Time, opts DateTimePickerOptions) DateTimeResult {
	loc := opts.Location
	if loc == nil {
		panic("mujica: DateTimePicker needs a Location")
	}
	if !opts.Min.IsZero() && !opts.Max.IsZero() && opts.Min.After(opts.Max) {
		panic(fmt.Sprintf("mujica: DateTimePicker Min %v is after Max %v", opts.Min, opts.Max))
	}
	cfg := timePickerConfig("DateTimePicker", TimePickerOptions{Seconds: opts.Seconds, Hour12: opts.Hour12, MinuteStep: opts.MinuteStep})
	cfg.label = msg(c, "date_time_picker.time")
	k := Tokens(c)
	d := Density(c)
	today := datetimeToday(c, opts.Today)
	var lo, hi Date
	if !opts.Min.IsZero() {
		lo = DateOf(opts.Min.In(loc))
	}
	if !opts.Max.IsZero() {
		hi = DateOf(opts.Max.In(loc))
	}
	col := ui.Column(c).Width(d.ControlHeight() * 9).Gap(d.Unit()).AlignItems(ui.Stretch).Disabled(opts.Disabled)
	open := ui.Local(col, "date-time-open", func() bool { return false })
	wasOpen := ui.Local(col, "date-time-was-open", func() bool { return false })
	day := ui.Local(col, "date-time-day", func() Date { return Date{} })
	clock := ui.Local(col, "date-time-clock", func() TimeOfDay { return TimeOfDay{} })
	cursor := ui.Local(col, "date-time-cursor", func() Date { return Date{} })
	problem := ui.Local(col, "date-time-error", func() string { return "" })
	format := func(t time.Time) string { return dateTimeText(c, t.In(loc), cfg) }
	changed, cancelled, done := false, false, false
	label := opts.Label
	if label == "" {
		label = msg(c, "date_time_picker.placeholder")
	}
	text := ""
	if !value.IsZero() {
		text = format(*value)
	}
	col.Children(func() {
		field := datetimeField(c, ui.Row(c), false)
		field.Children(func() {
			t := datetimeTrigger(c, text, msg(c, "date_time_picker.placeholder"), "calendar").Label(label).Description(text)
			if t.Clicked() || !*open && t.Shortcut(0, ui.KeyDown) {
				*open = !*open
				if *open {
					*day, *clock, *problem = today, TimeOfDay{}, ""
					if !value.IsZero() {
						w := value.In(loc)
						*day, *clock = DateOf(w), TimeOfDay{w.Hour(), w.Minute(), w.Second()}
					}
					*cursor = *day
					slog.Debug("mujica: date_time_picker open")
				}
			}
		})
		datetimePopup(c, field, open, func() {
			picked, _ := calendarMonth(c, cursor, calendarConfig{
				chosen:    func(x Date) bool { return x == *day },
				enabled:   calendarEnabled(lo, hi, nil),
				today:     today,
				weekStart: opts.WeekStart,
				autoFocus: true,
			})
			if !picked.IsZero() {
				*day, *problem = picked, ""
			}
			ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
				if _, edited := timeEditor(c, clock, cfg); edited {
					*problem = ""
				}
				ui.Text(c, loc.String()).FontSize(theme.CaptionSize).TextColor(k.TextMuted)
			})
			datetimeError(c, *problem)
			datetimeFooter(c, func() {
				if Button(c, msg(c, "calendar.cancel"), ButtonOptions{Variant: Secondary}).Clicked() {
					*open = false
				}
				if Button(c, msg(c, "calendar.confirm"), ButtonOptions{}).Clicked() {
					at, p := dateTimeConfirm(c, *day, *clock, loc, opts, cfg)
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
		cancelled = true
	}
	*wasOpen = *open
	return DateTimeResult{Element: col, changed: changed, cancelled: cancelled}
}

// dateTimeConfirm turns the draft into an instant, or explains why it cannot.
func dateTimeConfirm(c *ui.Context, d Date, t TimeOfDay, loc *time.Location, opts DateTimePickerOptions, cfg timeConfig) (time.Time, string) {
	if p := timePickerProblem(c, t, cfg); p != "" {
		return time.Time{}, p
	}
	wall := fmt.Sprintf("%v %s", d, timeText(c, t, cfg))
	at, n := dateTimeResolve(d, t, loc)
	switch {
	case n == 0:
		slog.Debug("mujica: date_time_picker gap", "wall", wall, "zone", loc.String())
		return at, fmt.Sprintf(msg(c, "date_time_picker.gap"), wall, loc.String())
	case n > 1:
		slog.Debug("mujica: date_time_picker ambiguous", "wall", wall, "zone", loc.String())
		return at, fmt.Sprintf(msg(c, "date_time_picker.ambiguous"), wall, loc.String())
	case !opts.Min.IsZero() && at.Before(opts.Min):
		return at, fmt.Sprintf(msg(c, "calendar.before_min"), dateTimeText(c, opts.Min.In(loc), cfg))
	case !opts.Max.IsZero() && at.After(opts.Max):
		return at, fmt.Sprintf(msg(c, "calendar.after_max"), dateTimeText(c, opts.Max.In(loc), cfg))
	}
	return at, ""
}

// dateTimeResolve finds the instants showing wall time d t in loc: none in a DST gap, two in an overlap.
func dateTimeResolve(d Date, t TimeOfDay, loc *time.Location) (time.Time, int) {
	wall := time.Date(d.Year, d.Month, d.Day, t.Hour, t.Minute, t.Second, 0, time.UTC)
	guess := time.Date(d.Year, d.Month, d.Day, t.Hour, t.Minute, t.Second, 0, loc)
	var hits []time.Time
	for _, probe := range []time.Time{guess.Add(-24 * time.Hour), guess, guess.Add(24 * time.Hour)} {
		_, off := probe.Zone()
		u := wall.Add(-time.Duration(off) * time.Second).In(loc)
		y, m, dd := u.Date()
		if (Date{y, m, dd}) != d || u.Hour() != t.Hour || u.Minute() != t.Minute || u.Second() != t.Second {
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

// dateTimeText formats a wall time with the picker's clock style.
func dateTimeText(c *ui.Context, w time.Time, cfg timeConfig) string {
	return DateOf(w).String() + " " + timeText(c, TimeOfDay{w.Hour(), w.Minute(), w.Second()}, cfg) + " " + w.Format("MST")
}
