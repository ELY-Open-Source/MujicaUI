package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// PasswordInputOptions configure a PasswordInput.
type PasswordInputOptions struct {
	Placeholder string
	// Label names the field for assistive technology.
	Label    string
	Disabled bool
	// Error marks the value invalid and shows the message below.
	Error string
}

// PasswordInput is a masked field with a show and hide toggle.
func PasswordInput(c *ui.Context, value *string, opts PasswordInputOptions) InputResult {
	if value == nil {
		panic("mujica: PasswordInput needs a value")
	}
	var r InputResult
	look := &textInputLook{disabled: opts.Disabled, err: opts.Error}
	r.Element, r.Input = textInputFrame(c, look, func(frame *ui.Element) *ui.Element {
		shown := ui.Local(frame, "shown", func() bool { return false })
		var in *ui.Element
		// MyGo's mask cannot be lifted, so each mode keeps its own editor.
		mode := "masked"
		if *shown {
			mode = "shown"
		}
		ui.Row(c).Key(mode).Grow(1).Children(func() {
			in = textInputEditor(c, frame, value, opts.Disabled, ui.TextInputBase).Grow(1).MinWidth(Density(c).Unit() * 10).
				Placeholder(opts.Placeholder).Label(opts.Label)
			if !*shown {
				in.Password()
			}
		})
		r.changed = in.Changed() && !opts.Disabled
		r.submitted = in.Submitted() && !opts.Disabled
		if passwordInputToggle(c, *shown) {
			*shown = !*shown
			slog.Debug("mujica: password_input reveal", "shown", *shown)
		}
		return in
	}, nil)
	return r
}

// passwordInputToggle builds the show or hide button and reports a click.
func passwordInputToggle(c *ui.Context, shown bool) bool {
	k := Tokens(c)
	name, icon := msg(c, "password_input.show"), "eye"
	if shown {
		name, icon = msg(c, "password_input.hide"), "eye-off"
	}
	b := ui.ButtonBase(c).Size(theme.MinHitSize, theme.MinHitSize).Radius(theme.ControlRadius).Label(name).Tooltip(name)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	fg := k.TextMuted
	if b.Pressed() {
		b.Background(k.Selection)
		fg = k.Text
	} else if b.Hovered() {
		b.Background(k.SurfaceHover)
		fg = k.Text
	}
	b.Children(func() { ui.Icon(c, icons.Must(icon)).Size(16, 16).TextColor(fg) })
	return b.Clicked()
}
