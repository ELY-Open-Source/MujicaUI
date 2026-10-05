package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// DataTableColumnPicker is a Columns button whose popover shows and hides columns.
// The last visible column cannot be hidden.
func DataTableColumnPicker[R any, K comparable](c *ui.Context, s *DataTableState[K], cols []DataColumn[R]) *ui.Element {
	k := Tokens(c)
	u := Density(c).Unit()
	b := Button(c, msg(c, "data_table.columns"), ButtonOptions{Variant: Secondary, Icon: icons.Must("eye")})
	open := ui.Local(b, "open", func() bool { return false })
	if b.Clicked() {
		*open = !*open
	}
	visible := 0
	for _, col := range cols {
		if !s.Hidden[col.ID] {
			visible++
		}
	}
	ui.PopoverBase(c, b, open, func(panel *ui.Element) {
		panel.Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).Radius(theme.PopupRadius).
			Padding(u * 2).Gap(u).MinWidth(200).MarginY(u)
		for _, col := range cols {
			ui.Column(c).Key(col.ID).Children(func() { dataTablePickerItem(c, s, col.ID, col.Title, visible) })
		}
	})
	return b
}

// dataTablePickerItem is one column's check box in the picker.
func dataTablePickerItem[K comparable](c *ui.Context, s *DataTableState[K], id, title string, visible int) {
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	on := !s.Hidden[id]
	last := on && visible == 1
	cb := ui.CheckboxBase(c, &on).Disabled(last).MinHeight(d.ControlHeight()).
		Padding(0, u*2).Gap(u * 2).AlignItems(ui.Center).Radius(theme.ControlRadius)
	internal.Ring(cb, k.Focus, theme.ControlRadius)
	if cb.Hovered() && !last {
		cb.Background(k.SurfaceHover)
	}
	if cb.Changed() && !last {
		if s.Hidden == nil {
			s.Hidden = map[string]bool{}
		}
		s.Hidden[id] = !on
		slog.Debug("mujica: data_table column", "id", id, "visible", on)
	}
	cb.Children(func() {
		dataTableCheckBox(c, on)
		ui.Text(c, title).SingleLine()
	})
}

// dataTableCheckBox draws a 16 DIP check box: OnAccent check on Accent when on.
func dataTableCheckBox(c *ui.Context, on bool) {
	k := Tokens(c)
	box := ui.Box(c).Size(16, 16).Shrink(0).Radius(theme.ControlRadius).Center()
	if on {
		box.Background(k.Accent).TextColor(k.OnAccent).Children(func() {
			ui.Icon(c, icons.Must("check")).Size(14, 14)
		})
	} else {
		box.Border(theme.BorderWidth, k.ControlBorder).Background(k.Surface)
	}
}
