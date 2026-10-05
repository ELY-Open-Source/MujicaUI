package mujica

import (
	"cmp"
	"fmt"
	"strconv"
	"testing"

	"github.com/egoist/mygo/ui"
)

type ttHarness struct {
	kids     map[string][]string
	size     map[string]int
	status   map[string]DataStatus
	loads    int
	st       TreeTableState[string]
	opts     TreeTableOptions[string]
	compares int
	sorts    int
	c        *ui.Context
}

func newTTHarness() *ttHarness {
	h := &ttHarness{
		kids:   map[string][]string{"Vestry": {"Stole", "Alb", "Cope"}, "Belfry": {}},
		size:   map[string]int{"Vestry": 9, "Belfry": 4, "Stole": 2, "Alb": 5, "Cope": 1, "Organ": 7},
		status: map[string]DataStatus{"Belfry": DataUnloaded},
	}
	text := func(n string) { ui.Text(h.c, n).SingleLine() }
	h.opts = TreeTableOptions[string]{
		Roots:      []string{"Vestry", "Belfry", "Organ"},
		Children:   func(n string) []string { return h.kids[n] },
		Label:      func(n string) string { return n },
		ItemStatus: func(n string) DataStatus { return h.status[n] },
		Load:       func(n string) { h.loads++; h.status[n] = DataLoading },
		Columns: []TreeColumn[string]{
			{ID: "name", Title: "Name", Cell: text,
				Compare: func(a, b string) int { h.compares++; return cmp.Compare(a, b) }},
			{ID: "size", Title: "Size", Width: 80, Align: ui.End,
				Cell:    func(n string) { ui.Text(h.c, strconv.Itoa(h.size[n])) },
				Compare: func(a, b string) int { h.compares++; return cmp.Compare(h.size[a], h.size[b]) }},
		},
	}
	return h
}

func (h *ttHarness) view(c *ui.Context) {
	Use(c, Settings{Mode: Light})
	h.c = c
	ui.Column(c).Padding(10).Children(func() {
		v := TreeTable(c, &h.st, h.opts)
		v.Size(360, 320)
		if v.SortChanged() {
			h.sorts++
		}
	})
}

func (h *ttHarness) order(tt *ui.Tester) string {
	var out []string
	for _, s := range tt.Texts() {
		if _, ok := h.size[s]; ok && (len(out) == 0 || out[len(out)-1] != s) {
			out = append(out, s)
		}
	}
	return fmt.Sprint(out)
}

func TestTreeTableSortsSiblings(t *testing.T) {
	h := newTTHarness()
	tt := ui.NewTester(h.view, 400, 360)
	h.st.SetOpen("Vestry", true)
	tt.Frame()
	if got := h.order(tt); got != "[Vestry Stole Alb Cope Belfry Organ]" {
		t.Fatalf("unsorted %s", got)
	}
	if err := tt.Click("Alb"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Size"); err != nil {
		t.Fatal(err)
	}
	if got := h.order(tt); got != "[Belfry Organ Vestry Cope Stole Alb]" || h.sorts != 1 {
		t.Fatalf("by size %s sorts %d", got, h.sorts)
	}
	if k, _ := h.st.Selected(); k != "Alb" {
		t.Fatalf("selection moved to %q", k)
	}
	n := h.compares
	tt.Frame()
	tt.Frame()
	if h.compares != n {
		t.Fatalf("re-sorted without a change: %d", h.compares-n)
	}
	_ = tt.Click("Size")
	if got := h.order(tt); got != "[Vestry Alb Stole Cope Organ Belfry]" {
		t.Fatalf("descending %s", got)
	}
}

func TestTreeTableLazyAndKeys(t *testing.T) {
	h := newTTHarness()
	tt := ui.NewTester(h.view, 400, 360)
	_ = tt.Click("Belfry")
	tt.Key(0, ui.KeyRight)
	if !h.st.IsOpen("Belfry") || h.loads != 1 || !tt.HasText("Loading") {
		t.Fatalf("open %v loads %d texts %q", h.st.IsOpen("Belfry"), h.loads, tt.Texts())
	}
	h.kids["Belfry"] = []string{"Rope"}
	h.size["Rope"] = 3
	h.status["Belfry"] = DataReady
	tt.Key(0, ui.KeyRight)
	if k, _ := h.st.Selected(); k != "Rope" {
		t.Fatalf("Right chose %q", k)
	}
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyLeft)
	if h.st.IsOpen("Belfry") || tt.HasText("Rope") {
		t.Fatal("Left did not close")
	}
}

func TestTreeTableLayout(t *testing.T) {
	h := newTTHarness()
	h.st.Layout = ui.TableLayout{Order: []string{"size", "name"}, Widths: map[string]float32{"size": 120}}
	tt := ui.NewTester(h.view, 400, 360)
	size, _ := tt.Find("Size")
	name, _ := tt.Find("Name")
	if size.X >= name.X || h.st.Layout.Widths["size"] != 120 {
		t.Fatalf("layout lost: size %v name %v %v", size.X, name.X, h.st.Layout)
	}
}

func TestTreeTableBadOptions(t *testing.T) {
	h := newTTHarness()
	h.opts.Columns[1].Cell = nil
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without Cell")
		}
	}()
	ui.NewTester(h.view, 300, 300)
}

func TestTreeTableLargeIsVirtualized(t *testing.T) {
	roots, kids := bigTree()
	var st TreeTableState[int]
	rows := map[int]bool{}
	cell := func(c *ui.Context) func(int) {
		return func(n int) { rows[n] = true; ui.Text(c, strconv.Itoa(n)) }
	}
	view := func(c *ui.Context) {
		Use(c, Settings{})
		clear(rows)
		TreeTable(c, &st, TreeTableOptions[int]{Roots: roots, Children: kids, Label: strconv.Itoa,
			Columns: []TreeColumn[int]{{ID: "n", Title: "N", Cell: cell(c), Compare: cmp.Compare[int]}, {ID: "m", Title: "M", Cell: cell(c)}}}).Size(320, 360)
	}
	tt := ui.NewTester(view, 400, 400)
	const bound = 26 // see TestTreeLargeIsVirtualized
	if len(rows) == 0 || len(rows) > bound || rows[1] {
		t.Fatalf("collapsed: %d rows", len(rows))
	}
	for _, r := range roots {
		st.SetOpen(r, true)
	}
	tt.Frame()
	if len(rows) > bound || st.tree.outline.Rows() != 100_000 {
		t.Fatalf("open: %d rows of %d", len(rows), st.tree.outline.Rows())
	}
	st.Sort = ui.SortOrder{Column: "n", Descending: true}
	tt.Frame()
	if len(rows) > bound || st.tree.outline.Item(0) != 99_000 || st.tree.outline.Item(1) != 99_999 {
		t.Fatalf("sorted: %d rows, first %d %d", len(rows), st.tree.outline.Item(0), st.tree.outline.Item(1))
	}
}

func TestTreeTableSnapshot(t *testing.T) {
	snapshot(t, "tree_table", 420, 240, func(c *ui.Context) {
		h := ui.Local(c.Root(), "h", func() *ttHarness {
			h := newTTHarness()
			h.st.SetOpen("Vestry", true)
			h.st.Sort = ui.SortOrder{Column: "name"}
			h.st.tree.started, h.st.tree.cursor, h.st.tree.asked = true, 2, map[string]bool{}
			h.st.tree.sel, h.st.tree.has = "Vestry", true
			return h
		})
		(*h).c = c
		ui.Column(c).Padding(16).Children(func() {
			TreeTable(c, &(*h).st, (*h).opts).Size(388, 208)
		})
	})
}

func TestTreeTableSelectionFollowsLargeSort(t *testing.T) {
	kids := make([]int, 3000)
	for i := range kids {
		kids[i] = i + 1
	}
	var st TreeTableState[int]
	view := func(c *ui.Context) {
		Use(c, Settings{})
		TreeTable(c, &st, TreeTableOptions[int]{Roots: []int{0}, Label: strconv.Itoa,
			Children: func(n int) []int {
				if n == 0 {
					return kids
				}
				return nil
			},
			Columns: []TreeColumn[int]{{ID: "n", Title: "N", Cell: func(n int) { ui.Text(c, "#"+strconv.Itoa(n)) }, Compare: cmp.Compare[int]}},
		}).Size(320, 360)
	}
	st.SetOpen(0, true)
	tt := ui.NewTester(view, 400, 400)
	if err := tt.Click("#3"); err != nil {
		t.Fatal(err)
	}
	st.Sort = ui.SortOrder{Column: "n", Descending: true}
	tt.Frame()
	if k, ok := st.Selected(); !ok || k != 3 || st.tree.outline.Item(st.tree.cursor) != 3 {
		t.Fatalf("after sort: selected %d %v at row %d", k, ok, st.tree.cursor)
	}
}

func TestTreeTableColumnMenuKeyboard(t *testing.T) {
	h := newTTHarness()
	tt := ui.NewTester(h.view, 420, 480)
	tabTo(t, tt, "Columns")
	tt.Key(0, ui.KeyEnter)
	tabTo(t, tt, "Sort Size ascending")
	tt.Key(0, ui.KeyEnter)
	if h.st.Sort != (ui.SortOrder{Column: "size"}) || h.sorts != 1 {
		t.Fatalf("sort %+v sorts %d", h.st.Sort, h.sorts)
	}
	tabTo(t, tt, "Move Name right")
	tt.Key(0, ui.KeyEnter)
	tabTo(t, tt, "Make Name narrower")
	tt.Key(0, ui.KeyEnter)
	if fmt.Sprint(h.st.Layout.Order) != "[size name]" || h.st.Layout.Widths["name"] != 104 {
		t.Fatalf("layout %+v", h.st.Layout)
	}
	tt.Key(0, ui.KeyEscape)
	if got := h.order(tt); got != "[Belfry Organ Vestry]" {
		t.Fatalf("rows not sorted by size: %s", got)
	}
}
