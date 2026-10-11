package git

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
)

func TestFileHistoryDiffAction(t *testing.T) {
	var s data.ListState[string]
	commits := gitHistory(3)
	diff := ""
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := FileHistory(c, &s, "commit_graph.go", commits, FileHistoryOptions{})
		r.Element.Height(240)
		if a, item, ok := r.Action(); ok && a == "diff" {
			diff = item
		}
	}, 700, 280)
	if !tt.HasText("commit_graph.go") || !tt.HasText("History") {
		t.Fatalf("header: %q", tt.Texts())
	}
	r, _ := tt.Find("Commit number 2")
	tt.ClickAt(700-30, r.Y+r.H/2)
	if diff != commits[1].Hash {
		t.Fatalf("diff %q", diff)
	}
	_ = tt.Click("Commit number 1")
	tt.Key(0, ui.KeyUp)
	if key, _ := s.Selected(); key != commits[1].Hash {
		t.Fatalf("selected %q", key)
	}
}

func TestFileHistoryEmptyAndInvalid(t *testing.T) {
	var s data.ListState[string]
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		FileHistory(c, &s, "a.go", nil, FileHistoryOptions{}).Element.Height(200)
	}, 400, 240)
	if !tt.HasText("No history for this file") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		var s data.ListState[string]
		FileHistory(c, &s, "", nil, FileHistoryOptions{})
	}, 300, 120, 1)
}

func TestFileHistorySnapshot(t *testing.T) {
	commits := gitHistory(3)
	snapshot(t, "file_history", 640, 200, func(c *ui.Context) {
		s := ui.Local(c.Root(), "list", func() *data.ListState[string] { return &data.ListState[string]{} })
		ui.Column(c).Padding(16).Children(func() {
			FileHistory(c, *s, "internal/git.go", commits, FileHistoryOptions{}).Element.Height(170)
		})
	})
}
