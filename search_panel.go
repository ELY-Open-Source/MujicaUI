package mujica

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// SearchMatch is a hit; Start and End index Text.
type SearchMatch struct {
	Line       int
	Text       string
	Start, End int
}

// SearchFileResult is a file and its hits.
type SearchFileResult struct {
	Path    string
	Matches []SearchMatch
}

// SearchPanelState is a SearchPanel's query, patterns, folded files and place.
type SearchPanelState struct {
	Query, Include, Exclude string
	// Collapsed holds the paths of folded files.
	Collapsed ui.Selection[string]
	rows      ui.ListState
	selected  int
	flat      []searchRow
	data      *SearchFileResult
	n, folded int
	started   bool
}

type searchRow struct{ file, match int }

// SearchPanelOptions configure a SearchPanel.
type SearchPanelOptions struct {
	// Status shows a search in progress or failed.
	Status   DataStatus
	Error    string
	Disabled bool
}

// SearchOpen is a hit the user opened.
type SearchOpen struct {
	Path string
	Line int
}

// SearchPanelResult is a SearchPanel and what happened in this frame.
type SearchPanelResult struct {
	Element *ui.Element
	// Changed asks the caller to search again.
	Changed bool
	Opened  *SearchOpen
}

// SearchPanel shows search results grouped by file.
func SearchPanel(c *ui.Context, results []SearchFileResult, s *SearchPanelState, opts SearchPanelOptions) SearchPanelResult {
	if s == nil {
		panic("mujica: SearchPanel needs a state")
	}
	if !s.started {
		s.started, s.selected = true, -1
	}
	d := Density(c)
	var res SearchPanelResult
	res.Element = codePanel(c).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	res.Element.Children(func() {
		ui.Column(c).Shrink(0).Padding(d.Unit() * 2).Gap(d.Unit() * 2).Children(func() {
			res.Changed = SearchInput(c, &s.Query, SearchInputOptions{Disabled: opts.Disabled}).Changed()
			ui.Row(c).Gap(d.Unit() * 2).Children(func() {
				for _, f := range []struct {
					v   *string
					key string
				}{{&s.Include, "search_panel.include"}, {&s.Exclude, "search_panel.exclude"}} {
					r := TextInput(c, f.v, TextInputOptions{Placeholder: msg(c, f.key), Label: msg(c, f.key), Disabled: opts.Disabled})
					r.Element.Grow(1).MinWidth(0)
					res.Changed = res.Changed || r.Changed()
				}
			})
		})
		if res.Changed {
			slog.Debug("mujica: search_panel query", "query", s.Query, "include", s.Include, "exclude", s.Exclude)
		}
		searchPanelFlatten(results, s)
		if len(results) == 0 || opts.Status != DataReady {
			text := msg(c, "search_panel.idle")
			if s.Query != "" {
				text = msg(c, "search_panel.none")
			}
			res.Changed = dataPlaceholder(c, opts.Status, opts.Error, text, text) || res.Changed
			return
		}
		total := 0
		for _, f := range results {
			total += len(f.Matches)
		}
		codeCaption(c, fmt.Sprintf(msg(c, "search_panel.summary"), total, len(results))).Padding(0, d.Unit()*3, d.Unit(), d.Unit()*3)
		s.rows.Selected = &s.selected
		s.rows.Key = func(i int) any { return s.flat[i] }
		s.rows.Label = func(i int) string { return searchPanelLabel(results, s.flat[i]) }
		list := codeRows(c, &s.rows, len(s.flat), func(i int) { searchPanelRow(c, results, s, s.flat[i], disabled) })
		if list.Submitted() && !disabled && s.selected >= 0 {
			r := s.flat[s.selected]
			f := results[r.file]
			if r.match < 0 {
				searchPanelFold(s, f.Path, !s.Collapsed.Has(f.Path))
			} else {
				res.Opened = &SearchOpen{Path: f.Path, Line: f.Matches[r.match].Line}
				slog.Debug("mujica: search_panel open", "line", res.Opened.Line)
			}
		}
		if !disabled && s.selected >= 0 && s.flat[s.selected].match < 0 {
			path := results[s.flat[s.selected].file].Path
			switch {
			case list.Shortcut(0, ui.KeyLeft):
				searchPanelFold(s, path, true)
			case list.Shortcut(0, ui.KeyRight):
				searchPanelFold(s, path, false)
			}
		}
	})
	return res
}

func searchPanelFold(s *SearchPanelState, path string, fold bool) {
	if fold {
		s.Collapsed.Add(path)
	} else {
		s.Collapsed.Remove(path)
	}
	s.folded = -1
	slog.Debug("mujica: search_panel fold", "folded", fold)
}

// searchPanelFlatten relists rows when results or folds change.
func searchPanelFlatten(results []SearchFileResult, s *SearchPanelState) {
	var first *SearchFileResult
	if len(results) > 0 {
		first = &results[0]
	}
	if s.data == first && s.n == len(results) && s.folded == s.Collapsed.Len() {
		return
	}
	s.data, s.n, s.folded = first, len(results), s.Collapsed.Len()
	s.flat = s.flat[:0]
	for fi, f := range results {
		s.flat = append(s.flat, searchRow{fi, -1})
		if !s.Collapsed.Has(f.Path) {
			for mi := range f.Matches {
				s.flat = append(s.flat, searchRow{fi, mi})
			}
		}
	}
	s.selected = min(s.selected, len(s.flat)-1)
}

func searchPanelLabel(results []SearchFileResult, r searchRow) string {
	f := results[r.file]
	if r.match < 0 {
		return f.Path
	}
	m := f.Matches[r.match]
	return fmt.Sprintf("%d %s", m.Line, strings.TrimSpace(m.Text))
}

func searchPanelRow(c *ui.Context, results []SearchFileResult, s *SearchPanelState, r searchRow, disabled bool) {
	k := Tokens(c)
	f := results[r.file]
	if r.match < 0 {
		chev, key := "chevron-down", "search_panel.fold"
		if s.Collapsed.Has(f.Path) {
			chev, key = "chevron-right", "search_panel.unfold"
		}
		if iconAction(c, icons.Must(chev), fmt.Sprintf(msg(c, key), f.Path)).Clicked() && !disabled {
			searchPanelFold(s, f.Path, !s.Collapsed.Has(f.Path))
		}
		ui.Icon(c, icons.Must("file-text")).Size(16, 16).Shrink(0).TextColor(k.TextMuted)
		ui.Text(c, f.Path).FontWeight(600).SingleLine().Grow(1).MinWidth(0)
		codeCaption(c, fmt.Sprint(len(f.Matches))).Shrink(0)
		return
	}
	m := f.Matches[r.match]
	if m.Start < 0 || m.End > len(m.Text) || m.Start > m.End {
		panic("mujica: SearchMatch range outside its Text")
	}
	codeMono(c, fmt.Sprint(m.Line)).Width(codeTextWidth(c, "99999")).TextColor(k.TextMuted).Shrink(0).Margin(0, 0, 0, 28)
	from := max(0, m.Start-40)
	before := strings.TrimLeft(m.Text[from:m.Start], " \t")
	if from > 0 {
		before = "…" + before
	}
	spans := []ui.Span{{Text: before}, {Text: m.Text[m.Start:m.End], Background: k.Selection, Underline: true, DecorationColor: k.AccentText, Weight: 600}, {Text: m.Text[m.End:]}}
	ui.RichText(c, spans...).SingleLine().Grow(1).MinWidth(0)
}
