package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestSystemPromptEditorEditSaveRestore(t *testing.T) {
	value, saved := "Be brief.", "Be brief."
	saves, restores := 0, 0
	tt := chatTester(520, 360, func(c *ui.Context) {
		v := SystemPromptEditor(c, &value, SystemPromptEditorOptions{Default: "You are helpful.", Saved: saved})
		if v.Save {
			saves++
			saved = value
		}
		if v.Restored {
			restores++
		}
	})
	if tt.HasText("Unsaved changes") || !tt.HasText("9 characters") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	tt.Type(" Be kind.")
	if !tt.HasText("Unsaved changes") || !tt.HasText("18 characters") {
		t.Fatalf("texts %q", tt.Texts())
	}
	tt.Key(chatMod(), ui.KeyS)
	if saves != 1 || saved != "Be brief. Be kind." || tt.HasText("Unsaved changes") {
		t.Fatalf("saves %d saved %q", saves, saved)
	}
	_ = tt.Click("Restore default")
	if value != "You are helpful." || restores != 1 || !tt.HasText("Unsaved changes") {
		t.Fatalf("restore gave %q", value)
	}
	_ = tt.Click("Save")
	if saves != 2 {
		t.Fatalf("saves %d", saves)
	}
}

func TestSystemPromptEditorLimitAndDisabled(t *testing.T) {
	value, saves := "abc", 0
	tt := chatTester(520, 360, func(c *ui.Context) {
		if SystemPromptEditor(c, &value, SystemPromptEditorOptions{MaxLength: 5}).Save {
			saves++
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnd)
	tt.Type("defgh")
	tt.Key(chatMod(), ui.KeyS)
	_ = tt.Click("Save")
	if saves != 0 || !tt.HasText("3 characters over the limit") {
		t.Fatalf("over the limit: saves %d, texts %q", saves, tt.Texts())
	}
	value = "abcde"
	tt = chatTester(520, 360, func(c *ui.Context) {
		if SystemPromptEditor(c, &value, SystemPromptEditorOptions{Saved: "x", Disabled: true}).Save {
			saves++
		}
	})
	_ = tt.Click("Save")
	tt.Key(chatMod(), ui.KeyS)
	_ = tt.Click("Restore default")
	if saves != 0 || value != "abcde" {
		t.Fatalf("disabled editor acted: %d %q", saves, value)
	}
}

func TestSystemPromptEditorSnapshot(t *testing.T) {
	snapshot(t, "system_prompt_editor", 520, 280, func(c *ui.Context) {
		v := "You are a courtly assistant."
		ui.Column(c).Padding(16).Children(func() {
			SystemPromptEditor(c, &v, SystemPromptEditorOptions{Default: "You are helpful.", Saved: "x", MaxLength: 2000})
		})
	})
}
