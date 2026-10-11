package chat

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var mentionSample = []ContextItem{
	{ID: "w1", Label: "Versailles — Wikipedia", Detail: "en.wikipedia.org", Kind: ContextWeb},
	{ID: "f1", Label: "guests.csv", Detail: "~/ball", Kind: ContextFile},
	{ID: "d1", Label: "Protocol", Detail: "Notion", Kind: ContextDoc},
	{ID: "f2", Label: "seating.xlsx", Detail: "~/ball", Kind: ContextFile},
}

// contextMentionTester supplies filtered results one frame later.
type contextMentionTester struct {
	text    string
	items   []ContextItem
	loading bool
	chosen  []string
	queries []string
}

func (m *contextMentionTester) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Column(c).Fill().Justify(ui.End).Padding(16).Children(func() {
		p := PromptComposer(c, &m.text, PromptComposerOptions{Label: "Prompt"})
		r := ContextMentionMenu(c, p, &m.text, m.items, ContextMentionOptions{Loading: m.loading})
		if r.QueryChanged() {
			m.queries = append(m.queries, r.Query())
			m.items = nil
			for _, it := range mentionSample {
				if strings.Contains(strings.ToLower(it.Label), strings.ToLower(r.Query())) {
					m.items = append(m.items, it)
				}
			}
		}
		if it, ok := r.Chosen(); ok {
			m.chosen = append(m.chosen, it.ID)
		}
	})
}

func TestContextMentionMenuGroupsAndChooses(t *testing.T) {
	m := &contextMentionTester{}
	tt := ui.NewTester(m.view, 600, 500)
	_ = tt.Click("Prompt")
	tt.Type("Read @")
	tt.Frame()
	files, _ := tt.Find("Files")
	docs, _ := tt.Find("Documents")
	web, ok := tt.Find("Web pages")
	if !ok || !(files.Y < docs.Y && docs.Y < web.Y) {
		t.Fatalf("groups not in file, doc, web order: %q", tt.Texts())
	}
	tt.Type("s")
	tt.Frame()
	if len(m.queries) < 2 || m.queries[len(m.queries)-1] != "s" {
		t.Fatalf("queries %q", m.queries)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(m.chosen) != 1 || m.text != "Read " {
		t.Fatalf("chosen %v text %q", m.chosen, m.text)
	}
}

func TestContextMentionMenuOnlyAfterSpaceAndLoading(t *testing.T) {
	m := &contextMentionTester{loading: true}
	tt := ui.NewTester(m.view, 600, 500)
	_ = tt.Click("Prompt")
	tt.Type("mail@host")
	if tt.HasText("Files") || len(m.queries) != 0 {
		t.Fatal("an @ inside a word opened the menu")
	}
	tt.Type(" @x")
	m.items = nil
	tt.Frame()
	if !tt.HasText("Searching…") {
		t.Fatalf("loading state missing: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if len(m.chosen) != 0 {
		t.Fatal("Enter chose from an empty list")
	}
}

func TestContextMentionMenuBadItem(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for an item without ID")
		}
	}()
	m := &contextMentionTester{items: []ContextItem{{Label: "x"}}}
	ui.NewTester(m.view, 300, 300)
}

func TestContextMentionMenuSnapshot(t *testing.T) {
	snapshotAfter(t, "context_mention_menu", 560, 420, func(c *ui.Context) {
		(*ui.Local(c.Root(), "m", func() *contextMentionTester { return &contextMentionTester{items: mentionSample} })).view(c)
	}, func(tt *ui.Tester) {
		_ = tt.Click("Prompt")
		tt.Type("Compare @")
		tt.Frame()
	})
}

func TestContextMentionMenuAsyncFirstChoiceAndDisabledQuery(t *testing.T) {
	text, off := "@", false
	var items []ContextItem
	picked, queries := "", 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		p := PromptComposer(c, &text, PromptComposerOptions{Label: "Prompt", Disabled: off})
		r := ContextMentionMenu(c, p, &text, items, ContextMentionOptions{Loading: items == nil})
		if r.QueryChanged() {
			queries++
		}
		if it, ok := r.Chosen(); ok {
			picked = it.ID
		}
	}, 400, 300)
	_ = tt.Click("Prompt")
	items = []ContextItem{{ID: "first", Label: "First", Kind: ContextDoc}}
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if picked != "first" || text != "" {
		t.Errorf("async first choice: %q text %q", picked, text)
	}
	off, text = true, "@other"
	before := queries
	tt.Frame()
	if queries != before {
		t.Errorf("disabled menu issued %d queries", queries-before)
	}
}
