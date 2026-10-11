package data

import (
	"log/slog"
	"slices"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TreeColumn is a column of a TreeTable; the first shows the arrows.
type TreeColumn[K comparable] struct {
	// ID identifies the column in sort and layout; required, unique.
	ID    string
	Title string
	// Width is in DIPs; 0 shares the room left. MinWidth bounds resizing.
	Width, MinWidth float32
	Align           ui.Align
	// Compare orders siblings, as cmp.Compare does; nil makes the column unsortable.
	Compare func(a, b K) int
	// Cell builds the item's cell in this column; required.
	Cell func(item K)
}

// TreeTableOptions configure a TreeTable of items of type K.
type TreeTableOptions[K comparable] struct {
	Columns []TreeColumn[K]
	Roots   []K
	// Children returns kids (nil leaf, empty unloaded); required; return a new slice on change.
	Children func(item K) []K
	// Label returns an item's text, read out and matched by typing; required.
	Label func(item K) string
	// ItemStatus returns the load state of an item's children; nil means loaded.
	ItemStatus func(item K) DataStatus
	// Load asks for an item's children when it opens DataUnloaded, and on retry.
	Load     func(item K)
	Multiple bool
	Empty    string
	Status   DataStatus
	Error    string
}

// TreeTableState holds a TreeTable's sort, layout, open items and selection.
type TreeTableState[K comparable] struct {
	// Sort orders siblings by a column; the header shows it.
	Sort ui.SortOrder
	// Layout is the column order and widths the user set.
	Layout ui.TableLayout
	// Selection holds the chosen items of a Multiple table.
	Selection ui.Selection[K]
	Tree      treeCore[K]
	sorted    map[K]treeSorted[K]
	roots     treeSorted[K]
}

// treeSorted caches one sibling list sorted, until its source or the sort changes.
type treeSorted[K comparable] struct {
	src  *K
	n    int
	sort ui.SortOrder
	out  []K
}

// Selected returns the chosen item, or the last chosen of a Multiple table.
func (s *TreeTableState[K]) Selected() (K, bool) { return s.Tree.Sel, s.Tree.Has }

// SetOpen opens or closes an item.
func (s *TreeTableState[K]) SetOpen(item K, open bool) { s.Tree.setOpen(item, open) }

// IsOpen reports whether an item shows its children.
func (s *TreeTableState[K]) IsOpen(item K) bool { return s.Tree.Outline.Open.Has(item) }

// TreeTable shows a hierarchy in sortable, resizable columns, building only rows in view.
func TreeTable[K comparable](c *ui.Context, s *TreeTableState[K], opts TreeTableOptions[K]) DataTableView {
	treeTableCheck(opts)
	k := core.Tokens(c)
	d := core.Density(c)
	t := &s.Tree
	t.begin(opts.Label, &s.Selection, opts.Multiple)
	if s.sorted == nil {
		s.sorted = map[K]treeSorted[K]{}
	}
	cols := make([]ui.TableColumn, len(opts.Columns))
	menu := make([]tableColumn, len(opts.Columns))
	for i, col := range opts.Columns {
		cols[i] = ui.TableColumn{ID: col.ID, Title: col.Title, Width: col.Width, MinWidth: col.MinWidth, Align: col.Align, Sortable: col.Compare != nil}
		menu[i] = tableColumn{ID: col.ID, Title: col.Title, Width: col.Width, MinWidth: col.MinWidth, Sortable: col.Compare != nil}
	}
	// The sort is read after the column menu, which may change it.
	var cmp func(a, b K) int
	var order ui.SortOrder
	sorted := func(src []K, cache *treeSorted[K]) []K {
		if cmp == nil || len(src) < 2 {
			return src
		}
		if cache.src == &src[0] && cache.n == len(src) && cache.sort == order {
			return cache.out
		}
		out := slices.Clone(src)
		slices.SortStableFunc(out, func(a, b K) int {
			if order.Descending {
				return cmp(b, a)
			}
			return cmp(a, b)
		})
		*cache = treeSorted[K]{src: &src[0], n: len(src), sort: order, out: out}
		return out
	}
	children := func(item K) []K {
		kids := opts.Children(item)
		if len(kids) < 2 || cmp == nil {
			return kids
		}
		cache := s.sorted[item]
		out := sorted(kids, &cache)
		s.sorted[item] = cache
		return out
	}
	before := s.Sort
	var v DataTableView
	v.Element = ui.Column(c).MinHeight(0).Gap(d.Unit() * 2)
	v.Element.Children(func() {
		if len(opts.Roots) == 0 {
			t.clear(&s.Selection)
			v.retry = DataPlaceholder(c, opts.Status, opts.Error, opts.Empty, core.Msg(c, "list.empty"))
			return
		}
		ui.Row(c).Justify(ui.End).Children(func() { tableColumnMenu(c, menu, &s.Sort, &s.Layout, nil) })
		order = s.Sort
		for _, col := range opts.Columns {
			if col.ID == order.Column {
				cmp = col.Compare
			}
		}
		roots := sorted(opts.Roots, &s.roots)
		head := s.Sort
		t.Outline.List.Sort, t.Outline.List.Columns = &s.Sort, s.Layout
		restore := DataSelectTheme(c, dataRowUnit(c))
		e := ui.OutlineTable(c, &t.Outline, cols, roots, children, func(item K, col int) {
			if col > 0 {
				opts.Columns[col].Cell(item)
				return
			}
			ui.Row(c).Grow(1).MinWidth(0).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
				ui.Row(c).Grow(1).MinWidth(0).AlignItems(ui.Center).Children(func() { opts.Columns[0].Cell(item) })
				t.itemStatus(c, item, opts.ItemStatus, opts.Load)
				if t.chosen(item, &s.Selection, opts.Multiple) {
					ui.Icon(c, icons.Must("check")).Size(16, 16).Shrink(0).TextColor(k.AccentText)
				}
			})
		})
		restore()
		e.Grow(1).MinHeight(0).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius)
		internal.Ring(e, k.Focus, theme.ControlRadius)
		s.Layout = t.Outline.List.Columns
		v.changed = e.Changed() && s.Sort == head
		v.submitted = e.Submitted()
		v.retry = dataFooter(c, opts.Status, opts.Error)
	})
	v.sorted = s.Sort != before
	if len(opts.Roots) > 0 {
		t.end(c, v.changed)
	}
	if v.sorted {
		slog.Debug("mujica: tree_table sort", "column", s.Sort.Column, "descending", s.Sort.Descending)
	}
	return v
}

func treeTableCheck[K comparable](opts TreeTableOptions[K]) {
	if opts.Children == nil || opts.Label == nil || len(opts.Columns) == 0 {
		panic("mujica: TreeTable needs Columns, Children and Label")
	}
	seen := map[string]bool{}
	for _, col := range opts.Columns {
		if col.ID == "" || seen[col.ID] || col.Cell == nil {
			panic("mujica: TreeTable column needs a unique ID and Cell: " + col.ID)
		}
		seen[col.ID] = true
	}
}
