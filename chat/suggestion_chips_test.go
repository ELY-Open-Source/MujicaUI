package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var chipTexts = []string{"Plan a ball", "Write a sonnet", "Explain the minuet"}

func TestSuggestionChipsFill(t *testing.T) {
	input := "draft"
	var last SuggestionChipsView
	tt := chatTester(500, 200, func(c *ui.Context) {
		if v := SuggestionChips(c, &input, chipTexts, SuggestionChipsOptions{Label: "Suggestions"}); v.Chosen != "" {
			last = v
		}
	})
	_ = tt.Click("Write a sonnet")
	if input != "Write a sonnet" || last.Sent {
		t.Fatalf("input %q sent %v", input, last.Sent)
	}
}

func TestSuggestionChipsSendAndKeys(t *testing.T) {
	input := "draft"
	var last SuggestionChipsView
	tt := chatTester(500, 200, func(c *ui.Context) {
		if v := SuggestionChips(c, &input, chipTexts, SuggestionChipsOptions{Send: true}); v.Chosen != "" {
			last = v
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyEnter)
	if last.Chosen != "Explain the minuet" || !last.Sent || input != "draft" {
		t.Fatalf("chose %+v input %q", last, input)
	}
}

func TestSuggestionChipsDisabled(t *testing.T) {
	input := ""
	tt := chatTester(500, 200, func(c *ui.Context) {
		SuggestionChips(c, &input, chipTexts, SuggestionChipsOptions{Disabled: true})
	})
	_ = tt.Click("Plan a ball")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if input != "" {
		t.Fatalf("disabled chips filled %q", input)
	}
}

func TestSuggestionChipsRejectEmpty(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	in := ""
	ui.Render(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SuggestionChips(c, &in, []string{""}, SuggestionChipsOptions{})
	}, 100, 100, 1)
}

func TestSuggestionChipsSnapshot(t *testing.T) {
	snapshot(t, "suggestion_chips", 520, 120, func(c *ui.Context) {
		in := ""
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			SuggestionChips(c, &in, chipTexts, SuggestionChipsOptions{})
			SuggestionChips(c, &in, chipTexts[:2], SuggestionChipsOptions{Send: true, Disabled: true})
		})
	})
}
