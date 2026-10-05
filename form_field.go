package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// FormFieldOptions configure a FormField.
type FormFieldOptions struct {
	// Required marks the label; the check itself belongs to the caller or Form.
	Required    bool
	Description string
	// Error shows below the control and marks it invalid for assistive technology.
	Error string
	// Horizontal puts the label beside the control instead of above it.
	Horizontal bool
}

// FormField labels the element control returns and shows help and error text for it.
func FormField(c *ui.Context, label string, opts FormFieldOptions, control func() *ui.Element) *ui.Element {
	if !opts.Horizontal {
		return formFieldBuild(c, label, opts, false, control)
	}
	var f *ui.Element
	formFieldGrid(c, func() { f = formFieldBuild(c, label, opts, true, control) })
	return f
}

// formFieldGrid holds fields whose labels sit in a column beside their controls.
func formFieldGrid(c *ui.Context, fields func()) *ui.Element {
	d := Density(c)
	return ui.Grid(c).ColumnTracks(ui.FitContent(), ui.Fr(1)).GapX(d.Unit() * 3).GapY(d.Unit() * 4).
		AlignItems(ui.Start).Children(fields)
}

// formFieldBuild builds the label and the control column; beside makes them two grid cells.
func formFieldBuild(c *ui.Context, label string, opts FormFieldOptions, beside bool, control func() *ui.Element) *ui.Element {
	if label == "" {
		panic("mujica: FormField needs a label")
	}
	k := Tokens(c)
	d := Density(c)
	var head, content, ctrl *ui.Element
	body := func() {
		head = ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center)
		if beside {
			head.MinHeight(d.ControlHeight()).Justify(ui.End)
		}
		var text *ui.Element
		head.Children(func() {
			text = ui.Text(c, label).FontWeight(500).SingleLine()
			if opts.Required {
				ui.Text(c, "*").TextColor(k.Danger).Role(ui.RoleNone)
			}
		})
		clicked := head.Clicked()
		content = ui.Column(c).Gap(d.Unit() * 1.5).MinWidth(0)
		content.Children(func() {
			ctrl = control()
			if ctrl == nil {
				panic("mujica: FormField " + label + ": control returned nil")
			}
			if opts.Description != "" {
				ui.Text(c, opts.Description).TextColor(k.TextMuted).FontSize(FontSize(c, theme.CaptionSize)).Role(ui.RoleNone)
			}
			if opts.Error != "" {
				ui.Text(c, opts.Error).TextColor(k.Danger).FontSize(FontSize(c, theme.CaptionSize)).Role(ui.RoleNone)
			}
		})
		ctrl.Label(label).Description(opts.Description).Error(opts.Error)
		if opts.Required {
			ctrl.Description(msg(c, "form.required"))
		}
		if ctrl.IsDisabled() {
			text.TextColor(k.TextMuted)
		} else if clicked {
			ctrl.Focus()
		}
	}
	if beside {
		body()
		return content
	}
	return ui.Column(c).Gap(d.Unit() * 1.5).Children(body)
}
