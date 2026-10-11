package git

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

func TestBranchSelectorSearchAndChoose(t *testing.T) {
	sel, changes := "main", 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if BranchSelector(c, &sel, branchSample(), BranchSelectorOptions{}).Changed() {
			changes++
		}
	}, 600, 500)
	if err := tt.Click("Branch: main"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Local") || !tt.HasText("Remote") || !tt.HasText("origin/release") {
		t.Fatalf("groups: %q", tt.Texts())
	}
	tt.Type("rel")
	if tt.HasText("fix/blame") || !tt.HasText("origin/release") {
		t.Fatalf("filter: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if sel != "origin/release" || changes != 1 || tt.HasText("Remote") {
		t.Fatalf("Enter chose %q changes %d", sel, changes)
	}
	_ = tt.Click("Branch: origin/release")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if sel != "fix/blame" {
		t.Fatalf("arrows chose %q", sel)
	}
	_ = tt.Click("Branch: fix/blame")
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Local") {
		t.Fatal("Escape did not close")
	}
}

func TestBranchSelectorCreateAndDisabled(t *testing.T) {
	sel, created, disabled := "main", "", false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := BranchSelector(c, &sel, branchSample(), BranchSelectorOptions{AllowCreate: true, Disabled: disabled})
		if name, ok := r.Created(); ok {
			created = name
		}
	}, 600, 500)
	_ = tt.Click("Branch: main")
	tt.Type("feature/new")
	if !tt.HasText("Create branch “feature/new”") {
		t.Fatalf("create row: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if created != "feature/new" || sel != "main" {
		t.Fatalf("created %q sel %q", created, sel)
	}
	_ = tt.Click("Branch: main")
	disabled = true
	tt.Frame()
	if tt.HasText("Local") {
		t.Fatal("popup stayed open on a disabled selector")
	}
	_ = tt.Click("Branch: main")
	if tt.HasText("Local") {
		t.Fatal("disabled selector opened")
	}
}

func TestBranchSelectorSnapshot(t *testing.T) {
	snapshotAfter(t, "branch_selector", 420, 420, func(c *ui.Context) {
		sel := ui.Local(c.Root(), "sel", func() string { return "main" })
		ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() {
			BranchSelector(c, sel, branchSample(), BranchSelectorOptions{AllowCreate: true})
		})
	}, func(tt *ui.Tester) {
		_ = tt.Click("Branch: main")
		tt.Type("ma")
	})
}
