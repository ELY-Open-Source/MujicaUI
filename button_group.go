package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ButtonGroupItem is one button of a ButtonGroup. Give Label, Icon or both; icon-only needs Name.
type ButtonGroupItem struct {
	Label string
	Icon  *ui.SVG
	// Name names an icon-only item for assistive technology.
	Name     string
	Disabled bool
}

// ButtonGroupOptions configure a ButtonGroup.
type ButtonGroupOptions struct {
	Vertical bool
	// Label names the group for assistive technology.
	Label    string
	Disabled bool
}

// ButtonGroupResult is the group element and the item clicked this frame.
type ButtonGroupResult struct {
	Element *ui.Element
	clicked int
}

// ClickedItem returns the index of the item clicked this frame, or -1.
func (r ButtonGroupResult) ClickedItem() int { return r.clicked }

// ButtonGroup joins buttons in one bordered strip; arrows move between them.
func ButtonGroup(c *ui.Context, items []ButtonGroupItem, opts ButtonGroupOptions) ButtonGroupResult {
	if len(items) == 0 {
		panic("mujica: ButtonGroup needs items")
	}
	k := Tokens(c)
	d := Density(c)
	orient, last := ui.Horizontal, len(items)-1
	g := ui.Row(c).AlignItems(ui.Stretch).Shrink(0)
	if opts.Vertical {
		orient = ui.Vertical
		g.Column()
	}
	g.FocusGroup(orient).Role(ui.RoleToolbar).Label(opts.Label).Disabled(opts.Disabled)
	r := ButtonGroupResult{Element: g, clicked: -1}
	g.Children(func() {
		for i, it := range items {
			if it.Label == "" && it.Icon == nil {
				panic("mujica: ButtonGroupItem needs a Label or an Icon")
			}
			if it.Label == "" && it.Name == "" {
				panic("mujica: icon-only ButtonGroupItem needs a Name")
			}
			b := ui.ButtonBase(c).Key(i).MinHeight(d.ControlHeight()).MinWidth(d.ControlHeight()).
				Padding(0, d.Unit()*3).Gap(d.Unit() * 2).Disabled(it.Disabled)
			internal.Ring(b, k.Focus, theme.ControlRadius)
			if it.Name != "" {
				b.Label(it.Name)
			}
			bg := k.Surface
			if b.Pressed() {
				bg = k.Selection
			} else if b.Hovered() {
				bg = k.SurfaceHover
			}
			b.Background(bg).TextColor(k.Text).Transition(colorFade)
			b.Radius(buttonGroupCorners(i, last, opts.Vertical)...)
			b.Border(theme.BorderWidth, k.ControlBorder).BorderWidth(buttonGroupEdges(i, opts.Vertical)...)
			b.Children(func() {
				if it.Icon != nil {
					ui.Icon(c, it.Icon).Size(16, 16)
				}
				if it.Label != "" {
					ui.Text(c, it.Label).SingleLine()
				}
			})
			if b.Clicked() {
				r.clicked = i
				slog.Debug("mujica: button_group click", "item", i)
			}
		}
	})
	return r
}

// buttonGroupCorners rounds only the outer corners of item i of 0..last.
func buttonGroupCorners(i, last int, vertical bool) []float32 {
	var tl, tr, br, bl float32
	r := theme.ControlRadius
	first, end := i == 0, i == last
	if vertical {
		if first {
			tl, tr = r, r
		}
		if end {
			br, bl = r, r
		}
	} else {
		if first {
			tl, bl = r, r
		}
		if end {
			tr, br = r, r
		}
	}
	return []float32{tl, tr, br, bl}
}

// buttonGroupEdges draws each shared edge once: items after the first skip the leading one.
func buttonGroupEdges(i int, vertical bool) []float32 {
	w := theme.BorderWidth
	switch {
	case i == 0:
		return []float32{w}
	case vertical:
		return []float32{0, w, w, w}
	default:
		return []float32{w, w, w, 0}
	}
}
