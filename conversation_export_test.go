package mujica

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestExportConversationFormats(t *testing.T) {
	msgs := chatMessages()
	md := ExportConversation("Banquet", msgs, ExportMarkdown)
	want := "# Banquet\n\n**User** · (2026-10-05 15:20)\n\nHow should guests be seated at the banquet?\n"
	if !strings.HasPrefix(md, want) || !strings.Contains(md, "**Lady Ann**") {
		t.Fatalf("markdown %q", md)
	}
	var doc struct {
		Title    string
		Messages []struct{ Role, Author, Text, Time string }
	}
	if err := json.Unmarshal([]byte(ExportConversation("Banquet", msgs, ExportJSON)), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.Title != "Banquet" || len(doc.Messages) != 3 || doc.Messages[1].Role != "assistant" || doc.Messages[2].Author != "Lady Ann" ||
		doc.Messages[0].Time != "2026-10-05T15:20:00Z" {
		t.Fatalf("json %+v", doc)
	}
	txt := ExportConversation("Banquet", msgs[:1], ExportText)
	if txt != "Banquet\n\nUser (2026-10-05 15:20):\nHow should guests be seated at the banquet?\n" {
		t.Fatalf("text %q", txt)
	}
	if got := ExportConversation("Empty", nil, ExportJSON); !strings.Contains(got, `"messages": []`) {
		t.Fatalf("empty json %q", got)
	}
}

func TestConversationExportRangeAndExport(t *testing.T) {
	var s ConversationExportState
	var last ConversationExportView
	selected := []int{2}
	tt := chatTester(560, 600, func(c *ui.Context) {
		if v := ConversationExport(c, &s, chatMessages(), ConversationExportOptions{Title: "Banquet", Selected: selected}); v.Exported {
			last = v
		}
	})
	_ = tt.Click("JSON")
	_ = tt.Click("Selected (1)")
	_ = tt.Click("Export")
	if !last.Exported || last.Format != ExportJSON || !strings.Contains(last.Content, "And the musicians?") ||
		strings.Contains(last.Content, "banquet?") {
		t.Fatalf("export %+v", last)
	}
	_ = tt.Click("Plain text")
	_ = tt.Click("Last")
	s.Last = 2
	tt.Frame()
	_ = tt.Click("Export")
	if last.Format != ExportText || strings.Contains(last.Content, "banquet?") || !strings.Contains(last.Content, "Seat by rank") {
		t.Fatalf("last-2 export %q", last.Content)
	}
}

func TestConversationExportDisabledAndChecks(t *testing.T) {
	var s ConversationExportState
	exported := false
	tt := chatTester(560, 600, func(c *ui.Context) {
		if ConversationExport(c, &s, chatMessages(), ConversationExportOptions{Title: "T", Disabled: true}).Exported {
			exported = true
		}
	})
	_ = tt.Click("JSON")
	_ = tt.Click("Export")
	if exported || s.Format != ExportMarkdown {
		t.Fatalf("disabled export acted: %v %v", exported, s.Format)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a bad selection")
		}
	}()
	ui.Render(func(c *ui.Context) {
		Use(c, Settings{})
		ConversationExport(c, &s, chatMessages(), ConversationExportOptions{Title: "T", Selected: []int{9}})
	}, 100, 100, 1)
}

func TestConversationExportSnapshot(t *testing.T) {
	snapshot(t, "conversation_export", 560, 470, func(c *ui.Context) {
		s := &ConversationExportState{}
		ui.Column(c).Padding(16).Children(func() {
			ConversationExport(c, s, chatMessages(), ConversationExportOptions{Title: "Banquet seating", Selected: []int{0, 2}})
		})
	})
}
