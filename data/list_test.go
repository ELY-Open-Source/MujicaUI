package data

import (
	"fmt"
	"image/color"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

type listHarness struct {
	ids      []int
	st       ListState[int]
	opts     ListOptions[int]
	built    int
	changes  int
	submits  int
	retries  int
	listSize [2]float32
}

func newListHarness(n int) *listHarness {
	h := &listHarness{listSize: [2]float32{300, 360}}
	for i := range n {
		h.ids = append(h.ids, i)
	}
	h.opts.Key = func(i int) int { return h.ids[i] }
	h.opts.Label = func(i int) string { return fmt.Sprintf("Item %d", h.ids[i]) }
	return h
}

func (h *listHarness) view(c *ui.Context) {
	core.Use(c, core.Settings{Mode: core.Light})
	h.built = 0
	ui.Column(c).Padding(10).Children(func() {
		v := List(c, &h.st, len(h.ids), h.opts, func(i int) {
			h.built++
			ui.Text(c, h.opts.Label(i)).SingleLine()
		})
		v.Element.Size(h.listSize[0], h.listSize[1])
		if v.Changed() {
			h.changes++
		}
		if v.Submitted() {
			h.submits++
		}
		if v.Retried() {
			h.retries++
		}
	})
}

func TestListSingleSelect(t *testing.T) {
	h := newListHarness(20)
	tt := ui.NewTester(h.view, 400, 400)
	if _, ok := h.st.Selected(); ok {
		t.Fatal("zero state has a selection")
	}
	if err := tt.Click("Item 3"); err != nil {
		t.Fatal(err)
	}
	if k, ok := h.st.Selected(); !ok || k != 3 || h.changes != 1 {
		t.Fatalf("selected %v %v, changes %d", k, ok, h.changes)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	if k, _ := h.st.Selected(); k != 5 {
		t.Fatalf("after Down×2 selected %d, want 5", k)
	}
	tt.Key(0, ui.KeyEnd)
	if k, _ := h.st.Selected(); k != 19 {
		t.Fatalf("End selected %d", k)
	}
	tt.Key(0, ui.KeyEnter)
	if h.submits != 1 {
		t.Fatalf("submits %d", h.submits)
	}
}

func TestListKeyFollowsInsert(t *testing.T) {
	h := newListHarness(10)
	tt := ui.NewTester(h.view, 400, 400)
	if err := tt.Click("Item 4"); err != nil {
		t.Fatal(err)
	}
	h.ids = append([]int{100, 101}, h.ids...)
	tt.Frame()
	if k, _ := h.st.Selected(); k != 4 {
		t.Fatalf("selection moved to %d after insert", k)
	}
	h.ids = append(h.ids[:5], h.ids[7:]...) // drops items 3 and 4
	tt.Frame()
	if k, ok := h.st.Selected(); ok {
		t.Fatalf("removed item still selected: %d", k)
	}
}

func TestListMultiple(t *testing.T) {
	h := newListHarness(10)
	h.opts.Multiple = true
	tt := ui.NewTester(h.view, 400, 400)
	if err := tt.Click("Item 1"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ClickWith(ui.Cmd, "Item 4"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ClickWith(ui.Shift, "Item 6"); err != nil {
		t.Fatal(err)
	}
	// Shift extends from the anchor click (MyGo keeps the pivot at item 1).
	for id := 1; id <= 6; id++ {
		if !h.st.Selection.Has(id) {
			t.Fatalf("item %d not chosen; %d chosen", id, h.st.Selection.Len())
		}
	}
	if h.st.Selection.Len() != 6 {
		t.Fatalf("chosen %d, want 6", h.st.Selection.Len())
	}
	tt.Key(ui.Cmd, ui.KeyA)
	if h.st.Selection.Len() != 10 {
		t.Fatalf("select all chose %d", h.st.Selection.Len())
	}
}

func TestListGridKeyboard(t *testing.T) {
	h := newListHarness(12)
	h.opts.Grid = true
	h.opts.ItemWidth, h.opts.ItemHeight = 80, 60
	h.listSize = [2]float32{380, 300}
	tt := ui.NewTester(h.view, 420, 340)
	if err := tt.Click("Item 0"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyRight)
	if k, _ := h.st.Selected(); k != 1 {
		t.Fatalf("Right chose %d", k)
	}
	tt.Key(0, ui.KeyDown)
	if k, _ := h.st.Selected(); k <= 1 {
		t.Fatalf("Down chose %d, want a later row", k)
	}
}

func TestListStatuses(t *testing.T) {
	h := newListHarness(0)
	tt := ui.NewTester(h.view, 400, 400)
	if !tt.HasText("No items") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	h.opts.Empty = "Nothing in the vault"
	tt.Frame()
	if !tt.HasText("Nothing in the vault") {
		t.Fatalf("custom empty: %q", tt.Texts())
	}
	h.opts.Status = DataLoading
	tt.Frame()
	if !tt.HasText("Loading") || tt.HasText("Nothing") {
		t.Fatalf("loading: %q", tt.Texts())
	}
	h.opts.Status, h.opts.Error = DataFailed, "Vault sealed"
	tt.Frame()
	if err := tt.Click("Retry"); err != nil || h.retries != 1 || !tt.HasText("Vault sealed") {
		t.Fatalf("failed: %v retries %d texts %q", err, h.retries, tt.Texts())
	}
	h.ids = []int{1, 2}
	h.opts.Status = DataLoadingMore
	tt.Frame()
	if !tt.HasText("Item 2") || !tt.HasText("Loading more") {
		t.Fatalf("loading more: %q", tt.Texts())
	}
}

func TestListLargeIsVirtualized(t *testing.T) {
	h := newListHarness(100_000)
	tt := ui.NewTester(h.view, 400, 400)
	// Rows are at least 36 DIPs: 360 DIPs show 10, plus overscan.
	const bound = 24
	if h.built == 0 || h.built > bound {
		t.Fatalf("first frame built %d items, bound %d", h.built, bound)
	}
	h.st.ScrollTo(73_000)
	tt.Frame()
	if h.built > bound || !tt.HasText("Item 73000") {
		t.Fatalf("after scroll built %d, texts %q", h.built, tt.Texts())
	}
	if err := tt.Click("Item 73001"); err != nil {
		t.Fatal(err)
	}
	if k, _ := h.st.Selected(); k != 73001 {
		t.Fatalf("selected %d", k)
	}
}

func TestListSnapshot(t *testing.T) {
	names := []string{"Matins", "Lauds", "Prime", "Terce", "Sext", "None"}
	snapshot(t, "list", 640, 260, func(c *ui.Context) {
		st := ui.Local(c.Root(), "st", func() *ListState[int] {
			s := &ListState[int]{}
			s.Selection.Add(1)
			s.Selection.Add(2)
			return s
		})
		g := ui.Local(c.Root(), "g", func() *ListState[int] { return &ListState[int]{started: true, cursor: 4} })
		opts := ListOptions[int]{Key: func(i int) int { return i }, Label: func(i int) string { return names[i] }, Multiple: true}
		ui.Row(c).Padding(16).Gap(16).Children(func() {
			List(c, *st, len(names), opts, func(i int) { ui.Text(c, names[i]) }).Element.Size(220, 228)
			gopts := opts
			gopts.Grid, gopts.Multiple, gopts.ItemWidth, gopts.ItemHeight = true, false, 100, 64
			List(c, *g, len(names), gopts, func(i int) {
				ui.Text(c, names[i]).FontWeight(600)
				ui.Text(c, "Office").TextColor(core.Tokens(c).TextMuted)
			}).Element.Size(360, 228)
		})
	})
}

func TestListKeyFollowsLargeReorder(t *testing.T) {
	h := newListHarness(5000)
	tt := ui.NewTester(h.view, 400, 400)
	if err := tt.Click("Item 3"); err != nil {
		t.Fatal(err)
	}
	slices.Reverse(h.ids)
	tt.Frame()
	if k, ok := h.st.Selected(); !ok || k != 3 || h.st.cursor != 4996 {
		t.Fatalf("after reversing: selected %d %v at row %d", k, ok, h.st.cursor)
	}
}

func TestListGridKeyFollowsInsert(t *testing.T) {
	h := newListHarness(12)
	h.opts.Grid = true
	h.opts.ItemWidth, h.opts.ItemHeight = 80, 60
	h.listSize = [2]float32{380, 300}
	tt := ui.NewTester(h.view, 420, 340)
	if err := tt.Click("Item 2"); err != nil {
		t.Fatal(err)
	}
	h.ids = append([]int{100}, h.ids...)
	tt.Frame()
	if k, _ := h.st.Selected(); k != 2 || h.st.cursor != 3 {
		t.Fatalf("after insert: selected %d at %d", k, h.st.cursor)
	}
}

func TestListSelectionColors(t *testing.T) {
	h := newListHarness(5)
	tt := ui.NewTester(h.view, 400, 400)
	if err := tt.Click("Item 1"); err != nil {
		t.Fatal(err)
	}
	r, _ := tt.Find("Item 1")
	got := color.RGBAModel.Convert(tt.Image().At(int(r.X+r.W-6), int(r.Y+r.H/2))).(color.RGBA)
	if want := toRGBA(theme.Light().Selection); got != want {
		t.Fatalf("selected row background %v, want Selection %v", got, want)
	}
}
