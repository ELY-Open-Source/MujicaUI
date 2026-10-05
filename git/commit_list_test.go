package git

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
)

var gitTestTime = time.Date(2026, 10, 5, 9, 30, 0, 0, time.UTC)

// gitTestHash is a stable, hash-like hex string for commit i.
func gitTestHash(i int) string { return fmt.Sprintf("%016x%024x", uint64(i)*0x9E3779B97F4A7C15, i) }

// gitHistory is a linear history of n commits, newest first.
func gitHistory(n int) []Commit {
	out := make([]Commit, n)
	for i := range out {
		out[i] = Commit{Hash: gitTestHash(n - i), Author: "Sakiko", When: gitTestTime.Add(-time.Duration(i) * time.Hour),
			Subject: fmt.Sprintf("Commit number %d", n-i)}
		if i+1 < n {
			out[i].Parents = []string{gitTestHash(n - i - 1)}
		}
	}
	return out
}

func TestCommitListChooseAndOpen(t *testing.T) {
	commits := gitHistory(5)
	commits[0].Refs = []string{"main"}
	var s data.ListState[string]
	changes, opens := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := CommitList(c, &s, commits, CommitListOptions{})
		r.Element.Height(300)
		if r.Changed() {
			changes++
		}
		if r.Submitted() {
			opens++
		}
	}, 700, 340)
	if !tt.HasText("main") || !tt.HasText("1715609") || !tt.HasText("2026-10-05 09:30") {
		t.Fatalf("row parts: %q", tt.Texts())
	}
	if err := tt.Click("Commit number 3"); err != nil {
		t.Fatal(err)
	}
	if key, _ := s.Selected(); key != commits[2].Hash || changes != 1 {
		t.Fatalf("selected %q changes %d", key, changes)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if key, _ := s.Selected(); key != commits[3].Hash || opens != 1 {
		t.Fatalf("keyboard: %q opens %d", key, opens)
	}
}

func TestCommitListEmptyAndInvalid(t *testing.T) {
	var s data.ListState[string]
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CommitList(c, &s, nil, CommitListOptions{})
	}, 400, 200)
	if !tt.HasText("No commits") {
		t.Fatalf("empty: %q", tt.Texts())
	}
	defer func() {
		if recover() == nil {
			t.Fatal("hashless commit did not panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		var s data.ListState[string]
		CommitList(c, &s, []Commit{{Subject: "x"}}, CommitListOptions{})
	}, 300, 100, 1)
}

func TestCommitListSnapshot(t *testing.T) {
	commits := gitHistory(4)
	commits[0].Refs = []string{"main", "v1.1"}
	snapshot(t, "commit_list", 640, 220, func(c *ui.Context) {
		s := ui.Local(c.Root(), "list", func() *data.ListState[string] { return &data.ListState[string]{} })
		ui.Column(c).Padding(16).Children(func() {
			CommitList(c, *s, commits, CommitListOptions{}).Element.Height(180)
		})
	})
}
