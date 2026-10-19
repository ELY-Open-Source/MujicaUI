package agent

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func permissionTester(open *bool, got *[]PermissionChoice) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Button(c, "Ask", input.ButtonOptions{}).Clicked() {
			*open = true
		}
		r := PermissionPrompt(c, open, PermissionPromptOptions{Title: "Write to the ledger?", Scope: "Write access to ~/court/ledger"})
		if r.Choice != PermissionNone {
			*got = append(*got, r.Choice)
		}
	}, 600, 400)
}

func TestPermissionPromptChoices(t *testing.T) {
	for _, tc := range []struct {
		tabs int
		want PermissionChoice
	}{{0, PermissionDeny}, {1, PermissionAllowAlways}, {2, PermissionAllowOnce}} {
		open := false
		var got []PermissionChoice
		tt := permissionTester(&open, &got)
		_ = tt.Click("Ask")
		if !tt.Focused("Deny") || !tt.HasText("Write access to ~/court/ledger") {
			t.Fatalf("Deny focused %v, texts %q", tt.Focused("Deny"), tt.Texts())
		}
		for range tc.tabs {
			tt.Key(0, ui.KeyTab)
		}
		tt.Key(0, ui.KeyEnter)
		if len(got) != 1 || got[0] != tc.want || open {
			t.Fatalf("tabs %d: got %v open %v", tc.tabs, got, open)
		}
	}
}

func TestPermissionPromptEscapeAndBackdrop(t *testing.T) {
	open := false
	var got []PermissionChoice
	tt := permissionTester(&open, &got)
	_ = tt.Click("Ask")
	tt.ClickAt(5, 395)
	if !open || len(got) != 0 {
		t.Fatal("backdrop acted")
	}
	tt.Key(0, ui.KeyEscape)
	if len(got) != 1 || got[0] != PermissionDeny {
		t.Fatalf("Escape gave %v", got)
	}
}

func TestPermissionPromptNeedsScope(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	open := true
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		PermissionPrompt(c, &open, PermissionPromptOptions{Title: "x"})
	}, 200, 200)
}

func TestPermissionPromptSnapshot(t *testing.T) {
	snapshot(t, "permission_prompt", 600, 340, func(c *ui.Context) {
		open := true
		PermissionPrompt(c, &open, PermissionPromptOptions{Title: "Write to the ledger?",
			Description: "The agent wants to record tonight's guests.", Scope: "Write access to ~/court/ledger and its subfolders"})
	})
}
