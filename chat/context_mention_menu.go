package chat

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
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

// ContextItem describes a file, document or webpage mention.
type ContextItem struct {
	// ID is the item's stable identity; required.
	ID     string
	Label  string
	Detail string
	Kind   ContextKind
}

// ContextMentionOptions configure a ContextMentionMenu.
type ContextMentionOptions struct {
	// Loading indicates caller-owned asynchronous fetching.
	Loading bool
}

// ContextMentionResult carries the popup, query and chosen item.
type ContextMentionResult struct {
	Element      ui.Element
	query        string
	queryChanged bool
	chosen       ContextItem
	ok           bool
}

// Query returns mention text; QueryChanged requests fresh results.
func (r ContextMentionResult) Query() string      { return r.query }
func (r ContextMentionResult) QueryChanged() bool { return r.queryChanged }

// Chosen returns the selected item and consumes its trigger.
func (r ContextMentionResult) Chosen() (ContextItem, bool) { return r.chosen, r.ok }

// ContextMentionMenu groups results for the active mention query.
// The caller filters results synchronously or asynchronously.
func ContextMentionMenu(c *ui.Context, composer input.InputResult, value *string, items []ContextItem, opts ContextMentionOptions) ContextMentionResult {
	if !composer.Element.Valid() || !composer.Input.Valid() || value == nil {
		panic("mujica: ContextMentionMenu needs composer and value")
	}
	data := ui.Local(composer.Element, "mention-data", func() mentionMenuData { return mentionMenuData{} })
	locale := core.CurrentLocale(c)
	if !data.built || locale != data.locale || !slices.Equal(data.items, items) {
		data.built, data.locale = true, locale
		data.items = append(data.items[:0], items...)
		data.rows, data.order = data.rows[:0], data.order[:0]

		for kind := ContextFile; kind <= ContextWeb; kind++ {
			for i, it := range items {
				if it.ID == "" || it.Label == "" {
					panic("mujica: ContextItem needs ID and Label")
				}
				if it.Kind == kind {
					data.rows = append(data.rows, slashMenuRow{label: it.Label, detail: it.Detail, icon: kind.icon(), group: core.Msg(c, kind.group())})
					data.order = append(data.order, i)
				}
			}
		}
		if len(data.order) != len(items) {
			panic("mujica: unknown ContextKind")
		}
	}
	rows, order := data.rows, data.order
	st := ui.Local(composer.Element, slashMenuKey{"mention"}, func() slashMenuState { return slashMenuState{} })
	at := contextMentionAt(*value)
	var r ContextMentionResult
	if at >= 0 {
		r.query = (*value)[at+1:]
	}
	// Opening or changing queries requests caller-owned results.
	r.queryChanged = at >= 0 && !composer.Element.IsDisabled() && composer.Input.Focused() && (r.query != st.query || !st.open && st.dismissed == nil)
	empty := core.Msg(c, "context_mention_menu.empty")
	if opts.Loading {
		empty = core.Msg(c, "context_mention_menu.loading")
	}
	i, _, panel := slashMenuShow(c, composer, st, r.query, at >= 0, rows, slashMenuOptions{empty: empty, label: core.Msg(c, "context_mention_menu.label"), loading: opts.Loading})
	r.Element = panel
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

type mentionMenuData struct {
	built  bool
	locale core.Locale
	items  []ContextItem
	rows   []slashMenuRow
	order  []int
}
