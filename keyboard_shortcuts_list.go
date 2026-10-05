package mujica

import (
	"reflect"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// ShortcutGroup is a titled group of bindings in a KeyboardShortcutsList.
type ShortcutGroup struct {
	Title    string
	Bindings []ShortcutBinding
}

// KeyboardShortcutsListOptions configure a KeyboardShortcutsList.
type KeyboardShortcutsListOptions struct {
	// Label names the search field; empty uses the built-in one.
	Label string
}

type shortcutsListRow struct{ group, binding int }

type shortcutsListState struct {
	groups []ShortcutGroup
	query  string
	rows   []shortcutsListRow
	ready  bool
}

// KeyboardShortcutsList displays searchable, grouped shortcuts using visible rows.
func KeyboardShortcutsList(c *ui.Context, query *string, groups []ShortcutGroup, opts KeyboardShortcutsListOptions) *ui.Element {
	k := Tokens(c)
	d := Density(c)
	label := opts.Label
	if label == "" {
		label = msg(c, "keyboard_shortcuts_list.search")
	}
	e := ui.Column(c).Gap(d.Unit() * 4).MinWidth(320)
	st := ui.Local(e, "shortcuts", func() shortcutsListState { return shortcutsListState{} })
	e.Children(func() {
		SearchInput(c, query, SearchInputOptions{Placeholder: label, Label: label})
		shortcutsListPlan(st, groups, strings.ToLower(strings.TrimSpace(*query)))
		if len(st.rows) == 0 {
			EmptyState(c, EmptyStateOptions{Kind: EmptyNoResults, Query: *query})
			return
		}
		accountRows(c, len(st.rows), func(i int) any { return st.rows[i] }, func(i int) {
			r := st.rows[i]
			g := groups[r.group]
			if r.binding < 0 {
				ui.Text(c, g.Title).Font(theme.SerifFont).FontWeight(600).TextColor(k.Text).Padding(d.Unit()*2, 0).Role(ui.RoleGroup).Label(g.Title)
				return
			}
			b := g.Bindings[r.binding]
			ui.Row(c).MinHeight(d.ControlHeight()).AlignItems(ui.Center).Gap(d.Unit()*4).
				BorderWidth(0, 0, theme.BorderWidth, 0).BorderColor(k.Border).Role(ui.RoleRow).
				Label(b.Action + ", " + b.Shortcut.String()).Children(func() {
				ui.Text(c, b.Action).Grow(1).MinWidth(0).Role(ui.RoleNone)
				shortcutKbd(c, b.Shortcut)
			})
		})
	})
	return e
}

func shortcutsListPlan(st *shortcutsListState, groups []ShortcutGroup, q string) {
	if st.ready && st.query == q && reflect.DeepEqual(st.groups, groups) {
		return
	}
	st.ready, st.query = true, q
	st.groups = slices.Clone(groups)
	st.rows = st.rows[:0]
	for gi, g := range groups {
		st.groups[gi].Bindings = slices.Clone(g.Bindings)
		header := false
		for bi, b := range g.Bindings {
			if b.Action == "" {
				panic("mujica: KeyboardShortcutsList binding needs an Action")
			}
			if q != "" && !strings.Contains(strings.ToLower(b.Action), q) && !strings.Contains(strings.ToLower(b.Shortcut.String()), q) {
				continue
			}
			if !header && g.Title != "" {
				st.rows = append(st.rows, shortcutsListRow{gi, -1})
				header = true
			}
			st.rows = append(st.rows, shortcutsListRow{gi, bi})
		}
	}
}

func accountRows(c *ui.Context, n int, key func(int) any, row func(int)) *ui.Element {
	if n <= 8 {
		return ui.Column(c).Children(func() {
			for i := range n {
				ui.Column(c).Key(key(i)).Children(func() { row(i) })
			}
		})
	}
	site := ui.Column(c).MinHeight(0).Shrink(1)
	st := ui.Local(site, "rows", func() ui.ListState { return ui.ListState{} })
	st.Key = key
	site.Children(func() { ui.List(c, st, n, row).Height(320).MinHeight(0).Shrink(1) })
	return site
}
