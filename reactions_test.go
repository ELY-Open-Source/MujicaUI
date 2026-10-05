package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

type reactionsProbe struct {
	items    []Reaction
	disabled bool
	toggled  []string
	added    []string
}

var reactionChoices = []string{"👍", "❤️", "🎉", "👀", "🙏", "😂", "🔥"}

func (p *reactionsProbe) view(c *ui.Context) {
	Use(c, Settings{})
	ui.Column(c).Padding(12).Children(func() {
		r := Reactions(c, p.items, ReactionsOptions{Choices: reactionChoices, Disabled: p.disabled})
		if r.Toggled != "" {
			p.toggled = append(p.toggled, r.Toggled)
		}
		if r.Added != "" {
			p.added = append(p.added, r.Added)
		}
	})
}

func newReactionsProbe() *reactionsProbe {
	return &reactionsProbe{items: []Reaction{{Emoji: "👍", Count: 3, Mine: true, People: []string{"Ada", "Blaise", "you"}}, {Emoji: "🎉", Count: 1}}}
}

func TestReactionsToggleAndAdd(t *testing.T) {
	p := newReactionsProbe()
	tt := ui.NewTester(p.view, 400, 300)
	if err := tt.Click("👍 3, including you"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("🎉 1"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Add reaction"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("🔥"); err != nil {
		t.Fatal(err)
	}
	if len(p.toggled) != 2 || p.toggled[0] != "👍" || len(p.added) != 1 || p.added[0] != "🔥" {
		t.Fatalf("toggled %v added %v", p.toggled, p.added)
	}
	if _, ok := tt.Find("❤️"); ok {
		t.Fatal("picker stayed open")
	}
}

func TestReactionsKeyboard(t *testing.T) {
	p := newReactionsProbe()
	tt := ui.NewTester(p.view, 400, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if !tt.HasText("❤️") {
		t.Fatal("Enter on add did not open the picker")
	}
	if len(p.toggled) != 1 {
		t.Fatalf("Space toggled %v", p.toggled)
	}
}

func TestReactionsDisabledAndInvalid(t *testing.T) {
	p := newReactionsProbe()
	p.disabled = true
	tt := ui.NewTester(p.view, 400, 300)
	_ = tt.Click("🎉 1")
	_ = tt.Click("Add reaction")
	if len(p.toggled) != 0 || tt.HasText("❤️") {
		t.Fatal("disabled reactions acted")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("zero count did not panic")
		}
	}()
	ui.NewTester((&reactionsProbe{items: []Reaction{{Emoji: "x"}}}).view, 400, 300)
}

func TestReactionsSnapshot(t *testing.T) {
	snapshot(t, "reactions", 300, 70, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() { Reactions(c, newReactionsProbe().items, ReactionsOptions{Choices: reactionChoices}) })
	})
}
