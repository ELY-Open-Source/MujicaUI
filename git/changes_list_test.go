package git

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func changesSample() []ChangedFile {
	return []ChangedFile{
		{Path: "commit_graph.go", Status: GitAdded, Staged: true},
		{Path: "internal/git.go", Status: GitModified},
		{Path: "docs/old.md", Status: GitDeleted},
		{Path: "notes.txt", Status: GitUntracked},
	}
}

func TestChangesListGroupsAndStage(t *testing.T) {
	var s ChangesListState
	files := changesSample()
	var acted [][2]string
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := ChangesList(c, &s, files, ChangesListOptions{})
		r.Element.Height(320)
		if a, item, ok := r.Action(); ok {
			acted = append(acted, [2]string{a, item})
		}
	}, 600, 360)
	if !tt.HasText("Staged changes (1)") || !tt.HasText("Changes (3)") {
		t.Fatalf("groups: %q", tt.Texts())
	}
	for _, name := range []string{"Stage internal/git.go", "Unstage commit_graph.go", "Stage all", "Unstage all"} {
		if err := tt.Click(name); err != nil {
			t.Fatal(err)
		}
	}
	want := [][2]string{{"stage", "internal/git.go"}, {"unstage", "commit_graph.go"}, {"stage_all", ""}, {"unstage_all", ""}}
	for i := range want {
		if i >= len(acted) || acted[i] != want[i] {
			t.Fatalf("actions %q, want %q", acted, want)
		}
	}
	if err := tt.Click("old.md"); err != nil {
		t.Fatal(err)
	}
	if got := s.Selected(); got != 2 {
		t.Fatalf("selected %d, want 2", got)
	}
	tt.Key(0, ui.KeyDown)
	if got := s.Selected(); got != 3 {
		t.Fatalf("Down selected %d", got)
	}
}

func TestChangesListDisabledAndEmptyGroups(t *testing.T) {
	var s ChangesListState
	files := []ChangedFile{{Path: "a.go", Status: GitModified}}
	acted := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := ChangesList(c, &s, files, ChangesListOptions{Disabled: true})
		r.Element.Height(200)
		_, _, ok := r.Action()
		acted = acted || ok || r.Changed()
	}, 500, 240)
	if !tt.HasText("Staged changes (0)") {
		t.Fatalf("empty group: %q", tt.Texts())
	}
	_ = tt.Click("Stage a.go")
	_ = tt.Click("a.go")
	if acted || s.Selected() != -1 {
		t.Fatal("disabled list acted")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		var s ChangesListState
		ChangesList(c, &s, []ChangedFile{{Status: GitAdded}}, ChangesListOptions{})
	}, 300, 100, 1)
}

func TestChangesListSnapshot(t *testing.T) {
	files := changesSample()
	snapshot(t, "changes_list", 520, 300, func(c *ui.Context) {
		s := ui.Local(c.Root(), "changes", func() *ChangesListState { return &ChangesListState{} })
		ui.Column(c).Padding(16).Children(func() {
			ChangesList(c, *s, files, ChangesListOptions{}).Element.Height(260)
		})
	})
}
