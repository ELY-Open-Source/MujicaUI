package data

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TreeOptions configure a Tree of items of type K, which are also their keys.
type TreeOptions[K comparable] struct {
	Roots []K
	// Children returns an item's children: nil for a leaf, empty while unloaded; required.
	Children func(item K) []K
	// Label returns an item's text, read out and matched by typing; required.
	Label func(item K) string
	// ItemStatus returns the load state of an item's children; nil means loaded.
	ItemStatus func(item K) DataStatus
	// Load asks for an item's children when it opens DataUnloaded, and on retry.
	Load     func(item K)
	Multiple bool
	// Empty replaces the built-in "No items" text.
	Empty  string
	Status DataStatus
	Error  string
}

// TreeState is a Tree's open items, selection and place.
type TreeState[K comparable] struct {
	// Selection holds the chosen items of a Multiple tree.
	Selection ui.Selection[K]
	Tree      treeCore[K]
}

// Selected returns the chosen item, or the last chosen of a Multiple tree.
func (s *TreeState[K]) Selected() (K, bool) { return s.Tree.Sel, s.Tree.Has }

// SetOpen opens or closes an item.
func (s *TreeState[K]) SetOpen(item K, open bool) { s.Tree.setOpen(item, open) }

// IsOpen reports whether an item shows its children.
func (s *TreeState[K]) IsOpen(item K) bool { return s.Tree.Outline.Open.Has(item) }

// treeCore is the outline state shared by Tree and TreeTable.
type treeCore[K comparable] struct {
	Outline ui.OutlineState[K]
	cursor  int
	started bool
	Sel     K
	Has     bool
	// asked holds the open items whose children Load was asked for.
	asked map[K]bool
}

func (t *treeCore[K]) setOpen(item K, open bool) {
	if open {
		t.Outline.Open.Add(item)
	} else {
		t.Outline.Open.Remove(item)
	}
}

// begin wires the outline's selection before it is built.
func (t *treeCore[K]) begin(label func(K) string, sel *ui.Selection[K], multiple bool) {
	if !t.started {
		t.started, t.cursor, t.asked = true, -1, map[K]bool{}
	}
	t.Outline.List.Selected = &t.cursor
	t.Outline.List.Label = func(i int) string { return label(t.Outline.Item(i)) }
	t.Outline.List.Selection = nil
	if multiple {
		t.Outline.List.Selection = sel
	}
}

// clear drops the choice and row state when the tree has no items.
func (t *treeCore[K]) clear(sel *ui.Selection[K]) {
	var zero K
	t.cursor, t.Sel, t.Has = -1, zero, false
	sel.Clear()
	clear(t.asked)
	t.Outline = ui.OutlineState[K]{Open: t.Outline.Open}
}

// end reads a new choice, or finds the chosen item again after rows moved.
func (t *treeCore[K]) end(c *ui.Context, changed bool) {
	rows := t.Outline.Rows()
	at := t.cursor >= 0 && t.cursor < rows
	if changed || !t.Has {
		t.Has = at
		if at {
			t.Sel = t.Outline.Item(t.cursor)
		}
		return
	}
	if at && t.Outline.Item(t.cursor) == t.Sel {
		return
	}
	t.cursor = -1
	for i := range rows {
		if t.Outline.Item(i) == t.Sel {
			t.cursor = i
			break
		}
	}
	if t.cursor < 0 {
		t.Has = false
		return
	}
	// The rows showed the old place: build them again at the new one.
	c.Invalidate()
	slog.Debug("mujica: tree selection moved", "row", t.cursor)
}

// chosen reports whether item shows as chosen.
func (t *treeCore[K]) chosen(item K, sel *ui.Selection[K], multiple bool) bool {
	if multiple {
		return sel.Has(item)
	}
	return t.Has && t.Sel == item
}

// itemStatus asks for unloaded children of an open item and shows loading or failure.
func (t *treeCore[K]) itemStatus(c *ui.Context, item K, status func(K) DataStatus, load func(K)) {
	if status == nil {
		return
	}
	st := status(item)
	open := t.Outline.Open.Has(item)
	if st != DataUnloaded || !open {
		delete(t.asked, item)
	}
	if st == DataUnloaded && open && !t.asked[item] {
		if load == nil {
			panic("mujica: tree item unloaded without Load")
		}
		t.asked[item] = true
		slog.Debug("mujica: tree load", "item", item)
		load(item)
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	switch st {
	case DataLoading:
		core.Spinner(c, core.SpinnerOptions{Size: 14, Label: core.Msg(c, "list.loading"), ShowLabel: true})
	case DataFailed:
		ui.Row(c).Gap(u).AlignItems(ui.Center).TextColor(k.Danger).Children(func() {
			ui.Icon(c, icons.Must("circle-alert")).Size(14, 14)
			ui.Text(c, core.Msg(c, "list.failed")).FontSize(core.FontSize(c, theme.CaptionSize)).SingleLine()
		})
		if input.Button(c, core.Msg(c, "list.retry"), input.ButtonOptions{Variant: input.Ghost}).MinHeight(theme.MinHitSize).Clicked() {
			slog.Debug("mujica: tree retry", "item", item)
			load(item)
		}
	}
}

// TreeView is a built Tree and what happened to it in this frame.
type TreeView struct {
	Element                   ui.Element
	changed, submitted, retry bool
}

// Changed reports a new selection.
func (v TreeView) Changed() bool { return v.changed }

// Submitted reports Enter or a double click on an item.
func (v TreeView) Submitted() bool { return v.submitted }

// Retried reports a click on the retry action of a failed first load.
func (v TreeView) Retried() bool { return v.retry }

// Tree shows items in a hierarchy, building only the rows in view.
func Tree[K comparable](c *ui.Context, s *TreeState[K], opts TreeOptions[K]) TreeView {
	if opts.Children == nil || opts.Label == nil {
		panic("mujica: Tree needs Children and Label")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	u := d.Unit()
	t := &s.Tree
	t.begin(opts.Label, &s.Selection, opts.Multiple)
	var v TreeView
	v.Element = ui.Column(c).MinHeight(0)
	v.Element.Children(func() {
		if len(opts.Roots) == 0 {
			t.clear(&s.Selection)
			v.retry = DataPlaceholder(c, opts.Status, opts.Error, opts.Empty, core.Msg(c, "list.empty"))
			return
		}
		// Outline rows pad by one unit above and below.
		h := d.ControlHeight() - 2*c.Theme().Space(1)
		restore := DataSelectTheme(c, 0)
		e := ui.Outline(c, &t.Outline, opts.Roots, opts.Children, func(item K) {
			ui.Row(c).Grow(1).MinWidth(0).MinHeight(h).Gap(u * 2).AlignItems(ui.Center).Children(func() {
				ui.Text(c, opts.Label(item)).SingleLine().Grow(1).MinWidth(0)
				t.itemStatus(c, item, opts.ItemStatus, opts.Load)
				if t.chosen(item, &s.Selection, opts.Multiple) {
					ui.Icon(c, icons.Must("check")).Size(16, 16).Shrink(0).TextColor(k.AccentText)
				}
			})
		})
		restore()
		e.Grow(1).MinHeight(0).Padding(u).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.ControlRadius)
		internal.Ring(e, k.Focus, theme.ControlRadius)
		v.changed, v.submitted = e.Changed(), e.Submitted()
		v.retry = dataFooter(c, opts.Status, opts.Error)
	})
	if len(opts.Roots) > 0 {
		t.end(c, v.changed)
	}
	if v.changed {
		slog.Debug("mujica: tree changed", "row", t.cursor, "count", s.Selection.Len())
	}
	return v
}
