package chat

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestThinkingBlockFoldsAndTimes(t *testing.T) {
	open, thinking := false, true
	started := time.Now().Add(-5 * time.Second)
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ThinkingBlock(c, &open, ThinkingBlockOptions{Thinking: thinking, Started: started, Elapsed: 12 * time.Second}, func() {
			ui.Text(c, "Weigh the guest list")
		})
	}, 400, 200)
	if !tt.HasText("Thinking… 5 s") || tt.HasText("Weigh the guest list") {
		t.Fatalf("collapsed thinking: %q", tt.Texts())
	}
	_ = tt.Click("Thinking… 5 s")
	if !open || !tt.HasText("Weigh the guest list") {
		t.Fatal("click did not open the reasoning")
	}
	thinking = false
	tt.Frame()
	if !tt.HasText("Thought for 12 s") {
		t.Fatalf("finished header %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if open {
		t.Fatal("Space did not fold it")
	}
}

func TestThinkingBlockDefaultOpen(t *testing.T) {
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ThinkingBlock(c, &open, ThinkingBlockOptions{Elapsed: time.Second}, func() { ui.Text(c, "Notes") })
	}, 400, 200)
	if !tt.HasText("Notes") {
		t.Fatal("default open block hides its content")
	}
}

func TestThinkingBlockNeedsStart(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for Thinking without Started")
		}
	}()
	open := false
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ThinkingBlock(c, &open, ThinkingBlockOptions{Thinking: true}, func() {})
	}, 100, 100)
}

func TestThinkingBlockSnapshot(t *testing.T) {
	open := true
	snapshot(t, "thinking_block", 420, 150, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			ThinkingBlock(c, &open, ThinkingBlockOptions{Elapsed: 14 * time.Second}, func() {
				ui.Text(c, "The guest list has forty names; the hall seats thirty-six. Four must dine in the antechamber.")
			})
		})
	})
}
