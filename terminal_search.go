package mujica

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"
)

// TerminalSearchState is a search through a terminal's scrollback and screen.
type TerminalSearchState struct {
	Query string
	// Current is the chosen match, from 0.
	Current int
	query   string
	version int
	ran     bool
}

// TerminalSearchOptions configure a TerminalSearch.
type TerminalSearchOptions struct {
	Disabled bool
}

// TerminalSearch finds Query in term and shows the current match.
func TerminalSearch(c *ui.Context, term *TerminalState, s *TerminalSearchState, opts TerminalSearchOptions) ChoiceResult {
	if term == nil || s == nil {
		panic("mujica: TerminalSearch needs a terminal and a state")
	}
	term.init()
	fresh := terminalSearchRun(term, s)
	r := SearchInput(c, &s.Query, SearchInputOptions{Find: true, Matches: len(term.matches), Current: &s.Current,
		Disabled: opts.Disabled, Label: msg(c, "terminal_search.label")})
	var res ChoiceResult
	res.Element = r.Element
	if r.Changed() {
		s.Current = 0
		terminalSearchRun(term, s)
		res.changed = true
	}
	if r.Moved() {
		res.changed = true
	}
	if res.changed || fresh {
		term.current = s.Current
		if s.Current < len(term.matches) {
			term.reveal(term.matches[s.Current].pos.Line)
		}
		slog.Debug("mujica: terminal_search move", "matches", len(term.matches), "current", s.Current)
	}
	return res
}

// terminalSearchRun refinds matches; it reports a new query.
func terminalSearchRun(term *TerminalState, s *TerminalSearchState) bool {
	if s.ran && s.query == s.Query && s.version == term.version {
		return false
	}
	fresh := !s.ran || s.query != s.Query
	s.ran, s.query, s.version = true, s.Query, term.version
	term.matches = term.matches[:0]
	q := []rune(strings.ToLower(s.Query))
	if len(q) > 0 {
		for i := range term.Lines() {
			line := []rune(strings.ToLower(term.LineText(i)))
			for col := 0; col+len(q) <= len(line); col++ {
				if string(line[col:col+len(q)]) == string(q) {
					term.matches = append(term.matches, terminalMatch{TerminalPos{i, col}, len(q)})
					col += len(q) - 1
				}
			}
		}
	}
	s.Current = max(0, min(s.Current, len(term.matches)-1))
	term.current = s.Current
	return fresh
}
