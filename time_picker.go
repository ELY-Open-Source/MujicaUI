package mujica

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

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
	// ReadOnly keeps the segments focusable and copyable but blocks edits and the popup.
	ReadOnly bool
	// Label names the field for assistive technology.
	Label string
}

// TimePicker edits *value in segments or a popup.
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
		*draft, *seen = *value, *value
	}
	changed := false
	col.Children(func() {
		cfg.problem = *problem
		_, edited, bad := timeEditor(c, draft, cfg)
		p := bad
		if p == "" {
			p = timePickerProblem(c, *draft, cfg)
		}
		if p != *problem {
			*problem = p
			c.Invalidate()
		}
		if edited && p == "" && *draft != *value {
			slog.Debug("mujica: time_picker change", "from", *value, "to", *draft)
			*value, *seen, changed = *draft, *draft, true
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
	// problem is the error shown last frame, attached to the segments.
	problem  string
	readOnly bool
}

// timePickerConfig checks opts and fills in defaults.
func timePickerConfig(name string, opts TimePickerOptions) timeConfig {
	cfg := timeConfig{seconds: opts.Seconds, hour12: opts.Hour12, minuteStep: max(opts.MinuteStep, 1),
		secondStep: max(opts.SecondStep, 1), min: opts.Min, max: opts.Max, label: opts.Label, readOnly: opts.ReadOnly}
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

// timeEditor draws segments and popup; reports edits and invalid entry.
func timeEditor(c *ui.Context, draft *TimeOfDay, cfg timeConfig) (field *ui.Element, edited bool, bad string) {
	k := Tokens(c)
	d := Density(c)
	field = datetimeField(c, ui.Row(c), cfg.problem != "", cfg.readOnly).Gap(d.Unit() * 2)
	open := ui.Local(field, "time-open", func() bool { return false })
	fresh := ui.Local(field, "time-fresh", func() bool { return false })
	field.Children(func() {
		segs := ui.Row(c).AlignItems(ui.Center).Gap(d.Unit() * 1.5).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(cfg.label)
		if segs.Shortcut(ui.Cmd, ui.KeyC) {
			c.WriteClipboard(timeText(c, *draft, cfg))
		}
		segs.Children(func() {
			focusNext := false
			for i, part := range timeParts(cfg) {
				if i > 0 && part != timePeriodPart {
					ui.Text(c, ":").TextColor(k.TextMuted).Role(ui.RoleNone)
				}
				if part == timePeriodPart {
					ui.Box(c).Width(d.Unit())
				}
				seg, done, segBad := timeSegment(c, draft, part, cfg, &edited)
				if focusNext {
					seg.Focus()
				}
				focusNext = done
				if bad == "" {
					bad = segBad
				}
			}
		})
		if cfg.readOnly {
			return
		}
		b := datetimeIconButton(c, "clock", msg(c, "time_picker.open")).Role(ui.RolePopUpButton)
		if b.Clicked() {
			*open, *fresh = !*open, true
		}
	})
	if !*open && !cfg.readOnly && field.Shortcut(ui.Alt, ui.KeyDown) {
		*open, *fresh = true, true
	}
	datetimePopup(c, field, open, cfg.readOnly, func() {
		if timeColumns(c, draft, cfg, *fresh) {
			edited = true
		}
		*fresh = false
	})
	return field, edited, bad
}

// timeLimits returns the lowest and highest value a segment accepts.
func timeLimits(part timePart, cfg timeConfig) (lo, hi int) {
	switch {
	case part != timeHourPart:
		return 0, 59
	case cfg.hour12:
		return 1, 12
	}
	return 0, 23
}

// timeSegment draws one segment; done reports a complete valid entry, bad an invalid one.
func timeSegment(c *ui.Context, draft *TimeOfDay, part timePart, cfg timeConfig, edited *bool) (seg *ui.Element, done bool, bad string) {
	k := Tokens(c)
	d := Density(c)
	name := msg(c, timePartNames[part])
	if cfg.label != "" {
		name = cfg.label + ": " + name
	}
	seg = ui.Box(c).Focusable().Role(ui.RoleStepper).Label(name).
		Padding(0, d.Unit()).Radius(theme.ControlRadius).MinWidth(theme.BodySize * 1.6).Shrink(0).Center()
	internal.Ring(seg, k.Focus, theme.ControlRadius)
	typed := ui.Local(seg, "time-typed", func() string { return "" })
	badText := ui.Local(seg, "time-bad", func() string { return "" })
	badFor := ui.Local(seg, "time-bad-for", func() TimeOfDay { return TimeOfDay{} })
	if !seg.Focused() {
		*typed = ""
	}
	if *badFor != *draft {
		*badText = ""
	}
	set := func(t TimeOfDay) {
		*badText = ""
		if t != *draft {
			*draft = t
			*edited = true
		}
	}
	t := *draft
	lo, hi := timeLimits(part, cfg)
	if !cfg.readOnly {
		switch {
		case seg.Shortcut(0, ui.KeyUp):
			set(timeStepPart(t, part, cfg, 1))
		case seg.Shortcut(0, ui.KeyDown):
			set(timeStepPart(t, part, cfg, -1))
		}
		t = *draft
		if part == timePeriodPart {
			if seg.Shortcut(0, ui.KeyA) && t.Hour >= 12 || seg.Shortcut(0, ui.KeyP) && t.Hour < 12 {
				t.Hour = (t.Hour + 12) % 24
				set(t)
			}
		} else {
			for n := range 10 {
				if !seg.Shortcut(0, ui.Key0+ui.Key(n)) {
					continue
				}
				v, ok, complete := timeType(typed, n, lo, hi)
				switch {
				case ok:
					set(timeWithPart(t, part, v, cfg))
					done = complete
				case complete:
					*badText, *badFor = fmt.Sprintf("%02d", v), *draft
					slog.Debug("mujica: time_picker invalid segment", "part", int(part), "typed", v)
				}
			}
		}
	}
	t = *draft
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
	if *badText != "" {
		text = *badText
		bad = fmt.Sprintf(msg(c, "time_picker.part_range"), msg(c, timePartNames[part]), lo, hi)
	}
	seg.Description(timeSegmentDescription(c, part, text, lo, hi)).Error(cfg.problem)
	fg := k.Text
	switch {
	case *badText != "":
		fg = k.Danger
	case seg.Focused() && !cfg.readOnly:
		seg.Background(k.Accent)
		fg = k.OnAccent
	case seg.Focused():
		seg.Background(k.Selection)
	case seg.Hovered():
		seg.Background(k.SurfaceHover)
	}
	seg.Children(func() { ui.Text(c, text).TextColor(fg).FontFeatures("tnum").SingleLine() })
	return seg, done, bad
}

// timeSegmentDescription gives assistive technology a segment's value and range.
func timeSegmentDescription(c *ui.Context, part timePart, text string, lo, hi int) string {
	if part == timePeriodPart {
		return text
	}
	return fmt.Sprintf(msg(c, "time_picker.segment"), strings.TrimSpace(text), lo, hi)
}

// timeStepPart steps one segment of t by dir, wrapping round its cycle.
func timeStepPart(t TimeOfDay, part timePart, cfg timeConfig, dir int) TimeOfDay {
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
	return t
}

// timeWithPart sets one segment of t to v, keeping the period on a 12-hour clock.
func timeWithPart(t TimeOfDay, part timePart, v int, cfg timeConfig) TimeOfDay {
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
	return t
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

// timeType appends digit n; reports value, validity and completion.
func timeType(typed *string, n, lo, hi int) (v int, ok, complete bool) {
	s := *typed + strconv.Itoa(n)
	if len(s) > 2 {
		s = strconv.Itoa(n)
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	complete = len(s) == 2 || v*10 > hi
	if complete {
		s = ""
	}
	*typed = s
	return v, v >= lo && v <= hi, complete
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
