package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func memorySample() []MemoryEntry {
	return []MemoryEntry{
		{ID: "1", Text: "Lady Sakiko prefers the east wing.", Source: "Chat, 3 Oct"},
		{ID: "2", Text: "The quartet needs two days' notice.", Source: "Plan, 1 Oct"},
		{ID: "3", Text: "Claret is served after nine.", Source: "Ledger"},
	}
}

func memoryTester(entries *[]MemoryEntry, opts MemoryPanelOptions, got *[]MemoryPanelResult) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		r := MemoryPanel(c, entries, opts)
		if r.Edited != "" || r.Deleted != "" {
			*got = append(*got, r)
		}
	}, 560, 560)
}

func TestMemoryPanelSearchEditDelete(t *testing.T) {
	entries := memorySample()
	var got []MemoryPanelResult
	tt := memoryTester(&entries, MemoryPanelOptions{}, &got)
	if !tt.HasText("3 entries") || !tt.HasText("Claret is served after nine.") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Search memory")
	tt.Type("quartet")
	if tt.HasText("Claret is served after nine.") || !tt.HasText("The quartet needs two days' notice.") {
		t.Fatalf("search: %q", tt.Texts())
	}
	_ = tt.Click("Edit")
	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("The quartet needs a week's notice.")
	_ = tt.Click("Save")
	if entries[1].Text != "The quartet needs a week's notice." || len(got) != 1 || got[0].Edited != "2" {
		t.Fatalf("edit: %+v %+v", entries[1], got)
	}
	_ = tt.Click("Delete")
	if !tt.HasText("Forget this memory?") || len(entries) != 3 {
		t.Fatal("delete did not ask first")
	}
	_ = tt.Click("Cancel")
	if len(entries) != 3 {
		t.Fatal("cancel deleted")
	}
	_ = tt.Click("Delete")
	tt.Key(0, ui.KeyTab) // Cancel has the focus; Tab to Delete
	tt.Key(0, ui.KeyEnter)
	if len(entries) != 2 || got[1].Deleted != "2" || !tt.HasText("No memories match") {
		t.Fatalf("delete: %+v %+v %q", entries, got, tt.Texts())
	}
}

func TestMemoryPanelDisabled(t *testing.T) {
	entries := memorySample()
	var got []MemoryPanelResult
	tt := memoryTester(&entries, MemoryPanelOptions{Disabled: true}, &got)
	_ = tt.Click("Edit")
	_ = tt.Click("Delete")
	if tt.HasText("Save") || tt.HasText("Forget this memory?") || len(got) != 0 {
		t.Fatal("disabled panel acted")
	}
}

func TestMemoryPanelManyEntriesVirtualized(t *testing.T) {
	entries := make([]MemoryEntry, 100000)
	for i := range entries {
		entries[i] = MemoryEntry{ID: fmt.Sprint(i), Text: fmt.Sprintf("fact %d", i), Source: "import"}
	}
	var got []MemoryPanelResult
	tt := memoryTester(&entries, MemoryPanelOptions{}, &got)
	if n := len(tt.Texts()); n > 300 {
		t.Fatalf("%d texts built for 100000 entries", n)
	}
	_ = tt.Click("Search memory")
	tt.Type("fact 99999")
	if !tt.HasText("fact 99999") || tt.HasText("fact 0") {
		t.Fatal("search over many entries failed")
	}
}

func TestMemoryPanelSnapshot(t *testing.T) {
	snapshot(t, "memory_panel", 520, 360, func(c *ui.Context) {
		entries := memorySample()
		ui.Column(c).Padding(16).Children(func() { MemoryPanel(c, &entries, MemoryPanelOptions{}) })
	})
}
