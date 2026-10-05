package input

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
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

// NumberInputSpec is NumberInputOptions resolved.
type NumberInputSpec struct {
	// lo and hi are the bounds on the precision grid; min and max the original ones.
	Lo, Hi, min, max, Step float64
	Decimals               int
}

func (s NumberInputSpec) format(v float64) string {
	if math.IsNaN(v) {
		return ""
	}
	return strconv.FormatFloat(v, 'f', s.Decimals, 64)
}

// Round rounds v to the precision.
func (s NumberInputSpec) Round(v float64) float64 {
	p := math.Pow(10, float64(s.Decimals))
	return math.Round(v*p)/p + 0 // + 0 turns -0 into 0
}

// clamp rounds v to the precision, then bounds it; the bounds sit on that grid.
func (s NumberInputSpec) clamp(v float64) float64 {
	return math.Max(s.Lo, math.Min(s.Hi, s.Round(v)))
}

// StepBy moves v n grid steps from Min, exactly checked against the original bounds.
func (s NumberInputSpec) StepBy(v float64, n int) float64 {
	if math.IsNaN(v) {
		// From empty, the first step lands next to zero.
		v = s.clamp(0)
	}
	base := s.Lo
	if math.IsInf(base, 0) {
		base = 0
	}
	at := func(k float64) float64 { return s.Round(base + k*s.Step) }
	k := math.Round((v - base) / s.Step)
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

func NumberInputResolve(opts NumberInputOptions) NumberInputSpec {
	s := NumberInputSpec{Lo: opts.Min, Hi: opts.Max, Step: opts.Step, Decimals: opts.Precision}
	if s.Lo > s.Hi {
		panic("mujica: NumberInput Min is above Max")
	}
	if s.Step < 0 || opts.Precision < 0 {
		panic("mujica: NumberInput Step and Precision cannot be negative")
	}
	if s.Lo == s.Hi {
		if opts.StepperOnly {
			panic("mujica: NumberInput StepperOnly needs Min < Max")
		}
		s.Lo, s.Hi = math.Inf(-1), math.Inf(1)
	}
	if s.Step == 0 {
		s.Step = 1
	}
	if s.Decimals == 0 {
		if t := strconv.FormatFloat(s.Step, 'f', -1, 64); strings.Contains(t, ".") {
			s.Decimals = len(t) - strings.IndexByte(t, '.') - 1
		}
	}
	// Bounds move inward onto the precision grid, compared exactly with the originals.
	s.min, s.max = s.Lo, s.Hi
	p := math.Pow(10, float64(s.Decimals))
	if s.Lo = s.Round(s.Lo); s.Lo < s.min {
		s.Lo = s.Round(s.Lo + 1/p)
	}
	if s.Hi = s.Round(s.Hi); s.Hi > s.max {
		s.Hi = s.Round(s.Hi - 1/p)
	}
	if s.Lo > s.Hi {
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
	s := NumberInputResolve(opts)
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
	stepBy := func(n int) { set(s.StepBy(*value, n)) }
	look := &textInputLook{disabled: opts.Disabled, err: opts.Error}
	r.Element, r.Input = textInputFrame(c, look, func(frame *ui.Element) *ui.Element {
		text := ui.Local(frame, "text", func() string { return s.format(*value) })
		invalid := ui.Local(frame, "invalid", func() string { return "" })
		seen := ui.Local(frame, "seen", func() float64 { return *value })
		if !numberInputSame(*seen, *value) {
			// Set from outside: the text follows, and a typing error goes.
			*text, *invalid = s.format(*value), ""
		}
		in := TextInputEditor(c, frame, text, look.disabled, ui.TextInputBase).Grow(1).MinWidth(core.Density(c).Unit() * 12).
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
			{core.Msg(c, "number_input.decrease"), "minus", -1, math.IsNaN(*value) || *value > s.Lo},
			{core.Msg(c, "number_input.increase"), "plus", 1, math.IsNaN(*value) || *value < s.Hi},
		} {
			if NumberInputButton(c, b.icon, b.name, b.can) {
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
func numberInputParse(c *ui.Context, s NumberInputSpec, text string, allowEmpty bool, set func(float64)) string {
	t := strings.TrimSpace(text)
	if t == "" && allowEmpty {
		set(math.NaN())
		return ""
	}
	v, err := strconv.ParseFloat(t, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return core.Msg(c, "number_input.invalid")
	}
	if v < s.Lo || v > s.Hi {
		return fmt.Sprintf(core.Msg(c, "number_input.range"), s.format(s.Lo), s.format(s.Hi))
	}
	set(v)
	return ""
}

// NumberInputButton builds a pointer-only step button; Up and Down are its keys.
func NumberInputButton(c *ui.Context, icon, name string, can bool) bool {
	k := core.Tokens(c)
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
func numberInputStepper(c *ui.Context, value *float64, s NumberInputSpec, opts NumberInputOptions) InputResult {
	k := core.Tokens(c)
	d := core.Density(c)
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
				set(s.StepBy(*value, 1))
			case st.Shortcut(0, ui.KeyDown), st.Shortcut(0, ui.KeyLeft):
				set(s.StepBy(*value, -1))
			case st.Shortcut(0, ui.KeyHome):
				set(s.Lo)
			case st.Shortcut(0, ui.KeyEnd):
				set(s.Hi)
			}
			st.Children(func() {
				if numberInputArrow(c, "chevron-up", *value < s.Hi, false) {
					set(s.StepBy(*value, 1))
					st.Focus()
				}
				if numberInputArrow(c, "chevron-down", *value > s.Lo, true) {
					set(s.StepBy(*value, -1))
					st.Focus()
				}
			})
			r.Input = st
		})
		if opts.Error != "" {
			ui.Text(c, opts.Error).TextColor(k.Danger).FontSize(core.FontSize(c, theme.CaptionSize))
		}
	})
	if r.changed {
		slog.Debug("mujica: number_input step", "value", *value)
	}
	return r
}

// numberInputArrow builds one half of the stepper and reports a click.
func numberInputArrow(c *ui.Context, icon string, can, lower bool) bool {
	k := core.Tokens(c)
	h := core.Density(c).ControlHeight()/2 - theme.BorderWidth
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
