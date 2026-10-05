package input

import (
	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

// IconButtonOptions configure an IconButton.
type IconButtonOptions struct {
	// Variant is the emphasis; the zero value is Primary, as for Button.
	Variant ButtonVariant
	// Tooltip shows on hover; empty uses the button's name.
	Tooltip  string
	Loading  bool
	Disabled bool
}

// IconButton is a square button showing only icon; label names its action.
func IconButton(c *ui.Context, icon *ui.SVG, label string, opts IconButtonOptions) *ui.Element {
	if icon == nil {
		panic("mujica: IconButton needs an icon")
	}
	if label == "" {
		panic("mujica: IconButton needs an action name")
	}
	h := core.Density(c).ControlHeight()
	b := Button(c, "", ButtonOptions{Variant: opts.Variant, Icon: icon, Loading: opts.Loading, Disabled: opts.Disabled, Label: label})
	b.Padding(0).Size(h, h)
	tip := opts.Tooltip
	if tip == "" {
		tip = label
	}
	return b.Tooltip(tip)
}
