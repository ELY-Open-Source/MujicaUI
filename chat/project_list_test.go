package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

type projectFixture struct {
	s        ProjectListState
	projects []ChatProject
	owner    map[string]string
	last     ProjectListView
	moves    int
}

func newProjectFixture() *projectFixture {
	return &projectFixture{
		projects: []ChatProject{{ID: "p1", Name: "Coronation", Folder: "Court"}, {ID: "p2", Name: "Gardens", Folder: "Court"}, {ID: "p3", Name: "Music"}},
		owner:    map[string]string{"a": "p1", "b": "p1", "c": "p3"},
	}
}

func (f *projectFixture) view(opts ProjectListOptions) func(c *ui.Context) {
	return func(c *ui.Context) {
		convs := func(id string) []ChatConversation {
			var out []ChatConversation
			for _, conv := range chatConversations() {
				if f.owner[conv.ID] == id {
					out = append(out, conv)
				}
			}
			return out
		}
		v := ProjectList(c, &f.s, f.projects, convs, opts)
		if v.Moved {
			f.owner[v.MoveID] = v.MoveTo
			f.moves++
		}
		if v.Changed || v.Opened != "" || v.Moved {
			f.last = v
		}
	}
}

func TestProjectListSelectAndExpand(t *testing.T) {
	f := newProjectFixture()
	tt := chatTester(420, 500, f.view(ProjectListOptions{Label: "Projects"}))
	if err := tt.Click("Gardens"); err != nil {
		t.Fatal(err)
	}
	if f.s.Selected != "p2" || !f.last.Changed {
		t.Fatalf("selected %q", f.s.Selected)
	}
	if tt.HasText("Court etiquette notes") {
		t.Fatal("conversations shown before expanding")
	}
	_ = tt.Click("Show conversations in Coronation")
	if !tt.HasText("Court etiquette notes") || !tt.HasText("Ballroom seating plan") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Court etiquette notes")
	if f.last.Opened != "a" {
		t.Fatalf("opened %q", f.last.Opened)
	}
	// The folder collapses its projects.
	_ = tt.Click("Court")
	if tt.HasText("Gardens") {
		t.Fatal("folder did not collapse")
	}
}

func TestProjectListKeyboard(t *testing.T) {
	f := newProjectFixture()
	tt := chatTester(420, 500, f.view(ProjectListOptions{}))
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Court") {
		t.Fatal("Tab did not enter at the folder")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("Court etiquette notes") {
		t.Fatal("Right did not expand")
	}
	tt.Key(0, ui.KeyEnter)
	if f.s.Selected != "p1" {
		t.Fatalf("Enter chose %q", f.s.Selected)
	}
	tt.Key(0, ui.KeyLeft)
	if tt.HasText("Court etiquette notes") {
		t.Fatal("Left did not collapse")
	}
}

func TestProjectListMoveByMenuAndDrag(t *testing.T) {
	f := newProjectFixture()
	tt := chatTester(420, 500, f.view(ProjectListOptions{}))
	_ = tt.Click("Show conversations in Coronation")
	if err := tt.Click("Move Court etiquette notes"); err != nil {
		t.Fatal(err)
	}
	// The menu opens on Gardens; Down reaches Music.
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if f.owner["a"] != "p3" || f.moves != 1 {
		t.Fatalf("menu move: owner %q moves %d", f.owner["a"], f.moves)
	}
	from, _ := tt.Find("Ballroom seating plan")
	to, _ := tt.Find("Gardens")
	tt.Press(from.X+5, from.Y+from.H/2)
	tt.Move(from.X+20, from.Y+from.H/2+5)
	tt.Move(to.X+10, to.Y+to.H/2)
	tt.Release(to.X+10, to.Y+to.H/2)
	if f.owner["b"] != "p2" || f.last.MoveTo != "p2" {
		t.Fatalf("drag move: owner %q", f.owner["b"])
	}
}

func TestProjectListDisabled(t *testing.T) {
	f := newProjectFixture()
	tt := chatTester(420, 500, f.view(ProjectListOptions{Disabled: true}))
	_ = tt.Click("Gardens")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if f.s.Selected != "" || f.moves != 0 {
		t.Fatalf("disabled list acted: %q", f.s.Selected)
	}
}

func TestProjectListSnapshot(t *testing.T) {
	snapshot(t, "project_list", 380, 280, func(c *ui.Context) {
		f := newProjectFixture()
		f.s.Selected = "p1"
		f.s.open = map[string]bool{"p1": true}
		f.s.closed = map[string]bool{}
		ui.Column(c).Padding(16).Children(func() { f.view(ProjectListOptions{})(c) })
	})
}
