package mujica

import (
	"fmt"
	"log/slog"
	"path"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
)

// StackFrame is one call of a thread's stack, innermost first.
type StackFrame struct {
	ID, Name, File string
	Line           int
}

// DebugThread is a thread and its stack.
type DebugThread struct {
	ID, Name string
	Frames   []StackFrame
}

// CallStackState is the current frame and folded threads.
type CallStackState struct {
	// Current is the chosen frame's ID.
	Current string
	// Collapsed holds folded thread IDs.
	Collapsed ui.Selection[string]
	rows      ui.ListState
	selected  int
	flat      []callRow
	seen      string
	folds     int
	started   bool
}

type callRow struct{ thread, frame int }

// CallStackOptions configure a CallStack.
type CallStackOptions struct {
	Disabled bool
}

// CallStack lists frames by thread; choosing sets Current.
func CallStack(c *ui.Context, threads []DebugThread, s *CallStackState, opts CallStackOptions) ChoiceResult {
	if s == nil {
		panic("mujica: CallStack needs a state")
	}
	s.flat = s.flat[:0]
	at := -1
	for ti, t := range threads {
		s.flat = append(s.flat, callRow{ti, -1})
		if s.Collapsed.Has(t.ID) {
			continue
		}
		for fi, f := range t.Frames {
			if f.ID == s.Current {
				at = len(s.flat)
			}
			s.flat = append(s.flat, callRow{ti, fi})
		}
	}
	// Choose the row again after outside changes.
	if !s.started || s.seen != s.Current || s.folds != s.Collapsed.Len() || s.selected >= len(s.flat) {
		s.started, s.seen, s.selected, s.folds = true, s.Current, at, s.Collapsed.Len()
	}
	d := Density(c)
	var res ChoiceResult
	res.Element = codePanel(c).Disabled(opts.Disabled)
	disabled := res.Element.IsDisabled()
	res.Element.Children(func() {
		if len(threads) == 0 {
			codeCaption(c, msg(c, "call_stack.empty")).Padding(d.Unit() * 4)
			return
		}
		was := s.selected
		s.rows.Selected = &s.selected
		s.rows.Key = func(i int) any { return s.flat[i] }
		s.rows.Label = func(i int) string { return callLabel(threads, s.flat[i]) }
		list := codeRows(c, &s.rows, len(s.flat), func(i int) { callRowView(c, threads, s, s.flat[i], disabled) })
		if disabled {
			s.selected = was
			return
		}
		if s.selected < 0 || s.selected >= len(s.flat) {
			return
		}
		r := s.flat[s.selected]
		t := threads[r.thread]
		switch {
		case r.frame >= 0 && s.selected != was:
			s.Current = t.Frames[r.frame].ID
			s.seen, res.changed = s.Current, true
			slog.Debug("mujica: call_stack frame", "thread", t.ID, "depth", r.frame)
		case r.frame < 0 && (list.Submitted() ||
			!s.Collapsed.Has(t.ID) && list.Shortcut(0, ui.KeyLeft) || s.Collapsed.Has(t.ID) && list.Shortcut(0, ui.KeyRight)):
			callFold(s, t.ID)
		}
	})
	return res
}

func callFold(s *CallStackState, id string) {
	if s.Collapsed.Has(id) {
		s.Collapsed.Remove(id)
	} else {
		s.Collapsed.Add(id)
	}
	slog.Debug("mujica: call_stack fold", "thread", id, "folded", s.Collapsed.Has(id))
}

func callLabel(threads []DebugThread, r callRow) string {
	t := threads[r.thread]
	if r.frame < 0 {
		return t.Name
	}
	f := t.Frames[r.frame]
	return fmt.Sprintf("%s %s:%d", f.Name, path.Base(f.File), f.Line)
}

func callRowView(c *ui.Context, threads []DebugThread, s *CallStackState, r callRow, disabled bool) {
	t := threads[r.thread]
	if r.frame < 0 {
		chev, key := "chevron-down", "search_panel.fold"
		if s.Collapsed.Has(t.ID) {
			chev, key = "chevron-right", "search_panel.unfold"
		}
		if iconAction(c, icons.Must(chev), fmt.Sprintf(msg(c, key), t.Name)).Clicked() && !disabled {
			callFold(s, t.ID)
		}
		ui.Text(c, t.Name).FontWeight(600).SingleLine().Grow(1).MinWidth(0)
		codeCaption(c, fmt.Sprintf(msg(c, "call_stack.frames"), len(t.Frames))).Shrink(0)
		return
	}
	f := t.Frames[r.frame]
	codeMono(c, f.Name).Grow(1).MinWidth(0).Margin(0, 0, 0, 20)
	codeCaption(c, fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)).Shrink(0)
}
