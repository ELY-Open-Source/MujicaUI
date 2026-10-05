package mujica

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

func promptTemplates() []PromptTemplate {
	return []PromptTemplate{
		{ID: "t1", Category: "Letters", Title: "Invitation", Body: "Write an invitation to {{guest}} for the {{event}}."},
		{ID: "t2", Category: "Music", Title: "Recital program", Body: "List pieces for a {{ instrument }} recital, then thank {{guest}}."},
		{ID: "t3", Category: "Letters", Title: "Thank-you note", Body: "Thank the hosts."},
	}
}

func TestPromptLibraryVariables(t *testing.T) {
	if got := PromptLibraryVariables("{{a}} and {{ b }} and {{a}}"); !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("variables %q", got)
	}
	if got := PromptLibraryFill("Hi {{ name }}, {{x}}!", map[string]string{"name": "Ann"}); got != "Hi Ann, !" {
		t.Fatalf("fill %q", got)
	}
}

func TestPromptLibraryBrowseFillInsert(t *testing.T) {
	var s PromptLibraryState
	templates := promptTemplates()
	inserted := ""
	tt := chatTester(940, 520, func(c *ui.Context) {
		if v := PromptLibrary(c, &s, templates, PromptLibraryOptions{Label: "Prompts"}); v.Inserted != "" {
			inserted = v.Inserted
		}
	})
	_ = tt.Click("Music")
	if tt.HasText("Invitation") || !tt.HasText("Recital program") {
		t.Fatalf("category filter: %q", tt.Texts())
	}
	_ = tt.Click("All")
	_ = tt.Click("Search prompts")
	tt.Type("thank")
	if tt.HasText("Invitation") || !tt.HasText("Thank-you note") || !tt.HasText("Recital program") {
		t.Fatalf("search covers titles and bodies: %q", tt.Texts())
	}
	_ = tt.Click("Recital program")
	if !tt.HasText("[instrument]") {
		t.Fatal("placeholder not previewed")
	}
	_ = tt.Click("Insert")
	if inserted != "" {
		t.Fatal("inserted with empty variables")
	}
	// Each variable's field sits under its label.
	for _, f := range [][2]string{{"instrument", "harp"}, {"guest", "the duke"}} {
		r, _ := tt.Find(f[0])
		tt.ClickAt(r.X+20, r.Y+r.H+18)
		tt.Type(f[1])
	}
	_ = tt.Click("Insert")
	if inserted != "List pieces for a harp recital, then thank the duke." {
		t.Fatalf("inserted %q", inserted)
	}
}

func TestPromptLibraryKeyboardAndDisabled(t *testing.T) {
	var s PromptLibraryState
	tt := chatTester(940, 520, func(c *ui.Context) { PromptLibrary(c, &s, promptTemplates(), PromptLibraryOptions{}) })
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if s.Category != "Letters" {
		t.Fatalf("keys chose category %q", s.Category)
	}
	var d PromptLibraryState
	tt = chatTester(940, 520, func(c *ui.Context) { PromptLibrary(c, &d, promptTemplates(), PromptLibraryOptions{Disabled: true}) })
	_ = tt.Click("Invitation")
	_ = tt.Click("Music")
	if d.Selected != "" || d.Category != "" {
		t.Fatalf("disabled library chose %q %q", d.Selected, d.Category)
	}
}

func TestPromptLibraryFiltersOnlyOnChange(t *testing.T) {
	var s PromptLibraryState
	templates := promptTemplates()
	tt := chatTester(940, 520, func(c *ui.Context) { PromptLibrary(c, &s, templates, PromptLibraryOptions{}) })
	s.shown.cats = append(s.shown.cats, "Marker")
	tt.Frame()
	if !slices.Contains(s.shown.cats, "Marker") {
		t.Fatal("refiltered without a change")
	}
	s.Query = "invit"
	tt.Frame()
	if !slices.Equal(s.shown.ids, []int{0}) {
		t.Fatalf("ids %v", s.shown.ids)
	}
}

func TestPromptLibrarySnapshot(t *testing.T) {
	snapshot(t, "prompt_library", 900, 360, func(c *ui.Context) {
		s := &PromptLibraryState{Selected: "t1", Values: map[string]string{"guest": "Lady Ann"}}
		ui.Column(c).Padding(16).Children(func() { PromptLibrary(c, s, promptTemplates(), PromptLibraryOptions{}) })
	})
}
