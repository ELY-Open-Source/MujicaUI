package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestWelcomeScreenGreetsAndChooses(t *testing.T) {
	input, got := "", ""
	tt := chatTester(900, 600, func(c *ui.Context) {
		v := WelcomeScreen(c, WelcomeScreenOptions{Greeting: "Good evening", Subtitle: "How may I serve?", Cards: chatCards(),
			Prompts: func() { SuggestionChips(c, &input, []string{"Plan a ball"}, SuggestionChipsOptions{}) }})
		if v.Example != "" {
			got = v.Example
		}
	})
	for _, s := range []string{"Good evening", "How may I serve?", "Research"} {
		if !tt.HasText(s) {
			t.Fatalf("%q missing from %q", s, tt.Texts())
		}
	}
	_ = tt.Click("“Explain the minuet”")
	if got != "Explain the minuet" {
		t.Fatalf("example %q", got)
	}
	_ = tt.Click("Plan a ball")
	if input != "Plan a ball" {
		t.Fatalf("slot prompt left %q", input)
	}
}

func TestWelcomeScreenDisabled(t *testing.T) {
	got := ""
	tt := chatTester(900, 600, func(c *ui.Context) {
		if v := WelcomeScreen(c, WelcomeScreenOptions{Greeting: "Hello", Cards: chatCards(), Disabled: true}); v.Example != "" {
			got = v.Example
		}
	})
	_ = tt.Click("“Explain the minuet”")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if got != "" {
		t.Fatalf("disabled screen chose %q", got)
	}
}

func TestWelcomeScreenNeedsGreeting(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		WelcomeScreen(c, WelcomeScreenOptions{})
	}, 100, 100, 1)
}

func TestWelcomeScreenSnapshot(t *testing.T) {
	snapshot(t, "welcome_screen", 900, 420, func(c *ui.Context) {
		in := ""
		WelcomeScreen(c, WelcomeScreenOptions{Greeting: "Good afternoon", Subtitle: "How may I serve the court?", Cards: chatCards(),
			Prompts: func() { SuggestionChips(c, &in, []string{"Plan a ball", "Write a sonnet"}, SuggestionChipsOptions{}) }})
	})
}
