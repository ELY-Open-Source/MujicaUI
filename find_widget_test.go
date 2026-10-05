package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func findTester(text *string, s *FindState, opts FindWidgetOptions, replaced *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		*replaced += FindWidget(c, text, s, opts).Replaced
	}, 640, 200)
}

func TestFindWidgetOptions(t *testing.T) {
	text := "Seal the seal; sealing wax. SEAL"
	s := &FindState{Query: "seal"}
	n := 0
	tt := findTester(&text, s, FindWidgetOptions{}, &n)
	if len(s.Matches()) != 4 || !tt.HasText("1 of 4") {
		t.Fatalf("plain matches %d", len(s.Matches()))
	}
	if err := tt.Click("Match case"); err != nil {
		t.Fatal(err)
	}
	if len(s.Matches()) != 2 {
		t.Fatalf("case-sensitive matches %d, want 2", len(s.Matches()))
	}
	_ = tt.Click("Whole word")
	if len(s.Matches()) != 1 || s.Matches()[0][0] != 9 {
		t.Fatalf("whole-word matches %v", s.Matches())
	}
	_ = tt.Click("Match case")
	_ = tt.Click("Whole word")
	_ = tt.Click("Regular expression")
	s.Query = `seal(ing)?`
	tt.Frame()
	if len(s.Matches()) != 4 || s.Matches()[2] != [2]int{15, 22} {
		t.Fatalf("regex matches %v", s.Matches())
	}
}

func TestFindWidgetInvalidRegex(t *testing.T) {
	text := "abc"
	s := &FindState{Query: "a(", Regex: true}
	n := 0
	tt := findTester(&text, s, FindWidgetOptions{}, &n)
	if len(s.Matches()) != 0 || !tt.HasText("Invalid pattern: error parsing regexp: missing closing ): `(?i)a(`") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestFindWidgetReplace(t *testing.T) {
	text := "one cat, two cat, red cat"
	s := &FindState{Query: "cat", Replacement: "owl", ShowReplace: true, Current: 1}
	n := 0
	tt := findTester(&text, s, FindWidgetOptions{}, &n)
	if err := tt.Click("Replace"); err != nil {
		t.Fatal(err)
	}
	if text != "one cat, two owl, red cat" || n != 1 {
		t.Fatalf("replace one gave %q (%d)", text, n)
	}
	if err := tt.Click("Replace all"); err != nil {
		t.Fatal(err)
	}
	if text != "one owl, two owl, red owl" || n != 3 || len(s.Matches()) != 0 {
		t.Fatalf("replace all gave %q (%d)", text, n)
	}
	text, s.Query, s.Regex, s.Replacement = "a1 b2", `([a-z])(\d)`, true, "$2$1"
	tt.Frame()
	_ = tt.Click("Replace all")
	if text != "1a 2b" {
		t.Fatalf("regex replace gave %q", text)
	}
}

func TestFindWidgetKeysAndReadOnly(t *testing.T) {
	text := "x y x y x"
	s := &FindState{}
	n := 0
	tt := findTester(&text, s, FindWidgetOptions{ReadOnly: true}, &n)
	if tt.HasText("Replace all") {
		t.Fatal("read-only widget offers replace")
	}
	tt.Key(0, ui.KeyTab)
	tt.Type("x")
	tt.Key(0, ui.KeyEnter)
	if s.Current != 1 || !tt.HasText("2 of 3") {
		t.Fatalf("Enter went to %d", s.Current)
	}
}

func TestFindWidgetDisabled(t *testing.T) {
	text := "cat"
	s := &FindState{Query: "cat", Replacement: "dog", ShowReplace: true}
	n := 0
	tt := findTester(&text, s, FindWidgetOptions{Disabled: true}, &n)
	_ = tt.Click("Replace all")
	_ = tt.Click("Match case")
	if text != "cat" || s.MatchCase {
		t.Fatalf("disabled widget changed %q case=%v", text, s.MatchCase)
	}
}

func TestFindWidgetSnapshot(t *testing.T) {
	snapshot(t, "find_widget", 640, 130, func(c *ui.Context) {
		st := ui.Local(c.Root(), "snap", func() *FindState {
			return &FindState{Query: "seal", Replacement: "signet", ShowReplace: true, WholeWord: true}
		})
		text := "the seal and the seal"
		ui.Column(c).Padding(12).Children(func() {
			FindWidget(c, &text, *st, FindWidgetOptions{})
		})
	})
}
