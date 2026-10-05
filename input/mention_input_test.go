package input

import (
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

var mentionPeople = []MentionCandidate{{ID: "u1", Name: "Alice"}, {ID: "u2", Name: "Albert"}, {ID: "u3", Name: "Bob"}}

func mentionFind(q string) []MentionCandidate {
	var out []MentionCandidate
	for _, p := range mentionPeople {
		if strings.HasPrefix(strings.ToLower(p.Name), strings.ToLower(q)) {
			out = append(out, p)
		}
	}
	return out
}

func TestMentionDocWholeTokens(t *testing.T) {
	d := mentionDoc{text: []rune("hi @al"), caret: 6, anchor: 6}
	at, q := d.query()
	if at != 3 || q != "al" {
		t.Fatalf("query %d %q", at, q)
	}
	d.choose(at, "u1", "Alice")
	d.insert("there")
	if string(d.text) != "hi @Alice there" || d.mentions[0] != (Mention{ID: "u1", Name: "Alice", Start: 3, End: 9}) {
		t.Fatalf("after choose: %q %+v", string(d.text), d.mentions)
	}
	d.moveTo(10, false)
	d.move(-1, false) // to the mention's end
	d.move(-1, false) // over the whole mention
	if d.caret != 3 {
		t.Fatalf("Left over a mention lands at %d", d.caret)
	}
	d.moveTo(9, false)
	d.erase(false)
	if string(d.text) != "hi  there" || len(d.mentions) != 0 || d.caret != 3 {
		t.Fatalf("Backspace after a mention: %q %+v %d", string(d.text), d.mentions, d.caret)
	}
	if at, _ := (&mentionDoc{text: []rune("a@b"), caret: 3, anchor: 3}).query(); at != -1 {
		t.Fatal("an email-like @ started a query")
	}
}

func mentionTester(v *string, ms *[]Mention, dis *bool, n *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(400).Children(func() {
			if MentionInput(c, v, ms, MentionInputOptions{Candidates: mentionFind, Label: "Comment", Disabled: *dis}).Changed() {
				*n++
			}
		})
	}, 500, 400)
}

func TestMentionInputCandidatesKeyboard(t *testing.T) {
	v, ms := "", []Mention(nil)
	dis, n := false, 0
	tt := mentionTester(&v, &ms, &dis, &n)
	tt.Key(0, ui.KeyTab)
	tt.Type("ping @al")
	if !tt.HasText("@Alice") || !tt.HasText("@Albert") || tt.HasText("@Bob") {
		t.Fatalf("candidates %q", tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if v != "ping @Albert " || len(ms) != 1 || ms[0].ID != "u2" || ms[0].Start != 5 || ms[0].End != 12 {
		t.Fatalf("chose: %q %+v", v, ms)
	}
	if tt.HasText("@Alice") {
		t.Fatal("popup stayed open")
	}
	tt.Type("ok")
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyBackspace)
	tt.Key(0, ui.KeyBackspace) // the space
	tt.Key(0, ui.KeyBackspace) // the whole mention
	if v != "ping " || len(ms) != 0 {
		t.Fatalf("whole delete: %q %+v", v, ms)
	}
	tt.Command("undo")
	if v != "ping @Albert" || len(ms) != 1 {
		t.Fatalf("undo: %q %+v", v, ms)
	}
	tt.Type(" @b")
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("@Bob") {
		t.Fatal("Escape did not close the candidates")
	}
}

func TestMentionInputClickCompositionDisabled(t *testing.T) {
	v, ms := "", []Mention(nil)
	dis, n := false, 0
	tt := mentionTester(&v, &ms, &dis, &n)
	tt.Key(0, ui.KeyTab)
	tt.Type("@b")
	if err := tt.Click("Bob"); err != nil {
		t.Fatal(err)
	}
	if v != "@Bob " || !slices.Equal(ms, []Mention{{ID: "u3", Name: "Bob", Start: 0, End: 4}}) {
		t.Fatalf("click: %q %+v", v, ms)
	}
	tt.Compose("你", 1)
	if v != "@Bob " {
		t.Fatalf("composition committed early: %q", v)
	}
	tt.Type("你好")
	if v != "@Bob 你好" {
		t.Fatalf("commit: %q", v)
	}
	before := n
	dis = true
	tt.Frame()
	tt.Type("x")
	tt.Key(0, ui.KeyBackspace)
	if v != "@Bob 你好" || n != before {
		t.Fatalf("disabled edit: %q", v)
	}
}

func TestMentionInputRejectsBadMentions(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	v, ms := "hi Bob", []Mention{{ID: "u3", Name: "Bob", Start: 3, End: 6}}
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		MentionInput(c, &v, &ms, MentionInputOptions{Candidates: mentionFind})
	}, 300, 100)
}

func TestMentionInputSnapshot(t *testing.T) {
	v, ms := "thanks @Alice, see @Bob", []Mention{{ID: "u1", Name: "Alice", Start: 7, End: 13}, {ID: "u3", Name: "Bob", Start: 19, End: 23}}
	e, none := "", []Mention(nil)
	snapshot(t, "mention_input", 420, 120, func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			MentionInput(c, &v, &ms, MentionInputOptions{Candidates: mentionFind})
			MentionInput(c, &e, &none, MentionInputOptions{Candidates: mentionFind, Placeholder: "Write a comment…"})
		})
	})
}
