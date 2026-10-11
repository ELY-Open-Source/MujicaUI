package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestQuoteReplyJumpsAndFlattensExcerpt(t *testing.T) {
	jumps, off := 0, false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(off).Children(func() {
			if QuoteReply(c, "Marie", "Let them\n  eat cake.").Clicked() {
				jumps++
			}
		})
	}, 400, 200)
	if !tt.HasText("Let them eat cake.") {
		t.Fatalf("excerpt %q", tt.Texts())
	}
	_ = tt.Click("Reply to Marie: Let them eat cake.")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if jumps != 2 {
		t.Fatalf("jumps %d", jumps)
	}
	off = true
	_ = tt.Click("Reply to Marie: Let them eat cake.")
	if jumps != 2 {
		t.Fatal("disabled quote jumped")
	}
}

func TestQuoteReplyNeedsAuthor(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic without an author")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		QuoteReply(c, "", "x")
	}, 100, 100)
}

func TestQuoteReplySnapshot(t *testing.T) {
	snapshot(t, "quote_reply", 420, 70, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			QuoteReply(c, "Herald", "The ball begins at nine in the Hall of Mirrors, and every guest must arrive by the east door.")
		})
	})
}
