package mujica

import (
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// ChoiceChipsOptions configure ChoiceChips.
type ChoiceChipsOptions struct {
	// Multiple lets several chips be on; otherwise one or none.
	Multiple bool
	Disabled bool
	// Label names the chip set for assistive technology.
	Label string
}

// ChoiceChips are wrapping filter chips toggling values in *value; arrows move between them.
func ChoiceChips[T comparable](c *ui.Context, value *[]T, items []ToggleItem[T], opts ChoiceChipsOptions) ChoiceResult {
	if value == nil {
		panic("mujica: ChoiceChips needs a value")
	}
	if !opts.Multiple && len(*value) > 1 {
		panic("mujica: single ChoiceChips holds more than one value")
	}
	k := Tokens(c)
	d := Density(c)
	g := ui.Row(c).Wrap().Gap(d.Unit()*2).Disabled(opts.Disabled).Role(ui.RoleGroup).FocusGroup(ui.Horizontal)
	if opts.Label != "" {
		g.Label(opts.Label)
	}
	changed := false
	g.Children(func() {
		for _, it := range items {
			on := slices.Contains(*value, it.Value)
			pressed := on
			b := ui.ToggleBase(c, &pressed).Disabled(it.Disabled)
			// The base took input before Disabled applied: drop a disabled chip's press.
			if b.Changed() && !b.IsDisabled() && toggleGroupApply(value, it.Value, pressed, ToggleGroupOptions{Multiple: opts.Multiple, AllowEmpty: true}) {
				changed, on = true, pressed
			}
			choiceChipLook(c, k, d, b, on, it)
		}
	})
	if changed {
		slog.Debug("mujica: choice_chips changed", "label", opts.Label, "value", *value)
	}
	return ChoiceResult{g, changed}
}

func choiceChipLook[T comparable](c *ui.Context, k theme.Tokens, d theme.Density, b *ui.Element, on bool, it ToggleItem[T]) {
	b.MinHeight(max(theme.MinHitSize, d.ControlHeight()-4)).MinWidth(theme.MinHitSize).Padding(0, d.Unit()*3).Gap(d.Unit()*1.5).
		Radius(theme.PopupRadius)
	internal.Ring(b, k.Focus, theme.PopupRadius)
	bg, fg, line := k.Surface, k.Text, k.ControlBorder
	switch {
	case on:
		bg, fg, line = k.Selection, k.AccentText, k.AccentText
	case b.Pressed():
		bg = k.Selection
	case b.Hovered():
		bg = k.SurfaceHover
	}
	b.Background(bg).TextColor(fg).Border(theme.BorderWidth, line).Transition(colorFade)
	if it.Name != "" {
		b.Label(it.Name)
	}
	b.Children(func() {
		switch {
		case on:
			ui.Icon(c, icons.Must("check")).Size(14, 14)
		case it.Icon != nil:
			ui.Icon(c, it.Icon).Size(14, 14)
		}
		if it.Label != "" {
			ui.Text(c, it.Label).SingleLine()
		}
	})
}
