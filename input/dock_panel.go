package input

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/navigation"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// DockZone is an area of a DockPanel.
type DockZone int

const (
	DockCenter DockZone = iota
	DockLeft
	DockRight
	DockBottom
)

// DockLayout is the caller-held arrangement: panel IDs per zone in tab order, and each zone's shown tab.
type DockLayout struct {
	Zones  [4][]string
	Active [4]string
}

// DockItem is one dockable panel.
type DockItem struct {
	ID, Title string
	Build     func()
}

// DockPanelOptions configure a DockPanel.
type DockPanelOptions struct {
	Disabled bool
}

// dockPanelDrag is a tab being dragged.
type dockPanelDrag struct{ id string }

// DockPanel arranges panels in left, center, right and bottom zones; tabs drag to an edge or zone to dock there.
func DockPanel(c *ui.Context, layout *DockLayout, items []DockItem, opts DockPanelOptions) core.ChoiceResult {
	if layout == nil || len(items) == 0 {
		panic("mujica: DockPanel needs a layout and items")
	}
	byID := map[string]DockItem{}
	for _, it := range items {
		if it.ID == "" || it.Title == "" || it.Build == nil {
			panic("mujica: DockPanel item needs an ID, Title and Build")
		}
		if _, dup := byID[it.ID]; dup {
			panic("mujica: DockPanel item ID repeated: " + it.ID)
		}
		byID[it.ID] = it
	}
	placed := map[string]bool{}
	for z, ids := range layout.Zones {
		for _, id := range ids {
			if _, ok := byID[id]; !ok || placed[id] {
				panic(fmt.Sprintf("mujica: DockPanel zone %d holds unknown or repeated %q", z, id))
			}
			placed[id] = true
		}
	}
	changed := false
	for _, it := range items {
		if !placed[it.ID] {
			layout.Zones[DockCenter] = append(layout.Zones[DockCenter], it.ID)
			changed = true
		}
	}
	for z, ids := range layout.Zones {
		if len(ids) > 0 && !slices.Contains(ids, layout.Active[z]) {
			layout.Active[z] = ids[0]
		}
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	root := ui.Box(c).Disabled(opts.Disabled).MinHeight(0)
	disabled := root.IsDisabled()
	// follow is the tab that moved; it keeps the focus in its new zone.
	follow := ui.Local(root, "dock_panel.follow", func() string { return "" })
	move := func(id string, to DockZone) {
		*follow = id
		for z := range layout.Zones {
			layout.Zones[z] = slices.DeleteFunc(slices.Clone(layout.Zones[z]), func(s string) bool { return s == id })
		}
		layout.Zones[to] = append(layout.Zones[to], id)
		layout.Active[to] = id
		changed = true
		slog.Debug("mujica: dock_panel docked", "panel", id, "zone", to)
	}
	dragging := false
	root.Children(func() {
		ui.Column(c).Fill().Children(func() {
			ui.Row(c).Grow(1).MinHeight(0).AlignItems(ui.Stretch).Children(func() {
				for _, z := range []DockZone{DockLeft, DockCenter, DockRight} {
					if dockPanelZone(c, k, u, layout, byID, z, disabled, move, &changed, follow) {
						dragging = true
					}
				}
			})
			if dockPanelZone(c, k, u, layout, byID, DockBottom, disabled, move, &changed, follow) {
				dragging = true
			}
		})
		// Edges stay one frame after the drag, when the drop arrives.
		was := ui.Local(root, "dock_panel.drag", func() bool { return false })
		if (dragging || *was) && !disabled {
			dockPanelEdges(c, k, u, move)
		}
		*was = dragging
	})
	return core.NewChoiceResult(root, changed)
}

// dockPanelZone builds zone z's tabs and shown panel; it reports a tab of it being dragged.
func dockPanelZone(c *ui.Context, k theme.Tokens, u float32, layout *DockLayout, byID map[string]DockItem, z DockZone, disabled bool, move func(string, DockZone), changed *bool, follow *string) bool {
	ids := layout.Zones[z]
	if len(ids) == 0 {
		return false
	}
	box := ui.Column(c.Key(z)).MinWidth(0).MinHeight(0).Background(k.Surface).Clip().Role(ui.RoleGroup).
		Label(core.Msg(c, fmt.Sprintf("dock_panel.zone%d", z)))
	switch z {
	case DockLeft, DockRight:
		box.Width(220).Shrink(0).Border(theme.BorderWidth, k.Border)
	case DockBottom:
		box.Height(160).Shrink(0).Border(theme.BorderWidth, k.Border)
	default:
		box.Grow(1)
	}
	if d, ok := ui.DragOver[dockPanelDrag](box); ok && !slices.Contains(ids, d.id) {
		box.Border(2, k.AccentText)
	}
	if d, ok := ui.Drop[dockPanelDrag](box); ok && !disabled {
		move(d.id, z)
	}
	dragging := false
	box.Children(func() {
		at := slices.Index(ids, layout.Active[z])
		base := ui.TabsBase(c, &at, len(ids))
		base.List.Gap(u).Padding(theme.FocusGap+theme.FocusWidth, u*2, 0, u*2).Children(func() {
			for i, id := range ids {
				tab := base.Tab(i).Drag(dockPanelDrag{id})
				navigation.TabsStyle(c, tab, navigation.Tab{Label: byID[id].Title}, i == at, false)
				tab.Description(core.Msg(c, "dock_panel.hint"))
				if *follow == id {
					tab.Focus()
					*follow = ""
				}
				dragging = dragging || tab.Dragging()
				for key, to := range map[ui.Key]DockZone{ui.KeyLeft: DockLeft, ui.KeyRight: DockRight, ui.KeyDown: DockBottom, ui.KeyUp: DockCenter} {
					if tab.Shortcut(ui.Shift|ui.Alt, key) && !disabled && to != z {
						move(id, to)
					}
				}
			}
		})
		if base.List.Changed() && !disabled && ids[at] != layout.Active[z] {
			layout.Active[z], *changed = ids[at], true
		}
		ui.Box(c).FillWidth().Height(theme.BorderWidth).Background(k.Border).Role(ui.RoleNone)
		ui.Scroll(c).Grow(1).MinHeight(0).Padding(u * 3).Children(byID[layout.Active[z]].Build)
	})
	return dragging
}

// dockPanelEdges shows drop strips on the four edges while a tab is dragged.
func dockPanelEdges(c *ui.Context, k theme.Tokens, u float32, move func(string, DockZone)) {
	for _, z := range []DockZone{DockLeft, DockRight, DockBottom, DockCenter} {
		e := ui.Box(c.Key(fmt.Sprint("edge", z))).Absolute().Role(ui.RoleNone)
		switch z {
		case DockLeft:
			e.Left(0).Top(0).Bottom(0).Width(12 * u)
		case DockRight:
			e.Right(0).Top(0).Bottom(0).Width(12 * u)
		case DockBottom:
			e.Left(12 * u).Right(12 * u).Bottom(0).Height(10 * u)
		default:
			e.Left(0).Right(0).Top(0).Height(4 * u)
		}
		_, over := ui.DragOver[dockPanelDrag](e)
		if over {
			e.Background(k.Selection).Border(2, k.AccentText)
		}
		if d, ok := ui.Drop[dockPanelDrag](e); ok {
			move(d.id, z)
		}
	}
}
