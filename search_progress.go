package mujica

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// SearchStage describes one caller-owned search step.
type SearchStage struct {
	Label string
	Done  bool
}

// SearchProgressOptions describe the caller's search.
type SearchProgressOptions struct {
	// Query is the keyword searched; required.
	Query string
	// Read and Total count read and discovered sources.
	Read, Total int
	Stages      []SearchStage
	// Done folds progress into a summary users can reopen.
	Done bool
}

// searchProgressState tracks completion and manual folding.
type searchProgressState struct {
	open, done bool
}

// SearchProgress displays keywords, read counts and collapsible stages.
func SearchProgress(c *ui.Context, opts SearchProgressOptions) ChoiceResult {
	if opts.Query == "" || opts.Read < 0 || opts.Total < 0 || opts.Read > opts.Total {
		panic("mujica: SearchProgress needs a Query and 0 <= Read <= Total")
	}
	r := ChoiceResult{Element: ui.Column(c).MinWidth(0)}
	st := ui.Local(r.Element, "search", func() searchProgressState { return searchProgressState{open: true} })
	if opts.Done && !st.done {
		// Finishing folds it once; the user may open it again.
		st.open = false
		slog.Debug("mujica: search_progress complete", "read", opts.Read)
	}
	if st.done && !opts.Done {
		st.open = true
		slog.Debug("mujica: search_progress restart")
	}
	st.done = opts.Done
	title := fmt.Sprintf(msg(c, "search_progress.searching"), opts.Query)
	if opts.Done {
		title = fmt.Sprintf(msg(c, "search_progress.done"), opts.Query, opts.Read)
	}
	r.Element.Children(func() { r.changed = searchProgressBody(c, st, title, opts) })
	if r.changed {
		slog.Debug("mujica: search_progress fold", "open", st.open)
	}
	return r
}

// searchProgressBody is the collapsible header and the stages.
func searchProgressBody(c *ui.Context, st *searchProgressState, title string, opts SearchProgressOptions) bool {
	k := Tokens(c)
	u := Density(c).Unit()
	small := FontSize(c, theme.CaptionSize)
	return Collapsible(c, &st.open, CollapsibleOptions{Title: title, Trigger: func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).TextColor(k.TextMuted).Children(func() {
			if opts.Done {
				ui.Icon(c, icons.Must("search")).Size(14, 14)
			} else {
				Spinner(c, SpinnerOptions{Size: 14, Label: title})
			}
			ui.Text(c, title).FontSize(small).SingleLine()
		})
	}}, func() {
		ui.Text(c, fmt.Sprintf(msg(c, "search_progress.read"), opts.Read, opts.Total)).FontSize(small).FontFeatures("tnum").TextColor(k.TextMuted)
		ui.Column(c).Gap(u).Role(ui.RoleList).Children(func() {
			for _, s := range opts.Stages {
				ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Role(ui.RoleListItem).Children(func() {
					if s.Done {
						ui.Icon(c, icons.Must("check")).Size(14, 14).TextColor(k.Success)
					} else {
						ui.Icon(c, icons.Must("circle")).Size(14, 14).TextColor(k.TextMuted)
					}
					t := ui.Text(c, s.Label).FontSize(small)
					if !s.Done {
						t.TextColor(k.TextMuted)
					}
				})
			}
		})
	}).Changed()
}
