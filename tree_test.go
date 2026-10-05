package mujica

import (
	"fmt"
	"image/color"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

type treeHarness struct {
	kids    map[string][]string
	status  map[string]DataStatus
	loads   []string
	st      TreeState[string]
	opts    TreeOptions[string]
	built   map[string]bool
	changes int
}

func newTreeHarness() *treeHarness {
	h := &treeHarness{
		kids: map[string][]string{
			"Nave":     {"Pulpit", "Pews"},
			"Crypt":    {}, // loaded lazily
			"Pulpit":   nil,
			"Pews":     nil,
			"Sacristy": nil,
		},
		status: map[string]DataStatus{"Crypt": DataUnloaded},
	}
	h.opts = TreeOptions[string]{
		Roots:      []string{"Nave", "Crypt", "Sacristy"},
		Children:   func(n string) []string { return h.kids[n] },
		Label:      func(n string) string { return n },
		ItemStatus: func(n string) DataStatus { return h.status[n] },
		Load: func(n string) {
			h.loads = append(h.loads, n)
			h.status[n] = DataLoading
		},
	}
	return h
}

func (h *treeHarness) view(c *ui.Context) {
	Use(c, Settings{Mode: Light})
	h.built = map[string]bool{}
	opts := h.opts
	label := opts.Label
	opts.Label = func(n string) string { h.built[n] = true; return label(n) }
	ui.Column(c).Padding(10).Children(func() {
		v := Tree(c, &h.st, opts)
		v.Size(320, 360)
		if v.Changed() {
			h.changes++
		}
	})
}

func TestTreeExpandAndKeyboard(t *testing.T) {
	h := newTreeHarness()
	tt := ui.NewTester(h.view, 400, 400)
	if tt.HasText("Pulpit") || h.built["Pulpit"] {
		t.Fatal("closed subtree built")
	}
	if err := tt.Click("Nave"); err != nil {
		t.Fatal(err)
	}
	if k, _ := h.st.Selected(); k != "Nave" || h.changes != 1 {
		t.Fatalf("selected %q changes %d", k, h.changes)
	}
	tt.Key(0, ui.KeyRight)
	if !h.st.IsOpen("Nave") || !tt.HasText("Pulpit") {
		t.Fatalf("Right did not open: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyRight)
	if k, _ := h.st.Selected(); k != "Pulpit" {
		t.Fatalf("Right to first child chose %q", k)
	}
	tt.Key(0, ui.KeyDown)
	if k, _ := h.st.Selected(); k != "Pews" {
		t.Fatalf("Down chose %q", k)
	}
	tt.Key(0, ui.KeyLeft)
	if k, _ := h.st.Selected(); k != "Nave" {
		t.Fatalf("Left to parent chose %q", k)
	}
	tt.Key(0, ui.KeyLeft)
	if h.st.IsOpen("Nave") || tt.HasText("Pews") {
		t.Fatal("Left did not close")
	}
}

func TestTreeLazyLoad(t *testing.T) {
	h := newTreeHarness()
	tt := ui.NewTester(h.view, 400, 400)
	if len(h.loads) != 0 {
		t.Fatalf("loaded before opening: %v", h.loads)
	}
	h.st.SetOpen("Crypt", true)
	tt.Frame()
	if fmt.Sprint(h.loads) != "[Crypt]" || !tt.HasText("Loading") {
		t.Fatalf("loads %v texts %q", h.loads, tt.Texts())
	}
	tt.Frame()
	if len(h.loads) != 1 {
		t.Fatalf("asked twice: %v", h.loads)
	}
	h.status["Crypt"] = DataFailed
	tt.Frame()
	if err := tt.Click("Retry"); err != nil {
		t.Fatal(err)
	}
	if len(h.loads) != 2 {
		t.Fatalf("retry did not load: %v", h.loads)
	}
	h.kids["Crypt"] = []string{"Ossuary"}
	h.status["Crypt"] = DataReady
	tt.Frame()
	if !tt.HasText("Ossuary") || tt.HasText("Loading") {
		t.Fatalf("children missing: %q", tt.Texts())
	}
}

func TestTreeMultipleAndEmpty(t *testing.T) {
	h := newTreeHarness()
	h.opts.Multiple = true
	tt := ui.NewTester(h.view, 400, 400)
	_ = tt.Click("Nave")
	_ = tt.ClickWith(ui.Cmd, "Sacristy")
	if h.st.Selection.Len() != 2 || !h.st.Selection.Has("Sacristy") {
		t.Fatalf("chose %d", h.st.Selection.Len())
	}
	h.opts.Roots = nil
	tt.Frame()
	if !tt.HasText("No items") {
		t.Fatalf("empty: %q", tt.Texts())
	}
}

func TestTreeBadOptions(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without Children")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		var s TreeState[int]
		Tree(c, &s, TreeOptions[int]{Roots: []int{1}})
	}, 100, 100)
}

// bigTree has 100 roots of 999 children each: 100,000 nodes.
func bigTree() (roots []int, kids func(int) []int) {
	for r := range 100 {
		roots = append(roots, r*1000)
	}
	cache := map[int][]int{}
	return roots, func(n int) []int {
		if n%1000 != 0 {
			return nil
		}
		if k, ok := cache[n]; ok {
			return k
		}
		var k []int
		for i := 1; i < 1000; i++ {
			k = append(k, n+i)
		}
		cache[n] = k
		return k
	}
}

func TestTreeLargeIsVirtualized(t *testing.T) {
	roots, kids := bigTree()
	var st TreeState[int]
	built := 0
	view := func(c *ui.Context) {
		Use(c, Settings{})
		built = 0
		Tree(c, &st, TreeOptions[int]{Roots: roots, Children: kids, Label: func(n int) string {
			built++
			return fmt.Sprintf("Node %d", n)
		}}).Size(320, 360)
	}
	tt := ui.NewTester(view, 400, 400)
	// Two Label calls per row; a jump builds 2×(10 rows + 2 overscan) + 2 edges.
	const bound = 2 * 26
	if built == 0 || built > bound {
		t.Fatalf("collapsed: label calls %d, bound %d", built, bound)
	}
	for _, r := range roots {
		st.SetOpen(r, true)
	}
	tt.Frame()
	if built > bound || st.tree.outline.Rows() != 100_000 {
		t.Fatalf("open: label calls %d rows %d", built, st.tree.outline.Rows())
	}
	st.tree.outline.List.ScrollTo(64_000, ui.Start)
	tt.Frame()
	if built > bound || !tt.HasText(fmt.Sprintf("Node %d", st.tree.outline.Item(64_000))) {
		t.Fatalf("scrolled: label calls %d texts %q", built, tt.Texts())
	}
}

func TestTreeSnapshot(t *testing.T) {
	snapshot(t, "tree", 400, 260, func(c *ui.Context) {
		h := ui.Local(c.Root(), "h", func() *treeHarness {
			h := newTreeHarness()
			h.st.SetOpen("Nave", true)
			h.st.SetOpen("Crypt", true)
			h.status["Crypt"] = DataLoading
			h.st.tree.started, h.st.tree.cursor, h.st.tree.asked = true, 1, map[string]bool{}
			h.st.tree.sel, h.st.tree.has = "Pulpit", true
			return h
		})
		ui.Column(c).Padding(16).Children(func() {
			Tree(c, &(*h).st, (*h).opts).Size(360, 228)
		})
	})
}

func TestTreeSelectionFollowsLargeReorder(t *testing.T) {
	roots := make([]int, 3000)
	for i := range roots {
		roots[i] = i
	}
	var st TreeState[int]
	view := func(c *ui.Context) {
		Use(c, Settings{Mode: Light})
		Tree(c, &st, TreeOptions[int]{Roots: roots, Children: func(int) []int { return nil },
			Label: func(n int) string { return fmt.Sprintf("Node %d", n) }}).Size(320, 360)
	}
	tt := ui.NewTester(view, 400, 400)
	if err := tt.Click("Node 2"); err != nil {
		t.Fatal(err)
	}
	r, _ := tt.Find("Node 2")
	got := color.RGBAModel.Convert(tt.Image().At(int(r.X+r.W-6), int(r.Y+r.H/2))).(color.RGBA)
	if want := toRGBA(theme.Light().Selection); got != want {
		t.Fatalf("selected row background %v, want Selection %v", got, want)
	}
	roots = slices.Clone(roots)
	slices.Reverse(roots)
	tt.Frame()
	if k, ok := st.Selected(); !ok || k != 2 || st.tree.cursor != 2997 {
		t.Fatalf("after reversing: selected %d %v at row %d", k, ok, st.tree.cursor)
	}
}
