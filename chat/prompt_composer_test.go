package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
)

func promptTester(text *string, sends *int, opts PromptComposerOptions) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if PromptComposer(c, text, opts).Submitted() {
			*sends++
		}
	}, 500, 400)
}

func TestPromptComposerEnterSendsShiftEnterBreaks(t *testing.T) {
	text, sends := "", 0
	tt := promptTester(&text, &sends, PromptComposerOptions{Label: "Prompt"})
	_ = tt.Click("Prompt")
	tt.Key(0, ui.KeyEnter)
	if sends != 0 {
		t.Fatal("blank prompt sent")
	}
	tt.Type("Bonjour")
	tt.Key(ui.Shift, ui.KeyEnter)
	tt.Type("Madame")
	tt.Key(0, ui.KeyEnter)
	if sends != 1 || text != "Bonjour\nMadame" {
		t.Fatalf("sends %d text %q", sends, text)
	}
}

func TestPromptComposerCmdEnterMode(t *testing.T) {
	text, sends := "", 0
	tt := promptTester(&text, &sends, PromptComposerOptions{Label: "Prompt", CmdEnterSends: true})
	_ = tt.Click("Prompt")
	tt.Type("a")
	tt.Key(0, ui.KeyEnter)
	tt.Type("b")
	tt.Key(ui.Cmd, ui.KeyEnter)
	if sends != 1 || text != "a\nb" {
		t.Fatalf("sends %d text %q", sends, text)
	}
}

func TestPromptComposerSlotsAndDisabled(t *testing.T) {
	text, sends := "hello", 0
	tt := promptTester(&text, &sends, PromptComposerOptions{Label: "Prompt", Disabled: true})
	_ = tt.Click("Prompt")
	tt.Type("x")
	tt.Key(0, ui.KeyEnter)
	if sends != 0 || text != "hello" {
		t.Fatalf("disabled composer: sends %d text %q", sends, text)
	}
	tt = ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		PromptComposer(c, &text, PromptComposerOptions{
			Attachments: func() { ui.Text(c, "report.pdf") },
			Tools:       func() { ui.Text(c, "Web search") },
			Actions:     func() { ui.Text(c, "Send it") },
		})
	}, 500, 400)
	for _, s := range []string{"report.pdf", "Web search", "Send it", "Message"} {
		if _, ok := tt.Find(s); !ok {
			t.Errorf("missing %q", s)
		}
	}
}

func TestPromptComposerSnapshot(t *testing.T) {
	text := "Draft the invitation for the spring ball,\nin the style of the court."
	snapshot(t, "prompt_composer", 520, 220, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			PromptComposer(c, &text, PromptComposerOptions{
				Attachments: func() { AttachmentChip(c, Attachment{Name: "guests.csv", Size: 4096}) },
				Tools:       func() { core.IconAction(c, icons.Must("paperclip"), "Attach") },
				Actions:     func() { SendButton(c, SendButtonOptions{Shortcut: "Enter"}) },
			})
		})
	})
}

func TestPromptComposerCompositionDoesNotSend(t *testing.T) {
	text, sends := "draft", 0
	composing := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if PromptComposer(c, &text, PromptComposerOptions{Label: "Prompt", Composing: composing}).Submitted() {
			sends++
		}
	}, 500, 400)
	_ = tt.Click("Prompt")
	composing = true
	tt.Frame()
	tt.Compose("邀请", 2)
	tt.Key(0, ui.KeyEnter)
	if sends != 0 || text != "draft" {
		t.Fatalf("composition submitted: sends %d text %q", sends, text)
	}
	tt.Type("邀请")
	composing = false
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if sends != 1 || text != "draft邀请" {
		t.Fatalf("committed input: sends %d text %q", sends, text)
	}
}
