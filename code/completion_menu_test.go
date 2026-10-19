package code

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/input"
)

func completionSample() []Completion {
	return []Completion{
		{Label: "Seal", Kind: SymbolMethod, Detail: "func(d Decree) error", Doc: "Seal stamps a decree."},
		{Label: "Search", Kind: SymbolFunction, Detail: "func(q string) []Decree"},
		{Label: "Season", Kind: SymbolConstant, Detail: "= 4"},
	}
}

type completionProbe struct {
	text     string
	open     bool
	s        CompletionMenuState
	disabled bool
	accepted []string
}

func (p *completionProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Column(c).Padding(16).Children(func() {
		in := input.TextInput(c, &p.text, input.TextInputOptions{Label: "Code"})
		if in.Changed() {
			p.open = strings.HasPrefix(strings.ToLower(p.text), "se")
		}
		r := CompletionMenu(c, in.Input, &p.open, completionSample(), &p.s, CompletionMenuOptions{Disabled: p.disabled})
		if r.Accepted != nil {
			p.accepted = append(p.accepted, r.Accepted.Label)
			p.text = r.Accepted.Label
		}
	})
}

func TestCompletionMenuKeys(t *testing.T) {
	p := &completionProbe{}
	tt := ui.NewTester(p.view, 480, 400)
	tt.Key(0, ui.KeyTab)
	tt.Type("se")
	if !p.open || !tt.HasText("Seal") || !tt.HasText("Seal stamps a decree.") {
		t.Fatalf("menu did not open: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	if p.s.Selected != 0 {
		t.Fatalf("Down wrapped to %d, want 0", p.s.Selected)
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyTab)
	if len(p.accepted) != 1 || p.accepted[0] != "Season" || p.open {
		t.Fatalf("accepted %v open %v", p.accepted, p.open)
	}
	if !tt.Focused("Code") {
		t.Fatal("accepting moved the focus away from the editor")
	}
}

func TestCompletionMenuEnterClickEscape(t *testing.T) {
	p := &completionProbe{}
	tt := ui.NewTester(p.view, 480, 400)
	tt.Key(0, ui.KeyTab)
	tt.Type("se")
	tt.Key(0, ui.KeyEnter)
	if len(p.accepted) != 1 || p.accepted[0] != "Seal" {
		t.Fatalf("Enter accepted %v", p.accepted)
	}
	p.open = true
	tt.Frame()
	if err := tt.Click("Search"); err != nil {
		t.Fatal(err)
	}
	if len(p.accepted) != 2 || p.accepted[1] != "Search" {
		t.Fatalf("click accepted %v", p.accepted)
	}
	p.open = true
	tt.Frame()
	tt.Key(0, ui.KeyEscape)
	if p.open || len(p.accepted) != 2 {
		t.Fatalf("Escape: open %v accepted %v", p.open, p.accepted)
	}
}

func TestCompletionMenuDisabledClosed(t *testing.T) {
	p := &completionProbe{open: true, disabled: true}
	tt := ui.NewTester(p.view, 480, 400)
	if p.open || tt.HasText("Search") {
		t.Fatal("disabled menu shows")
	}
}

func TestCompletionMenuSnapshot(t *testing.T) {
	snapshot(t, "completion_menu", 420, 260, func(c *ui.Context) {
		p := ui.Local(c.Root(), "snap", func() *completionProbe { return &completionProbe{text: "se", open: true} })
		(*p).view(c)
	})
}
