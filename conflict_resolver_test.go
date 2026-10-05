package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

const conflictText = "intro\n<<<<<<< HEAD\nours one\n=======\ntheirs one\n>>>>>>> feature\nmiddle\n<<<<<<< HEAD\nours two\n||||||| base\nold two\n=======\ntheirs two\n>>>>>>> feature\nend"

func TestConflictParse(t *testing.T) {
	parts, bad := conflictParse([]string{"a", "<<<<<<< HEAD", "x", "=======", "y", ">>>>>>> br", "b"})
	if bad != 0 || len(parts) != 3 || !parts[1].conflict || parts[1].oursRef != "HEAD" || parts[1].theirsRef != "br" || parts[1].line != 2 {
		t.Fatalf("parts %+v bad %d", parts, bad)
	}
	parts, bad = conflictParse([]string{"a", "<<<<<<< HEAD", "x", "======="})
	if bad != 2 || len(parts) != 1 || parts[0].conflict {
		t.Fatalf("unclosed: %+v bad %d", parts, bad)
	}
}

func TestConflictResolverAcceptAndCount(t *testing.T) {
	var s ConflictResolverState
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Scroll(c).Fill().Children(func() {
			if ConflictResolver(c, &s, conflictText, ConflictResolverOptions{}).Changed() {
				changes++
			}
		})
	}, 800, 900)
	if !tt.HasText("0 of 2 conflicts resolved") || !tt.HasText("Conflict 2 · line 8") {
		t.Fatalf("start: %q", tt.Texts())
	}
	if err := tt.Click("Accept incoming"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("1 of 2 conflicts resolved") || changes != 1 {
		t.Fatalf("after one: %q", tt.Texts())
	}
	if err := tt.Click("Accept both"); err != nil {
		t.Fatal(err)
	}
	if got := s.Result(); got != "intro\nours one\ntheirs one\nmiddle\n<<<<<<< HEAD\nours two\n||||||| base\nold two\n=======\ntheirs two\n>>>>>>> feature\nend" {
		t.Fatalf("result %q", got)
	}
	if err := tt.Click("Undo"); err != nil {
		t.Fatal(err)
	}
	if d, n := s.Resolved(); d != 0 || n != 2 || changes != 3 {
		t.Fatalf("undo: %d/%d changes %d", d, n, changes)
	}
}

func TestConflictResolverBrokenAndDisabled(t *testing.T) {
	var s ConflictResolverState
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ConflictResolver(c, &s, "a\n<<<<<<< HEAD\nx", ConflictResolverOptions{})
	}, 600, 300)
	if !tt.HasText("Conflict markers are not closed after line 2") {
		t.Fatalf("broken: %q", tt.Texts())
	}
	if s.Result() != "a\n<<<<<<< HEAD\nx" {
		t.Fatalf("broken result %q", s.Result())
	}
	var d ConflictResolverState
	tt = ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ConflictResolver(c, &d, conflictText, ConflictResolverOptions{Disabled: true})
	}, 800, 900)
	_ = tt.Click("Accept current")
	if done, _ := d.Resolved(); done != 0 {
		t.Fatal("disabled resolver accepted")
	}
}

func TestConflictResolverSnapshot(t *testing.T) {
	snapshotAfter(t, "conflict_resolver", 640, 560, func(c *ui.Context) {
		s := ui.Local(c.Root(), "s", func() *ConflictResolverState { return &ConflictResolverState{} })
		ui.Column(c).Padding(16).Children(func() {
			ConflictResolver(c, *s, conflictText, ConflictResolverOptions{})
		})
	}, func(tt *ui.Tester) { _ = tt.Click("Accept incoming") })
}
