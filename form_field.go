package mujica

import (
	"github.com/egoist/mygo/ui"
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

// FormField labels the control that control builds and adds help and error text.
// The label names the control and a click on it focuses the control.
func FormField(c *ui.Context, label string, opts FormFieldOptions, control func()) *ui.Element {
	if label == "" {
		panic("mujica: FormField needs a label")
	}
	if opts.Horizontal {
		var f *ui.Element
		ui.Form(c, func() { f = formFieldBuild(c, label, opts, control) })
		return f
	}
	return formFieldBuild(c, label, opts, control)
}

// formFieldBuild builds the field in the layout the caller set up.
func formFieldBuild(c *ui.Context, label string, opts FormFieldOptions, control func()) *ui.Element {
	if opts.Required {
		label += " *"
	}
	return ui.Field(c, label, control).Description(opts.Description).Error(opts.Error)
}
