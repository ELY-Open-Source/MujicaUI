package mujica

import (
	"cmp"
	"fmt"
	"image/color"
	"strconv"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

type relic struct {
	ID   int
	Name string
	Age  int
}

type tableHarness struct {
	rows     []relic
	st       DataTableState[int]
	opts     DataTableOptions[relic, int]
	compares int
	changes  int
	sorts    int
	c        *ui.Context
	density  theme.Density
	seen     map[string]int
}

func newTableHarness(rows []relic) *tableHarness {
	h := &tableHarness{rows: rows}
	h.opts = DataTableOptions[relic, int]{
		Key: func(r relic) int { return r.ID },
		Columns: []DataColumn[relic]{
			{ID: "name", Title: "Name", Sortable: true, Text: func(r relic) string { return r.Name },
				Compare: func(a, b relic) int { h.compares++; return cmp.Compare(a.Name, b.Name) }},
			{ID: "age", Title: "Age", Width: 80, Align: ui.End, Sortable: true, Text: func(r relic) string { return strconv.Itoa(r.Age) },
				Compare: func(a, b relic) int { h.compares++; return cmp.Compare(a.Age, b.Age) }},
		},
	}
	return h
}

func (h *tableHarness) view(c *ui.Context) {
	Use(c, Settings{Mode: Light, Density: h.density})
	h.c, h.seen = c, map[string]int{}
	ui.Column(c).Padding(10).Gap(8).Children(func() {
		cols := h.opts
		for i := range cols.Columns {
			text := cols.Columns[i].Text
			cols.Columns[i].Text = func(r relic) string { h.seen[r.Name]++; return text(r) }
		}
		v := DataTable(c, &h.st, h.rows, cols)
		v.Element.Size(360, 300)
		if v.Changed() {
			h.changes++
		}
		if v.SortChanged() {
			h.sorts++
		}
	})
}

func relics() []relic {
	return []relic{{1, "Censer", 300}, {2, "Altar", 120}, {3, "Bell", 80}, {4, "Diptych", 450}}
}

func (h *tableHarness) order(tt *ui.Tester) []string {
	var names []string
	for _, s := range tt.Texts() {
		for _, r := range h.rows {
			// Rows are labeled with their name too: skip the repeat.
			if r.Name == s && (len(names) == 0 || names[len(names)-1] != s) {
				names = append(names, s)
			}
		}
	}
	return names
}

func TestDataTableSortsOnce(t *testing.T) {
	h := newTableHarness(relics())
	tt := ui.NewTester(h.view, 400, 360)
	if got := fmt.Sprint(h.order(tt)); got != "[Censer Altar Bell Diptych]" {
		t.Fatalf("initial order %s", got)
	}
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprint(h.order(tt)); got != "[Altar Bell Censer Diptych]" || h.st.Sort != (ui.SortOrder{Column: "name"}) || h.sorts != 1 {
		t.Fatalf("ascending %s sort %+v sorts %d", got, h.st.Sort, h.sorts)
	}
	n := h.compares
	tt.Frame()
	tt.Frame()
	if h.compares != n {
		t.Fatalf("sorted again without a change: %d compares", h.compares-n)
	}
	if err := tt.Click("Age"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Age"); err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprint(h.order(tt)); got != "[Diptych Censer Altar Bell]" || !h.st.Sort.Descending {
		t.Fatalf("age descending %s", got)
	}
}

func TestDataTableSelectionFollowsRow(t *testing.T) {
	h := newTableHarness(relics())
	tt := ui.NewTester(h.view, 400, 360)
	if err := tt.Click("Bell"); err != nil {
		t.Fatal(err)
	}
	if k, ok := h.st.Selected(); !ok || k != 3 || h.changes != 1 {
		t.Fatalf("selected %d %v changes %d", k, ok, h.changes)
	}
	_ = tt.Click("Name")
	_ = tt.Click("Name")
	if k, _ := h.st.Selected(); k != 3 || h.rows[h.st.cache.view[h.st.cursor]].Name != "Bell" {
		t.Fatalf("after sort selected %d at row %d", k, h.st.cursor)
	}
	tt.Key(0, ui.KeyDown)
	if k, _ := h.st.Selected(); k != 2 { // descending: Diptych Censer Bell Altar
		t.Fatalf("Down chose %d", k)
	}
	h.st.Filter = "ce"
	tt.Frame()
	if h.st.Rows() != 1 || h.st.cursor != -1 {
		t.Fatalf("filter rows %d cursor %d", h.st.Rows(), h.st.cursor)
	}
	h.st.Filter = ""
	tt.Frame()
	if k, ok := h.st.Selected(); !ok || k != 2 || h.rows[h.st.cache.view[h.st.cursor]].Name != "Altar" {
		t.Fatalf("selection lost through filter: %d %v", k, ok)
	}
	h.rows = []relic{{1, "Censer", 300}}
	tt.Frame()
	if _, ok := h.st.Selected(); ok {
		t.Fatal("removed row stays selected")
	}
}

func TestDataTableMultipleMarks(t *testing.T) {
	h := newTableHarness(relics())
	h.opts.Multiple = true
	tt := ui.NewTester(h.view, 400, 360)
	table, _ := tt.Find("Name")
	for _, name := range []string{"Altar", "Diptych"} {
		r, _ := tt.Find(name)
		tt.ClickAt(table.X-28, r.Y+r.H/2) // the mark column sits before Name
	}
	if h.st.Selection.Len() != 2 || !h.st.Selection.Has(2) || !h.st.Selection.Has(4) || h.changes != 2 {
		t.Fatalf("marks chose %d, changes %d", h.st.Selection.Len(), h.changes)
	}
	if _, ok := tt.Find("Select row"); !ok {
		t.Fatal("row check box has no name")
	}
	_ = tt.Click("Age")
	if h.st.Selection.Len() != 2 || !h.st.Selection.Has(4) {
		t.Fatal("selection changed by sorting")
	}
	h.rows = h.rows[:2] // Censer, Altar
	tt.Frame()
	if h.st.Selection.Len() != 1 || !h.st.Selection.Has(2) {
		t.Fatalf("prune kept %d", h.st.Selection.Len())
	}
}

// tabTo presses Tab until the element labeled label has the focus.
func tabTo(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	for range 60 {
		if tt.Focused(label) {
			return
		}
		tt.Key(0, ui.KeyTab)
	}
	t.Fatalf("Tab never reached %q", label)
}

func TestDataTableColumnMenuKeyboard(t *testing.T) {
	h := newTableHarness(relics())
	tt := ui.NewTester(h.view, 420, 480)
	tabTo(t, tt, "Columns")
	tt.Key(0, ui.KeyEnter)
	tabTo(t, tt, "Sort Age descending")
	tt.Key(0, ui.KeyEnter)
	if h.st.Sort != (ui.SortOrder{Column: "age", Descending: true}) || h.sorts != 1 {
		t.Fatalf("sort %+v sorts %d", h.st.Sort, h.sorts)
	}
	tabTo(t, tt, "Move Age left")
	tt.Key(0, ui.KeyEnter)
	if fmt.Sprint(h.st.Layout.Order) != "[age name]" {
		t.Fatalf("order %v", h.st.Layout.Order)
	}
	tabTo(t, tt, "Make Age wider")
	tt.Key(0, ui.KeyEnter)
	if h.st.Layout.Widths["age"] != 96 {
		t.Fatalf("widths %v", h.st.Layout.Widths)
	}
	tabTo(t, tt, "Name")
	tt.Key(0, ui.KeySpace)
	if !h.st.Hidden["name"] {
		t.Fatalf("name not hidden: %v", h.st.Hidden)
	}
	tabTo(t, tt, "Age")
	tt.Key(0, ui.KeySpace)
	if h.st.Hidden["age"] {
		t.Fatal("hid the last visible column")
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Censer") || !tt.HasText("Age") {
		t.Fatalf("table did not follow the menu: %q", tt.Texts())
	}
}

func TestDataTableCustomCellAndControlled(t *testing.T) {
	h := newTableHarness(relics())
	h.opts.Controlled = true
	h.opts.Columns[0].Cell = func(r relic) { ui.Text(h.c, "◆ "+r.Name) }
	tt := ui.NewTester(h.view, 400, 360)
	if !tt.HasText("◆ Altar") || tt.HasText("Altar") && !tt.HasText("◆ Altar") {
		t.Fatalf("custom cell missing: %q", tt.Texts())
	}
	_ = tt.Click("Name")
	first := ""
	for _, s := range tt.Texts() {
		if len(s) > 4 && s[:4] == "◆ " {
			first = s
			break
		}
	}
	if h.sorts != 1 || h.compares != 0 || first != "◆ Censer" {
		t.Fatalf("controlled sorted locally: sorts %d compares %d first %q", h.sorts, h.compares, first)
	}
}

func TestDataTableBadOptions(t *testing.T) {
	for name, mutate := range map[string]func(o *DataTableOptions[relic, int]){
		"no key":       func(o *DataTableOptions[relic, int]) { o.Key = nil },
		"no compare":   func(o *DataTableOptions[relic, int]) { o.Columns[0].Compare = nil },
		"duplicate id": func(o *DataTableOptions[relic, int]) { o.Columns[1].ID = "name" },
	} {
		h := newTableHarness(relics())
		mutate(&h.opts)
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: no panic", name)
				}
			}()
			ui.NewTester(h.view, 300, 300)
		}()
	}
	h := newTableHarness([]relic{{1, "A", 1}, {1, "B", 2}})
	defer func() {
		if recover() == nil {
			t.Error("duplicate row keys: no panic")
		}
	}()
	ui.NewTester(h.view, 300, 300)
}

func TestDataTableStatuses(t *testing.T) {
	h := newTableHarness(nil)
	tt := ui.NewTester(h.view, 400, 360)
	if !tt.HasText("No items") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	h.rows = relics()
	h.st.Filter = "zzz"
	tt.Frame()
	if !tt.HasText("No matching rows") {
		t.Fatalf("no match: %q", tt.Texts())
	}
	h.st.Filter = ""
	h.opts.Status = DataFailed
	tt.Frame()
	if !tt.HasText("Could not load") || !tt.HasText("Retry") || !tt.HasText("Bell") {
		t.Fatalf("failed append: %q", tt.Texts())
	}
}

func TestDataTableLargeIsVirtualized(t *testing.T) {
	rows := make([]relic, 100_000)
	for i := range rows {
		rows[i] = relic{ID: i, Name: fmt.Sprintf("Relic %06d", (i*7919)%100_000), Age: i % 997}
	}
	h := newTableHarness(rows)
	for i := 3; i <= 10; i++ {
		id := fmt.Sprintf("c%d", i)
		h.opts.Columns = append(h.opts.Columns, DataColumn[relic]{ID: id, Title: id, Width: 60, Text: func(r relic) string { return id }})
	}
	tt := ui.NewTester(h.view, 400, 360)
	// Rows whose cells ran in the last frame: 300 DIPs show about 7 rows of 36 DIPs, plus overscan.
	const bound = 16
	if len(h.seen) == 0 || len(h.seen) > bound {
		t.Fatalf("built %d rows, bound %d", len(h.seen), bound)
	}
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	if len(h.seen) > bound || !tt.HasText("Relic 000000") || !tt.HasText("Relic 000001") {
		t.Fatalf("after sort built %d rows, texts %q", len(h.seen), tt.Texts())
	}
	if err := tt.Click("Relic 000001"); err != nil {
		t.Fatal(err)
	}
	n := h.compares
	_ = tt.Click("Name") // descending: the row moves 99,998 places
	tt.Frame()
	if k, _ := h.st.Selected(); k != rowWithName(rows, "Relic 000001") {
		t.Fatalf("selection lost: %d", k)
	}
	m := h.compares
	if m == n {
		t.Fatal("descending did not sort")
	}
	tt.Frame()
	tt.Frame()
	if h.compares != m {
		t.Fatal("re-sorted each frame")
	}
}

func rowWithName(rows []relic, name string) int {
	for _, r := range rows {
		if r.Name == name {
			return r.ID
		}
	}
	panic(name)
}

func TestDataTableSnapshot(t *testing.T) {
	snapshot(t, "data_table", 640, 300, func(c *ui.Context) {
		h := ui.Local(c.Root(), "h", func() *tableHarness {
			h := newTableHarness(relics())
			h.opts.Multiple = true
			h.st.Sort = ui.SortOrder{Column: "name"}
			h.st.Selection.Add(3)
			return h
		})
		ui.Column(c).Padding(16).Children(func() {
			DataTable(c, &(*h).st, (*h).rows, (*h).opts).Element.Size(600, 260)
		})
	})
}

func TestDataTableSelectionColors(t *testing.T) {
	h := newTableHarness(relics())
	tt := ui.NewTester(h.view, 400, 360)
	if err := tt.Click("Bell"); err != nil {
		t.Fatal(err)
	}
	r, _ := tt.Find("Bell")
	got := color.RGBAModel.Convert(tt.Image().At(int(r.X+r.W-6), int(r.Y+r.H/2))).(color.RGBA)
	if want := toRGBA(theme.Light().Selection); got != want {
		t.Fatalf("selected row background %v, want Selection %v", got, want)
	}
}

func TestDataTableCompactLargeText(t *testing.T) {
	h := newTableHarness(relics())
	h.opts.Multiple = true
	h.density = theme.Compact
	tt := ui.NewTester(h.view, 600, 700)
	tt.SetPreferences(ui.Preferences{TextScale: 2})
	head, _ := tt.Find("Name")
	row, _ := tt.Find("Censer") // the first row, labeled by its first column
	if head.Y+head.H > row.Y-1 {
		t.Fatalf("header text %v runs past the header into row %v", head, row)
	}
	box, ok := tt.Find("Select row")
	if !ok || box.W < theme.MinHitSize || box.H < theme.MinHitSize || box.X < row.X {
		t.Fatalf("check box hit area %v in row %v", box, row)
	}
	if name, _ := tt.Find("Name"); box.X+box.W > name.X {
		t.Fatalf("check box %v overlaps the Name column at %v", box, name.X)
	}
}

func TestDataTableCheckBoxHitWidth(t *testing.T) {
	for _, d := range []theme.Density{theme.Compact, theme.Default} {
		h := newTableHarness(relics())
		h.opts.Multiple = true
		h.density = d
		tt := ui.NewTester(h.view, 400, 360)
		box, _ := tt.Find("Select row")
		if box.W < theme.MinHitSize || box.H < theme.MinHitSize {
			t.Fatalf("%v: check box hit area %v", d, box)
		}
		tt.ClickAt(box.X+1, box.Y+box.H/2) // its left edge still toggles
		if h.st.Selection.Len() != 1 {
			t.Fatalf("%v: left edge of the hit area missed", d)
		}
	}
}
