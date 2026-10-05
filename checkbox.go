package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// CheckboxOptions configure a Checkbox.
type CheckboxOptions struct {
	// Indeterminate shows the mixed mark; a click then checks the box.
	Indeterminate bool
	Disabled      bool
}

// checkboxSize is the box edge in DIPs.
const checkboxSize float32 = 16

// Checkbox toggles *value; a click on the box or its label toggles, as does Space.
func Checkbox(c *ui.Context, value *bool, label string, opts CheckboxOptions) *ui.Element {
	if value == nil {
		panic("mujica: Checkbox needs a value")
	}
	k := Tokens(c)
	d := Density(c)
	// Disabled on a wrapper blocks the base's click before it toggles.
	var row *ui.Element
	wrap := ui.Row(c).Shrink(0).Disabled(opts.Disabled)
	wrap.Children(func() {
		row = ui.CheckboxBase(c, value).Gap(d.Unit() * 2).MinHeight(d.ControlHeight()).FocusRing(false)
	})
	if wrap.IsDisabled() && wrap.FocusWithin() {
		wrap.FocusRing(false).Focus()
	}
	if row.Changed() && opts.Indeterminate && !*value {
		// From mixed, a click always checks.
		*value = true
	}
	if opts.Indeterminate {
		row.Description(msg(c, "checkbox.mixed"))
	}
	on, mixed := *value, opts.Indeterminate
	hover, pressed := row.Hovered(), row.Pressed()
	row.Children(func() {
		bg, border := checkboxColors(k, on || mixed, hover, pressed)
		box := ui.Box(c).Size(checkboxSize, checkboxSize).Shrink(0).Radius(theme.ControlRadius).
			Background(bg).Border(theme.BorderWidth, border)
		box.Transition(textInputFade)
		box.DrawOver(func(p *ui.Painter, r ui.Rect) {
			switch {
			case mixed:
				p.Fill(ui.Rect{X: r.X + r.W*0.25, Y: r.Y + r.H/2 - 1, W: r.W * 0.5, H: 2}, k.OnAccent, 1)
			case on:
				var path ui.Path
				path.MoveTo(r.X+r.W*0.24, r.Y+r.H*0.52).LineTo(r.X+r.W*0.42, r.Y+r.H*0.7).LineTo(r.X+r.W*0.77, r.Y+r.H*0.3)
				p.StrokePath(&path, 2, k.OnAccent)
			}
			if row.FocusVisible() {
				internal.FocusRing(p, r, theme.ControlRadius, k.Focus)
			}
		})
		if label != "" {
			ui.Text(c, label)
		}
	})
	if row.Changed() {
		slog.Debug("mujica: checkbox toggle", "label", label, "checked", *value)
	}
	return row
}

// checkboxColors gives a check box or radio mark its fill and edge.
func checkboxColors(k theme.Tokens, on, hover, pressed bool) (bg, border ui.Color) {
	switch {
	case on && pressed:
		return k.AccentPressed, k.Text
	case on && hover:
		return k.AccentHover, k.Text
	case on:
		return k.Accent, k.AccentText
	case pressed:
		return k.Selection, k.AccentText
	case hover:
		return k.SurfaceHover, k.AccentText
	}
	return k.Surface, k.ControlBorder
}

// CheckboxItem is one box of a CheckboxGroup.
type CheckboxItem struct {
	Label    string
	Value    *bool
	Disabled bool
}

// CheckboxGroupOptions configure a CheckboxGroup.
type CheckboxGroupOptions struct {
	Disabled bool
}

// CheckboxGroup shows a select-all box over items: checked, mixed or clear.
func CheckboxGroup(c *ui.Context, label string, items []CheckboxItem, opts CheckboxGroupOptions) InputResult {
	if len(items) == 0 {
		panic("mujica: CheckboxGroup needs items")
	}
	d := Density(c)
	on, enabled, enabledOn := 0, 0, 0
	for _, it := range items {
		if it.Value == nil {
			panic("mujica: CheckboxItem needs a Value")
		}
		if *it.Value {
			on++
		}
		if !it.Disabled {
			enabled++
			if *it.Value {
				enabledOn++
			}
		}
	}
	var r InputResult
	r.Element = ui.Column(c).Shrink(0).Disabled(opts.Disabled || enabled == 0)
	r.Element.Children(func() {
		all := on == len(items)
		head := Checkbox(c, &all, label, CheckboxOptions{Indeterminate: on > 0 && !all})
		r.Input = head
		if head.Changed() {
			// Select-all leaves disabled items alone.
			target := enabledOn < enabled
			for _, it := range items {
				if !it.Disabled {
					*it.Value = target
				}
			}
			r.changed = true
			slog.Debug("mujica: checkbox_group select_all", "label", label, "checked", target)
		}
		ui.Column(c).Padding(0, 0, 0, checkboxSize+d.Unit()*2).Role(ui.RoleGroup).Label(label).Children(func() {
			for _, it := range items {
				if Checkbox(c, it.Value, it.Label, CheckboxOptions{Disabled: it.Disabled}).Changed() {
					r.changed = true
				}
			}
		})
	})
	return r
}
