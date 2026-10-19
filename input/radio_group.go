package input

import (
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
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
func RadioGroup[T comparable](c *ui.Context, value *T, items []RadioItem[T], opts RadioGroupOptions) core.ChoiceResult {
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
	k := core.Tokens(c)
	d := core.Density(c)
	var r core.ChoiceResult
	// Disabled on a wrapper blocks the bases' clicks before they choose.
	wrap := ui.Box(c).Shrink(0).Disabled(opts.Disabled)
	// Keys by first appearance of each value: fmt keys would merge equal-looking values.
	known := ui.Local(wrap, "values", func() []T { return nil })
	wrap.Children(func() {
		r.Element = ui.RadioGroup(c, func() {
			for _, it := range items {
				if radioGroupItem(c, k, d, value, it, stableKey(known, it.Value)) {
					r = core.NewChoiceResult(r.Element, true)
					slog.Debug("mujica: radio_group choose", "label", it.Label)
				}
			}
		}).Label(opts.Label).Gap(0)
	})
	if opts.Horizontal {
		r.Element.Row().Gap(d.Unit() * 6)
	}
	if wrap.IsDisabled() && wrap.FocusWithin() {
		wrap.FocusRing(false).Focus()
	}
	return r
}

// radioGroupItem builds one radio and reports whether it was chosen.
func radioGroupItem[T comparable](c *ui.Context, k theme.Tokens, d theme.Density, value *T, it RadioItem[T], key int) bool {
	// Keyed by value so focus follows the item, not its position.
	var row ui.Element
	item := ui.Box(c.Key(key)).Disabled(it.Disabled)
	item.Children(func() {
		row = ui.RadioBase(c, value, it.Value).Gap(d.Unit() * 2).MinHeight(d.ControlHeight()).FocusRing(false)
	})
	if item.IsDisabled() && item.FocusWithin() {
		item.FocusRing(false).Focus()
	}
	on := *value == it.Value
	hover, pressed := row.Hovered(), row.Pressed()
	row.Children(func() {
		bg, border := checkboxColors(k, on, hover, pressed)
		dot := ui.Box(c).Size(checkboxSize, checkboxSize).Shrink(0).Radius(checkboxSize/2).
			Background(bg).Border(theme.BorderWidth, border)
		dot.Transition(core.ColorFade)
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

// stableKey keys v by its first appearance in *known, so equal-looking values stay apart.
func stableKey[T comparable](known *[]T, v T) int {
	i := slices.Index(*known, v)
	if i < 0 {
		i = len(*known)
		*known = append(*known, v)
	}
	return i
}
