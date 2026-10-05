package mujica

import (
	"log/slog"
	"reflect"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// SettingEntry identifies a searchable setting and its destination.
type SettingEntry struct {
	ID, Page    string
	Label       string
	Description string
	// Keywords match too, as synonyms.
	Keywords []string
}

// SettingsSearchOptions configure a SettingsSearch.
type SettingsSearchOptions struct {
	// Limit caps the results shown; zero means 8.
	Limit int
}

// SettingsSearchResult carries the selected destination.
type SettingsSearchResult struct {
	Element *ui.Element
	chosen  *SettingEntry
}

// Chosen returns the selected entry and reveals its row.
func (r SettingsSearchResult) Chosen() (SettingEntry, bool) {
	if r.chosen == nil {
		return SettingEntry{}, false
	}
	return *r.chosen, true
}

type settingsSearchState struct {
	cursor int
	source []SettingEntry
	query  string
	limit  int
	hits   []int
	ready  bool
}

// SettingsSearch highlights matches and navigates to settings.
func SettingsSearch(c *ui.Context, query *string, entries []SettingEntry, opts SettingsSearchOptions) SettingsSearchResult {
	if opts.Limit < 0 {
		panic("mujica: SettingsSearch Limit cannot be negative")
	}
	limit := opts.Limit
	if limit == 0 {
		limit = 8
	}
	d := Density(c)
	var r SettingsSearchResult
	r.Element = ui.Column(c).Gap(d.Unit() * 2).MinWidth(280)
	st := ui.Local(r.Element, "search", func() settingsSearchState { return settingsSearchState{} })
	r.Element.Children(func() {
		res := SearchInput(c, query, SearchInputOptions{Placeholder: msg(c, "settings_search.placeholder"), Label: msg(c, "settings_search.placeholder")})
		in := res.Input
		q := strings.ToLower(strings.TrimSpace(*query))
		settingsSearchPlan(st, entries, q, limit)
		hits := st.hits
		st.cursor = max(0, min(st.cursor, len(hits)-1))
		if res.Changed() {
			st.cursor = 0
		}
		disabled := in.IsDisabled()
		before := st.cursor
		switch {
		case disabled || len(hits) == 0:
		case in.Shortcut(0, ui.KeyDown):
			st.cursor = (st.cursor + 1) % len(hits)
		case in.Shortcut(0, ui.KeyUp):
			st.cursor = (st.cursor + len(hits) - 1) % len(hits)
		}
		if res.Submitted() && len(hits) > 0 && !disabled {
			r.chosen = &entries[hits[st.cursor]]
		}
		if q != "" && len(hits) == 0 {
			ui.Text(c, msg(c, "settings_search.none")).TextColor(Tokens(c).TextMuted).Role(ui.RoleStatus)
		}
		list := accountRows(c, len(hits), func(n int) any { return entries[hits[n]].ID }, func(n int) {
			i := hits[n]
			if settingsSearchHit(c, entries[i], q, n == st.cursor) {
				r.chosen = &entries[i]
			}
		})
		list.Role(ui.RoleList)
		if len(hits) > 8 && (st.cursor != before || res.Changed()) {
			ui.Local(list, "rows", func() ui.ListState { return ui.ListState{} }).ScrollIntoView(st.cursor)
		}
	})
	if r.chosen != nil {
		settingsReveal(c, r.chosen.ID)
		*query, st.cursor = "", 0
		slog.Debug("mujica: settings_search chosen", "id", r.chosen.ID)
	}
	return r
}

func settingsSearchPlan(st *settingsSearchState, entries []SettingEntry, q string, limit int) {
	if st.ready && st.query == q && st.limit == limit && reflect.DeepEqual(st.source, entries) {
		return
	}
	st.ready, st.query, st.limit = true, q, limit
	st.source = slices.Clone(entries)
	st.hits = st.hits[:0]
	ids := map[string]bool{}
	for i, e := range entries {
		if e.ID == "" || e.Page == "" || e.Label == "" || ids[e.ID] {
			panic("mujica: SettingEntry needs a unique ID, Page and Label")
		}
		ids[e.ID] = true
		st.source[i].Keywords = slices.Clone(e.Keywords)
		text := strings.ToLower(e.Label + " " + e.Description + " " + strings.Join(e.Keywords, " "))
		if q != "" && strings.Contains(text, q) && len(st.hits) < limit {
			st.hits = append(st.hits, i)
		}
	}
}

// settingsSearchHit builds a marked, clickable search result.
func settingsSearchHit(c *ui.Context, e SettingEntry, q string, current bool) bool {
	k := Tokens(c)
	d := Density(c)
	b := ui.ButtonBase(c).Key(e.ID).MinHeight(d.ControlHeight()).Padding(d.Unit(), d.Unit()*3).Justify(ui.Start).
		Radius(theme.ControlRadius).Role(ui.RoleListItem).Label(e.Label).Description(e.Description).Selected(current).Gap(d.Unit() * 2)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case current:
		b.Background(k.Selection)
	case b.Hovered():
		b.Background(k.SurfaceHover)
	}
	b.Children(func() {
		ui.Box(c).Size(8, 8).Shrink(0).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			if current {
				p.FillPath(internal.Diamond(internal.Snap(p, r.X+4), internal.Snap(p, r.Y+4), 8), k.AccentText)
			}
		})
		ui.Column(c).Gap(2).Children(func() {
			ui.RichText(c, settingsSearchSpans(c, e.Label, q, theme.BodySize)...)
			if e.Description != "" {
				ui.RichText(c, settingsSearchSpans(c, e.Description, q, theme.CaptionSize)...).TextColor(k.TextMuted)
			}
		})
	})
	return b.Clicked() && !b.IsDisabled()
}

// settingsSearchSpans highlights matches while preserving original Unicode text.
func settingsSearchSpans(c *ui.Context, s, q string, size float32) []ui.Span {
	k := Tokens(c)
	size = FontSize(c, size)
	plain := func(t string) ui.Span { return ui.Span{Text: t, Size: size, Color: k.Text} }
	if q == "" {
		return []ui.Span{plain(s)}
	}
	var out []ui.Span
	low := strings.ToLower(s)
	original := []rune(s)
	for {
		i := strings.Index(low, q)
		if i < 0 {
			break
		}
		start, count := utf8.RuneCountInString(low[:i]), utf8.RuneCountInString(q)
		if start > 0 {
			out = append(out, plain(string(original[:start])))
		}
		out = append(out, ui.Span{Text: string(original[start : start+count]), Size: size, Weight: 700, Background: k.Selection, Color: k.Text})
		original, low = original[start+count:], low[i+len(q):]
	}
	if len(original) > 0 {
		out = append(out, plain(string(original)))
	}
	return out
}
