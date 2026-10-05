package mujica

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// TimeOfDay is a wall-clock time with no date and no location.
type TimeOfDay struct {
	Hour, Minute, Second int
}

// Valid reports whether t is a time between 00:00:00 and 23:59:59.
func (t TimeOfDay) Valid() bool {
	return t.Hour >= 0 && t.Hour < 24 && t.Minute >= 0 && t.Minute < 60 && t.Second >= 0 && t.Second < 60
}

// String formats t as 15:04:05.
func (t TimeOfDay) String() string { return fmt.Sprintf("%02d:%02d:%02d", t.Hour, t.Minute, t.Second) }

// Compare returns -1, 0 or 1 as t is before, equal to or after o.
func (t TimeOfDay) Compare(o TimeOfDay) int {
	a, b := t.seconds(), o.seconds()
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func (t TimeOfDay) seconds() int { return t.Hour*3600 + t.Minute*60 + t.Second }

// TimePickerOptions configure a TimePicker.
type TimePickerOptions struct {
	// Seconds shows and edits the seconds.
	Seconds bool
	// Hour12 shows a 12-hour clock with AM and PM.
	Hour12 bool
	// MinuteStep and SecondStep must divide 60; zero means 1.
	MinuteStep, SecondStep int
	// Min and Max bound the time; a zero Max means 23:59:59.
	Min, Max TimeOfDay
	Disabled bool
	// Label names the field for assistive technology.
	Label string
}

// TimePicker edits *value by typing into hour, minute and second segments or by choosing in a popup.
// Times outside the range or off the step are reported and not written.
func TimePicker(c *ui.Context, value *TimeOfDay, opts TimePickerOptions) DateTimeResult {
	cfg := timePickerConfig("TimePicker", opts)
	if !value.Valid() {
		panic(fmt.Sprintf("mujica: TimePicker value %v is not a time of day", *value))
	}
	col := ui.Column(c).Gap(Density(c).Unit()).AlignItems(ui.Start).Disabled(opts.Disabled)
	draft := ui.Local(col, "time-picker-draft", func() TimeOfDay { return *value })
	seen := ui.Local(col, "time-picker-seen", func() TimeOfDay { return *value })
	problem := ui.Local(col, "time-picker-error", func() string { return "" })
	if *seen != *value {
		*draft, *seen, *problem = *value, *value, ""
	}
	changed := false
	col.Children(func() {
		cfg.invalid = *problem != ""
		if _, edited := timeEditor(c, draft, cfg); edited {
			*problem = timePickerProblem(c, *draft, cfg)
			if *problem == "" && *draft != *value {
				slog.Debug("mujica: time_picker change", "from", *value, "to", *draft)
				*value, *seen, changed = *draft, *draft, true
			}
		}
		datetimeError(c, *problem)
	})
	return DateTimeResult{Element: col, changed: changed}
}

// timeConfig is the resolved behavior of a time editor.
type timeConfig struct {
	seconds, hour12        bool
	minuteStep, secondStep int
	min, max               TimeOfDay
	label                  string
	invalid                bool
}

// timePickerConfig checks opts and fills in defaults.
func timePickerConfig(name string, opts TimePickerOptions) timeConfig {
	cfg := timeConfig{seconds: opts.Seconds, hour12: opts.Hour12, minuteStep: max(opts.MinuteStep, 1),
		secondStep: max(opts.SecondStep, 1), min: opts.Min, max: opts.Max, label: opts.Label}
	for _, s := range []int{opts.MinuteStep, opts.SecondStep} {
		if s < 0 || s > 0 && 60%s != 0 {
			panic(fmt.Sprintf("mujica: %s step %d does not divide 60", name, s))
		}
	}
	if cfg.max == (TimeOfDay{}) {
		cfg.max = TimeOfDay{23, 59, 59}
	}
	if !cfg.min.Valid() || !cfg.max.Valid() || cfg.min.Compare(cfg.max) > 0 {
		panic(fmt.Sprintf("mujica: %s range %v..%v is not valid", name, cfg.min, cfg.max))
	}
	return cfg
}

// timePickerProblem explains why t cannot be written, or returns "".
func timePickerProblem(c *ui.Context, t TimeOfDay, cfg timeConfig) string {
	switch {
	case t.Minute%cfg.minuteStep != 0:
		return fmt.Sprintf(msg(c, "time_picker.minute_step"), cfg.minuteStep)
	case cfg.seconds && t.Second%cfg.secondStep != 0:
		return fmt.Sprintf(msg(c, "time_picker.second_step"), cfg.secondStep)
	case t.Compare(cfg.min) < 0 || t.Compare(cfg.max) > 0:
		return fmt.Sprintf(msg(c, "time_picker.out_of_range"), timeText(c, cfg.min, cfg), timeText(c, cfg.max, cfg))
	}
	return ""
}

// timeText formats t for messages, honoring the 12-hour clock and seconds.
func timeText(c *ui.Context, t TimeOfDay, cfg timeConfig) string {
	s := fmt.Sprintf("%02d:%02d", t.Hour, t.Minute)
	if cfg.hour12 {
		s = fmt.Sprintf("%d:%02d", timeHour12(t.Hour), t.Minute)
	}
	if cfg.seconds {
		s += fmt.Sprintf(":%02d", t.Second)
	}
	if cfg.hour12 {
		s = fmt.Sprintf(msg(c, "time_picker.with_period"), s, timePeriod(c, t.Hour))
	}
	return s
}

func timeHour12(h int) int {
	if h%12 == 0 {
		return 12
	}
	return h % 12
}

func timePeriod(c *ui.Context, h int) string {
	if h < 12 {
		return msg(c, "time_picker.am")
	}
	return msg(c, "time_picker.pm")
}

// timePart is one segment of a time editor.
type timePart int

const (
	timeHourPart timePart = iota
	timeMinutePart
	timeSecondPart
	timePeriodPart
)

var timePartNames = [...]string{"time_picker.hours", "time_picker.minutes", "time_picker.seconds", "time_picker.period"}

func timeParts(cfg timeConfig) []timePart {
	parts := []timePart{timeHourPart, timeMinutePart}
	if cfg.seconds {
		parts = append(parts, timeSecondPart)
	}
	if cfg.hour12 {
		parts = append(parts, timePeriodPart)
	}
	return parts
}

// timeEditor draws a segmented time field with a popup and reports whether the user edited *draft.
func timeEditor(c *ui.Context, draft *TimeOfDay, cfg timeConfig) (field *ui.Element, edited bool) {
	k := Tokens(c)
	d := Density(c)
	field = datetimeField(c, ui.Row(c), cfg.invalid).Gap(d.Unit() * 2)
	open := ui.Local(field, "time-open", func() bool { return false })
	fresh := ui.Local(field, "time-fresh", func() bool { return false })
	field.Children(func() {
		segs := ui.Row(c).AlignItems(ui.Center).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(cfg.label)
		segs.Children(func() {
			focusNext := false
			for i, part := range timeParts(cfg) {
				if i > 0 && part != timePeriodPart {
					ui.Text(c, ":").TextColor(k.TextMuted).Role(ui.RoleNone)
				}
				if part == timePeriodPart {
					ui.Box(c).Width(d.Unit())
				}
				seg, done := timeSegment(c, draft, part, cfg, &edited)
				if focusNext {
					seg.Focus()
					focusNext = false
				}
				focusNext = done
			}
		})
		b := datetimeIconButton(c, "clock", msg(c, "time_picker.open"))
		if b.Clicked() {
			*open, *fresh = !*open, true
		}
	})
	if !*open && field.Shortcut(ui.Alt, ui.KeyDown) {
		*open, *fresh = true, true
	}
	datetimePopup(c, field, open, func() {
		if timeColumns(c, draft, cfg, *fresh) {
			edited = true
		}
		*fresh = false
	})
	return field, edited
}

// timeSegment draws one focusable segment; done reports a fully typed value.
func timeSegment(c *ui.Context, draft *TimeOfDay, part timePart, cfg timeConfig, edited *bool) (seg *ui.Element, done bool) {
	k := Tokens(c)
	d := Density(c)
	name := msg(c, timePartNames[part])
	if cfg.label != "" {
		name = cfg.label + ": " + name
	}
	seg = ui.Box(c).Focusable().FocusRing(false).Role(ui.RoleStepper).Label(name).
		Padding(d.Unit()/2, d.Unit()).Radius(theme.ControlRadius).MinWidth(theme.BodySize * 1.6).Center()
	typed := ui.Local(seg, "time-typed", func() string { return "" })
	if !seg.Focused() {
		*typed = ""
	}
	set := func(t TimeOfDay) {
		if t != *draft {
			*draft = t
			*edited = true
		}
	}
	t := *draft
	step := func(dir int) {
		switch part {
		case timeHourPart:
			t.Hour = (t.Hour + dir + 24) % 24
		case timeMinutePart:
			t.Minute = timeStep(t.Minute, cfg.minuteStep, dir)
		case timeSecondPart:
			t.Second = timeStep(t.Second, cfg.secondStep, dir)
		case timePeriodPart:
			t.Hour = (t.Hour + 12) % 24
		}
		set(t)
	}
	switch {
	case seg.Shortcut(0, ui.KeyUp):
		step(1)
	case seg.Shortcut(0, ui.KeyDown):
		step(-1)
	}
	if part == timePeriodPart {
		if seg.Shortcut(0, ui.KeyA) && t.Hour >= 12 || seg.Shortcut(0, ui.KeyP) && t.Hour < 12 {
			t.Hour = (t.Hour + 12) % 24
			set(t)
		}
	} else {
		for n := range 10 {
			if seg.Shortcut(0, ui.Key0+ui.Key(n)) {
				done = timeType(typed, n, part, cfg, &t)
				set(t)
			}
		}
	}
	text := ""
	switch part {
	case timeHourPart:
		text = fmt.Sprintf("%02d", t.Hour)
		if cfg.hour12 {
			text = fmt.Sprintf("%2d", timeHour12(t.Hour))
		}
	case timeMinutePart:
		text = fmt.Sprintf("%02d", t.Minute)
	case timeSecondPart:
		text = fmt.Sprintf("%02d", t.Second)
	case timePeriodPart:
		text = timePeriod(c, t.Hour)
	}
	fg := k.Text
	switch {
	case seg.Focused():
		seg.Background(k.Accent)
		fg = k.OnAccent
	case seg.Hovered():
		seg.Background(k.SurfaceHover)
	}
	seg.Children(func() { ui.Text(c, text).TextColor(fg).FontFeatures("tnum").SingleLine() })
	return seg, done
}

// timeStep moves v to the next multiple of step in dir, wrapping within an hour.
func timeStep(v, step, dir int) int {
	if dir > 0 {
		return (v/step*step + step) % 60
	}
	if v%step != 0 {
		return v / step * step
	}
	return (v - step + 60) % 60
}

// timeType adds digit n to a segment's typed digits and reports when the segment is complete.
func timeType(typed *string, n int, part timePart, cfg timeConfig, t *TimeOfDay) bool {
	lo, hi := 0, 59
	if part == timeHourPart {
		hi = 23
		if cfg.hour12 {
			lo, hi = 1, 12
		}
	}
	s := *typed + fmt.Sprint(n)
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	if len(s) > 2 || v > hi {
		s, v = fmt.Sprint(n), n
	}
	*typed = s
	if v >= lo {
		switch part {
		case timeHourPart:
			if cfg.hour12 {
				v = v%12 + t.Hour/12*12
			}
			t.Hour = v
		case timeMinutePart:
			t.Minute = v
		case timeSecondPart:
			t.Second = v
		}
	}
	if len(s) == 2 || n*10 > hi && len(s) == 1 {
		*typed = ""
		return true
	}
	return false
}

// timeColumns lists choosable values per segment and reports whether a choice edited *draft.
func timeColumns(c *ui.Context, draft *TimeOfDay, cfg timeConfig, fresh bool) (edited bool) {
	k := Tokens(c)
	d := Density(c)
	h := d.ControlHeight() - 4
	within := func(lo, hi TimeOfDay) bool { return hi.Compare(cfg.min) >= 0 && lo.Compare(cfg.max) <= 0 }
	ui.Row(c).Gap(d.Unit()).Children(func() {
		for _, part := range timeParts(cfg) {
			t := *draft
			var options []TimeOfDay
			switch part {
			case timeHourPart:
				for v := range 24 {
					if !cfg.hour12 || v/12 == t.Hour/12 {
						options = append(options, TimeOfDay{v, t.Minute, t.Second})
					}
				}
			case timeMinutePart:
				for v := 0; v < 60; v += cfg.minuteStep {
					options = append(options, TimeOfDay{t.Hour, v, t.Second})
				}
			case timeSecondPart:
				for v := 0; v < 60; v += cfg.secondStep {
					options = append(options, TimeOfDay{t.Hour, t.Minute, v})
				}
			case timePeriodPart:
				options = []TimeOfDay{{t.Hour % 12, t.Minute, t.Second}, {t.Hour%12 + 12, t.Minute, t.Second}}
			}
			name := msg(c, timePartNames[part])
			ui.Scroll(c).Height(h * 6).Width(d.ControlHeight() * 1.5).Children(func() {
				ui.Column(c).FocusGroup(ui.Vertical).Role(ui.RoleRadioGroup).Label(name).Children(func() {
					for _, o := range options {
						var lo, hi TimeOfDay
						text := ""
						switch part {
						case timeHourPart:
							lo, hi = TimeOfDay{o.Hour, 0, 0}, TimeOfDay{o.Hour, 59, 59}
							text = fmt.Sprintf("%02d", o.Hour)
							if cfg.hour12 {
								text = fmt.Sprintf("%d", timeHour12(o.Hour))
							}
						case timeMinutePart:
							lo, hi = TimeOfDay{o.Hour, o.Minute, 0}, TimeOfDay{o.Hour, o.Minute, 59}
							text = fmt.Sprintf("%02d", o.Minute)
						case timeSecondPart:
							lo, hi = o, o
							text = fmt.Sprintf("%02d", o.Second)
						case timePeriodPart:
							lo, hi = TimeOfDay{o.Hour / 12 * 12, 0, 0}, TimeOfDay{o.Hour/12*12 + 11, 59, 59}
							text = timePeriod(c, o.Hour)
						}
						sel := t
						item := ui.RadioBase(c, &sel, o).Height(h).FillWidth().Center().Radius(theme.ControlRadius).
							Disabled(!within(lo, hi)).Label(text)
						fg := k.Text
						switch {
						case o == t:
							item.Background(k.Accent)
							fg = k.OnAccent
							if fresh {
								item.ScrollIntoView()
							}
						case item.Pressed():
							item.Background(k.Selection)
						case item.Hovered():
							item.Background(k.SurfaceHover)
						}
						internal.Ring(item, k.Focus, theme.ControlRadius)
						item.Children(func() { ui.Text(c, text).TextColor(fg).FontFeatures("tnum") })
						if item.Clicked() && o != *draft {
							*draft, edited = o, true
						}
					}
				})
			})
		}
	})
	return edited
}
