package mujica

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// NumberInputOptions configure a NumberInput.
type NumberInputOptions struct {
	// Min and Max bound the value; equal bounds (the zero value) mean no bounds.
	Min, Max float64
	// Step is what Up, Down and the buttons add; zero means 1.
	Step float64
	// Precision is the decimals kept and shown; zero takes the decimals of Step.
	Precision int
	// AllowEmpty lets the field be empty, which sets the value to NaN.
	AllowEmpty bool
	// StepperOnly shows the value with step buttons and no typing; it needs bounds.
	StepperOnly bool
	Placeholder string
	// Label names the field for assistive technology.
	Label    string
	Disabled bool
	// Error marks the value invalid and shows the message below.
	Error string
}

// numberInputSpec is NumberInputOptions resolved.
type numberInputSpec struct {
	lo, hi, step float64
	decimals     int
}

func (s numberInputSpec) format(v float64) string {
	if math.IsNaN(v) {
		return ""
	}
	return strconv.FormatFloat(v, 'f', s.decimals, 64)
}

// clamp rounds v to the precision, then bounds it; the bounds sit on that grid.
func (s numberInputSpec) clamp(v float64) float64 {
	p := math.Pow(10, float64(s.decimals))
	return math.Max(s.lo, math.Min(s.hi, math.Round(v*p)/p))
}

// numberInputSame reports equal values, counting NaN (empty) as equal.
func numberInputSame(a, b float64) bool { return a == b || math.IsNaN(a) && math.IsNaN(b) }

func numberInputResolve(opts NumberInputOptions) numberInputSpec {
	s := numberInputSpec{lo: opts.Min, hi: opts.Max, step: opts.Step, decimals: opts.Precision}
	if s.lo > s.hi {
		panic("mujica: NumberInput Min is above Max")
	}
	if s.step < 0 || opts.Precision < 0 {
		panic("mujica: NumberInput Step and Precision cannot be negative")
	}
	if s.lo == s.hi {
		if opts.StepperOnly {
			panic("mujica: NumberInput StepperOnly needs Min < Max")
		}
		s.lo, s.hi = math.Inf(-1), math.Inf(1)
	}
	if s.step == 0 {
		s.step = 1
	}
	if s.decimals == 0 {
		if t := strconv.FormatFloat(s.step, 'f', -1, 64); strings.Contains(t, ".") {
			s.decimals = len(t) - strings.IndexByte(t, '.') - 1
		}
	}
	// Bounds move inward onto the precision grid, so clamped values keep both.
	p := math.Pow(10, float64(s.decimals))
	s.lo = math.Ceil(math.Round(s.lo*p*1e6)/1e6) / p
	s.hi = math.Floor(math.Round(s.hi*p*1e6)/1e6) / p
	if s.lo > s.hi {
		panic("mujica: NumberInput has no value in range at this precision")
	}
	return s
}

// NumberInput edits *value as a number; Up and Down step it.
func NumberInput(c *ui.Context, value *float64, opts NumberInputOptions) InputResult {
	if value == nil {
		panic("mujica: NumberInput needs a value")
	}
	if opts.StepperOnly && opts.AllowEmpty {
		panic("mujica: NumberInput StepperOnly cannot be empty")
	}
	if math.IsNaN(*value) && !opts.AllowEmpty {
		panic("mujica: NumberInput value is NaN without AllowEmpty")
	}
	s := numberInputResolve(opts)
	if opts.StepperOnly {
		return numberInputStepper(c, value, s, opts)
	}
	var r InputResult
	set := func(v float64) {
		if !math.IsNaN(v) {
			v = s.clamp(v)
		}
		if v != *value && !(math.IsNaN(v) && math.IsNaN(*value)) {
			*value = v
			r.changed = true
		}
	}
	stepBy := func(n float64) {
		base := *value
		if math.IsNaN(base) {
			base = math.Max(s.lo, math.Min(s.hi, 0))
		}
		set(base + n*s.step)
	}
	look := &textInputLook{disabled: opts.Disabled, err: opts.Error}
	r.Element, r.Input = textInputFrame(c, look, func(frame *ui.Element) *ui.Element {
		text := ui.Local(frame, "text", func() string { return s.format(*value) })
		invalid := ui.Local(frame, "invalid", func() string { return "" })
		seen := ui.Local(frame, "seen", func() float64 { return *value })
		if !numberInputSame(*seen, *value) {
			// Set from outside: the text follows, and a typing error goes.
			*text, *invalid = s.format(*value), ""
		}
		in := textInputEditor(c, frame, text, opts.Disabled, ui.TextInputBase).Grow(1).MinWidth(Density(c).Unit() * 12).
			Placeholder(opts.Placeholder).Label(opts.Label)
		if in.Changed() && !opts.Disabled {
			*invalid = numberInputParse(c, s, *text, opts.AllowEmpty, set)
		}
		switch {
		case opts.Disabled:
		case in.Shortcut(0, ui.KeyUp):
			stepBy(1)
			*text, *invalid = s.format(*value), ""
		case in.Shortcut(0, ui.KeyDown):
			stepBy(-1)
			*text, *invalid = s.format(*value), ""
		}
		if !in.Focused() && *invalid == "" {
			// Leaving shows the value rounded; invalid text stays to be fixed.
			*text = s.format(*value)
		}
		defer func() { *seen = *value }()
		for _, b := range []struct {
			name  string
			icon  string
			delta float64
			can   bool
		}{
			{msg(c, "number_input.decrease"), "minus", -1, math.IsNaN(*value) || *value > s.lo},
			{msg(c, "number_input.increase"), "plus", 1, math.IsNaN(*value) || *value < s.hi},
		} {
			if numberInputButton(c, b.icon, b.name, b.can) {
				stepBy(b.delta)
				*text, *invalid = s.format(*value), ""
				in.Focus()
			}
		}
		if look.err == "" {
			look.err = *invalid
		}
		return in
	}, nil)
	if r.changed {
		slog.Debug("mujica: number_input change", "value", *value)
	}
	return r
}

// numberInputParse applies typed text and returns why it is invalid, if it is.
func numberInputParse(c *ui.Context, s numberInputSpec, text string, allowEmpty bool, set func(float64)) string {
	t := strings.TrimSpace(text)
	if t == "" && allowEmpty {
		set(math.NaN())
		return ""
	}
	v, err := strconv.ParseFloat(t, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return msg(c, "number_input.invalid")
	}
	if v < s.lo || v > s.hi {
		return fmt.Sprintf(msg(c, "number_input.range"), s.format(s.lo), s.format(s.hi))
	}
	set(v)
	return ""
}

// numberInputButton builds a pointer-only step button; Up and Down are its keys.
func numberInputButton(c *ui.Context, icon, name string, can bool) bool {
	k := Tokens(c)
	b := ui.Box(c).Size(theme.MinHitSize, theme.MinHitSize).Shrink(0).Center().Radius(theme.ControlRadius).
		Role(ui.RoleButton).Label(name).Cursor(ui.CursorPointer).Disabled(!can)
	fg := k.TextMuted
	if b.Pressed() {
		b.Background(k.Selection)
		fg = k.Text
	} else if b.Hovered() {
		b.Background(k.SurfaceHover)
		fg = k.Text
	}
	b.Children(func() { ui.Icon(c, icons.Must(icon)).Size(14, 14).TextColor(fg) })
	return b.Clicked()
}

// numberInputStepper shows the value beside MyGo's stepper, restyled.
func numberInputStepper(c *ui.Context, value *float64, s numberInputSpec, opts NumberInputOptions) InputResult {
	k := Tokens(c)
	d := Density(c)
	var r InputResult
	r.Element = ui.Column(c).Gap(d.Unit()).Shrink(0).Disabled(opts.Disabled)
	r.Element.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).MinHeight(d.ControlHeight()).Children(func() {
			ui.Text(c, s.format(*value)).FontFeatures("tnum").MinWidth(d.Unit() * 8).TextAlign(ui.End)
			// MyGo's stepper arrows ignore Disabled: give them a copy.
			target := value
			if opts.Disabled {
				held := *value
				target = &held
			}
			st := ui.Stepper(c, target, s.lo, s.hi, s.step).Label(opts.Label).Width(d.Unit()*6).
				Radius(theme.ControlRadius).Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder)
			if opts.Error != "" {
				st.Error(opts.Error).Border(theme.BorderWidth, k.Danger)
			}
			internal.Ring(st, k.Focus, theme.ControlRadius)
			r.Input = st
			r.changed = st.Changed() && !opts.Disabled
		})
		if opts.Error != "" {
			ui.Text(c, opts.Error).TextColor(k.Danger).FontSize(theme.CaptionSize)
		}
	})
	if r.changed {
		*value = s.clamp(*value)
		slog.Debug("mujica: number_input step", "value", *value)
	}
	return r
}
