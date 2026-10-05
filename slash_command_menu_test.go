package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

var slashCommandsSample = []SlashCommand{
	{Name: "summarize", Description: "Summarize the thread"},
	{Name: "search", Description: "Search the web"},
	{Name: "translate", Description: "Translate to French"},
}

// slashTester is a composer with a slash menu; it records chosen commands and sends.
type slashTester struct {
	text   string
	chosen []string
	sends  int
	off    bool
}

func (s *slashTester) view(c *ui.Context) {
	Use(c, Settings{})
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
