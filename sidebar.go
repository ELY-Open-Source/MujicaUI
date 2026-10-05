package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// SidebarItem is a destination of a Sidebar, or with Items a nested group.
type SidebarItem struct {
	// ID is what choosing the item stores; unused for groups.
	ID    string
	Label string
	Icon  *ui.SVG
	// Count shows beside the label when positive.
	Count int
	// Items nest under the item, which then expands and collapses.
	Items []SidebarItem
}

// SidebarSection is a titled group of a Sidebar; its title collapses it.
type SidebarSection struct {
	Title string
	Items []SidebarItem
}

// SidebarOptions configure a Sidebar.
type SidebarOptions struct {
	// Compact shows only icons, labels as tooltips; every item needs an Icon.
	Compact bool
	// Width is 240 DIPs when zero, 56 when Compact.
	Width float32
}

// SidebarResult is the built Sidebar.
type SidebarResult struct {
	*ui.Element
	changed bool
}

// Changed reports a new choice in this frame.
func (r SidebarResult) Changed() bool { return r.changed }

// Sidebar is a source list choosing an item ID into *selected, built on MyGo's Sidebar.
func Sidebar(c *ui.Context, selected *string, sections []SidebarSection, opts SidebarOptions) SidebarResult {
	k := Tokens(c)
	width := opts.Width
	if width == 0 {
		width = 240
		if opts.Compact {
			width = 56
		}
	}
	wrap := ui.Column(c).Width(width).FillHeight().AlignItems(ui.Stretch).Shrink(0)
	closed := ui.Local(wrap, "closed", func() map[string]*bool { return map[string]*bool{} })
	res := SidebarResult{Element: wrap}
	prev := *selected
	wrap.Children(func() {
		// MyGo draws icons in Accent: swap in readable state colors for the sidebar.
		saved := c.Theme()
		t := *saved
		t.Accent, t.AccentText, t.Surface = k.AccentText, k.AccentText, k.Surface
		c.SetTheme(&t)
		sb := ui.Sidebar(c, selected, func() {
			for si, s := range sections {
				if opts.Compact {
					if si > 0 {
						sidebarRule(c)
					}
					sidebarFlat(c, *selected, s.Items)
					continue
				}
				if s.Title == "" {
					panic("mujica: SidebarSection needs a Title")
				}
				ui.SidebarSection(c, s.Title, sidebarOpen(*closed, "s:"+s.Title), func() {
					sidebarItems(c, *selected, s.Items, *closed, 0)
				})
			}
		}).FillHeight().FillWidth()
		sb.DrawOver(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X + r.W - 1, Y: r.Y, W: 1, H: r.H}, k.Border, 0)
		})
		c.SetTheme(saved)
		if sb.Changed() {
			res.changed = true
			slog.Debug("mujica: sidebar choose", "from", prev, "to", *selected)
		}
	})
	return res
}

// sidebarOpen returns the shown state of a collapsible part, open at first.
func sidebarOpen(closed map[string]*bool, key string) *bool {
	o, ok := closed[key]
	if !ok {
		o = new(bool)
		*o = true
		closed[key] = o
	}
	return o
}

func sidebarItems(c *ui.Context, selected string, items []SidebarItem, closed map[string]*bool, depth int) {
	u := Density(c).Unit()
	for _, it := range items {
		if it.Label == "" {
			panic("mujica: SidebarItem needs a Label")
		}
		if it.Items != nil {
			ui.Column(c).Padding(0, 0, 0, float32(depth)*u*4).Children(func() {
				ui.SidebarSection(c, it.Label, sidebarOpen(closed, fmt.Sprintf("g%d:%s", depth, it.Label)), func() {
					sidebarItems(c, selected, it.Items, closed, depth+1)
				})
			})
			continue
		}
		ui.Column(c).Padding(0, 0, 0, float32(depth)*u*4).Children(func() {
			sidebarItem(c, selected, it, false)
		})
	}
}

// sidebarFlat lists every leaf for the compact rail.
func sidebarFlat(c *ui.Context, selected string, items []SidebarItem) {
	for _, it := range items {
		if it.Items != nil {
			sidebarFlat(c, selected, it.Items)
			continue
		}
		if it.Icon == nil {
			panic("mujica: compact Sidebar items need an Icon")
		}
		sidebarItem(c, selected, it, true)
	}
}

func sidebarItem(c *ui.Context, selected string, it SidebarItem, compact bool) {
	if it.ID == "" {
		panic("mujica: SidebarItem needs an ID")
	}
	k := Tokens(c)
	d := Density(c)
	on := it.ID == selected
	label := it.Label
	if compact {
		label = ""
	}
	e := ui.SidebarItem(c, it.ID, it.Icon, label).MinHeight(d.ControlHeight()).Radius(theme.ControlRadius)
	name := it.Label
	if it.Count > 0 {
		name = fmt.Sprintf("%s (%d)", it.Label, it.Count)
	}
	if compact {
		e.Tooltip(name).Label(name)
	} else if it.Count > 0 {
		e.Label(name).Children(func() {
			ui.Text(c, sidebarCount(it.Count)).FontSize(theme.CaptionSize).FontFeatures("tnum").
				TextColor(k.TextMuted).Padding(0, d.Unit()*1.5).Radius(theme.ControlRadius).Background(k.Background)
		})
	}
	if on {
		e.Background(k.Selection).TextColor(k.Text).FontWeight(600)
		e.Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + 4, W: 3, H: r.H - 8}, k.AccentText, 1.5)
		})
	}
}

func sidebarCount(n int) string {
	if n > 999 {
		return "999+"
	}
	return fmt.Sprint(n)
}

func sidebarRule(c *ui.Context) {
	ui.Box(c).Height(1).Margin(Density(c).Unit()*2, 4).Background(Tokens(c).Ornament).Role(ui.RoleNone)
}
