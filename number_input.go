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
	// lo and hi are the bounds on the precision grid; min and max the original ones.
	lo, hi, min, max, step float64
	decimals               int
}

func (s numberInputSpec) format(v float64) string {
	if math.IsNaN(v) {
		return ""
	}
	return strconv.FormatFloat(v, 'f', s.decimals, 64)
}

// round rounds v to the precision.
func (s numberInputSpec) round(v float64) float64 {
	p := math.Pow(10, float64(s.decimals))
	return math.Round(v*p)/p + 0 // + 0 turns -0 into 0
}

// clamp rounds v to the precision, then bounds it; the bounds sit on that grid.
func (s numberInputSpec) clamp(v float64) float64 {
	return math.Max(s.lo, math.Min(s.hi, s.round(v)))
}

// stepBy moves v n grid steps from Min, exactly checked against the original bounds.
func (s numberInputSpec) stepBy(v float64, n int) float64 {
	if math.IsNaN(v) {
		// From empty, the first step lands next to zero.
		v = s.clamp(0)
	}
	base := s.lo
	if math.IsInf(base, 0) {
		base = 0
	}
	at := func(k float64) float64 { return s.round(base + k*s.step) }
	k := math.Round((v - base) / s.step)
	// Off the grid, the nearest point in the direction is the first step.
	if n > 0 && at(k) <= v || n < 0 && at(k) >= v {
		k += float64(n)
	}
	r := at(k)
	switch {
	case r > s.max:
		r = at(k - 1)
	case r < s.min:
		r = at(k + 1)
	}
	if r > s.max || r < s.min {
		return v
	}
	return r
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
	// Bounds move inward onto the precision grid, compared exactly with the originals.
	s.min, s.max = s.lo, s.hi
	p := math.Pow(10, float64(s.decimals))
	if s.lo = s.round(s.lo); s.lo < s.min {
		s.lo = s.round(s.lo + 1/p)
	}
	if s.hi = s.round(s.hi); s.hi > s.max {
		s.hi = s.round(s.hi - 1/p)
	}
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
	stepBy := func(n int) { set(s.stepBy(*value, n)) }
	look := &textInputLook{disabled: opts.Disabled, err: opts.Error}
	r.Element, r.Input = textInputFrame(c, look, func(frame *ui.Element) *ui.Element {
		text := ui.Local(frame, "text", func() string { return s.format(*value) })
		invalid := ui.Local(frame, "invalid", func() string { return "" })
		seen := ui.Local(frame, "seen", func() float64 { return *value })
		if !numberInputSame(*seen, *value) {
			// Set from outside: the text follows, and a typing error goes.
			*text, *invalid = s.format(*value), ""
		}
		in := textInputEditor(c, frame, text, look.disabled, ui.TextInputBase).Grow(1).MinWidth(Density(c).Unit() * 12).
			Placeholder(opts.Placeholder).Label(opts.Label)
		if in.Changed() && !look.disabled {
			*invalid = numberInputParse(c, s, *text, opts.AllowEmpty, set)
		}
		switch {
		case look.disabled:
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
			delta int
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

// numberInputStepper shows the value beside up and down arrows; it steps itself.
func numberInputStepper(c *ui.Context, value *float64, s numberInputSpec, opts NumberInputOptions) InputResult {
	k := Tokens(c)
	d := Density(c)
	var r InputResult
	r.Element = ui.Column(c).Gap(d.Unit()).Shrink(0).Disabled(opts.Disabled)
	disabled := r.Element.IsDisabled()
	set := func(v float64) {
		if v != *value && !disabled {
			*value = v
			r.changed = true
		}
	}
	r.Element.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).MinHeight(d.ControlHeight()).Children(func() {
			ui.Text(c, s.format(*value)).FontFeatures("tnum").MinWidth(d.Unit() * 8).TextAlign(ui.End)
			border := k.ControlBorder
			if opts.Error != "" {
				border = k.Danger
			}
			st := ui.Column(c).Focusable().Shrink(0).Width(d.Unit()*6).Radius(theme.ControlRadius).Clip().
				Background(k.Surface).Border(theme.BorderWidth, border).Role(ui.RoleStepper).Label(opts.Label).
				Description(s.format(*value)).Error(opts.Error)
			internal.Ring(st, k.Focus, theme.ControlRadius)
			switch {
			case st.Shortcut(0, ui.KeyUp), st.Shortcut(0, ui.KeyRight):
				set(s.stepBy(*value, 1))
			case st.Shortcut(0, ui.KeyDown), st.Shortcut(0, ui.KeyLeft):
				set(s.stepBy(*value, -1))
			case st.Shortcut(0, ui.KeyHome):
				set(s.lo)
			case st.Shortcut(0, ui.KeyEnd):
				set(s.hi)
			}
			st.Children(func() {
				if numberInputArrow(c, "chevron-up", *value < s.hi, false) {
					set(s.stepBy(*value, 1))
					st.Focus()
				}
				if numberInputArrow(c, "chevron-down", *value > s.lo, true) {
					set(s.stepBy(*value, -1))
					st.Focus()
				}
			})
			r.Input = st
		})
		if opts.Error != "" {
			ui.Text(c, opts.Error).TextColor(k.Danger).FontSize(FontSize(c, theme.CaptionSize))
		}
	})
	if r.changed {
		slog.Debug("mujica: number_input step", "value", *value)
	}
	return r
}

// numberInputArrow builds one half of the stepper and reports a click.
func numberInputArrow(c *ui.Context, icon string, can, lower bool) bool {
	k := Tokens(c)
	h := Density(c).ControlHeight()/2 - theme.BorderWidth
	a := ui.Box(c).Height(h).Center().Role(ui.RoleNone).Cursor(ui.CursorPointer).Disabled(!can)
	if lower {
		a.BorderWidth(theme.BorderWidth, 0, 0, 0).BorderColor(k.Border)
	}
	fg := k.TextMuted
	if a.Pressed() {
		a.Background(k.Selection)
		fg = k.Text
	} else if a.Hovered() {
		a.Background(k.SurfaceHover)
		fg = k.Text
	}
	a.Children(func() { ui.Icon(c, icons.Must(icon)).Size(12, 12).TextColor(fg) })
	return a.Clicked()
}
