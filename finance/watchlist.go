package finance

import (
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/navigation"
	"github.com/ZacharyZhang-NY/MujicaUI/overlay"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// WatchItem is a watched symbol with its latest quote.
type WatchItem struct {
	Symbol, Name  string
	Price         float64
	Change        float64
	ChangePercent float64
	Volume        float64
}

// WatchGroup is a named list of watched symbols.
type WatchGroup struct {
	Name  string
	Items []WatchItem
}

// WatchColumn is an optional Watchlist column.
type WatchColumn int

const (
	WatchPrice WatchColumn = iota
	WatchChange
	WatchVolume
)

var watchColumns = [3]string{"watchlist.price", "watchlist.change", "watchlist.volume"}

type WatchlistState struct {
	Group    int
	Selected string
	Hidden   map[WatchColumn]bool
	rows     ui.ListState
	cursor   int
}

// WatchlistOptions configure a Watchlist.
type WatchlistOptions struct {
	// Decimals of prices; zero means 2.
	Decimals int
	Disabled bool
}

// WatchlistView is a built Watchlist and this frame's events.
type WatchlistView struct {
	Element            *ui.Element
	Changed, Submitted bool
	Moved              bool
	Rows               []int
	To                 int
}

func Watchlist(c *ui.Context, s *WatchlistState, groups []WatchGroup, opts WatchlistOptions) WatchlistView {
	if s == nil || len(groups) == 0 {
		panic("mujica: Watchlist needs a state and a group")
	}
	dec := core.PriceDecimals("Watchlist", opts.Decimals)

	if s.Group < 0 || s.Group >= len(groups) {
		panic("mujica: Watchlist Group out of range")
	}
	d := core.Density(c)
	var v WatchlistView
	v.Element = ui.Column(c).Gap(d.Unit() * 2).Grow(1).MinHeight(0).Disabled(opts.Disabled)
	opts.Disabled = v.Element.IsDisabled()
	if opts.Disabled {
		copy := *s
		s = &copy
	}
	if s.Hidden == nil {
		s.Hidden = map[WatchColumn]bool{}
	}
	v.Element.Children(func() {
		ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
			tabs := make([]navigation.Tab, len(groups))
			for i, g := range groups {
				tabs[i] = navigation.Tab{Label: g.Name}
			}
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				navigation.Tabs(c, &s.Group, tabs, navigation.TabsOptions{Label: core.Msg(c, "watchlist.groups")})
			})
			overlay.DropdownMenu(c, "", overlay.DropdownMenuOptions{Icon: icons.Must("settings"), Label: core.Msg(c, "watchlist.columns"), Disabled: opts.Disabled}, func(m *overlay.PopupMenu) {
				for col, key := range watchColumns {
					on := !s.Hidden[WatchColumn(col)]
					if m.Check(core.Msg(c, key), &on, overlay.PopupMenuItemOptions{KeepOpen: true}) {
						s.Hidden[WatchColumn(col)] = !on
					}
				}
			})
		})
		items := groups[s.Group].Items
		v.Changed, v.Submitted = watchlistRows(c, s, items, dec, &v, opts.Disabled)
	})
	if v.Element.IsDisabled() {
		v.Changed, v.Submitted, v.Moved = false, false, false
	}
	switch {
	case v.Moved:
		slog.Debug("mujica: watchlist move", "rows", v.Rows, "to", v.To)
	case v.Changed:
		slog.Debug("mujica: watchlist select", "symbol", s.Selected)
	}
	return v
}

func watchlistRows(c *ui.Context, s *WatchlistState, items []WatchItem, dec int, v *WatchlistView, disabled bool) (changed, submitted bool) {
	k := core.Tokens(c)
	d := core.Density(c)
	s.cursor = slices.IndexFunc(items, func(it WatchItem) bool { return it.Symbol == s.Selected })
	before := s.cursor
	cursor := &s.cursor
	if disabled {
		cp := s.cursor
		cursor = &cp
	}
	s.rows.Key = func(i int) any { return items[i].Symbol }
	s.rows.Label = func(i int) string { return items[i].Symbol }
	s.rows.Selected = cursor
	s.rows.Reorder = nil
	if !disabled {
		s.rows.Reorder = func(rows []int, to int) { v.Moved, v.Rows, v.To = true, slices.Clone(rows), to }
	}
	if len(items) == 0 {
		ui.Text(c, core.Msg(c, "watchlist.empty")).TextColor(k.TextMuted).Padding(d.Unit() * 3)
		return false, false
	}
	restore := data.DataSelectTheme(c, 0)
	list := ui.List(c, &s.rows, len(items), func(i int) {
		it := items[i]
		ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).Gap(d.Unit() * 3).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icons.Must("grip-vertical")).Size(14, 14).TextColor(k.TextMuted)
			ui.Column(c).Grow(1).MinWidth(80).Children(func() {
				ui.Text(c, it.Symbol).Font(theme.MonoFont).FontWeight(700).SingleLine()
				if it.Name != "" {
					ui.Text(c, it.Name).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
				}
			})
			if !s.Hidden[WatchPrice] {
				ui.Row(c).Width(96).Justify(ui.End).Children(func() { PriceText(c, it.Price, PriceTextOptions{Decimals: dec}) })
			}
			if !s.Hidden[WatchChange] {
				ui.Row(c).Width(150).Justify(ui.End).Children(func() {
					PriceChangeBadge(c, it.Change, it.ChangePercent, PriceChangeBadgeOptions{Decimals: dec})
				})
			}
			if !s.Hidden[WatchVolume] {
				ui.Text(c, data.StatisticFormat(it.Volume, 0)).Width(90).TextAlign(ui.End).FontFeatures("tnum").TextColor(k.TextMuted)
			}
			if items[i].Symbol == s.Selected {
				ui.Icon(c, icons.Must("check")).Size(14, 14).TextColor(k.AccentText)
			}
		})
	}).Grow(1).MinHeight(0).Padding(d.Unit()).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius)
	restore()
	internal.Ring(list, k.Focus, theme.ControlRadius)
	if !disabled && s.cursor != before && s.cursor >= 0 {
		s.Selected, changed = items[s.cursor].Symbol, true
	}
	return changed, list.Submitted() && !disabled
}
