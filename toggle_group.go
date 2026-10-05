package mujica

import (
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ChoiceResult is a choice control's element and whether its value changed this frame.
type ChoiceResult struct {
	*ui.Element
	changed bool
}

// Changed reports whether the user changed the value in this frame.
func (r ChoiceResult) Changed() bool { return r.changed }

// ToggleItem is one button of a ToggleGroup.
type ToggleItem[T comparable] struct {
	Value T
	Label string
	Icon  *ui.SVG
	// Name names an icon-only item for assistive technology.
	Name     string
	Disabled bool
}

// ToggleGroupOptions configure a ToggleGroup.
type ToggleGroupOptions struct {
	// Multiple lets several items be pressed at once.
	Multiple bool
	// AllowEmpty lets the user release the last pressed item.
	AllowEmpty bool
	Vertical   bool
	Disabled   bool
	// Label names the group for assistive technology.
	Label string
}

// ToggleGroup is a joined row of toggles choosing one item of *value, or several with Multiple.
func ToggleGroup[T comparable](c *ui.Context, value *[]T, items []ToggleItem[T], opts ToggleGroupOptions) ChoiceResult {
	if !opts.Multiple && len(*value) > 1 {
		panic("mujica: single ToggleGroup holds more than one value")
	}
	k := Tokens(c)
	g := ui.Row(c).Shrink(0).AlignSelf(ui.Start).AlignItems(ui.Stretch).Radius(theme.ControlRadius).
		Border(theme.BorderWidth, k.ControlBorder).Background(k.Surface).Disabled(opts.Disabled).Role(ui.RoleGroup)
	if opts.Vertical {
		g.Column().FocusGroup(ui.Vertical)
	} else {
		g.FocusGroup(ui.Horizontal)
	}
	if opts.Label != "" {
		g.Label(opts.Label)
	}
	g.Dividers(theme.BorderWidth, k.ControlBorder)
	changed := false
	g.Children(func() {
		for _, it := range items {
			on := slices.Contains(*value, it.Value)
			pressed := on
			b := ui.ToggleBase(c, &pressed).Disabled(it.Disabled)
			// The base took input before Disabled applied: drop a disabled item's press.
			if b.Changed() && !it.Disabled && !opts.Disabled && toggleGroupApply(value, it.Value, pressed, opts) {
				changed = true
				on = pressed
			}
			toggleGroupLook(c, k, b, on, it)
		}
	})
	if changed {
		slog.Debug("mujica: toggle group changed", "label", opts.Label, "value", *value)
	}
	return ChoiceResult{g, changed}
}

// toggleGroupApply presses or releases v by the group's rules and reports a change.
func toggleGroupApply[T comparable](value *[]T, v T, pressed bool, opts ToggleGroupOptions) bool {
	switch {
	case pressed && opts.Multiple:
		*value = append(*value, v)
	case pressed:
		*value = []T{v}
	case len(*value) == 1 && !opts.AllowEmpty:
		return false
	default:
		i := slices.Index(*value, v)
		*value = slices.Delete(slices.Clone(*value), i, i+1)
	}
	return true
}

func toggleGroupLook[T comparable](c *ui.Context, k theme.Tokens, b *ui.Element, on bool, it ToggleItem[T]) {
	d := Density(c)
	b.MinHeight(d.ControlHeight()-2).MinWidth(d.ControlHeight()).Padding(0, d.Unit()*3).Gap(d.Unit() * 2)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	bg, fg := ui.Transparent, k.Text
	switch {
	case on:
		bg, fg = k.Selection, k.AccentText
	case b.Pressed():
		bg = k.Selection
	case b.Hovered():
		bg = k.SurfaceHover
	}
	b.Background(bg).TextColor(fg)
	if on {
		toggleMark(b, k)
	}
	if it.Name != "" {
		b.Label(it.Name)
	}
	b.Children(func() {
		if it.Icon != nil {
			ui.Icon(c, it.Icon).Size(16, 16)
		}
		if it.Label != "" {
			ui.Text(c, it.Label).SingleLine()
		}
	})
}
