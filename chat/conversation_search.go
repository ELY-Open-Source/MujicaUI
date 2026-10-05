package chat

import (
	"log/slog"
	"unicode"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ConversationSearchHit is one conversation the caller found for the query.
type ConversationSearchHit struct {
	ID, Title string
	// Snippet is matching message text, shown under the title.
	Snippet string
}

// ConversationSearchOptions configure a ConversationSearch.
type ConversationSearchOptions struct {
	// Loading shows the caller's search is under way.
	Loading bool
	// Error describes a failed search, shown with a retry button.
	Error       string
	Placeholder string
	Label       string
	Disabled    bool
}

// ConversationSearchView is a built ConversationSearch and this frame's events.
type ConversationSearchView struct {
	Element *ui.Element
	// QueryChanged asks the caller to search for the new query.
	QueryChanged bool
	// Chosen is the hit opened this frame, or empty.
	Chosen string
	Retry  bool
}

type conversationSearchState struct {
	cursor int
	rows   ui.ListState
	enter  bool
}

type conversationSearchKey struct{}

// ConversationSearch shows caller-supplied hits for a query, highlighting matches.
func ConversationSearch(c *ui.Context, query *string, hits []ConversationSearchHit, opts ConversationSearchOptions) ConversationSearchView {
	if query == nil {
		panic("mujica: ConversationSearch needs a query")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	col := ui.Column(c).Gap(d.Unit() * 2).MinHeight(0).Disabled(opts.Disabled)
	v := ConversationSearchView{Element: col}
	st := ui.Local(col, conversationSearchKey{}, func() *conversationSearchState { return &conversationSearchState{cursor: -1} })
	s := *st
	s.cursor = min(s.cursor, len(hits)-1)
	col.Children(func() {
		in := input.SearchInput(c, query, input.SearchInputOptions{Placeholder: opts.Placeholder, Label: opts.Label, Disabled: opts.Disabled})
		v.QueryChanged = in.Changed()
		if in.Input.Shortcut(0, ui.KeyDown) && !in.Input.IsDisabled() && len(hits) > 0 {
			s.enter, s.cursor = true, max(s.cursor, 0)
		}
		switch {
		case opts.Error != "":
			ui.Row(c).Gap(d.Unit() * 2).AlignItems(ui.Center).TextColor(k.Danger).Children(func() {
				ui.Icon(c, icons.Must("circle-alert")).Size(16, 16)
				ui.Text(c, opts.Error).Role(ui.RoleStatus).Grow(1)
				v.Retry = input.Button(c, core.Msg(c, "conversation_search.retry"), input.ButtonOptions{Variant: input.Secondary, Icon: icons.Must("refresh-cw")}).Clicked()
			})
		case opts.Loading && len(hits) == 0:
			core.Spinner(c, core.SpinnerOptions{Label: core.Msg(c, "conversation_search.searching"), ShowLabel: true})
		case len(hits) == 0 && *query != "":
			ui.Text(c, core.Msg(c, "conversation_search.no_results")).TextColor(k.TextMuted).Role(ui.RoleStatus)
		case len(hits) > 0:
			v.Chosen = conversationSearchList(c, s, *query, hits, opts.Loading)
		}
	})
	if v.QueryChanged {
		s.cursor = -1
		slog.Debug("mujica: conversation_search query", "len", len(*query))
	}
	if v.Chosen != "" {
		slog.Debug("mujica: conversation_search choose", "id", v.Chosen)
	}
	return v
}

func conversationSearchList(c *ui.Context, s *conversationSearchState, q string, hits []ConversationSearchHit, loading bool) string {
	k := core.Tokens(c)
	d := core.Density(c)
	chosen := ""
	s.rows.Key = func(i int) any { return hits[i].ID }
	s.rows.Label = func(i int) string { return hits[i].Title }
	s.rows.Selected = &s.cursor
	// Rows share one height, so the list sizes itself early.
	rowH := d.ControlHeight() + core.FontSize(c, theme.CaptionSize)*theme.LineHeight + 2
	restore := data.DataSelectTheme(c, 0)
	list := ui.List(c, &s.rows, len(hits), func(i int) {
		h := hits[i]
		row := ui.Column(c).Height(rowH).Padding(d.Unit(), d.Unit()*2).Gap(2).Justify(ui.Center).ClipY()
		if i == s.cursor {
			row.DrawOver(func(p *ui.Painter, r ui.Rect) {
				p.FillPath(internal.Diamond(r.X+3, r.Y+r.H/2, 6), k.AccentText)
			})
		}
		row.Children(func() {
			ui.RichText(c, conversationSearchSpans(h.Title, q, k)...).SingleLine()
			if h.Snippet != "" {
				ui.RichText(c, conversationSearchSpans(h.Snippet, q, k)...).SingleLine().TextColor(k.TextMuted).FontSize(core.FontSize(c, theme.CaptionSize))
			}
		})
		if !row.IsDisabled() && row.Clicked() {
			chosen = h.ID
		}
	}).Height(min(360, float32(len(hits))*rowH+2*d.Unit()+2)).Padding(d.Unit()).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Radius(theme.ControlRadius).Label(core.Msg(c, "conversation_search.results"))
	restore()
	internal.Ring(list, k.Focus, theme.ControlRadius)
	if s.enter {
		s.enter = false
		list.Focus()
	}
	if loading {
		core.Spinner(c, core.SpinnerOptions{Size: 12, Label: core.Msg(c, "conversation_search.searching")})
	}
	if list.Submitted() && !list.IsDisabled() && s.cursor >= 0 && s.cursor < len(hits) {
		chosen = hits[s.cursor].ID
	}
	return chosen
}

// conversationSearchSpans marks each case-insensitive match of q in text.
func conversationSearchSpans(text, q string, k theme.Tokens) []ui.Span {
	r := []rune(text)
	var spans []ui.Span
	at := 0
	for _, m := range conversationSearchMatches(text, q) {
		if m[0] > at {
			spans = append(spans, ui.Span{Text: string(r[at:m[0]])})
		}
		spans = append(spans, ui.Span{Text: string(r[m[0]:m[1]]), Weight: 700, Background: k.Selection, Underline: true})
		at = m[1]
	}
	if at < len(r) || len(spans) == 0 {
		spans = append(spans, ui.Span{Text: string(r[at:])})
	}
	return spans
}

// conversationSearchMatches returns rune ranges where q matches, ignoring case.
func conversationSearchMatches(text, q string) [][2]int {
	t, p := []rune(text), []rune(q)
	if len(p) == 0 {
		return nil
	}
	var out [][2]int
	for i := 0; i+len(p) <= len(t); {
		ok := true
		for j := range p {
			if unicode.ToLower(t[i+j]) != unicode.ToLower(p[j]) {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, [2]int{i, i + len(p)})
			i += len(p)
			continue
		}
		i++
	}
	return out
}
