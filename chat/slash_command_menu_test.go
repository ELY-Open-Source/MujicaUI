package chat

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var slashCommandsSample = []SlashCommand{
	{Name: "summarize", Description: "Summarize the thread"},
	{Name: "search", Description: "Search the web"},
	{Name: "translate", Description: "Translate to French"},
}

// slashTester records command choices and prompt submissions.
type slashTester struct {
	text   string
	chosen []string
	sends  int
	off    bool
}

func (s *slashTester) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Column(c).Fill().Justify(ui.End).Padding(16).Children(func() {
		p := PromptComposer(c, &s.text, PromptComposerOptions{Label: "Prompt", Disabled: s.off})
		if p.Submitted() {
			s.sends++
		}
		if cmd, ok := SlashCommandMenu(c, p, &s.text, slashCommandsSample).Chosen(); ok {
			s.chosen = append(s.chosen, cmd.Name)
		}
	})
}

func TestSlashCommandMenuFilterAndKeys(t *testing.T) {
	s := &slashTester{}
	tt := ui.NewTester(s.view, 500, 400)
	_ = tt.Click("Prompt")
	tt.Type("/s")
	sr, _ := tt.Find("Search the web")
	tr, ok := tt.Find("Translate to French")
	if !ok || !tt.HasText("/summarize") || tr.Y <= sr.Y {
		t.Fatalf("filter /s does not list prefix matches first: %q", tt.Texts())
	}
	tt.Type("e")
	if tt.HasText("/summarize") || !tt.HasText("/search") {
		t.Fatalf("filter /se: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if s.text != "/search " || len(s.chosen) != 1 || s.chosen[0] != "search" || s.sends != 0 {
		t.Fatalf("Enter: text %q chosen %v sends %d", s.text, s.chosen, s.sends)
	}
	if tt.HasText("/summarize") {
		t.Fatal("menu stays open after choosing")
	}
}

func TestSlashCommandMenuTabCompletesAndEscapeCloses(t *testing.T) {
	s := &slashTester{}
	tt := ui.NewTester(s.view, 500, 400)
	_ = tt.Click("Prompt")
	tt.Type("/tr")
	tt.Key(0, ui.KeyTab)
	if s.text != "/translate " || len(s.chosen) != 0 {
		t.Fatalf("Tab: text %q chosen %v", s.text, s.chosen)
	}
	s.text = ""
	tt.Frame()
	tt.Type("/")
	if !tt.HasText("/summarize") {
		t.Fatal("menu did not open on /")
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("/summarize") {
		t.Fatal("Escape did not close the menu")
	}
	tt.Type("t")
	if !tt.HasText("/translate") {
		t.Fatal("menu did not reopen as the filter changed")
	}
}

func TestSlashCommandMenuClickAndEmpty(t *testing.T) {
	s := &slashTester{}
	tt := ui.NewTester(s.view, 500, 400)
	_ = tt.Click("Prompt")
	tt.Type("/zz")
	if !tt.HasText("No matching commands") {
		t.Fatalf("empty state: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyBackspace)
	if err := tt.Click("/translate"); err != nil {
		t.Fatal(err)
	}
	if len(s.chosen) != 1 || s.chosen[0] != "translate" || !tt.Focused("Prompt") {
		t.Fatalf("click chose %v, composer focused %v", s.chosen, tt.Focused("Prompt"))
	}
}

func TestSlashCommandMenuDisabledComposer(t *testing.T) {
	s := &slashTester{text: "/s", off: true}
	tt := ui.NewTester(s.view, 500, 400)
	_ = tt.Click("Prompt")
	tt.Key(0, ui.KeyEnter)
	if tt.HasText("/search") || len(s.chosen) != 0 {
		t.Fatal("disabled composer shows or uses the menu")
	}
}

func TestSlashCommandMenuSnapshot(t *testing.T) {
	snapshotAfter(t, "slash_command_menu", 520, 360, func(c *ui.Context) {
		(*ui.Local(c.Root(), "s", func() *slashTester { return &slashTester{} })).view(c)
	}, func(tt *ui.Tester) {
		_ = tt.Click("Prompt")
		tt.Type("/")
		tt.Key(0, ui.KeyDown)
	})
}

func TestSlashCommandMenuNarrowAndLargeData(t *testing.T) {
	cmds := make([]SlashCommand, 100000)
	for i := range cmds {
		cmds[i] = SlashCommand{Name: fmt.Sprintf("cmd%d", i)}
	}
	text := "/"
	var panel *ui.Element
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		p := PromptComposer(c, &text, PromptComposerOptions{Label: "Prompt"})
		panel = SlashCommandMenu(c, p, &text, cmds).Element
	}, 180, 120)
	_ = tt.Click("Prompt")
	tt.Frame()
	if panel == nil {
		t.Fatal("menu missing")
	}
	b := panel.Bounds()
	if b.X < 0 || b.Y < 0 || b.X+b.W > 180 || b.Y+b.H > 120 {
		t.Errorf("popup outside window: %v", b)
	}
	if n := len(tt.Texts()); n > 150 {
		t.Errorf("built %d labels for 100000 commands, want <= 150", n)
	}
}

func TestSlashCommandMenuReopensAfterTriggerRemoved(t *testing.T) {
	s := &slashTester{}
	tt := ui.NewTester(s.view, 500, 400)
	_ = tt.Click("Prompt")
	tt.Type("/")
	tt.Key(0, ui.KeyEscape)
	tt.Key(0, ui.KeyBackspace)
	tt.Type("/")
	if !tt.HasText("/summarize") {
		t.Fatal("retyped slash stayed dismissed")
	}
}

func TestSlashCommandMenuQueryRevealsHighlight(t *testing.T) {
	commands := make([]SlashCommand, 100)
	for i := range commands {
		commands[i] = SlashCommand{Name: fmt.Sprintf("cmd%03d", i)}
	}
	text := "/"
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Fill().Justify(ui.End).Children(func() {
			p := PromptComposer(c, &text, PromptComposerOptions{Label: "Prompt"})
			SlashCommandMenu(c, p, &text, commands)
		})
	}, 500, 400)
	_ = tt.Click("Prompt")
	tt.Scroll(100, 150, 0, 1500)
	tt.Type("c")
	b, ok := tt.Find("/cmd000")
	if !ok || b.H == 0 || b.Y < 0 || b.Y+b.H > 400 {
		t.Fatalf("highlight outside view: %v found %v", b, ok)
	}
}
