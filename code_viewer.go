package mujica

import (
	"fmt"
	"log/slog"
	"runtime"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// CodeViewerOptions configure a CodeViewer.
type CodeViewerOptions struct {
	// Language picks the highlighting; empty is plain.
	Language string
	// Label names the view for assistive technology.
	Label    string
	Disabled bool
}

// CodeViewerState is a CodeViewer's current line, selection, folds and place.
type CodeViewerState struct {
	// Selection holds the selected line indices.
	Selection ui.Selection[int]
	// Folded holds the first lines of folded blocks.
	Folded  ui.Selection[int]
	rows    ui.ListState
	cursor  int
	shown   []int
	inBlock []bool
	ends    []int
	data    *string
	n       int
	lang    string
	folds   int
	started bool
}

// Line returns the current line index, or -1.
func (s *CodeViewerState) Line() int {
	if s.cursor < 0 || s.cursor >= len(s.shown) {
		return -1
	}
	return s.shown[s.cursor]
}

// Visible returns the first and last lines in view.
func (s *CodeViewerState) Visible() (first, last int) {
	a, b := s.rows.Visible()
	if len(s.shown) == 0 || a < 0 {
		return 0, -1
	}
	return s.shown[min(a, len(s.shown)-1)], s.shown[min(max(b, a), len(s.shown)-1)]
}

// ScrollTo brings line to the top.
func (s *CodeViewerState) ScrollTo(line int) {
	row, _ := slices.BinarySearch(s.shown, line)
	s.rows.ScrollTo(row, ui.Start)
}

// GoTo makes line current, unfolding and scrolling to it.
func (s *CodeViewerState) GoTo(line int) {
	for start := range s.Folded.All() {
		if start < line && line <= s.ends[start] {
			s.Folded.Remove(start)
		}
	}
	s.folds = -1
	s.Selection.Clear()
	s.Selection.Add(line)
	s.started, s.cursor = true, -2-line // resolved to a row once the rows are known
}

// CodeViewer is a read-only, foldable, virtualized code view.
func CodeViewer(c *ui.Context, lines []string, s *CodeViewerState, opts CodeViewerOptions) ChoiceResult {
	if s == nil {
		panic("mujica: CodeViewer needs a state")
	}
	if !s.started {
		s.started, s.cursor = true, -1
	}
	codeViewerPrepare(lines, s, opts.Language)
	k := Tokens(c)
	d := Density(c)
	var res ChoiceResult
	res.Element = codePanel(c).Disabled(opts.Disabled)
	res.Element.Children(func() {
		if len(lines) == 0 {
			codeCaption(c, msg(c, "code_viewer.empty")).Padding(d.Unit() * 4)
			return
		}
		digits := len(fmt.Sprint(len(lines)))
		gutter := codeTextWidth(c, strings.Repeat("9", digits))
		s.rows.Selected = &s.cursor
		s.rows.Selection = &s.Selection
		s.rows.Key = func(i int) any { return s.shown[i] }
		s.rows.Label = func(i int) string { return fmt.Sprintf("%d %s", s.shown[i]+1, lines[s.shown[i]]) }
		restore := dataSelectTheme(c, 0)
		list := ui.List(c, &s.rows, len(s.shown), func(i int) {
			li := s.shown[i]
			ui.Row(c).MinHeight(d.ControlHeight()).Padding(0, d.Unit()*2, 0, 0).Gap(d.Unit() * 2).AlignItems(ui.Center).Children(func() {
				codeSelectMark(c, i == s.cursor)
				if s.ends[li] > li {
					label := msg(c, "code_viewer.fold")
					if s.Folded.Has(li) {
						label = msg(c, "code_viewer.unfold")
					}
					chev := "chevron-down"
					if s.Folded.Has(li) {
						chev = "chevron-right"
					}
					if iconAction(c, icons.Must(chev), fmt.Sprintf(label, li+1)).Clicked() && !res.Element.IsDisabled() {
						codeViewerToggle(s, li)
					}
				} else {
					ui.Box(c).Size(theme.MinHitSize, theme.MinHitSize).Shrink(0)
				}
				codeMono(c, fmt.Sprint(li+1)).Width(gutter).TextAlign(ui.End).TextColor(k.TextMuted).Shrink(0)
				codeViewerLine(c, k, opts.Language, lines[li], s.inBlock[li])
			})
		})
		restore()
		list.Grow(1).MinHeight(d.ControlHeight() * 3).Label(opts.Label)
		internal.Ring(list, k.Focus, theme.ControlRadius)
		if list.Changed() {
			res.changed = true
			slog.Debug("mujica: code_viewer line", "line", s.Line())
		}
		codeViewerKeys(c, list, lines, s)
	})
	return res
}

// codeViewerKeys copies the selection and folds with Left and Right.
func codeViewerKeys(c *ui.Context, list *ui.Element, lines []string, s *CodeViewerState) {
	if list.IsDisabled() {
		return
	}
	copyMods := ui.Cmd
	if runtime.GOOS != "darwin" {
		copyMods = ui.Ctrl
	}
	line := s.Line()
	switch {
	case list.Shortcut(copyMods, ui.KeyC):
		picked := slices.Sorted(s.Selection.All())
		if len(picked) == 0 && line >= 0 {
			picked = []int{line}
		}
		out := make([]string, len(picked))
		for i, l := range picked {
			out[i] = lines[l]
		}
		c.WriteClipboard(strings.Join(out, "\n"))
		slog.Debug("mujica: code_viewer copy", "lines", len(out))
	case list.Shortcut(0, ui.KeyLeft) && line >= 0 && s.ends[line] > line && !s.Folded.Has(line):
		codeViewerToggle(s, line)
	case list.Shortcut(0, ui.KeyRight) && line >= 0 && s.Folded.Has(line):
		codeViewerToggle(s, line)
	}
}

func codeViewerToggle(s *CodeViewerState, line int) {
	if s.Folded.Has(line) {
		s.Folded.Remove(line)
	} else {
		s.Folded.Add(line)
	}
	s.folds = -1
	slog.Debug("mujica: code_viewer fold", "line", line, "folded", s.Folded.Has(line))
}

// codeViewerLine draws one line's syntax tokens.
func codeViewerLine(c *ui.Context, k theme.Tokens, lang, line string, inBlock bool) {
	toks, _ := internal.Highlight(lang, line, inBlock)
	spans := make([]ui.Span, len(toks))
	for i, t := range toks {
		spans[i] = ui.Span{Text: t.Text, Color: codeTokenColor(k, t.Kind), Italic: t.Kind == internal.TokComment}
	}
	ui.RichText(c, spans...).Font(theme.MonoFont).FontSize(FontSize(c, codeSize)).SingleLine().Grow(1).MinWidth(0)
}

// codeViewerPrepare rescans code and folds when they change.
func codeViewerPrepare(lines []string, s *CodeViewerState, lang string) {
	var first *string
	if len(lines) > 0 {
		first = &lines[0]
	}
	if s.data != first || s.n != len(lines) || s.lang != lang {
		s.data, s.n, s.lang, s.folds = first, len(lines), lang, -1
		s.inBlock = slices.Grow(s.inBlock[:0], len(lines))[:len(lines)]
		in := false
		for i, l := range lines {
			s.inBlock[i] = in
			_, in = internal.Highlight(lang, l, in)
		}
		s.ends = codeViewerFolds(lines, s.ends)
	}
	if s.folds == s.Folded.Len() {
		return
	}
	s.folds = s.Folded.Len()
	s.shown = s.shown[:0]
	for i := 0; i < len(lines); i++ {
		s.shown = append(s.shown, i)
		if s.Folded.Has(i) {
			i = s.ends[i]
		}
	}
	if s.cursor <= -2 {
		s.cursor = max(slices.Index(s.shown, -2-s.cursor), 0)
		s.rows.ScrollIntoView(s.cursor)
	}
	s.cursor = min(s.cursor, len(s.shown)-1)
}

// codeViewerFolds returns each line's indented block end.
func codeViewerFolds(lines []string, ends []int) []int {
	ends = slices.Grow(ends[:0], len(lines))[:len(lines)]
	indent := func(l string) int {
		if strings.TrimSpace(l) == "" {
			return -1
		}
		return len(l) - len(strings.TrimLeft(l, " \t"))
	}
	type open struct{ line, indent int }
	var stack []open
	last := -1
	for i, l := range lines {
		ends[i] = i
		in := indent(l)
		if in < 0 {
			continue
		}
		for len(stack) > 0 && stack[len(stack)-1].indent >= in {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			ends[top.line] = last
		}
		stack = append(stack, open{i, in})
		last = i
	}
	for _, o := range stack {
		ends[o.line] = last
	}
	return ends
}
