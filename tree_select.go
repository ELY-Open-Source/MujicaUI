package mujica

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// TreeSelectNode is one node of a TreeSelect. IDs must be unique in the tree.
type TreeSelectNode struct {
	ID, Label string
	Children  []TreeSelectNode
	Disabled  bool
}

// TreeSelectRule says how checking a parent affects its children in multiple mode.
type TreeSelectRule int

const (
	// TreeSelectCascade checks a parent's enabled leaves with it; the value holds leaf IDs.
	TreeSelectCascade TreeSelectRule = iota
	// TreeSelectIndependent checks each node on its own; the value holds any IDs.
	TreeSelectIndependent
)

// TreeSelectOptions configure a TreeSelect.
type TreeSelectOptions struct {
	Multiple bool
	Rule     TreeSelectRule
	// Placeholder shows while nothing is chosen; empty uses "Select…".
	Placeholder string
	// MaxTags is how many chosen nodes show as tags before "+N"; zero is 3.
	MaxTags  int
	Disabled bool
	// Label names the field for assistive technology.
	Label string
}

type treeSelectState struct {
	open     bool
	query    string
	expanded map[string]bool
	// focus is the ID of the row to focus as it is built; "\x00" is the first row.
	focus string
}

type treeSelectKey struct{}

// treeSelectRow is a visible row: its node, depth and parent.
type treeSelectRow struct {
	node   *TreeSelectNode
	depth  int
	parent string
}

// TreeSelect chooses one node, or several with Multiple, from a searchable tree in a popup.
func TreeSelect(c *ui.Context, selected *[]string, nodes []TreeSelectNode, opts TreeSelectOptions) ChoiceResult {
	if !opts.Multiple && len(*selected) > 1 {
		panic("mujica: single TreeSelect holds more than one ID")
	}
	k := Tokens(c)
	before := slices.Clone(*selected)
	tr := ui.ButtonBase(c).Disabled(opts.Disabled).Role(ui.RolePopUpButton)
	if opts.Label != "" {
		tr.Label(opts.Label)
	}
	st := *ui.Local(tr, treeSelectKey{}, func() *treeSelectState { return &treeSelectState{expanded: map[string]bool{}} })
	if opts.Disabled {
		st.open = false
	}
	if tr.Clicked() || (!st.open && tr.Shortcut(0, ui.KeyDown) && !tr.IsDisabled()) {
		st.open, st.query = !st.open, ""
	}
	labels := map[string]string{}
	treeSelectWalk(nodes, func(n *TreeSelectNode) { labels[n.ID] = n.Label })
	selectField(c, tr, st.open, func() {
		switch {
		case len(*selected) == 0:
			ui.Text(c, selectOr(opts.Placeholder, msg(c, "select.placeholder"))).SingleLine().TextColor(k.TextMuted).Grow(1)
		case !opts.Multiple:
			ui.Text(c, labels[(*selected)[0]]).SingleLine().Grow(1).Shrink(1)
		default:
			options := make([]SelectOption[string], 0, len(labels))
			for _, id := range *selected {
				options = append(options, SelectOption[string]{Value: id, Label: labels[id]})
			}
			multiSelectTags(c, selected, options, MultiSelectOptions{MaxTags: opts.MaxTags, Disabled: opts.Disabled})
		}
	})
	ui.PopoverBase(c, tr, &st.open, func(panel *ui.Element) {
		selectPanel(c, panel)
		panel.MinWidth(max(tr.Bounds().W, 260))
		ui.Column(c).Gap(4).Children(func() {
			if selectSearch(c, &st.query, msg(c, "tree_select.search")).Shortcut(0, ui.KeyDown) {
				st.focus = "\x00"
			}
			rows := treeSelectVisible(nodes, st)
			if len(rows) == 0 {
				selectNote(c, msg(c, "tree_select.no_results"))
				return
			}
			shown := make(map[string]bool, len(rows))
			for _, r := range rows {
				shown[r.node.ID] = true
			}
			ui.Scroll(c).MaxHeight(300).Padding(4).Children(func() {
				ui.Column(c).FocusGroup(ui.Vertical).Children(func() {
					entry := slices.IndexFunc(rows, func(r treeSelectRow) bool { return !r.node.Disabled })
					for i, r := range rows {
						treeSelectRowView(c, st, selected, r, i == entry, shown, opts)
					}
				})
			})
		})
	})
	changed := !slices.Equal(before, *selected)
	if changed {
		slog.Debug("mujica: tree select changed", "label", opts.Label, "selected", *selected)
	}
	return ChoiceResult{tr, changed}
}

func treeSelectWalk(nodes []TreeSelectNode, fn func(n *TreeSelectNode)) {
	for i := range nodes {
		fn(&nodes[i])
		treeSelectWalk(nodes[i].Children, fn)
	}
}

// treeSelectVisible lists the rows to show: expanded ones, or matches and their ancestors.
func treeSelectVisible(nodes []TreeSelectNode, st *treeSelectState) []treeSelectRow {
	q := strings.ToLower(strings.TrimSpace(st.query))
	var rows []treeSelectRow
	var visit func(ns []TreeSelectNode, depth int, parent string) bool
	visit = func(ns []TreeSelectNode, depth int, parent string) bool {
		found := false
		for i := range ns {
			n := &ns[i]
			at := len(rows)
			rows = append(rows, treeSelectRow{n, depth, parent})
			self := q == "" || strings.Contains(strings.ToLower(n.Label), q)
			below := false
			if q != "" || st.expanded[n.ID] {
				below = visit(n.Children, depth+1, n.ID)
			}
			if q != "" && !self && !below {
				rows = rows[:at]
				continue
			}
			found = true
		}
		return found
	}
	visit(nodes, 0, "")
	return rows
}

// treeSelectCheck returns a node's check state: its own in single mode, else by the rule.
func treeSelectCheck(n *TreeSelectNode, selected []string, opts TreeSelectOptions) internal.CheckState {
	if !opts.Multiple || opts.Rule == TreeSelectIndependent || len(n.Children) == 0 {
		if slices.Contains(selected, n.ID) {
			return internal.Checked
		}
		return internal.Unchecked
	}
	on, all := 0, 0
	treeSelectWalk(n.Children, func(x *TreeSelectNode) {
		if len(x.Children) == 0 && !x.Disabled {
			all++
			if slices.Contains(selected, x.ID) {
				on++
			}
		}
	})
	switch {
	case on == 0:
		return internal.Unchecked
	case on == all:
		return internal.Checked
	}
	return internal.Mixed
}

// treeSelectToggle checks or unchecks n by the rule.
func treeSelectToggle(n *TreeSelectNode, selected []string, on bool, opts TreeSelectOptions) []string {
	ids := []string{n.ID}
	if opts.Rule == TreeSelectCascade && len(n.Children) > 0 {
		ids = nil
		treeSelectWalk(n.Children, func(x *TreeSelectNode) {
			if len(x.Children) == 0 && !x.Disabled {
				ids = append(ids, x.ID)
			}
		})
	}
	for _, id := range ids {
		selected = multiSelectSet(selected, id, false)
		if on {
			selected = append(selected, id)
		}
	}
	return selected
}

// treeSelectFirstChild returns the first child of n that shows and is enabled, or "".
func treeSelectFirstChild(n *TreeSelectNode, shown map[string]bool) string {
	for _, ch := range n.Children {
		if shown[ch.ID] && !ch.Disabled {
			return ch.ID
		}
	}
	return ""
}

// treeSelectDescribe says what MyGo cannot set on custom rows (ui/element.go:282).
func treeSelectDescribe(c *ui.Context, b *ui.Element, branch, open bool, state internal.CheckState) {
	var parts []string
	if branch {
		parts = append(parts, msg(c, map[bool]string{true: "tree_select.expanded", false: "tree_select.collapsed"}[open]))
	}
	if state == internal.Mixed {
		parts = append(parts, msg(c, "tree_select.mixed"))
	}
	if len(parts) > 0 {
		b.Description(strings.Join(parts, ", "))
	}
}

func treeSelectRowView(c *ui.Context, st *treeSelectState, selected *[]string, r treeSelectRow, first bool, shown map[string]bool, opts TreeSelectOptions) {
	k := Tokens(c)
	d := Density(c)
	n := r.node
	branch := len(n.Children) > 0
	open := branch && (st.expanded[n.ID] || strings.TrimSpace(st.query) != "")
	ui.Row(c).Key(n.ID).AlignItems(ui.Center).Gap(2).Padding(0, 0, 0, float32(r.depth)*18).Children(func() {
		arrow := ui.Box(c).Size(theme.MinHitSize, theme.MinHitSize).Shrink(0).Radius(theme.ControlRadius)
		if branch {
			arrow.Role(ui.RoleButton).Label(fmt.Sprintf(msg(c, map[bool]string{true: "tree_select.collapse", false: "tree_select.expand"}[open]), n.Label))
			if arrow.Clicked() {
				st.expanded[n.ID] = !st.expanded[n.ID]
			}
			if arrow.Hovered() {
				arrow.Background(k.SurfaceHover)
			}
			arrow.Draw(func(p *ui.Painter, rc ui.Rect) {
				in := (rc.W - 16) / 2
				rc = ui.Rect{X: rc.X + in, Y: rc.Y + in, W: 16, H: 16}
				internal.Mark(p, rc, map[bool]string{true: "chevron-down", false: "chevron-right"}[open], k.TextMuted)
			})
		}
		state := treeSelectCheck(n, *selected, opts)
		var b *ui.Element
		if opts.Multiple {
			on := state == internal.Checked
			b = multiSelectCheck(c, k, &on, state, n.Label, n.Disabled)
			if b.Changed() {
				*selected = treeSelectToggle(n, *selected, on, opts)
			}
		} else {
			chosen := state == internal.Checked
			ui.Row(c).Grow(1).Disabled(n.Disabled).Children(func() {
				b = ui.ButtonBase(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2).Gap(d.Unit() * 2).
					Justify(ui.Start).Radius(theme.ControlRadius).Grow(1)
			})
			internal.Ring(b, k.Focus, theme.ControlRadius)
			if b.Clicked() {
				*selected, st.open = []string{n.ID}, false
			}
			b.Role(ui.RoleTreeItem).Selected(chosen).Background(ui.Transparent).TextColor(k.Text)
			switch {
			case chosen:
				b.Background(k.Selection).TextColor(k.AccentText)
			case b.Hovered():
				b.Background(k.SurfaceHover)
			}
			b.Children(func() {
				ui.Text(c, n.Label).SingleLine().Grow(1)
				selectCheck(c, chosen)
			})
		}
		treeSelectDescribe(c, b, branch, open, state)
		if st.focus == n.ID || first && st.focus == "\x00" {
			b.Focus()
			st.focus = ""
		}
		right, left := b.Shortcut(0, ui.KeyRight) && !b.IsDisabled(), b.Shortcut(0, ui.KeyLeft) && !b.IsDisabled()
		switch {
		case right && branch && !open:
			st.expanded[n.ID] = true
		case right && branch:
			st.focus = treeSelectFirstChild(n, shown)
		case left && open && st.expanded[n.ID]:
			st.expanded[n.ID] = false
		case left && r.parent != "":
			st.focus = r.parent
		}
	})
}
