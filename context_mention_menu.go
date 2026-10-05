package mujica

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"
)

// ContextKind is what a piece of context is.
type ContextKind int

const (
	ContextFile ContextKind = iota
	ContextDoc
	ContextWeb
)

func (k ContextKind) icon() string {
	switch k {
	case ContextFile:
		return "file"
	case ContextDoc:
		return "file-text"
	case ContextWeb:
		return "globe"
	}
	panic("mujica: unknown ContextKind")
}

func (k ContextKind) group() string {
	return [...]string{"context_mention_menu.file", "context_mention_menu.doc", "context_mention_menu.web"}[k]
}

// ContextItem is a file, document or web page the user can mention.
type ContextItem struct {
	// ID is the item's stable identity; required.
	ID     string
	Label  string
	Detail string
	Kind   ContextKind
}

// ContextMentionOptions configure a ContextMentionMenu.
type ContextMentionOptions struct {
	// Loading shows that the caller is fetching items for the query.
	Loading bool
}

// ContextMentionResult is the menu's popup, nil while closed, its query and the item chosen this frame.
type ContextMentionResult struct {
	Element      *ui.Element
	query        string
	queryChanged bool
	chosen       ContextItem
	ok           bool
}

// Query returns the text after "@"; QueryChanged reports a new one to fetch items for.
func (r ContextMentionResult) Query() string      { return r.query }
func (r ContextMentionResult) QueryChanged() bool { return r.queryChanged }

// Chosen returns the item picked by Enter, Tab or a click; its "@query" was taken out of the text.
func (r ContextMentionResult) Chosen() (ContextItem, bool) { return r.chosen, r.ok }

// ContextMentionMenu pops up items grouped by kind while the composer's text ends in "@query".
// The caller filters items for Query, possibly asynchronously, and attaches the chosen one.
func ContextMentionMenu(c *ui.Context, composer InputResult, value *string, items []ContextItem, opts ContextMentionOptions) ContextMentionResult {
	var rows []slashMenuRow
	var order []int
	for kind := ContextFile; kind <= ContextWeb; kind++ {
		for i, it := range items {
			if it.ID == "" || it.Label == "" {
				panic("mujica: ContextItem needs ID and Label")
			}
			if it.Kind == kind {
				rows = append(rows, slashMenuRow{label: it.Label, detail: it.Detail, icon: kind.icon(), group: msg(c, kind.group())})
				order = append(order, i)
			}
		}
	}
	if len(order) != len(items) {
		panic("mujica: unknown ContextKind")
	}
	st := ui.Local(composer.Element, slashMenuKey{"mention"}, func() slashMenuState { return slashMenuState{} })
	at := contextMentionAt(*value)
	var r ContextMentionResult
	if at >= 0 {
		r.query = (*value)[at+1:]
	}
	// A new query, or the menu opening, asks the caller for items.
	r.queryChanged = at >= 0 && (r.query != st.query || !st.open && st.dismissed == nil)
	empty := msg(c, "context_mention_menu.empty")
	if opts.Loading {
		empty = msg(c, "context_mention_menu.loading")
	}
	i, _ := slashMenuShow(c, composer, st, r.query, at >= 0, rows, slashMenuOptions{empty: empty, label: msg(c, "context_mention_menu.label"), loading: opts.Loading})
	r.Element = st.panel
	if i >= 0 {
		r.chosen, r.ok = items[order[i]], true
		*value = (*value)[:at]
		slog.Debug("mujica: context_mention_menu choose", "id", r.chosen.ID)
	}
	return r
}

// contextMentionAt returns where the trailing "@query" starts, or -1.
func contextMentionAt(s string) int {
	at := strings.LastIndexByte(s, '@')
	if at < 0 || strings.ContainsAny(s[at+1:], " \t\n") || at > 0 && !strings.ContainsRune(" \t\n", rune(s[at-1])) {
		return -1
	}
	return at
}
