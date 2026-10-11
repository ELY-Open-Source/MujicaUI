package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// PopupMenuItemOptions configure one item of a PopupMenu.
type PopupMenuItemOptions struct {
	Disabled bool
	// KeepOpen keeps the menu open after the item is chosen.
	KeepOpen bool
}

// PopupMenu builds the items of a DropdownMenu or ContextMenu, reporting choices as they happen.
type PopupMenu struct {
	c         *ui.Context
	root      *popupMenuRoot
	level     *popupMenuLevel
	subs      []popupMenuSub
	seq, rows int
	autofocus bool
}

// popupMenuRoot carries what every level of one open menu shares in a frame.
type popupMenuRoot struct{ close bool }

// popupMenuLevel is what one menu level keeps while it shows.
type popupMenuLevel struct {
	open  int  // seq of the open submenu, 0 for none
	byKey bool // the open submenu was opened from the keyboard
}

type popupMenuSub struct {
	item  *ui.Element
	seq   int
	build func(m *PopupMenu)
}

// Item adds an item and reports its choice, which closes the menu unless KeepOpen.
func (m *PopupMenu) Item(label string, opts PopupMenuItemOptions) bool {
	it := m.row(ui.ButtonBase(m.c), opts.Disabled, false)
	m.content(it, label, false, false)
	if !it.Clicked() {
		return false
	}
	slog.Debug("mujica: popup_menu choose", "item", label)
	m.root.close = m.root.close || !opts.KeepOpen
	return true
}

// Check adds an item toggling *checked and reports whether it changed.
func (m *PopupMenu) Check(label string, checked *bool, opts PopupMenuItemOptions) bool {
	// The base toggles as it is built, so a disabled wrapper comes first.
	var base *ui.Element
	ui.Column(m.c).Shrink(0).Disabled(opts.Disabled).Children(func() { base = ui.CheckboxBase(m.c, checked) })
	it := m.row(base, opts.Disabled, false)
	m.content(it, label, *checked, false)
	if opts.Disabled || !it.Changed() {
		return false
	}
	slog.Debug("mujica: popup_menu check", "item", label, "checked", *checked)
	m.root.close = m.root.close || !opts.KeepOpen
	return true
}

// Submenu adds an item opening a nested menu that build fills.
func (m *PopupMenu) Submenu(label string, opts PopupMenuItemOptions, build func(m *PopupMenu)) {
	if build == nil {
		panic("mujica: PopupMenu Submenu needs a build function")
	}
	it := m.row(ui.ButtonBase(m.c), opts.Disabled, true)
	seq := m.seq
	m.content(it, label, false, true)
	if opts.Disabled {
		return
	}
	if it.Clicked() || it.Shortcut(0, ui.KeyRight) {
		m.level.open, m.level.byKey = seq, it.FocusVisible()
		slog.Debug("mujica: popup_menu submenu", "item", label)
	}
	m.subs = append(m.subs, popupMenuSub{it, seq, build})
}

// Group adds a captioned group of items, separated from what comes before.
func (m *PopupMenu) Group(label string, fn func()) {
	k := core.Tokens(m.c)
	u := core.Density(m.c).Unit()
	if m.rows > 0 {
		m.Separator()
	}
	if label != "" {
		ui.Text(m.c, label).FontSize(core.FontSize(m.c, theme.CaptionSize)).TextColor(k.TextMuted).Padding(u, u*3, u/2, u*3+16+u*2).Shrink(0)
	}
	fn()
}

// Separator adds a line between items.
func (m *PopupMenu) Separator() {
	u := core.Density(m.c).Unit()
	ui.Box(m.c).Height(theme.BorderWidth).Margin(u, u*2).Background(core.Tokens(m.c).Border).Shrink(0)
}

// row styles an item element and runs the hover, focus and highlight logic.
func (m *PopupMenu) row(it *ui.Element, disabled, sub bool) *ui.Element {
	c := m.c
	k := core.Tokens(c)
	d := core.Density(c)
	m.seq++
	m.rows++
	seq := m.seq
	it.Justify(ui.Start).MinHeight(max(theme.MinHitSize, d.ControlHeight()-4)).Padding(0, d.Unit()*3).
		Gap(d.Unit() * 2).Radius(theme.ControlRadius).Disabled(disabled).TextColor(k.Text)
	internal.Ring(it, k.Focus, theme.ControlRadius)
	if disabled {
		return it
	}
	if m.autofocus {
		it.AutoFocus()
		m.autofocus = false
	}
	hovered := it.Hovered()
	was := ui.Local(it, "mujica.menu.hover", func() bool { return false })
	if hovered && !*was {
		it.Focus() // focus follows the pointer
		if sub {
			m.level.open, m.level.byKey = seq, false
		} else {
			m.level.open = 0
		}
	}
	*was = hovered
	if hovered || it.FocusVisible() || m.level.open == seq {
		it.Background(k.Selection).Draw(func(p *ui.Painter, r ui.Rect) {
			p.Fill(ui.Rect{X: r.X, Y: r.Y + 4, W: 2, H: r.H - 8}, k.AccentText, 1)
		})
	}
	return it
}

// content lays out the check column, label and submenu chevron of an item.
func (m *PopupMenu) content(it *ui.Element, label string, checked, sub bool) {
	c := m.c
	k := core.Tokens(c)
	it.Children(func() {
		mark := ui.Box(c).Size(16, 16).Shrink(0)
		if checked {
			mark.Children(func() { ui.Icon(c, icons.Must("check")).Size(16, 16).TextColor(k.AccentText) })
		}
		ui.Text(c, label).SingleLine().Grow(1)
		if sub {
			ui.Icon(c, icons.Must("chevron-right")).Size(14, 14).TextColor(k.TextMuted).Shrink(0)
		}
	})
}

// PopupMenuOpen shows the menu build fills below anchor while *open is true.
func PopupMenuOpen(c *ui.Context, anchor *ui.Element, open *bool, build func(m *PopupMenu)) *ui.Element {
	if !*open {
		return nil
	}
	core.OverlayLayer(c, open)
	b := anchor.Bounds()
	root := &popupMenuRoot{}
	panel := ui.PopoverBase(c, anchor, open, func(panel *ui.Element) {
		core.OverlayStyle(c, panel, OverlayAnchored(c, panel, b))
		panel.Margin(4, 0, 0, 0).MinWidth(PopoverWide(c, max(b.W, 200))).MaxWidth(PopoverWide(c, 360)).MaxHeight(PopoverRoom(c, b))
		popupMenuLevelBuild(c, panel, root, build, true)
		// The panel holds every level, so Tab in a submenu closes the whole menu too.
		if panel.Shortcut(0, ui.KeyTab) || panel.Shortcut(ui.Shift, ui.KeyTab) {
			root.close = true
		}
	})
	if root.close {
		*open = false
	}
	if !*open {
		slog.Debug("mujica: popup_menu close")
	}
	return panel
}

// popupMenuLevelBuild builds one level's items in outer, then its open submenu.
func popupMenuLevelBuild(c *ui.Context, outer *ui.Element, root *popupMenuRoot, build func(m *PopupMenu), autofocus bool) {
	u := core.Density(c).Unit()
	// The padding leaves room for the items' focus rings inside the clipping scroll.
	items := ui.Scroll(c).Shrink(1).Padding(u + theme.FocusGap + theme.FocusWidth).FocusGroup(ui.Vertical)
	m := &PopupMenu{c: c, root: root, autofocus: autofocus}
	m.level = ui.Local(items, "mujica.menu.level", func() popupMenuLevel { return popupMenuLevel{} })
	outer.Role(ui.RolePopup)
	items.Children(func() { build(m) })
	for _, s := range m.subs {
		if m.level.open == s.seq {
			popupMenuSubmenu(c, outer, root, m.level, s)
		}
	}
}

// popupMenuSubmenu builds a submenu beside its item, flipped to stay in the window.
func popupMenuSubmenu(c *ui.Context, outer *ui.Element, root *popupMenuRoot, lv *popupMenuLevel, s popupMenuSub) {
	u := core.Density(c).Unit()
	ob, ib := outer.Bounds(), s.item.Bounds()
	ww, wh := c.Size()
	box := ui.Column(c).Absolute().Key(s.seq).MinWidth(min(180, ww-8)).MaxWidth(min(360, ww-8))
	sb := box.Bounds()
	left, top := ib.X+ib.W+u-ob.X, ib.Y-u-ob.Y
	if ob.X+left+sb.W > ww-4 {
		left = ib.X - u - sb.W - ob.X
	}
	// Flipped or not, the submenu stays inside the window.
	left = max(4-ob.X, min(left, ww-4-sb.W-ob.X))
	top = max(4-ob.Y, min(top, wh-4-sb.H-ob.Y))
	box.Left(left).Top(top).MaxHeight(wh - 8)
	p := core.OverlayEnter(c, box)
	if sb.W == 0 || ob.W == 0 {
		p = 0 // not measured yet
		c.AnimationFrame()
	}
	core.OverlayStyle(c, box, p)
	box.Children(func() { popupMenuLevelBuild(c, box, root, s.build, lv.byKey) })
	if box.Shortcut(0, ui.KeyLeft) || box.Shortcut(0, ui.KeyEscape) {
		lv.open = 0
		s.item.Focus()
	}
}

// popupMenuTrigger styles a menu trigger button with a trailing chevron.
func popupMenuTrigger(c *ui.Context, b *ui.Element, label string, icon *ui.SVG, open bool) {
	k := core.Tokens(c)
	d := core.Density(c)
	b.MinHeight(d.ControlHeight()).MinWidth(theme.MinHitSize).Padding(0, d.Unit()*3, 0, d.Unit()*4).
		Gap(d.Unit() * 2).Radius(theme.ControlRadius)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	bg, fg, border := input.ButtonColors(k, input.Secondary, b.Hovered() && !b.IsDisabled(), b.Pressed() || open)
	b.Background(bg).TextColor(fg).Border(theme.BorderWidth, border)
	b.Children(func() {
		if icon != nil {
			ui.Icon(c, icon).Size(16, 16)
		}
		if label != "" {
			ui.Text(c, label).SingleLine()
		}
		ui.Icon(c, icons.Must("chevron-down")).Size(14, 14).TextColor(k.TextMuted)
	})
}
