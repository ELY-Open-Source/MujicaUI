package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// TextOptions configure a Text. The zero value wraps every line.
type TextOptions struct {
	// MaxLines cuts with an ellipsis after n lines; zero never.
	MaxLines int
	// Selectable lets the user select and copy the text.
	Selectable bool
}

// Text shows body text wrapping at its width.
func Text(c *ui.Context, s string, opts TextOptions) *ui.Element {
	if opts.MaxLines < 0 {
		panic("mujica: Text MaxLines must not be negative")
	}
	e := ui.Text(c, s).LineHeight(theme.LineHeight)
	switch {
	case opts.MaxLines == 1:
		e.SingleLine()
	case opts.MaxLines > 1:
		e.MaxLines(opts.MaxLines)
	}
	if opts.Selectable {
		e.Selectable().Focusable()
		internal.Ring(e, Tokens(c).Focus, theme.ControlRadius)
	}
	return e
}

// EditableTextOptions configure the editable mode of Text.
type EditableTextOptions struct {
	// Saving shows the caller's save in progress and blocks editing.
	Saving bool
	// Error is the caller's save failure, shown below the text.
	Error string
	// Label names the text for assistive technology.
	Label string
}

// EditableTextView is a shown EditableText and this frame's outcome.
type EditableTextView struct {
	Element *ui.Element
	// Canceled reports Escape discarded an edit.
	Canceled bool
	changed  bool
}

// Changed reports a confirmed new value, already in *value.
func (v EditableTextView) Changed() bool { return v.changed }

type editableTextState struct {
	editing, fresh bool
	draft          string
}

type editableTextKey struct{}

// EditableText is Text's editable mode with confirm and cancel.
func EditableText(c *ui.Context, value *string, opts EditableTextOptions) EditableTextView {
	k := Tokens(c)
	u := Density(c).Unit()
	line := k.ControlBorder
	if opts.Error != "" {
		line = k.Danger
	}
	var v EditableTextView
	ui.Column(c).Gap(u).AlignItems(ui.Start).Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			v = editableTextField(c, k, value, opts, line)
			if opts.Saving {
				Spinner(c, SpinnerOptions{Size: 12, Label: msg(c, "editable_text.saving")})
			}
		})
		if opts.Error != "" {
			ui.Text(c, opts.Error).TextColor(k.Danger).FontSize(FontSize(c, theme.CaptionSize))
		}
	})
	switch {
	case v.changed:
		slog.Debug("mujica: editabletext changed", "len", len(*value))
	case v.Canceled:
		slog.Debug("mujica: editabletext canceled")
	}
	return v
}

func editableTextField(c *ui.Context, k theme.Tokens, value *string, opts EditableTextOptions, line ui.Color) EditableTextView {
	pad := Density(c).Unit() / 2
	box := ui.Row(c).AlignItems(ui.Center).MinWidth(theme.MinHitSize * 2).MinHeight(theme.MinHitSize).
		Radius(theme.ControlRadius).Disabled(opts.Saving)
	v := EditableTextView{Element: box}
	st := ui.Local(box, editableTextKey{}, func() editableTextState { return editableTextState{} })
	if opts.Saving {
		st.editing = false
	}
	if !st.editing {
		box.Focusable().Padding(pad+1, pad*2+1).BorderWidth(0, 0, 1, 0).BorderColor(line).BorderStyle(ui.BorderDashed).
			Label(opts.Label).Error(opts.Error)
		internal.Ring(box, k.Focus, theme.ControlRadius)
		if box.Hovered() {
			box.Background(k.SurfaceHover)
		}
		box.Transition(colorFade)
		if !box.IsDisabled() && (box.DoubleClicked() || box.Shortcut(0, ui.KeyEnter)) {
			st.editing, st.fresh, st.draft = true, true, *value
		} else {
			box.Children(func() { ui.Text(c, *value).SingleLine() })
			return v
		}
	}
	var in *ui.Element
	box.Children(func() {
		in = ui.TextInputBase(c, &st.draft).Grow(1).MinWidth(theme.MinHitSize*3).Padding(pad, pad*2).
			Radius(theme.ControlRadius).Background(k.Surface).Border(theme.BorderWidth, k.Accent).TextColor(k.Text).
			Label(opts.Label).Error(opts.Error)
		internal.Ring(in, k.Focus, theme.ControlRadius)
	})
	done := func() {
		st.editing = false
		box.Focus()
	}
	switch {
	case st.fresh:
		st.fresh = false
		in.Focus()
	case in.Submitted():
		v.changed = st.draft != *value
		*value = st.draft
		done()
	case in.Shortcut(0, ui.KeyEscape):
		v.Canceled = true
		done()
	case !in.Focused():
		v.changed = st.draft != *value
		*value = st.draft
		st.editing = false
	}
	return v
}
