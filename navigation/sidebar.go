package navigation

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
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

// Sidebar is a source list choosing an item ID into *selected, built on MyGo's Sidebar.
func Sidebar(c *ui.Context, selected *string, sections []SidebarSection, opts SidebarOptions) core.ChoiceResult {
	k := core.Tokens(c)
	width := opts.Width
	if width == 0 {
		width = 240
		if opts.Compact {
			width = 56
		}
	}
	wrap := ui.Column(c).Width(width).FillHeight().AlignItems(ui.Stretch).Shrink(0)
	closed := ui.Local(wrap, "closed", func() map[string]*bool { return map[string]*bool{} })
	res := core.ChoiceResult{Element: wrap}
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
				path := s.Title
				sidebarGroup(c, s.Title, sidebarOpen(*closed, path), true, 0, func() {
					sidebarItems(c, *selected, s.Items, *closed, path, 0)
				})
			}
		}).FillHeight().FillWidth()
		sb.DrawOver(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X + r.W - 1, Y: r.Y, W: 1, H: r.H}, k.Border, 0)
		})
		c.SetTheme(saved)
		if sb.Changed() {
			res = core.NewChoiceResult(res.Element, true)
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

func sidebarItems(c *ui.Context, selected string, items []SidebarItem, closed map[string]*bool, parent string, depth int) {
	u := core.Density(c).Unit()
	for _, it := range items {
		if it.Label == "" {
			panic("mujica: SidebarItem needs a Label")
		}
		if it.Items != nil {
			// The key holds the whole path, so same-named groups stay apart.
			path := parent + "\x00" + it.Label
			sidebarGroup(c, it.Label, sidebarOpen(closed, path), false, depth, func() {
				sidebarItems(c, selected, it.Items, closed, path, depth+1)
			})
			continue
		}
		ui.Column(c).Padding(0, 0, 0, float32(depth)*u*4).Children(func() {
			sidebarItem(c, selected, it, false)
		})
	}
}

// sidebarGroup is a section title or nested group that Enter, Space or a click opens and closes.
func sidebarGroup(c *ui.Context, title string, open *bool, section bool, depth int, fn func()) {
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	ui.Column(c).Gap(u*0.5).Shrink(0).Padding(0, 0, 0, float32(depth)*u*4).Children(func() {
		was := *open
		tr := ui.CollapsibleBase(c, open).Trigger
		if *open != was {
			slog.Debug("mujica: sidebar group toggle", "group", title, "open", *open)
		}
		tr.AlignItems(ui.Center).Gap(u*2).Padding(u, u*2).Radius(theme.ControlRadius).Label(title).Cursor(ui.CursorPointer)
		if section {
			tr.Margin(u*2, 0, 0, 0).FontWeight(600).FontSize(core.FontSize(c, theme.CaptionSize))
		} else {
			tr.MinHeight(d.ControlHeight())
		}
		internal.Ring(tr, k.Focus, theme.ControlRadius)
		if tr.Hovered() {
			tr.Background(k.SurfaceHover)
		}
		fg := k.Text
		if section {
			fg = k.TextMuted
		}
		internal.Ease(tr, core.Motion(c, theme.StateDuration), fg)
		turn := tr.Animate("mujica.turn", sidebarTurn(*open), core.Motion(c, theme.StateDuration))
		tr.Children(func() {
			ui.Icon(c, icons.Must("chevron-right")).Size(12, 12).Shrink(0).Rotate(turn).TextColor(k.TextMuted)
			ui.Text(c, title).SingleLine().Grow(1)
		})
		if *open {
			fn()
		}
	})
}

func sidebarTurn(open bool) float32 {
	if open {
		return 90
	}
	return 0
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
	k := core.Tokens(c)
	d := core.Density(c)
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
			ui.Text(c, sidebarCount(it.Count)).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum").
				TextColor(k.TextMuted).Padding(0, d.Unit()*1.5).Radius(theme.ControlRadius).Background(k.Background)
		})
	}
	internal.Ease(e, core.Motion(c, theme.StateDuration), k.Text)
	if on {
		e.Background(k.Selection).FontWeight(600)
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
	ui.Box(c).Height(1).Margin(core.Density(c).Unit()*2, 4).Background(core.Tokens(c).Ornament).Role(ui.RoleNone)
}
