package mujica

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// HistoryEntry is one command a shell ran.
type HistoryEntry struct {
	Command  string
	Time     time.Time
	ExitCode int
}

// CommandHistoryState is the filter and place of a CommandHistory.
type CommandHistoryState struct {
	Query    string
	rows     ui.ListState
	selected int
	shown    []int
	query    string
	n        int
	ran      bool
}

// CommandHistoryOptions configure a CommandHistory.
type CommandHistoryOptions struct {
	Disabled bool
}

// CommandHistoryResult carries the entry to replay.
type CommandHistoryResult struct {
	Element *ui.Element
	// Replayed is the entry to run again, or nil.
	Replayed *HistoryEntry
}

// CommandHistory lists commands newest first with search and replay.
func CommandHistory(c *ui.Context, entries []HistoryEntry, s *CommandHistoryState, opts CommandHistoryOptions) CommandHistoryResult {
	if s == nil {
		panic("mujica: CommandHistory needs a state")
	}
	var res CommandHistoryResult
	res.Element = codePanel(c).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	replay := func(i int) {
		if disabled {
			return
		}
		e := entries[i]
		res.Replayed = &e
		slog.Debug("mujica: command_history replay", "index", i)
	}
	res.Element.Children(func() {
		codeHeader(c, func() {
			SearchInput(c, &s.Query, SearchInputOptions{Disabled: opts.Disabled, Placeholder: msg(c, "command_history.search")}).Element.Grow(1)
		})
		commandHistoryFilter(entries, s)
		if len(s.shown) == 0 {
			text := msg(c, "command_history.empty")
			if len(entries) > 0 {
				text = msg(c, "command_history.none")
			}
			codeCaption(c, text).Padding(Density(c).Unit() * 4)
			return
		}
		s.rows.Selected = &s.selected
		s.rows.Key = func(i int) any { return s.shown[i] }
		s.rows.Label = func(i int) string { return entries[s.shown[i]].Command }
		list := codeRows(c, &s.rows, len(s.shown), func(i int) {
			e := entries[s.shown[i]]
			commandHistoryRow(c, e)
			if iconAction(c, icons.Must("refresh-cw"), fmt.Sprintf(msg(c, "command_history.replay"), e.Command)).Clicked() {
				replay(s.shown[i])
			}
		})
		if list.Submitted() && s.selected >= 0 && s.selected < len(s.shown) {
			replay(s.shown[s.selected])
		}
	})
	return res
}

// commandHistoryFilter refilters when entries or query change.
func commandHistoryFilter(entries []HistoryEntry, s *CommandHistoryState) {
	if s.ran && s.n == len(entries) && s.query == s.Query {
		return
	}
	s.ran, s.n, s.query = true, len(entries), s.Query
	s.shown = s.shown[:0]
	q := strings.ToLower(strings.TrimSpace(s.Query))
	for i := len(entries) - 1; i >= 0; i-- {
		if q == "" || strings.Contains(strings.ToLower(entries[i].Command), q) {
			s.shown = append(s.shown, i)
		}
	}
	s.selected = min(s.selected, len(s.shown)-1)
}

func commandHistoryRow(c *ui.Context, e HistoryEntry) {
	k := Tokens(c)
	ui.Column(c).Grow(1).MinWidth(0).Children(func() {
		codeMono(c, e.Command)
		ui.Row(c).Gap(Density(c).Unit() * 2).Children(func() {
			if !e.Time.IsZero() {
				codeCaption(c, e.Time.Format("15:04:05"))
			}
			if e.ExitCode != 0 {
				codeCaption(c, fmt.Sprintf(msg(c, "command_history.exit"), e.ExitCode)).TextColor(k.Danger)
			}
		})
	})
}
