package git

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func reviewSample() ReviewThread {
	return ReviewThread{File: "commit_graph.go", Line: 88, Messages: []ReviewMessage{
		{Author: "Uika", When: gitTestTime, Body: "Should merge lanes reuse the parent's color?"},
		{Author: "Sakiko", When: gitTestTime, Body: "Yes, the first parent keeps it."},
	}}
}

func TestReviewCommentReplyResolveReopen(t *testing.T) {
	thread := reviewSample()
	draft := ""
	var replies []string
	resolves, reopens := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := ReviewComment(c, thread, &draft, ReviewCommentOptions{})
		if s, ok := r.Replied(); ok {
			replies = append(replies, s)
		}
		if r.Resolved() {
			resolves++
			thread.Resolved = true
		}
		if r.Reopened() {
			reopens++
			thread.Resolved = false
		}
	}, 600, 500)
	if !tt.HasText("commit_graph.go:88") || !tt.HasText("Yes, the first parent keeps it.") {
		t.Fatalf("thread: %q", tt.Texts())
	}
	_ = tt.Click("Reply")
	if len(replies) != 0 {
		t.Fatal("empty reply sent")
	}
	if err := tt.Click("Reply text"); err != nil {
		t.Fatal(err)
	}
	tt.Type("  Agreed.  ")
	tt.Key(ui.Cmd, ui.KeyEnter)
	if len(replies) != 1 || replies[0] != "Agreed." || draft != "" {
		t.Fatalf("replies %q draft %q", replies, draft)
	}
	_ = tt.Click("Resolve")
	if resolves != 1 || !tt.HasText("Resolved") || tt.HasText("Reply text") {
		t.Fatalf("resolve: %d %q", resolves, tt.Texts())
	}
	_ = tt.Click("Reopen")
	if reopens != 1 || !tt.HasText("Reply") {
		t.Fatalf("reopen: %d", reopens)
	}
}

func TestReviewCommentDisabledAndInvalid(t *testing.T) {
	draft := "Text"
	acted := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := ReviewComment(c, reviewSample(), &draft, ReviewCommentOptions{Disabled: true})
		_, replied := r.Replied()
		acted = acted || replied || r.Resolved()
	}, 600, 500)
	_ = tt.Click("Resolve")
	_ = tt.Click("Reply")
	if acted || draft != "Text" {
		t.Fatal("disabled thread acted")
	}
	for _, th := range []ReviewThread{{Line: 1, Messages: reviewSample().Messages}, {File: "a", Messages: reviewSample().Messages}, {File: "a", Line: 1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%+v did not panic", th)
				}
			}()
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ReviewComment(c, th, &draft, ReviewCommentOptions{})
			}, 300, 100, 1)
		}()
	}
}

func TestReviewCommentSnapshot(t *testing.T) {
	snapshot(t, "review_comment", 560, 420, func(c *ui.Context) {
		draft := ui.Local(c.Root(), "draft", func() string { return "" })
		ui.Column(c).Padding(16).Children(func() {
			ReviewComment(c, reviewSample(), draft, ReviewCommentOptions{})
		})
	})
}
