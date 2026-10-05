package mujica

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// DataColumn is a column of a DataTable showing rows of type R.
type DataColumn[R any] struct {
	// ID identifies the column in sort, layout and visibility; required, unique.
	ID    string
	Title string
	// Width is in DIPs; 0 shares the room left. MinWidth bounds resizing.
	Width, MinWidth float32
	Align           ui.Align
	// Sortable lets a header click sort by the column; local mode needs Compare.
	Sortable bool
	// Compare orders two rows, as cmp.Compare does.
	Compare func(a, b R) int
	// Text is the cell's plain text, for filtering and the default cell; required.
	Text func(r R) string
	// Cell builds a custom cell; nil shows Text.
	Cell func(r R)
}

// DataTableOptions configure a DataTable.
type DataTableOptions[R any, K comparable] struct {
	Columns []DataColumn[R]
	// Key returns a row's stable ID; required.
	Key      func(r R) K
	Multiple bool
	// Controlled shows rows as given: the caller sorts and filters on SortChanged.
	Controlled bool
	// Version must change when rows change in place; a new slice is noticed.
	Version int
	// Empty replaces the built-in "No items" text.
	Empty  string
	Status DataStatus
	Error  string
}

// DataTableState holds a DataTable's sort, filter, layout and selection.
type DataTableState[K comparable] struct {
	Sort ui.SortOrder
	// Filter keeps rows whose cells contain it, ignoring case, in local mode.
	Filter string
	// Hidden names the columns not shown, by ID.
	Hidden map[string]bool
	// Layout is the column order and widths the user set.
	Layout ui.TableLayout
	// Selection holds the chosen row IDs of a Multiple table.
	Selection ui.Selection[K]
	rows      ui.ListState
	cursor    int
	started   bool
	sel       K
	has       bool
	toggled   bool
	cache     dataTableCache[K]
}

// dataTableCache is the sorted, filtered view, rebuilt only when its inputs change.
type dataTableCache[K comparable] struct {
	valid   bool
	first   any
	n       int
	version int
	sort    ui.SortOrder
	filter  string
	view    []int
	at      map[K]int
}

// Selected returns the chosen row's ID, or the last chosen of a Multiple table.
func (s *DataTableState[K]) Selected() (K, bool) { return s.sel, s.has }

// Rows returns how many rows the table shows after filtering.
func (s *DataTableState[K]) Rows() int { return len(s.cache.view) }

// DataTableView is a built DataTable and what happened to it in this frame.
type DataTableView struct {
	*ui.Element
	changed, submitted, sorted, retry bool
}

// Changed reports a new selection.
func (v DataTableView) Changed() bool { return v.changed }

// Submitted reports Enter or a double click on a row.
func (v DataTableView) Submitted() bool { return v.submitted }

// SortChanged reports a header click or menu choice that changed State.Sort.
func (v DataTableView) SortChanged() bool { return v.sorted }

// Retried reports a click on the retry action of a failed load.
func (v DataTableView) Retried() bool { return v.retry }

const dataTableMarkID = "mujica.select"

// DataTable shows rows in columns, building only those in view.
func DataTable[R any, K comparable](c *ui.Context, s *DataTableState[K], rows []R, opts DataTableOptions[R, K]) DataTableView {
	dataTableCheck(opts)
	if !s.started {
		s.started, s.cursor = true, -1
	}
	k := Tokens(c)
	d := Density(c)
	before := s.Sort
	var v DataTableView
	v.Element = ui.Column(c).MinHeight(0).Gap(d.Unit() * 2)
	v.Children(func() {
		if len(rows) > 0 {
			ui.Row(c).Justify(ui.End).Children(func() {
				tableColumnMenu(c, dataTableMenu(opts), &s.Sort, &s.Layout, &s.Hidden)
			})
		}
		dataTableRefresh(s, rows, opts)
		if s.has {
			if i, ok := s.cache.at[s.sel]; ok {
				s.cursor = i
			} else {
				s.cursor = -1
			}
		}
		view := s.cache.view
		if len(view) == 0 {
			empty := opts.Empty
			if len(rows) > 0 {
				empty = msg(c, "data_table.noMatch")
			}
			v.retry = dataPlaceholder(c, opts.Status, opts.Error, empty, msg(c, "list.empty"))
			return
		}
		vis, cols := dataTableColumns(c, s, opts)
		s.rows.Key = func(i int) any { return opts.Key(rows[view[i]]) }
		s.rows.Label = func(i int) string { return opts.Columns[vis[0]].Text(rows[view[i]]) }
		s.rows.Selected, s.rows.Sort, s.rows.Columns = &s.cursor, &s.Sort, s.Layout
		s.rows.Selection = nil
		if opts.Multiple {
			s.rows.Selection = &s.Selection
		}
		// Rows reach the data-row height of the density.
		restore := dataSelectTheme(c, d.ControlHeight()/8)
		head := s.Sort
		e := ui.Table(c, &s.rows, cols, len(view), func(row, col int) {
			r := rows[view[row]]
			if col == 0 {
				dataTableMark(c, s, opts, opts.Key(r), row)
				return
			}
			dc := opts.Columns[vis[col-1]]
			if dc.Cell != nil {
				dc.Cell(r)
			} else {
				ui.Text(c, dc.Text(r)).SingleLine()
			}
		})
		restore()
		e.Grow(1).MinHeight(0).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius)
		internal.Ring(e, k.Focus, theme.ControlRadius)
		s.Layout = s.rows.Columns
		v.changed = e.Changed() && s.Sort == head || s.toggled
		s.toggled = false
		v.submitted = e.Submitted()
		v.retry = dataFooter(c, opts.Status, opts.Error)
		if s.cursor >= 0 && s.cursor < len(view) {
			s.sel, s.has = opts.Key(rows[view[s.cursor]]), true
		}
	})
	v.sorted = s.Sort != before
	if v.sorted {
		slog.Debug("mujica: data_table sort", "column", s.Sort.Column, "descending", s.Sort.Descending)
	}
	if v.changed {
		slog.Debug("mujica: data_table select", "count", s.Selection.Len())
	}
	return v
}

func dataTableCheck[R any, K comparable](opts DataTableOptions[R, K]) {
	if opts.Key == nil || len(opts.Columns) == 0 {
		panic("mujica: DataTable needs Key and Columns")
	}
	seen := map[string]bool{}
	for _, col := range opts.Columns {
		if col.ID == "" || seen[col.ID] || col.ID == dataTableMarkID || col.Text == nil {
			panic("mujica: DataTable column needs a unique ID and Text: " + col.ID)
		}
		if col.Sortable && !opts.Controlled && col.Compare == nil {
			panic("mujica: sortable DataTable column needs Compare: " + col.ID)
		}
		seen[col.ID] = true
	}
}

// dataTableRefresh rebuilds the view only when rows, sort or filter changed.
func dataTableRefresh[R any, K comparable](s *DataTableState[K], rows []R, opts DataTableOptions[R, K]) {
	cc := &s.cache
	var first any
	if len(rows) > 0 {
		first = &rows[0]
	}
	sort, filter := s.Sort, s.Filter
	if opts.Controlled {
		sort, filter = ui.SortOrder{}, ""
	}
	data := !cc.valid || cc.first != first || cc.n != len(rows) || cc.version != opts.Version
	if !data && cc.sort == sort && cc.filter == filter {
		return
	}
	var view []int
	if opts.Controlled {
		view = make([]int, len(rows))
		for i := range view {
			view[i] = i
		}
	} else {
		view = DataTableRows(rows, opts.Columns, sort, filter)
	}
	at := make(map[K]int, len(view))
	for i, j := range view {
		at[opts.Key(rows[j])] = i
	}
	if len(at) != len(view) {
		panic("mujica: DataTable rows with duplicate keys")
	}
	*cc = dataTableCache[K]{valid: true, first: first, n: len(rows), version: opts.Version, sort: sort, filter: filter, view: view, at: at}
	if data {
		dataTablePrune(s, rows, opts)
	}
	slog.Debug("mujica: data_table view", "rows", len(view), "of", len(rows))
}

// dataTablePrune drops chosen IDs whose rows left the data.
func dataTablePrune[R any, K comparable](s *DataTableState[K], rows []R, opts DataTableOptions[R, K]) {
	all := s.cache.at
	if len(s.cache.view) != len(rows) {
		all = make(map[K]int, len(rows))
		for i, r := range rows {
			all[opts.Key(r)] = i
		}
	}
	for key := range s.Selection.All() {
		if _, ok := all[key]; !ok {
			s.Selection.Remove(key)
		}
	}
	if _, ok := all[s.sel]; s.has && !ok {
		s.has = false
	}
}

// dataTableMenu lists the columns for the column menu.
func dataTableMenu[R any, K comparable](opts DataTableOptions[R, K]) []tableColumn {
	cols := make([]tableColumn, len(opts.Columns))
	for i, col := range opts.Columns {
		cols[i] = tableColumn{ID: col.ID, Title: col.Title, Width: col.Width, MinWidth: col.MinWidth, Sortable: col.Sortable}
	}
	return cols
}

func dataTableColumns[R any, K comparable](c *ui.Context, s *DataTableState[K], opts DataTableOptions[R, K]) ([]int, []ui.TableColumn) {
	w := Density(c).ControlHeight()
	cols := []ui.TableColumn{{ID: dataTableMarkID, Width: w, Fixed: true, Align: ui.Center}}
	var vis []int
	for i, col := range opts.Columns {
		if s.Hidden[col.ID] {
			continue
		}
		vis = append(vis, i)
		cols = append(cols, ui.TableColumn{ID: col.ID, Title: col.Title, Width: col.Width, MinWidth: col.MinWidth, Align: col.Align, Sortable: col.Sortable})
	}
	if len(vis) == 0 {
		panic("mujica: DataTable with every column hidden")
	}
	return vis, cols
}

// dataTableMark draws the selection column: a check box, or a diamond on the chosen row.
func dataTableMark[R any, K comparable](c *ui.Context, s *DataTableState[K], opts DataTableOptions[R, K], key K, row int) {
	k := Tokens(c)
	if !opts.Multiple {
		if row == s.cursor {
			ui.Box(c).Size(8, 8).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
				p.FillPath(internal.Diamond(internal.Snap(p, r.X+r.W/2), internal.Snap(p, r.Y+r.H/2), r.W), k.AccentText)
			})
		}
		return
	}
	on := s.Selection.Has(key)
	cb := ui.CheckboxBase(c, &on).Label(msg(c, "data_table.selectRow")).Size(theme.MinHitSize, theme.MinHitSize).
		Justify(ui.Center).AlignItems(ui.Center).Radius(theme.ControlRadius)
	internal.Ring(cb, k.Focus, theme.ControlRadius)
	cb.Children(func() { dataTableCheckBox(c, on) })
	if cb.Changed() {
		if on {
			s.Selection.Add(key)
		} else {
			s.Selection.Remove(key)
		}
		s.toggled = true
		slog.Debug("mujica: data_table toggle", "on", on)
	}
}

// DataTableRows sorts and filters rows as a local DataTable does; it returns indices.
func DataTableRows[R any](rows []R, cols []DataColumn[R], sort ui.SortOrder, filter string) []int {
	view := make([]int, 0, len(rows))
	q := strings.ToLower(filter)
	for i, r := range rows {
		if q == "" || dataTableMatch(r, cols, q) {
			view = append(view, i)
		}
	}
	for _, col := range cols {
		if col.ID == sort.Column && col.Compare != nil {
			cmp := col.Compare
			slices.SortStableFunc(view, func(a, b int) int {
				if sort.Descending {
					return cmp(rows[b], rows[a])
				}
				return cmp(rows[a], rows[b])
			})
		}
	}
	return view
}

func dataTableMatch[R any](r R, cols []DataColumn[R], q string) bool {
	for _, col := range cols {
		if strings.Contains(strings.ToLower(col.Text(r)), q) {
			return true
		}
	}
	return false
}
