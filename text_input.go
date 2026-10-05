package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// InputResult is a value input's outer element and its events.
type InputResult struct {
	Element *ui.Element
	// Input is the text field inside, for focus and keys.
	Input     *ui.Element
	changed   bool
	submitted bool
}

// Changed reports that the user changed the value this frame.
func (r InputResult) Changed() bool { return r.changed }

// Submitted reports Enter, or the submit shortcut of a TextArea.
func (r InputResult) Submitted() bool { return r.submitted }

// TextInputOptions configure a TextInput.
type TextInputOptions struct {
	Placeholder string
	// Label names the field for assistive technology.
	Label string
	// Leading and Trailing build content inside the field's edges.
	Leading  func()
	Trailing func()
	// Clearable shows a clear button while there is text; Escape clears too.
	Clearable bool
	Disabled  bool
	// ReadOnly keeps the text selectable and copyable but not editable.
	ReadOnly bool
	// Error marks the value invalid and shows the message below.
	Error string
}

// TextInput is a single-line field editing *value.
func TextInput(c *ui.Context, value *string, opts TextInputOptions) InputResult {
	if value == nil {
		panic("mujica: TextInput needs a value")
	}
	var r InputResult
	look := &textInputLook{disabled: opts.Disabled, readOnly: opts.ReadOnly, err: opts.Error}
	r.Element, r.Input = textInputFrame(c, look, func(frame *ui.Element) *ui.Element {
		if opts.Leading != nil {
			opts.Leading()
		}
		in := textInputEditor(c, frame, value, opts.ReadOnly || look.disabled, ui.TextInputBase).Grow(1).MinWidth(Density(c).Unit() * 10).
			Placeholder(opts.Placeholder).Label(opts.Label)
		r.changed = in.Changed() && !opts.ReadOnly && !look.disabled
		r.submitted = in.Submitted() && !look.disabled
		if opts.Clearable && *value != "" && !opts.ReadOnly && !look.disabled {
			if textInputClear(c, in) {
				*value = ""
				r.changed = true
				slog.Debug("mujica: text_input clear")
			}
		}
		if opts.Trailing != nil {
			opts.Trailing()
		}
		return in
	}, nil)
	return r
}

// textInputFade eases hover, focus and error colors.
var textInputFade = ui.ElementTransition{Colors: true, Duration: theme.StateDuration}

// textInputLook is the state a field frame shows; disabled gains any disabled ancestor.
type textInputLook struct {
	disabled, readOnly, area bool
	err                      string
}

// textInputFrame builds the bordered field around what build returns as the input.
func textInputFrame(c *ui.Context, look *textInputLook, build func(frame *ui.Element) *ui.Element, footer func()) (outer, in *ui.Element) {
	k := Tokens(c)
	d := Density(c)
	outer = ui.Column(c).Gap(d.Unit()).Shrink(0).Disabled(look.disabled)
	// One disabled state, the field's own or an ancestor's, for build to obey.
	look.disabled = outer.IsDisabled()
	outer.Children(func() {
		frame := ui.Row(c).Radius(theme.ControlRadius).Gap(d.Unit() * 2).Cursor(ui.CursorText)
		if look.area {
			frame.AlignItems(ui.Stretch)
		} else {
			frame.MinHeight(d.ControlHeight()).Padding(0, d.Unit()*3)
		}
		frame.Children(func() { in = build(frame) })
		in.FocusRing(false)
		if look.disabled && outer.FocusWithin() {
			// A field disabled while focused gives the focus up.
			outer.FocusRing(false).Focus()
		}
		if look.err != "" {
			in.Error(look.err)
		}
		if frame.Clicked() {
			in.Focus()
		}
		bg, border := k.Surface, k.ControlBorder
		switch {
		case look.readOnly:
			bg, border = k.Background, k.Border
		case frame.Hovered():
			border = k.ControlBorder.Mix(k.Text, 0.35)
		}
		switch {
		case look.err != "":
			border = k.Danger
		case in.Focused() && !look.readOnly:
			border = k.Focus
		}
		frame.Background(bg).Border(theme.BorderWidth, border).Transition(textInputFade)
		frame.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if in.FocusVisible() {
				internal.FocusRing(p, r, theme.ControlRadius, k.Focus)
			}
		})
		if look.err == "" && footer == nil {
			return
		}
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Start).Children(func() {
			if look.err != "" {
				ui.Text(c, look.err).TextColor(k.Danger).FontSize(FontSize(c, theme.CaptionSize)).Grow(1)
			} else {
				ui.Spacer(c)
			}
			if footer != nil {
				footer()
			}
		})
	})
	return outer, in
}

// textInputEditor builds the MyGo editor; a locked one edits a copy, undone at once.
func textInputEditor(c *ui.Context, frame *ui.Element, value *string, locked bool, base func(*ui.Context, *string) *ui.Element) *ui.Element {
	if !locked {
		return base(c, value)
	}
	shown := ui.Local(frame, "locked", func() string { return *value })
	*shown = *value
	in := base(c, shown)
	*shown = *value
	return in
}

// textInputClear builds the clear button and reports a clear by click or Escape.
func textInputClear(c *ui.Context, in *ui.Element) bool {
	k := Tokens(c)
	b := ui.Box(c).Size(theme.MinHitSize, theme.MinHitSize).Shrink(0).Center().Radius(theme.ControlRadius).
		Role(ui.RoleButton).Label(msg(c, "text_input.clear")).Cursor(ui.CursorPointer)
	fg := k.TextMuted
	if b.Hovered() {
		b.Background(k.SurfaceHover)
		fg = k.Text
	}
	b.Children(func() { ui.Icon(c, icons.Must("x")).Size(14, 14).TextColor(fg) })
	if b.Clicked() || in.Shortcut(0, ui.KeyEscape) {
		in.Focus()
		return true
	}
	return false
}
