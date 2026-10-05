package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
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
	h := Density(c).ControlHeight()
	b := Button(c, "", ButtonOptions{Variant: opts.Variant, Icon: icon, Loading: opts.Loading, Disabled: opts.Disabled, Label: label})
	b.Padding(0).Size(h, h)
	tip := opts.Tooltip
	if tip == "" {
		tip = label
	}
	return b.Tooltip(tip)
}

// iconAction is a square, borderless button showing icon, named label.
func iconAction(c *ui.Context, icon *ui.SVG, label string) *ui.Element {
	k := Tokens(c)
	b := ui.ButtonBase(c).Size(theme.MinHitSize, theme.MinHitSize).Radius(theme.ControlRadius).Label(label).TextColor(k.TextMuted)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case b.Pressed():
		b.Background(k.Selection).TextColor(k.Text)
	case b.Hovered():
		b.Background(k.SurfaceHover).TextColor(k.Text)
	}
	b.Children(func() { ui.Icon(c, icon).Size(16, 16) })
	return b
}
