package devtools

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/code"
	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/data"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// SchemaKind classifies nodes; zero means SchemaDatabase.
type SchemaKind int

const (
	SchemaDatabase SchemaKind = iota
	SchemaTable
	SchemaView
	SchemaColumn
	SchemaIndex
)

// SchemaNode describes an object; Detail holds type or index columns.
type SchemaNode struct {
	Name     string
	Kind     SchemaKind
	Detail   string
	Children []SchemaNode
}

// SchemaTreeState is a SchemaTree's open nodes and choice.
type SchemaTreeState struct {
	outline  ui.OutlineState[string]
	selected int
	nodes    map[string]*SchemaNode
	kids     map[string][]string
	roots    []string
	src      []SchemaNode
}

// Selected returns the chosen path and whether it exists.
func (s *SchemaTreeState) Selected() (string, bool) {
	if s.selected < 0 || s.selected >= s.outline.Rows() {
		return "", false
	}
	return s.outline.Item(s.selected), true
}

// SchemaTreeOptions configure a SchemaTree.
type SchemaTreeOptions struct {
	// Height of the tree; zero means 360.
	Height float32
}

// SchemaTree shows database objects in a virtualized keyboard tree.
func SchemaTree(c *ui.Context, s *SchemaTreeState, roots []SchemaNode, opts SchemaTreeOptions) core.ChoiceResult {
	if s == nil || opts.Height < 0 {
		panic("mujica: SchemaTree needs a state and Height >= 0")
	}
	h := opts.Height
	if h == 0 {
		h = 360
	}
	schemaIndex(s, roots)
	k := core.Tokens(c)
	d := core.Density(c)
	var r core.ChoiceResult
	r.Element = code.CodePanel(c).Height(h).Label(core.Msg(c, "schema_tree.label"))
	locked := r.Element.IsDisabled()
	if locked {
		copy := *s
		copy.outline.Open = ui.Selection[string]{}
		for path := range s.outline.Open.All() {
			copy.outline.Open.Add(path)
		}
		s = &copy
		for _, mods := range []ui.Modifiers{0, ui.Alt, ui.Shift} {
			r.Element.Shortcut(mods, ui.KeyLeft)
			r.Element.Shortcut(mods, ui.KeyRight)
		}
	}
	r.Element.Children(func() {
		if len(s.roots) == 0 {
			code.CodeCaption(c, core.Msg(c, "schema_tree.empty")).Padding(d.Unit() * 4)
			return
		}
		was := s.selected
		s.outline.List.Selected = &s.selected
		s.outline.List.Label = func(i int) string { return s.nodes[s.outline.Item(i)].Name }
		restore := data.DataSelectTheme(c, theme.MinHitSize/4)
		e := ui.Outline(c, &s.outline, s.roots, func(p string) []string { return s.kids[p] }, func(p string) {
			n := s.nodes[p]
			ui.Row(c).Grow(1).MinWidth(0).MinHeight(d.ControlHeight() - 2*d.Unit()).Gap(d.Unit() * 2).AlignItems(ui.Center).
				Label(fmt.Sprintf("%s %s %s", schemaWord(c, n.Kind), n.Name, n.Detail)).Children(func() {
				sel, _ := s.Selected()
				code.CodeSelectMark(c, sel == p)
				letter, col := schemaBadge(k, n.Kind)
				code.CodeBadge(c, letter, col)
				ui.Text(c, n.Name).SingleLine().Shrink(0)
				if n.Detail != "" {
					code.CodeMono(c, n.Detail).TextColor(k.TextMuted).MinWidth(0)
				}
			})
		}).Grow(1)
		restore()
		internal.Ring(e, k.Focus, theme.ControlRadius)
		r = core.NewChoiceResult(r.Element, !locked && s.selected != was)
	})
	if r.Changed() {
		p, _ := s.Selected()
		slog.Debug("mujica: schema_tree select", "path", p)
	}
	return r
}

// schemaIndex maps paths to nodes once per new schema.
func schemaIndex(s *SchemaTreeState, roots []SchemaNode) {
	if s.nodes != nil && len(roots) == len(s.src) && (len(roots) == 0 || &roots[0] == &s.src[0]) {
		return
	}
	if s.nodes == nil {
		s.selected = -1 // nothing chosen at first
	}
	s.src, s.nodes, s.kids, s.roots = roots, map[string]*SchemaNode{}, map[string][]string{}, nil
	var walk func(list []SchemaNode, parent string) []string
	walk = func(list []SchemaNode, parent string) []string {
		var paths []string
		for i := range list {
			n := &list[i]
			if n.Name == "" || n.Kind < SchemaDatabase || n.Kind > SchemaIndex {
				panic(fmt.Sprintf("mujica: SchemaTree node %q needs a Name and a known Kind", n.Name))
			}
			p := n.Name
			if parent != "" {
				p = parent + "/" + n.Name
			}
			if _, dup := s.nodes[p]; dup {
				panic("mujica: SchemaTree has two nodes at " + p)
			}
			s.nodes[p] = n
			paths = append(paths, p)
			if len(n.Children) > 0 {
				s.kids[p] = walk(n.Children, p)
			}
		}
		return paths
	}
	s.roots = walk(roots, "")
	if len(s.roots) == 0 {
		s.outline = ui.OutlineState[string]{}
		s.selected = -1
	}
}

func schemaBadge(k theme.Tokens, kind SchemaKind) (string, ui.Color) {
	switch kind {
	case SchemaDatabase:
		return "D", k.AccentText
	case SchemaTable:
		return "T", k.Info
	case SchemaView:
		return "V", k.Success
	case SchemaIndex:
		return "I", k.Warning
	}
	return "C", k.TextMuted
}

func schemaWord(c *ui.Context, kind SchemaKind) string {
	return core.Msg(c, [...]string{"schema_tree.database", "schema_tree.table", "schema_tree.view", "schema_tree.column", "schema_tree.index"}[kind])
}
