package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// RadioItem is one choice of a RadioGroup; Value is its stable identity.
type RadioItem[T comparable] struct {
	Value    T
	Label    string
	Disabled bool
}

// RadioGroupOptions configure a RadioGroup.
type RadioGroupOptions struct {
	// Horizontal lays the choices in a row; the default is a column.
	Horizontal bool
	// Label names the group for assistive technology.
	Label    string
	Disabled bool
}

// RadioGroup chooses one item's Value into *value; arrows move the choice.
func RadioGroup[T comparable](c *ui.Context, value *T, items []RadioItem[T], opts RadioGroupOptions) InputResult {
	if value == nil {
		panic("mujica: RadioGroup needs a value")
	}
	if len(items) == 0 {
		panic("mujica: RadioGroup needs items")
	}
	seen := make(map[T]bool, len(items))
	for _, it := range items {
		if seen[it.Value] {
			panic("mujica: RadioGroup has a duplicate value")
		}
		seen[it.Value] = true
	}
	k := Tokens(c)
	d := Density(c)
	var r InputResult
	r.Element = ui.RadioGroup(c, func() {
		for _, it := range items {
			if radioGroupItem(c, k, d, value, it) {
				r.changed = true
				slog.Debug("mujica: radio_group choose", "label", it.Label)
			}
		}
	}).Label(opts.Label).Disabled(opts.Disabled).Gap(0)
	if opts.Horizontal {
		r.Element.Row().Gap(d.Unit() * 6)
	}
	return r
}

// radioGroupItem builds one radio and reports whether it was chosen.
func radioGroupItem[T comparable](c *ui.Context, k theme.Tokens, d theme.Density, value *T, it RadioItem[T]) bool {
	row := ui.RadioBase(c, value, it.Value).Gap(d.Unit() * 2).MinHeight(d.ControlHeight()).FocusRing(false).Disabled(it.Disabled)
	on := *value == it.Value
	hover, pressed := row.Hovered(), row.Pressed()
	row.Children(func() {
		dot := ui.Box(c).Size(checkboxSize, checkboxSize).Shrink(0).Radius(checkboxSize / 2)
		switch {
		case on:
			dot.Background(k.Accent).Border(theme.BorderWidth, k.AccentText)
		case pressed:
			dot.Background(k.Selection).Border(theme.BorderWidth, k.AccentText)
		case hover:
			dot.Background(k.SurfaceHover).Border(theme.BorderWidth, k.AccentText)
		default:
			dot.Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder)
		}
		dot.Transition(textInputFade)
		dot.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if on {
				s := r.W * 0.375
				p.Fill(ui.Rect{X: r.X + (r.W-s)/2, Y: r.Y + (r.H-s)/2, W: s, H: s}, k.OnAccent, s/2)
			}
			if row.FocusVisible() {
				internal.FocusRing(p, r, checkboxSize/2, k.Focus)
			}
		})
		ui.Text(c, it.Label)
	})
	return row.Changed()
}
