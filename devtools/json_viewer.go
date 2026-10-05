package devtools

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/code"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// JsonViewerState is a JsonViewer's open nodes, search and choice.
type JsonViewerState struct {
	// Query finds keys and values; matches open their ancestors.
	Query string
	// Current is the chosen match.
	Current  int
	open     map[string]bool
	tree     jsonTree
	src      string
	err      error
	rows     []int // visible node indexes
	stale    bool  // rows need rebuilding
	searched string
	matches  []int
	sel      int
	list     ui.ListState
}

// JsonViewerOptions configure a JsonViewer.
type JsonViewerOptions struct {
	// Height of the tree; zero means 360.
	Height float32
	// Expanded opens every node at first.
	Expanded bool
}

// JsonViewerResult reports this frame's events.
type JsonViewerResult struct {
	Element *ui.Element
	// Copied is the path written to the clipboard, else "".
	Copied string
}

// JsonViewer shows searchable JSON trees with folding and path copying.
func JsonViewer(c *ui.Context, data string, s *JsonViewerState, opts JsonViewerOptions) JsonViewerResult {
	if s == nil || opts.Height < 0 {
		panic("mujica: JsonViewer needs a state and Height >= 0")
	}
	h := opts.Height
	if h == 0 {
		h = 360
	}
	if s.open == nil || data != s.src {
		s.src, s.open, s.stale, s.searched = data, map[string]bool{"$": true}, true, "\x00"
		s.tree, s.err = parseJSONTree(data)
		if opts.Expanded {
			for _, n := range s.tree.nodes {
				s.open[n.path] = true
			}
		}
	}
	u := core.Density(c).Unit()
	var r JsonViewerResult
	r.Element = code.CodePanel(c).Height(h)
	if r.Element.IsDisabled() {
		for _, key := range []ui.Key{ui.KeyDown, ui.KeyUp, ui.KeyHome, ui.KeyEnd, ui.KeyPageDown, ui.KeyPageUp} {
			r.Element.Shortcut(0, key)
		}
	}
	r.Element.Children(func() {
		if s.err != nil {
			ui.Row(c).Padding(u * 3).Gap(u * 2).TextColor(core.Tokens(c).Danger).Role(ui.RoleStatus).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(16, 16)
				ui.Text(c, fmt.Sprintf(core.Msg(c, "json_viewer.invalid"), s.err))
			})
			return
		}
		jsonViewerSearch(s)
		code.CodeHeader(c, func() {
			moved := input.SearchInput(c, &s.Query, input.SearchInputOptions{Find: true, Matches: len(s.matches), Current: &s.Current}).Moved()
			if moved && len(s.matches) > 0 {
				s.sel = s.matches[s.Current]
				s.stale = true
			}
			path := ""
			if s.sel >= 0 && s.sel < len(s.tree.nodes) {
				path = s.tree.nodes[s.sel].path
			}
			code.CodeMono(c, path).Grow(1).MinWidth(0).SingleLine()
			if input.Button(c, core.Msg(c, "json_viewer.copy"), input.ButtonOptions{Variant: input.Ghost, Icon: icons.Must("copy"), Disabled: path == ""}).Clicked() && path != "" {
				c.WriteClipboard(path)
				r.Copied = path
				slog.Debug("mujica: json_viewer copy path", "path", path)
			}
		})
		jsonViewerRows(s)
		at := -1
		for i, n := range s.rows {
			if n == s.sel {
				at = i
			}
		}
		was := at
		s.list.Selected = &at
		s.list.Label = func(i int) string { return s.tree.nodes[s.rows[i]].label() }
		list := code.CodeRows(c, &s.list, len(s.rows), func(i int) { jsonViewerRow(c, s, s.rows[i]) })
		if !list.IsDisabled() && at != was && at >= 0 {
			s.sel = s.rows[at]
		}
		if !list.IsDisabled() && at >= 0 && at < len(s.rows) {
			n := s.tree.nodes[s.rows[at]]
			if n.container() && (list.Shortcut(0, ui.KeyRight) && !s.open[n.path] || list.Shortcut(0, ui.KeyLeft) && s.open[n.path] || list.Submitted()) {
				s.open[n.path] = !s.open[n.path]
				s.stale = true
			}
		}
	})
	return r
}

// jsonViewerSearch finds matches and opens their ancestors.
func jsonViewerSearch(s *JsonViewerState) {
	if s.searched == s.Query {
		return
	}
	s.searched = s.Query
	s.matches = s.matches[:0]
	q := strings.ToLower(s.Query)
	if q != "" {
		for i, n := range s.tree.nodes {
			if strings.Contains(strings.ToLower(n.key), q) || strings.Contains(strings.ToLower(n.value), q) {
				s.matches = append(s.matches, i)
				for p := n.parent; p >= 0; p = s.tree.nodes[p].parent {
					s.open[s.tree.nodes[p].path] = true
				}
			}
		}
	}
	s.Current = max(0, min(s.Current, len(s.matches)-1))
	if len(s.matches) > 0 {
		s.sel = s.matches[s.Current]
	}
	s.stale = true
}

// jsonViewerRows rebuilds visible indices when nodes change.
func jsonViewerRows(s *JsonViewerState) {
	if !s.stale {
		return
	}
	s.stale = false
	s.rows = s.rows[:0]
	for i := 0; i < len(s.tree.nodes); {
		n := s.tree.nodes[i]
		s.rows = append(s.rows, i)
		if n.container() && !s.open[n.path] {
			i = n.end
			continue
		}
		i++
	}
	if i := s.list.Selected; i != nil && *i >= 0 {
		for at, n := range s.rows {
			if n == s.sel {
				s.list.ScrollIntoView(at)
			}
		}
	}
}

// jsonViewerRow shows disclosure, key and value.
func jsonViewerRow(c *ui.Context, s *JsonViewerState, i int) {
	k := core.Tokens(c)
	n := s.tree.nodes[i]
	ui.Box(c).Width(float32(n.depth) * 16).Shrink(0)
	mark := ui.Box(c).Size(theme.MinHitSize, theme.MinHitSize).Center().Shrink(0)
	if n.container() {
		name := "chevron-right"
		if s.open[n.path] {
			name = "chevron-down"
		}
		mark.Role(ui.RoleDisclosure).Label(n.path).Cursor(ui.CursorPointer)
		mark.Children(func() { ui.Icon(c, icons.Must(name)).Size(14, 14).TextColor(k.TextMuted) })
		if !mark.IsDisabled() && mark.Clicked() {
			s.open[n.path], s.stale = !s.open[n.path], true
		}
	}
	var spans []ui.Span
	if n.key != "" {
		for _, sp := range code.CodeMatchSpans(c, n.key, s.Query) {
			sp.Color, sp.Weight = k.AccentText, max(sp.Weight, 600)
			spans = append(spans, sp)
		}
		spans = append(spans, ui.Span{Text: ": ", Color: k.TextMuted})
	}
	switch {
	case n.container() && s.open[n.path]:
		spans = append(spans, ui.Span{Text: n.value[:1], Color: k.TextMuted})
	case n.container():
		spans = append(spans, ui.Span{Text: fmt.Sprintf(core.Msg(c, "json_viewer.items"), n.value[:1], n.count, n.value[1:]), Color: k.TextMuted})
	default:
		tok, _ := internal.Highlight("json", n.value, false)
		col := k.Text
		if len(tok) > 0 {
			col = code.CodeTokenColor(k, tok[0].Kind)
		}
		for _, sp := range code.CodeMatchSpans(c, n.value, s.Query) {
			sp.Color = col
			spans = append(spans, sp)
		}
	}
	ui.RichText(c, spans...).Font(theme.MonoFont).FontSize(core.FontSize(c, code.CodeSize)).SingleLine().Grow(1).MinWidth(0)
}
