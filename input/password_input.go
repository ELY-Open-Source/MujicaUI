package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
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
	r.Element, r.Input = textInputFrame(c, look, func(frame ui.Element) ui.Element {
		shown := ui.Local(frame, "shown", func() bool { return false })
		var in ui.Element
		// MyGo's mask cannot be lifted, so each mode keeps its own editor.
		mode := "masked"
		if *shown {
			mode = "shown"
		}
		ui.Row(c.Key(mode)).Grow(1).Children(func() {
			in = TextInputEditor(c, frame, value, look.disabled, ui.TextInputBase).Grow(1).MinWidth(core.Density(c).Unit() * 10).
				Placeholder(opts.Placeholder).Label(opts.Label)
			if !*shown {
				in.Password()
			}
		})
		r.changed = in.Changed() && !look.disabled
		r.submitted = in.Submitted() && !look.disabled
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
	name, icon := core.Msg(c, "password_input.show"), "eye"
	if shown {
		name, icon = core.Msg(c, "password_input.hide"), "eye-off"
	}
	return core.IconAction(c, icons.Must(icon), name).Tooltip(name).Clicked()
}
