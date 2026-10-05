package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// MenubarItemKind is what a menu item does.
type MenubarItemKind int

const (
	// MenubarAction runs a command.
	MenubarAction MenubarItemKind = iota
	// MenubarCheck toggles a setting the caller keeps in Checked.
	MenubarCheck
	// MenubarRadio picks one of a group; the caller sets Checked on one.
	MenubarRadio
	// MenubarSeparator draws a line between groups of items.
	MenubarSeparator
)

// MenubarShortcut is the key hint shown beside an item; Key zero shows none.
type MenubarShortcut struct {
	Mods ui.Modifiers
	Key  ui.Key
}

// MenubarItem is an item of a menu, or with Items a submenu.
type MenubarItem struct {
	// ID is reported by Chosen; required for actions, checks and radios.
	ID       string
	Label    string
	Kind     MenubarItemKind
	Checked  bool
	Disabled bool
	// Shortcut is only shown; handle the key with c.Shortcut.
	Shortcut MenubarShortcut
	Items    []MenubarItem
}

// MenubarMenu is a title of the bar and the menu it opens.
type MenubarMenu struct {
	Label string
	Items []MenubarItem
}

// MenubarOptions configure a Menubar.
type MenubarOptions struct {
	// Label names the bar; empty uses the built-in name.
	Label string
}

// MenubarResult is the built Menubar.
type MenubarResult struct {
	Element *ui.Element
	chosen  string
}

// Chosen returns the ID of the item the user chose in this frame.
func (r MenubarResult) Chosen() (string, bool) { return r.chosen, r.chosen != "" }

type menubarState struct {
	open    int
	hl      []int
	focusTo int
	ptr     [2]float32
	// reveal scrolls the highlighted item into view once.
	reveal bool
}

// Menubar is an in-window menu bar drawn by MujicaUI, with submenus and full keyboard control.
func Menubar(c *ui.Context, menus []MenubarMenu, opts MenubarOptions) MenubarResult {
	if len(menus) == 0 {
		panic("mujica: Menubar needs menus")
	}
	for _, m := range menus {
		if m.Label == "" {
			panic("mujica: MenubarMenu needs a Label")
		}
		menubarCheck(m.Items)
	}
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	label := opts.Label
	if label == "" {
		label = msg(c, "menubar.label")
	}
	bar := ui.Row(c).AlignItems(ui.Center).Gap(u * 0.5).Padding(u).FocusGroup(ui.Horizontal).
		Role(ui.RoleToolbar).Label(label).Background(k.Surface)
	bar.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - 1, W: r.W, H: 1}, k.Border, 0)
	})
	st := ui.Local(bar, "state", func() *menubarState { return &menubarState{open: -1, focusTo: -1} })
	m := &menubarRun{c: c, menus: menus, st: *st}
	res := MenubarResult{Element: bar}
	titles := make([]ui.Rect, len(menus))
	bar.Children(func() {
		for i, menu := range menus {
			t := m.title(i, menu.Label)
			titles[i] = t.Bounds()
		}
	})
	if m.st.open >= 0 && bar.Shortcut(0, ui.KeyEscape) {
		m.close()
	}
	if m.st.open >= 0 {
		m.overlay(titles)
	}
	// Presses in the menu refocus the title, so lost focus means Tab left the menubar.
	if m.st.open >= 0 && m.st.focusTo < 0 && !bar.FocusWithin() {
		slog.Debug("mujica: menubar close on focus out")
		m.st.open, m.st.hl = -1, nil
	}
	// Focus stays on the title, so the highlighted item is read out.
	if name := m.current(); m.moved && name != "" {
		c.Announce(name)
	}
	res.chosen = m.chosen
	return res
}

// menubarCheck panics on items that cannot work.
func menubarCheck(items []MenubarItem) {
	for _, it := range items {
		if it.Kind == MenubarSeparator {
			continue
		}
		if it.Label == "" {
			panic("mujica: MenubarItem needs a Label")
		}
		if it.Items != nil {
			if it.Kind != MenubarAction {
				panic("mujica: a submenu item must be of kind MenubarAction")
			}
			menubarCheck(it.Items)
		} else if it.ID == "" {
			panic("mujica: MenubarItem " + it.Label + " needs an ID")
		}
	}
}

// menubarRun is one frame of a Menubar.
type menubarRun struct {
	c      *ui.Context
	menus  []MenubarMenu
	st     *menubarState
	chosen string
	moved  bool
}

// current names the highlighted item of the open menu, "" for none.
func (m *menubarRun) current() string {
	st := m.st
	if st.open < 0 || len(st.hl) == 0 {
		return ""
	}
	l := len(st.hl) - 1
	if h := st.hl[l]; h >= 0 {
		return menubarName(m.c, m.items(l)[h])
	}
	return ""
}

func (m *menubarRun) title(i int, label string) *ui.Element {
	c, st := m.c, m.st
	k := Tokens(c)
	d := Density(c)
	open := st.open == i
	t := ui.ButtonBase(c).MinHeight(d.ControlHeight()-4).Padding(0, d.Unit()*3).Radius(theme.ControlRadius).
		Role(ui.RoleMenuButton).Label(label).TextColor(k.Text)
	internal.Ring(t, k.Focus, theme.ControlRadius)
	switch {
	case open:
		t.Background(k.Selection)
	case t.Hovered():
		t.Background(k.SurfaceHover)
	}
	internal.Ease(t, Motion(c, theme.StateDuration), k.Text)
	if st.focusTo == i {
		t.Focus()
		st.focusTo = -1
	}
	if st.open < 0 {
		switch {
		case t.Shortcut(0, ui.KeyRight):
			m.focus((i + 1) % len(m.menus))
		case t.Shortcut(0, ui.KeyLeft):
			m.focus((i - 1 + len(m.menus)) % len(m.menus))
		case t.Shortcut(0, ui.KeyDown), t.Shortcut(0, ui.KeyUp):
			m.openMenu(i, true)
		}
	} else if open {
		m.keys(t)
		t.Description(m.current())
	}
	if t.Clicked() {
		if open {
			m.close()
		} else {
			m.openMenu(i, t.FocusVisible())
		}
	}
	t.Children(func() { ui.Text(c, label).SingleLine() })
	return t
}

// keys drives an open menu from the focused title.
func (m *menubarRun) keys(t *ui.Element) {
	st := m.st
	l := len(st.hl) - 1
	items := m.items(l)
	h := st.hl[l]
	m.moved, st.reveal = true, true
	switch {
	case t.Shortcut(0, ui.KeyDown):
		st.hl[l] = menubarStep(items, h, 1)
	case t.Shortcut(0, ui.KeyUp):
		st.hl[l] = menubarStep(items, h, -1)
	case t.Shortcut(0, ui.KeyHome):
		st.hl[l] = menubarStep(items, -1, 1)
	case t.Shortcut(0, ui.KeyEnd):
		st.hl[l] = menubarStep(items, len(items), -1)
	case t.Shortcut(0, ui.KeyRight):
		if h >= 0 && items[h].Items != nil {
			st.hl = append(st.hl, menubarStep(items[h].Items, -1, 1))
		} else {
			m.openMenu((st.open+1)%len(m.menus), true)
		}
	case t.Shortcut(0, ui.KeyLeft):
		if l > 0 {
			st.hl = st.hl[:l]
		} else {
			m.openMenu((st.open-1+len(m.menus))%len(m.menus), true)
		}
	case t.Shortcut(0, ui.KeyEnter), t.Shortcut(0, ui.KeySpace):
		if h >= 0 {
			m.activate(l, h)
		}
	case t.Shortcut(0, ui.KeyEscape):
		if l > 0 {
			st.hl = st.hl[:l]
		} else {
			m.close()
		}
	default:
		m.moved, st.reveal = false, false
	}
}

// items returns the items shown at level l of the open menu.
func (m *menubarRun) items(l int) []MenubarItem {
	items := m.menus[m.st.open].Items
	for i := 0; i < l; i++ {
		items = items[m.st.hl[i]].Items
	}
	return items
}

func (m *menubarRun) focus(i int) {
	m.st.focusTo = i
	m.c.Invalidate()
}

func (m *menubarRun) openMenu(i int, highlight bool) {
	first := -1
	if highlight {
		first = menubarStep(m.menus[i].Items, -1, 1)
	}
	m.st.open, m.st.hl = i, []int{first}
	m.moved = highlight
	m.focus(i)
	slog.Debug("mujica: menubar open", "menu", m.menus[i].Label)
}

func (m *menubarRun) close() {
	if m.st.open >= 0 {
		m.focus(m.st.open)
	}
	m.st.open, m.st.hl = -1, nil
	slog.Debug("mujica: menubar close")
}

// activate chooses item h of level l, or opens it when it is a submenu.
func (m *menubarRun) activate(l, h int) {
	it := m.items(l)[h]
	if it.Disabled || it.Kind == MenubarSeparator {
		return
	}
	if it.Items != nil {
		m.st.hl = append(m.st.hl[:l+1], menubarStep(it.Items, -1, 1))
		m.focus(m.st.open)
		return
	}
	m.chosen = it.ID
	slog.Debug("mujica: menubar choose", "id", it.ID)
	m.close()
}

// menubarStep returns the next enabled item from h in direction dir, wrapping.
func menubarStep(items []MenubarItem, h, dir int) int {
	n := len(items)
	for k := 1; k <= n; k++ {
		j := ((h+dir*k)%n + n) % n
		if it := items[j]; it.Kind != MenubarSeparator && !it.Disabled {
			return j
		}
	}
	return -1
}
