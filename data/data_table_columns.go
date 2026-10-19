package data

import (
	"fmt"
	"log/slog"
	"maps"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// tableColumn is what the column menu knows of a column.
type tableColumn struct {
	ID, Title       string
	Width, MinWidth float32
	Sortable        bool
}

// tableWidthStep is how far one Narrower or Wider press resizes a column.
const tableWidthStep float32 = 16

// tableColumnMenu is the keyboard path to sort, move, resize and, with hidden, show columns.
func tableColumnMenu(c *ui.Context, cols []tableColumn, sort *ui.SortOrder, layout *ui.TableLayout, hidden *map[string]bool) {
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	b := input.Button(c, core.Msg(c, "data_table.columns"), input.ButtonOptions{Variant: input.Secondary, Icon: icons.Must("settings")})
	open := ui.Local(b, "open", func() bool { return false })
	if b.Clicked() {
		*open = !*open
	}
	order := tableOrder(cols, *layout)
	visible := len(cols)
	if hidden != nil {
		for _, col := range cols {
			if (*hidden)[col.ID] {
				visible--
			}
		}
	}
	ui.PopoverBase(c, b, open, func(panel ui.Element) {
		panel.Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).Radius(theme.PopupRadius).
			Padding(u * 2).Gap(u).MinWidth(240).MarginY(u)
		for at, id := range order {
			col := cols[slices.IndexFunc(cols, func(col tableColumn) bool { return col.ID == id })]
			ui.Row(c.Key(id)).Gap(u).AlignItems(ui.Center).Children(func() {
				if hidden != nil {
					tableVisibility(c, col, hidden, visible)
				} else {
					ui.Text(c, col.Title).SingleLine().Grow(1).Padding(0, u*2)
				}
				tableColumnActions(c, col, at, order, sort, layout)
			})
		}
	})
}

// tableColumnActions are one column's sort, move and width buttons.
func tableColumnActions(c *ui.Context, col tableColumn, at int, order []string, sort *ui.SortOrder, layout *ui.TableLayout) {
	act := func(icon, key string, disabled bool) bool {
		return input.Button(c, "", input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must(icon), Label: fmt.Sprintf(core.Msg(c, key), col.Title), Disabled: disabled}).Clicked()
	}
	if col.Sortable {
		for _, desc := range []bool{false, true} {
			icon, key := "arrow-up", "data_table.sortAsc"
			if desc {
				icon, key = "arrow-down", "data_table.sortDesc"
			}
			if act(icon, key, *sort == ui.SortOrder{Column: col.ID, Descending: desc}) {
				*sort = ui.SortOrder{Column: col.ID, Descending: desc}
				slog.Debug("mujica: table sort", "column", col.ID, "descending", desc)
			}
		}
	}
	for _, step := range []int{-1, 1} {
		icon, key := "chevron-left", "data_table.moveLeft"
		if step > 0 {
			icon, key = "chevron-right", "data_table.moveRight"
		}
		to := at + step
		if act(icon, key, to < 0 || to >= len(order)) {
			ids := slices.Clone(order)
			ids[at], ids[to] = ids[to], ids[at]
			layout.Order = ids
			slog.Debug("mujica: table move", "column", col.ID, "to", to)
		}
	}
	w := tableWidth(col, *layout)
	least := max(col.MinWidth, tableWidthStep)
	if act("minus", "data_table.narrower", w <= least) {
		tableSetWidth(layout, col.ID, max(least, w-tableWidthStep))
	}
	if act("plus", "data_table.wider", false) {
		tableSetWidth(layout, col.ID, w+tableWidthStep)
	}
}

// tableVisibility is a column's show/hide check box; the last visible one cannot hide.
func tableVisibility(c *ui.Context, col tableColumn, hidden *map[string]bool, visible int) {
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	on := !(*hidden)[col.ID]
	last := on && visible == 1
	cb := ui.CheckboxBase(c, &on).Label(col.Title).Disabled(last).MinHeight(d.ControlHeight()).Grow(1).
		Padding(0, u*2).Gap(u * 2).AlignItems(ui.Center).Radius(theme.ControlRadius)
	internal.Ring(cb, k.Focus, theme.ControlRadius)
	if cb.Hovered() && !last {
		cb.Background(k.SurfaceHover)
	}
	if cb.Changed() && !last {
		m := maps.Clone(*hidden)
		if m == nil {
			m = map[string]bool{}
		}
		m[col.ID] = !on
		*hidden = m
		slog.Debug("mujica: table column", "id", col.ID, "visible", on)
	}
	cb.Children(func() {
		dataTableCheckBox(c, on)
		ui.Text(c, col.Title).SingleLine()
	})
}

// tableOrder lists column IDs as the layout orders them, the rest after.
func tableOrder(cols []tableColumn, l ui.TableLayout) []string {
	var ids []string
	for _, id := range l.Order {
		if slices.ContainsFunc(cols, func(col tableColumn) bool { return col.ID == id }) && !slices.Contains(ids, id) {
			ids = append(ids, id)
		}
	}
	for _, col := range cols {
		if !slices.Contains(ids, col.ID) {
			ids = append(ids, col.ID)
		}
	}
	return ids
}

// tableWidth is a column's set width; a sharing column starts from its least, at least 120.
func tableWidth(col tableColumn, l ui.TableLayout) float32 {
	if w, ok := l.Widths[col.ID]; ok {
		return w
	}
	if col.Width > 0 {
		return col.Width
	}
	return max(col.MinWidth, 120)
}

func tableSetWidth(l *ui.TableLayout, id string, w float32) {
	m := maps.Clone(l.Widths)
	if m == nil {
		m = map[string]float32{}
	}
	m[id] = w
	l.Widths = m
	slog.Debug("mujica: table width", "column", id, "width", w)
}

// dataTableCheckBox draws a 16 DIP check box: OnAccent check on Accent when on.
func dataTableCheckBox(c *ui.Context, on bool) {
	k := core.Tokens(c)
	box := ui.Box(c).Size(16, 16).Shrink(0).Radius(theme.ControlRadius).Center()
	if on {
		box.Background(k.Accent).TextColor(k.OnAccent).Children(func() {
			ui.Icon(c, icons.Must("check")).Size(14, 14)
		})
	} else {
		box.Border(theme.BorderWidth, k.ControlBorder).Background(k.Surface)
	}
}
