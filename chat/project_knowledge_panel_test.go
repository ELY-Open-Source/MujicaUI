package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func knowledgeFiles() []KnowledgeFile {
	return []KnowledgeFile{{ID: "1", Name: "etiquette.pdf", Tokens: 120000}, {ID: "2", Name: "menu.md", Tokens: 50000}}
}

func TestProjectKnowledgePanelAddRemove(t *testing.T) {
	files := knowledgeFiles()
	adds := 0
	tt := chatTester(420, 400, func(c *ui.Context) {
		v := ProjectKnowledgePanel(c, files, ProjectKnowledgePanelOptions{Capacity: 400000})
		if v.Add {
			adds++
		}
		if v.Removed != "" {
			for i, f := range files {
				if f.ID == v.Removed {
					files = append(files[:i:i], files[i+1:]...)
					break
				}
			}
		}
	})
	if !tt.HasText("170000 of 400000 tokens (42%)") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Add files")
	if adds != 1 {
		t.Fatalf("adds %d", adds)
	}
	if err := tt.Click("Remove menu.md"); err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].ID != "1" || !tt.HasText("120000 of 400000 tokens (30%)") {
		t.Fatalf("files %v texts %q", files, tt.Texts())
	}
}

func TestProjectKnowledgePanelCapacityStates(t *testing.T) {
	for _, tc := range []struct {
		capacity int
		note     string
	}{{200000, "Nearly full"}, {150000, "Over capacity"}} {
		tt := chatTester(420, 400, func(c *ui.Context) {
			ProjectKnowledgePanel(c, knowledgeFiles(), ProjectKnowledgePanelOptions{Capacity: tc.capacity})
		})
		if !tt.HasText(tc.note) {
			t.Errorf("%d: texts %q", tc.capacity, tt.Texts())
		}
	}
	tt := chatTester(420, 400, func(c *ui.Context) {
		ProjectKnowledgePanel(c, nil, ProjectKnowledgePanelOptions{Capacity: 10})
	})
	if !tt.HasText("No files yet") {
		t.Fatal("empty state missing")
	}
}

func TestProjectKnowledgePanelDisabledAndKeys(t *testing.T) {
	removed := ""
	tt := chatTester(420, 400, func(c *ui.Context) {
		if v := ProjectKnowledgePanel(c, knowledgeFiles(), ProjectKnowledgePanelOptions{Capacity: 400000, Disabled: true}); v.Removed != "" {
			removed = v.Removed
		}
	})
	_ = tt.Click("Remove menu.md")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if removed != "" {
		t.Fatalf("disabled panel removed %q", removed)
	}
	tt = chatTester(420, 400, func(c *ui.Context) {
		if v := ProjectKnowledgePanel(c, knowledgeFiles(), ProjectKnowledgePanelOptions{Capacity: 400000}); v.Removed != "" {
			removed = v.Removed
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if removed != "2" {
		t.Fatalf("keys removed %q", removed)
	}
}

func TestProjectKnowledgePanelRejectsBadCapacity(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ProjectKnowledgePanel(c, nil, ProjectKnowledgePanelOptions{})
	}, 100, 100, 1)
}

func TestProjectKnowledgePanelSnapshot(t *testing.T) {
	snapshot(t, "project_knowledge_panel", 420, 240, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ProjectKnowledgePanel(c, knowledgeFiles(), ProjectKnowledgePanelOptions{Capacity: 200000})
		})
	})
}
