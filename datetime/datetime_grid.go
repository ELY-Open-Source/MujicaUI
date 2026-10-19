package datetime

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// datetimeGridConfig describes a page of month or year cells, addressed by ordinal.
type datetimeGridConfig struct {
	// pageOf returns the first ordinal and count of the page holding o.
	pageOf  func(o int) (start, n int)
	cols    int
	page    int
	title   func(start int) string
	prev    string
	next    string
	text    func(o int) string
	label   func(o int) string
	chosen  func(o int) bool
	current func(o int) bool
	enabled func(o int) bool
	// muted greys cells of the page that belong to a neighbor page.
	muted     func(o, start int) bool
	autoFocus bool
}

// datetimeGrid draws the page of *cursor and returns the ordinal picked, if any.
func datetimeGrid(c *ui.Context, cursor *int, cfg datetimeGridConfig) (picked int, ok bool) {
	k := core.Tokens(c)
	d := core.Density(c)
	root := ui.Column(c).Gap(d.Unit() * 2).Shrink(0).Width(d.ControlHeight() * 7)
	keyed := ui.Local(root, "datetime-grid-keys", func() bool { return false })
	root.Children(func() {
		start, _ := cfg.pageOf(*cursor)
		ui.Row(c).AlignItems(ui.Center).Gap(d.Unit()).BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).
			PaddingY(d.Unit()).Children(func() {
			if core.IconAction(c, icons.Must("chevron-left"), cfg.prev).Clicked() {
				*cursor, *keyed = *cursor-cfg.page, false
			}
			ui.Text(c, cfg.title(start)).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.BodySize+2)).Bold().
				TextColor(k.Text).Grow(1).TextAlign(ui.Center).SingleLine()
			if core.IconAction(c, icons.Must("chevron-right"), cfg.next).Clicked() {
				*cursor, *keyed = *cursor+cfg.page, false
			}
		})
		cells := ui.Column(c).Gap(d.Unit()).FocusGroup(ui.Horizontal | ui.Vertical)
		moved := true
		switch {
		case cells.Shortcut(0, ui.KeyLeft):
			*cursor--
		case cells.Shortcut(0, ui.KeyRight):
			*cursor++
		case cells.Shortcut(0, ui.KeyUp):
			*cursor -= cfg.cols
		case cells.Shortcut(0, ui.KeyDown):
			*cursor += cfg.cols
		case cells.Shortcut(0, ui.KeyPageUp):
			*cursor -= cfg.page
		case cells.Shortcut(0, ui.KeyPageDown):
			*cursor += cfg.page
		case cells.Shortcut(0, ui.KeyHome):
			*cursor, _ = datetimeGridBounds(*cursor, cfg)
		case cells.Shortcut(0, ui.KeyEnd):
			_, *cursor = datetimeGridBounds(*cursor, cfg)
		default:
			moved = false
		}
		if moved {
			*keyed = true
		}
		start, n := cfg.pageOf(*cursor)
		w := d.ControlHeight() * 7 / float32(cfg.cols)
		cells.Children(func() {
			for r := 0; r < n; r += cfg.cols {
				ui.Row(c).Children(func() {
					for o := start + r; o < start+min(r+cfg.cols, n); o++ {
						if datetimeGridCell(c, o, start, *cursor, cfg, w, moved, *keyed) {
							picked, ok = o, true
							*cursor, *keyed = o, false
						}
					}
				})
			}
		})
	})
	return picked, ok
}

// datetimeGridCell draws one month or year cell and reports whether it was picked.
func datetimeGridCell(c *ui.Context, o, start, cursor int, cfg datetimeGridConfig, w float32, focus, keyed bool) bool {
	k := core.Tokens(c)
	d := core.Density(c)
	enabled, chosen := cfg.enabled(o), cfg.chosen(o)
	sel := o - 1
	if chosen {
		sel = o
	}
	cell := ui.RadioBase(c, &sel, o).Size(w, d.ControlHeight()).Center().Radius(theme.ControlRadius).
		Label(cfg.label(o)).Disabled(!enabled).FocusRing(false)
	if o == cursor && enabled {
		if focus {
			cell.Focus()
		} else if cfg.autoFocus {
			core.AutoFocus(c, cell)
		}
	}
	fg := k.Text
	if cfg.muted(o, start) {
		fg = k.TextMuted
	}
	switch {
	case chosen:
		cell.Background(k.Accent)
		fg = k.OnAccent
		if cell.Hovered() {
			cell.Background(k.AccentHover)
		}
	case cell.Pressed():
		cell.Background(k.Selection)
	case cell.Hovered():
		cell.Background(k.SurfaceHover)
	}
	current := cfg.current(o)
	ring := cell.FocusVisible() || cell.Focused() && keyed
	cell.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if chosen {
			core.DatetimeSelectedMark(p, r, k.OnAccent)
		}
		if current {
			p.FillPath(internal.Diamond(r.X+r.W/2, r.Y+r.H-6, 5), fg)
		}
		if ring {
			internal.FocusRing(p, r, theme.ControlRadius, k.Focus)
		}
	})
	cell.Children(func() {
		t := ui.Text(c, cfg.text(o)).TextColor(fg).FontFeatures("tnum").SingleLine()
		if current {
			t.Bold()
		}
		if !enabled {
			t.Strikethrough()
		}
	})
	return cell.Clicked() && enabled
}

// datetimeGridBounds returns the page's own first and last ordinals.
func datetimeGridBounds(o int, cfg datetimeGridConfig) (first, last int) {
	start, n := cfg.pageOf(o)
	first, last = start, start+n-1
	for first < last && cfg.muted(first, start) {
		first++
	}
	for last > first && cfg.muted(last, start) {
		last--
	}
	return first, last
}

// datetimeGridField is the field and popup shared by MonthPicker and YearPicker.
type datetimeGridField struct {
	text, placeholder, label string
	clearable, hasValue      bool
	disabled, readOnly       bool
	// home is the ordinal the popup opens on.
	home int
}

// datetimeGridPicker draws f and its popup; reports a pick or clear.
func datetimeGridPicker(c *ui.Context, f datetimeGridField, cfg datetimeGridConfig) (col ui.Element, picked int, ok, cleared bool) {
	d := core.Density(c)
	col = ui.Column(c).Width(d.ControlHeight() * 6).Gap(d.Unit()).AlignItems(ui.Stretch).Disabled(f.disabled)
	open := ui.Local(col, "datetime-grid-open", func() bool { return false })
	cursor := ui.Local(col, "datetime-grid-cursor", func() int { return 0 })
	col.Children(func() {
		field := core.DatetimeField(c, ui.Row(c), false, f.readOnly)
		field.Children(func() {
			t := core.DatetimeTrigger(c, f.text, f.placeholder, "calendar").Label(f.label).Description(f.text)
			if !f.readOnly && (t.Clicked() || !*open && t.Shortcut(0, ui.KeyDown)) {
				*open, *cursor = !*open, f.home
			}
			if f.clearable && f.hasValue && !f.readOnly && core.IconAction(c, icons.Must("x"), core.Msg(c, "calendar.clear")).Clicked() {
				cleared = true
			}
		})
		core.DatetimePopup(c, field, open, f.readOnly, func() {
			cfg.autoFocus = true
			if picked, ok = datetimeGrid(c, cursor, cfg); ok {
				*open = false
			}
		})
	})
	return col, picked, ok, cleared
}
