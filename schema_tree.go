package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// SchemaKind is what a schema node is. The zero value is SchemaDatabase.
type SchemaKind int

const (
	SchemaDatabase SchemaKind = iota
	SchemaTable
	SchemaView
	SchemaColumn
	SchemaIndex
)

// SchemaNode is one database object; Detail is a column type or index columns.
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

// Selected returns the chosen node's path, as "shop/orders/id", and whether one is chosen.
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

// SchemaTree shows databases, tables, views, columns and indexes as a keyboard tree, virtualized.
func SchemaTree(c *ui.Context, s *SchemaTreeState, roots []SchemaNode, opts SchemaTreeOptions) ChoiceResult {
	if s == nil || opts.Height < 0 {
		panic("mujica: SchemaTree needs a state and Height >= 0")
	}
	h := opts.Height
	if h == 0 {
		h = 360
	}
	schemaIndex(s, roots)
	k := Tokens(c)
	d := Density(c)
	var r ChoiceResult
	r.Element = codePanel(c).Height(h).Label(msg(c, "schema_tree.label"))
	r.Element.Children(func() {
		if len(s.roots) == 0 {
			codeCaption(c, msg(c, "schema_tree.empty")).Padding(d.Unit() * 4)
			return
		}
		was := s.selected
		s.outline.List.Selected = &s.selected
		s.outline.List.Label = func(i int) string { return s.nodes[s.outline.Item(i)].Name }
		restore := dataSelectTheme(c, 0)
		e := ui.Outline(c, &s.outline, s.roots, func(p string) []string { return s.kids[p] }, func(p string) {
			n := s.nodes[p]
			ui.Row(c).Grow(1).MinWidth(0).MinHeight(d.ControlHeight() - 2*d.Unit()).Gap(d.Unit() * 2).AlignItems(ui.Center).
				Label(fmt.Sprintf("%s %s %s", schemaWord(c, n.Kind), n.Name, n.Detail)).Children(func() {
				sel, _ := s.Selected()
				codeSelectMark(c, sel == p)
				letter, col := schemaBadge(k, n.Kind)
				codeBadge(c, letter, col)
				ui.Text(c, n.Name).SingleLine().Shrink(0)
				if n.Detail != "" {
					codeMono(c, n.Detail).TextColor(k.TextMuted).MinWidth(0)
				}
			})
		}).Grow(1)
		restore()
		internal.Ring(e, k.Focus, theme.ControlRadius)
		r.changed = s.selected != was
	})
	if r.changed {
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
	return msg(c, [...]string{"schema_tree.database", "schema_tree.table", "schema_tree.view", "schema_tree.column", "schema_tree.index"}[kind])
}
