package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func messageEditorTester(draft *string, saves, cancels *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := MessageEditor(c, draft, MessageEditorOptions{Original: "Hello"})
		if r.Saved() {
			*saves++
		}
		if r.Canceled() {
			*cancels++
		}
	}, 400, 300)
}

func TestMessageEditorSaveRules(t *testing.T) {
	draft, saves, cancels := "Hello", 0, 0
	tt := messageEditorTester(&draft, &saves, &cancels)
	_ = tt.Click("Save")
	_ = tt.Click("Edit message")
	tt.Key(ui.Cmd, ui.KeyEnter)
	if saves != 0 {
		t.Fatal("saved unchanged text")
	}
	tt.Type(" there")
	tt.Key(ui.Cmd, ui.KeyEnter)
	_ = tt.Click("Save")
	if saves != 2 || draft != "Hello there" {
		t.Fatalf("saves %d draft %q", saves, draft)
	}
	_ = tt.Click("Edit message")
	tt.Command("selectAll")
	tt.Type("   ")
	_ = tt.Click("Save")
	if saves != 2 {
		t.Fatal("saved blank text")
	}
}

func TestMessageEditorCancel(t *testing.T) {
	draft, saves, cancels := "Hello", 0, 0
	tt := messageEditorTester(&draft, &saves, &cancels)
	_ = tt.Click("Cancel")
	_ = tt.Click("Edit message")
	tt.Key(0, ui.KeyEscape)
	if cancels != 2 || saves != 0 {
		t.Fatalf("cancels %d saves %d", cancels, saves)
	}
}

func TestMessageEditorSnapshot(t *testing.T) {
	draft := "Could you seat the ambassador nearer the throne?"
	snapshot(t, "message_editor", 420, 200, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { MessageEditor(c, &draft, MessageEditorOptions{Original: "Seat the ambassador."}) })
	})
}
