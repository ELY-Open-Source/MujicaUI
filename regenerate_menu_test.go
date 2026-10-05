package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func regenerateTester(opts RegenerateMenuOptions, got *[]RegenerateChoice) *ui.Tester {
	return chatTester(500, 420, func(c *ui.Context) {
		if v := RegenerateMenu(c, opts); v.Chosen {
			*got = append(*got, v.Choice)
		}
	})
}

func TestRegenerateMenuChoices(t *testing.T) {
	var got []RegenerateChoice
	opts := RegenerateMenuOptions{Models: []string{"Opus", "Sonnet"}, Current: "Opus"}
	tt := regenerateTester(opts, &got)
	_ = tt.Click("Regenerate")
	_ = tt.Click("Try again")
	_ = tt.Click("Regenerate")
	_ = tt.Click("Shorter")
	if len(got) != 2 || got[0].Kind != RegenerateAgain || got[1] != (RegenerateChoice{Kind: RegenerateVariant, Variant: "Shorter"}) {
		t.Fatalf("choices %+v", got)
	}
	_ = tt.Click("Regenerate")
	_ = tt.Click("Try another model")
	if tt.HasText("Opus") || !tt.HasText("Sonnet") {
		t.Fatalf("model submenu %q", tt.Texts())
	}
	_ = tt.Click("Sonnet")
	if len(got) != 3 || got[2] != (RegenerateChoice{Kind: RegenerateModel, Model: "Sonnet"}) {
		t.Fatalf("choices %+v", got)
	}
}

func TestRegenerateMenuKeyboardAndVariants(t *testing.T) {
	var got []RegenerateChoice
	tt := regenerateTester(RegenerateMenuOptions{Variants: []string{"In verse"}}, &got)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(got) != 1 || got[0].Kind != RegenerateAgain {
		t.Fatalf("keys chose %+v", got)
	}
	_ = tt.Click("Regenerate")
	if !tt.HasText("In verse") || tt.HasText("Longer") {
		t.Fatalf("variants %q", tt.Texts())
	}
}

func TestRegenerateMenuDisabled(t *testing.T) {
	var got []RegenerateChoice
	tt := regenerateTester(RegenerateMenuOptions{Disabled: true}, &got)
	_ = tt.Click("Regenerate")
	_ = tt.Click("Try again")
	if len(got) != 0 || tt.HasText("Try again") {
		t.Fatalf("disabled menu chose %+v", got)
	}
}

func TestRegenerateMenuSnapshot(t *testing.T) {
	snapshotAfter(t, "regenerate_menu", 360, 330, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			RegenerateMenu(c, RegenerateMenuOptions{Models: []string{"Opus", "Sonnet"}, Current: "Opus"})
		})
	}, func(tt *ui.Tester) { _ = tt.Click("Regenerate") })
}
