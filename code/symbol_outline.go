package code

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// SymbolKind is what a Symbol declares.
type SymbolKind int

const (
	SymbolFunction SymbolKind = iota
	SymbolMethod
	SymbolType
	SymbolInterface
	SymbolField
	SymbolVariable
	SymbolConstant
	SymbolModule
)

// Symbol is a declaration spanning Line to EndLine, from 1.
type Symbol struct {
	Name, Detail  string
	Kind          SymbolKind
	Line, EndLine int
	Children      []Symbol
}

// SymbolOutlineState is a SymbolOutline's filter, open symbols and choice.
type SymbolOutlineState struct {
	Query   string
	outline ui.OutlineState[int]
	cursor  int
	nodes   []symbolNode
	roots   []int
	data    *Symbol
	n       int
	line    int
	started bool
}

type symbolNode struct {
	sym    *Symbol
	parent int
	kids   []int
}

// SymbolOutlineOptions configure a SymbolOutline.
type SymbolOutlineOptions struct {
	// CursorLine chooses its symbol when it changes; 0 is off.
	CursorLine int
	Disabled   bool
}

// SymbolOutlineResult carries the symbol to jump to.
type SymbolOutlineResult struct {
	Element *ui.Element
	Jumped  *Symbol
}

var symbolBadges = [...]string{"f", "m", "T", "I", "p", "v", "c", "M"}

// SymbolOutline shows a filterable symbol tree.
func SymbolOutline(c *ui.Context, symbols []Symbol, s *SymbolOutlineState, opts SymbolOutlineOptions) SymbolOutlineResult {
	if s == nil {
		panic("mujica: SymbolOutline needs a state")
	}
	if !s.started {
		s.started, s.cursor = true, -1
	}
	symbolIndex(symbols, s)
	k := core.Tokens(c)
	d := core.Density(c)
	var res SymbolOutlineResult
	res.Element = CodePanel(c).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	res.Element.Children(func() {
		CodeHeader(c, func() {
			input.SearchInput(c, &s.Query, input.SearchInputOptions{Placeholder: core.Msg(c, "symbol_outline.filter"), Disabled: opts.Disabled}).Element.Grow(1)
		})
		q := strings.ToLower(strings.TrimSpace(s.Query))
		keep := symbolFilter(s, q)
		roots := symbolKept(s.roots, keep)
		if len(roots) == 0 {
			CodeCaption(c, core.Msg(c, "symbol_outline.empty")).Padding(d.Unit() * 4)
			return
		}
		follow := opts.CursorLine != s.line
		s.line = opts.CursorLine
		if follow {
			symbolFollow(s, opts.CursorLine, keep)
		}
		s.outline.List.Selected = &s.cursor
		s.outline.List.Label = func(i int) string { return s.nodes[s.outline.Item(i)].sym.Name }
		children := func(i int) []int {
			if len(s.nodes[i].kids) == 0 {
				return nil
			}
			if q != "" {
				s.outline.Open.Add(i)
			}
			return symbolKept(s.nodes[i].kids, keep)
		}
		restore := data.DataSelectTheme(c, 0)
		e := ui.Outline(c, &s.outline, roots, children, func(i int) {
			sym := s.nodes[i].sym
			if sym.Kind < SymbolFunction || sym.Kind > SymbolModule {
				panic("mujica: Symbol kind out of range")
			}
			ui.Row(c).MinHeight(d.ControlHeight() - 2*c.Theme().Space(1)).Gap(d.Unit() * 2).AlignItems(ui.Center).Grow(1).MinWidth(0).Children(func() {
				CodeBadge(c, symbolBadges[sym.Kind], [...]ui.Color{k.AccentText, k.AccentText, k.Info, k.Info, k.TextMuted, k.TextMuted, k.Warning, k.Success}[sym.Kind])
				ui.RichText(c, CodeMatchSpans(c, sym.Name, q)...).SingleLine().MinWidth(0)
				if sym.Detail != "" {
					CodeCaption(c, sym.Detail).MinWidth(0)
				}
				ui.Spacer(c)
				CodeSelectMark(c, s.cursor >= 0 && s.cursor < s.outline.Rows() && s.outline.Item(s.cursor) == i)
			})
		})
		restore()
		e.Grow(1).MinHeight(d.ControlHeight() * 3)
		internal.Ring(e, k.Focus, theme.ControlRadius)
		if follow && s.cursor >= 0 {
			s.outline.List.ScrollIntoView(s.cursor)
		}
		if !disabled && s.cursor >= 0 && s.cursor < s.outline.Rows() && e.Submitted() {
			sym := *s.nodes[s.outline.Item(s.cursor)].sym
			res.Jumped = &sym
			slog.Debug("mujica: symbol_outline jump", "line", sym.Line)
		}
	})
	return res
}

// symbolIndex numbers the symbols depth first when the tree changes.
func symbolIndex(symbols []Symbol, s *SymbolOutlineState) {
	var first *Symbol
	if len(symbols) > 0 {
		first = &symbols[0]
	}
	if s.data == first && s.n == len(symbols) {
		return
	}
	s.data, s.n, s.nodes, s.roots = first, len(symbols), s.nodes[:0], s.roots[:0]
	var walk func(list []Symbol, parent int) []int
	walk = func(list []Symbol, parent int) []int {
		var ids []int
		for i := range list {
			id := len(s.nodes)
			s.nodes = append(s.nodes, symbolNode{sym: &list[i], parent: parent})
			s.nodes[id].kids = walk(list[i].Children, id)
			ids = append(ids, id)
		}
		return ids
	}
	s.roots = walk(symbols, -1)
}

// symbolFilter marks matches and ancestors; nil keeps all.
func symbolFilter(s *SymbolOutlineState, q string) []bool {
	if q == "" {
		return nil
	}
	keep := make([]bool, len(s.nodes))
	for i, n := range s.nodes {
		if strings.Contains(strings.ToLower(n.sym.Name), q) {
			for j := i; j >= 0 && !keep[j]; j = s.nodes[j].parent {
				keep[j] = true
			}
		}
	}
	return keep
}

func symbolKept(ids []int, keep []bool) []int {
	if keep == nil {
		return ids
	}
	var out []int
	for _, id := range ids {
		if keep[id] {
			out = append(out, id)
		}
	}
	return out
}

// symbolFollow opens and chooses the innermost symbol spanning line.
func symbolFollow(s *SymbolOutlineState, line int, keep []bool) {
	best := -1
	for i, n := range s.nodes {
		if keep != nil && !keep[i] {
			continue
		}
		if n.sym.Line <= line && line <= max(n.sym.EndLine, n.sym.Line) {
			best = i
		}
	}
	if best < 0 {
		return
	}
	for p := s.nodes[best].parent; p >= 0; p = s.nodes[p].parent {
		s.outline.Open.Add(p)
	}
	s.cursor = symbolRow(s, best, keep)
	slog.Debug("mujica: symbol_outline follow", "line", line)
}

// symbolRow returns the row the shown node id is at.
func symbolRow(s *SymbolOutlineState, id int, keep []bool) int {
	row := 0
	var walk func(ids []int) bool
	walk = func(ids []int) bool {
		for _, n := range symbolKept(ids, keep) {
			if n == id {
				return true
			}
			row++
			if s.outline.Open.Has(n) && walk(s.nodes[n].kids) {
				return true
			}
		}
		return false
	}
	if walk(s.roots) {
		return row
	}
	return -1
}
