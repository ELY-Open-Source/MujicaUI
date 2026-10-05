package git

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestCommitInputCommitOnceWhileBusy(t *testing.T) {
	m := CommitMessage{}
	busy := false
	commits, amends := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := CommitInput(c, &m, CommitInputOptions{Busy: busy})
		if r.Committed() {
			commits++
		}
		if r.Amended() {
			amends++
		}
	}, 600, 400)
	_ = tt.Click("Commit")
	if commits != 0 {
		t.Fatal("empty title committed")
	}
	if err := tt.Click("Summary"); err != nil {
		t.Fatal(err)
	}
	tt.Type("Add git components")
	tt.Key(0, ui.KeyEnter)
	_ = tt.Click("Commit")
	if commits != 1 {
		t.Fatalf("commits %d, want 1 before Busy clears", commits)
	}
	busy = true
	tt.Frame()
	_ = tt.Click("Commit")
	_ = tt.Click("Amend")
	busy = false
	tt.Frame()
	if commits != 1 || amends != 0 {
		t.Fatalf("busy: commits %d amends %d", commits, amends)
	}
	_ = tt.Click("Amend")
	if amends != 1 {
		t.Fatalf("amends %d", amends)
	}
}

func TestCommitInputTitleHint(t *testing.T) {
	m := CommitMessage{Title: strings.Repeat("x", 55)}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CommitInput(c, &m, CommitInputOptions{})
	}, 600, 400)
	if _, ok := tt.Find("Summary is 5 characters over the suggested length"); !ok {
		t.Fatal("over-length hint not named")
	}
	m.Title = "short"
	tt.Frame()
	if _, ok := tt.Find("5 of 50 characters"); !ok {
		t.Fatal("count not named")
	}
}

func TestCommitInputDisabledAndInvalid(t *testing.T) {
	m := CommitMessage{Title: "x"}
	acted := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := CommitInput(c, &m, CommitInputOptions{Disabled: true})
		acted = acted || r.Committed() || r.Amended()
	}, 600, 400)
	_ = tt.Click("Commit")
	_ = tt.Click("Amend")
	if acted {
		t.Fatal("disabled input committed")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		CommitInput(c, &m, CommitInputOptions{TitleLimit: -1})
	}, 300, 100, 1)
}

func TestCommitInputSnapshot(t *testing.T) {
	snapshot(t, "commit_input", 520, 280, func(c *ui.Context) {
		m := ui.Local(c.Root(), "msg", func() CommitMessage {
			return CommitMessage{Title: "Add CommitGraph lanes with merge lines", Body: "Lanes are laid out once per history."}
		})
		ui.Column(c).Padding(16).Children(func() {
			CommitInput(c, m, CommitInputOptions{})
		})
	})
}
