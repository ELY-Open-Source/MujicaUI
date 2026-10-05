package mujica

import (
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// chatSearchBackend answers queries one frame late, as an async caller does.
type chatSearchBackend struct {
	query, asked string
	hits         []ConversationSearchHit
	loading      bool
	queries      []string
	chosen       string
	err          string
	retries      int
}

func (b *chatSearchBackend) view(opts ConversationSearchOptions) func(c *ui.Context) {
	return func(c *ui.Context) {
		if b.loading && b.err == "" {
			b.loading, b.hits = false, nil
			for _, conv := range chatConversations() {
				if strings.Contains(strings.ToLower(conv.Title), strings.ToLower(b.asked)) {
					b.hits = append(b.hits, ConversationSearchHit{ID: conv.ID, Title: conv.Title, Snippet: "about " + conv.Title})
				}
			}
		}
		opts.Loading, opts.Error = b.loading, b.err
		v := ConversationSearch(c, &b.query, b.hits, opts)
		if v.QueryChanged {
			b.queries = append(b.queries, b.query)
			b.asked, b.loading = b.query, true
			c.Invalidate()
		}
		if v.Chosen != "" {
			b.chosen = v.Chosen
		}
		if v.Retry {
			b.retries++
		}
	}
}

func TestConversationSearchAsyncAndChoose(t *testing.T) {
	b := &chatSearchBackend{}
	tt := chatTester(420, 400, b.view(ConversationSearchOptions{Label: "Search conversations"}))
	tt.Key(0, ui.KeyTab)
	tt.Type("garden")
	tt.Frame()
	if len(b.queries) == 0 || b.queries[len(b.queries)-1] != "garden" {
		t.Fatalf("queries %q", b.queries)
	}
	if !tt.HasText("Garden maze") {
		t.Fatalf("results %q", tt.Texts())
	}
	if err := tt.Click("Garden maze"); err != nil {
		t.Fatal(err)
	}
	if b.chosen != "f" {
		t.Fatalf("chose %q", b.chosen)
	}
}

func TestConversationSearchKeyboard(t *testing.T) {
	b := &chatSearchBackend{}
	tt := chatTester(420, 400, b.view(ConversationSearchOptions{Label: "Search conversations"}))
	tt.Key(0, ui.KeyTab)
	tt.Type("o")
	tt.Frame()
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if b.chosen == "" || b.chosen != b.hits[1].ID {
		t.Fatalf("keys chose %q of %v", b.chosen, b.hits)
	}
}

func TestConversationSearchStates(t *testing.T) {
	b := &chatSearchBackend{query: "zzz", asked: "zzz"}
	tt := chatTester(420, 400, b.view(ConversationSearchOptions{}))
	if !tt.HasText("No conversations match") {
		t.Fatalf("texts %q", tt.Texts())
	}
	b.err = "Search is down"
	tt.Frame()
	if !tt.HasText("Search is down") {
		t.Fatal("error missing")
	}
	_ = tt.Click("Retry")
	if b.retries != 1 {
		t.Fatalf("retries %d", b.retries)
	}
}

func TestConversationSearchHighlights(t *testing.T) {
	got := conversationSearchMatches("Ballroom SEATING and seats", "seat")
	if !slices.Equal(got, [][2]int{{9, 13}, {21, 25}}) {
		t.Fatalf("matches %v", got)
	}
	if got := conversationSearchMatches("宫廷礼仪与宫廷音乐", "宫廷"); !slices.Equal(got, [][2]int{{0, 2}, {5, 7}}) {
		t.Fatalf("CJK matches %v", got)
	}
	if conversationSearchMatches("abc", "") != nil {
		t.Fatal("empty query matched")
	}
}

func TestConversationSearchDisabled(t *testing.T) {
	b := &chatSearchBackend{query: "o", hits: []ConversationSearchHit{{ID: "x", Title: "Royal"}}}
	tt := chatTester(420, 400, b.view(ConversationSearchOptions{Disabled: true}))
	_ = tt.Click("Royal")
	tt.Key(0, ui.KeyTab)
	tt.Type("q")
	if b.chosen != "" || b.query != "o" || len(b.queries) != 0 {
		t.Fatalf("disabled search acted: %q %q %v", b.chosen, b.query, b.queries)
	}
}

func TestConversationSearchSnapshot(t *testing.T) {
	snapshot(t, "conversation_search", 420, 260, func(c *ui.Context) {
		q := "seat"
		ui.Column(c).Padding(16).Children(func() {
			ConversationSearch(c, &q, []ConversationSearchHit{
				{ID: "b", Title: "Ballroom seating plan", Snippet: "…where to seat the duke…"},
				{ID: "c", Title: "Theatre seats", Snippet: "Box seats for the opera"},
			}, ConversationSearchOptions{})
		})
	})
}
