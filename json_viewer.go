package mujica

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
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

// JsonViewer shows JSON as a tree with folding, path copy and search, virtualized.
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
	u := Density(c).Unit()
	var r JsonViewerResult
	r.Element = codePanel(c).Height(h)
	r.Element.Children(func() {
		if s.err != nil {
			ui.Row(c).Padding(u * 3).Gap(u * 2).TextColor(Tokens(c).Danger).Role(ui.RoleStatus).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(16, 16)
				ui.Text(c, fmt.Sprintf(msg(c, "json_viewer.invalid"), s.err))
			})
			return
		}
		jsonViewerSearch(s)
		codeHeader(c, func() {
			moved := SearchInput(c, &s.Query, SearchInputOptions{Find: true, Matches: len(s.matches), Current: &s.Current}).Moved()
			if moved && len(s.matches) > 0 {
				s.sel = s.matches[s.Current]
				s.stale = true
			}
			path := ""
			if s.sel >= 0 && s.sel < len(s.tree.nodes) {
				path = s.tree.nodes[s.sel].path
			}
			codeMono(c, path).Grow(1).MinWidth(0).SingleLine()
			if Button(c, msg(c, "json_viewer.copy"), ButtonOptions{Variant: Ghost, Icon: icons.Must("copy"), Disabled: path == ""}).Clicked() && path != "" {
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
		list := codeRows(c, &s.list, len(s.rows), func(i int) { jsonViewerRow(c, s, s.rows[i]) })
		if at != was && at >= 0 {
			s.sel = s.rows[at]
		}
		if at >= 0 && at < len(s.rows) {
			n := s.tree.nodes[s.rows[at]]
			if n.container() && (list.Shortcut(0, ui.KeyRight) && !s.open[n.path] || list.Shortcut(0, ui.KeyLeft) && s.open[n.path] || list.Submitted()) {
				s.open[n.path] = !s.open[n.path]
				s.stale = true
			}
		}
	})
	return r
}

// jsonViewerSearch finds matches and opens their ancestors when the query changes.
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

// jsonViewerRows lists the visible nodes, skipping closed subtrees, when anything changed.
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

// jsonViewerRow shows one node: fold mark, key, and value or size.
func jsonViewerRow(c *ui.Context, s *JsonViewerState, i int) {
	k := Tokens(c)
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
		if mark.Clicked() {
			s.open[n.path], s.stale = !s.open[n.path], true
		}
	}
	var spans []ui.Span
	if n.key != "" {
		for _, sp := range codeMatchSpans(c, n.key, s.Query) {
			sp.Color, sp.Weight = k.AccentText, max(sp.Weight, 600)
			spans = append(spans, sp)
		}
		spans = append(spans, ui.Span{Text: ": ", Color: k.TextMuted})
	}
	switch {
	case n.container() && s.open[n.path]:
		spans = append(spans, ui.Span{Text: n.value[:1], Color: k.TextMuted})
	case n.container():
		spans = append(spans, ui.Span{Text: fmt.Sprintf(msg(c, "json_viewer.items"), n.value[:1], n.count, n.value[1:]), Color: k.TextMuted})
	default:
		tok, _ := internal.Highlight("json", n.value, false)
		col := k.Text
		if len(tok) > 0 {
			col = codeTokenColor(k, tok[0].Kind)
		}
		for _, sp := range codeMatchSpans(c, n.value, s.Query) {
			sp.Color = col
			spans = append(spans, sp)
		}
	}
	ui.RichText(c, spans...).Font(theme.MonoFont).FontSize(FontSize(c, codeSize)).SingleLine().Grow(1).MinWidth(0)
}
