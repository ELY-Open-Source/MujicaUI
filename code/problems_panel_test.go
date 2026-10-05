package code

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func problemsSample() []Problem {
	return []Problem{
		{ProblemWarning, "court/seal.go", 30, 2, "unused variable wax", "vet"},
		{ProblemError, "court/seal.go", 12, 5, "undefined: Decree", "compiler"},
		{ProblemInfo, "docs/ledger.md", 4, 1, "line is long", "lint"},
	}
}

type problemsProbe struct {
	problems []Problem
	s        ProblemsPanelState
	disabled bool
	opened   []Problem
}

func (p *problemsProbe) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := ProblemsPanel(c, p.problems, &p.s, ProblemsPanelOptions{Disabled: p.disabled})
	r.Element.Height(360)
	if r.Opened != nil {
		p.opened = append(p.opened, *r.Opened)
	}
}

func TestProblemsPanelGroupsSortedByLine(t *testing.T) {
	p := &problemsProbe{problems: problemsSample()}
	tt := ui.NewTester(p.view, 640, 380)
	if !tt.HasText("Errors 1") || !tt.HasText("Warnings 1") || !tt.HasText("Info 1") {
		t.Fatalf("counts missing: %q", tt.Texts())
	}
	a, _ := tt.Find("undefined: Decree")
	b, _ := tt.Find("unused variable wax")
	if a.Y >= b.Y {
		t.Fatal("problems are not sorted by line within the file")
	}
	if len(p.s.flat) != 5 || p.s.flat[0].file != "court/seal.go" || p.s.flat[3].file != "docs/ledger.md" {
		t.Fatalf("rows %v", p.s.flat)
	}
}

func TestProblemsPanelFilterAndOpen(t *testing.T) {
	p := &problemsProbe{problems: problemsSample()}
	tt := ui.NewTester(p.view, 640, 380)
	if err := tt.Click("Warnings 1"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("unused variable wax") || !tt.HasText("undefined: Decree") {
		t.Fatalf("filter kept warnings: %q", tt.Texts())
	}
	r, _ := tt.Find("undefined: Decree")
	tt.ClickAt(r.X+10, r.Y+r.H/2)
	tt.ClickAt(r.X+10, r.Y+r.H/2)
	if len(p.opened) != 1 || p.opened[0].Line != 12 || p.opened[0].Col != 5 {
		t.Fatalf("opened %v", p.opened)
	}
}

func TestProblemsPanelKeyboardFold(t *testing.T) {
	p := &problemsProbe{problems: problemsSample()}
	tt := ui.NewTester(p.view, 640, 380)
	for range 4 {
		tt.Key(0, ui.KeyTab)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyLeft)
	if !p.s.Collapsed.Has("court/seal.go") || tt.HasText("undefined: Decree") {
		t.Fatal("Left did not fold the file")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(p.opened) != 1 || p.opened[0].Message != "undefined: Decree" {
		t.Fatalf("Enter opened %v", p.opened)
	}
}

func TestProblemsPanelDisabledAndEmpty(t *testing.T) {
	p := &problemsProbe{problems: problemsSample(), disabled: true}
	tt := ui.NewTester(p.view, 640, 380)
	_ = tt.Click("Errors 1")
	_ = tt.Click("Fold court/seal.go")
	if p.s.Hidden[ProblemError] || p.s.Collapsed.Len() != 0 {
		t.Fatal("disabled panel changed")
	}
	tt = ui.NewTester((&problemsProbe{}).view, 640, 380)
	if !tt.HasText("No problems") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestProblemsPanelSnapshot(t *testing.T) {
	snapshot(t, "problems_panel", 600, 240, func(c *ui.Context) {
		s := ui.Local(c.Root(), "snap", func() *ProblemsPanelState { return &ProblemsPanelState{} })
		ui.Column(c).Padding(12).Children(func() {
			ProblemsPanel(c, problemsSample(), *s, ProblemsPanelOptions{}).Element.Height(216)
		})
	})
}
