package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

var switcherSpaces = []Workspace{{ID: "ave", Name: "Ave Mujica", Detail: "5 members"}, {ID: "mygo", Name: "MyGO!!!!!", Detail: "5 members"}, {ID: "crychic", Name: "CRYCHIC", Detail: "Archived"}}

func TestWorkspaceSwitcherSearchAndSwitch(t *testing.T) {
	cur, n := "ave", 0
	tt := accountTester(500, 500, func(c *ui.Context) {
		if WorkspaceSwitcher(c, &cur, switcherSpaces, WorkspaceSwitcherOptions{}).Changed() {
			n++
		}
	})
	if err := tt.Click("Ave Mujica"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Find a workspace") || !tt.HasText("CRYCHIC") {
		t.Fatalf("open: %q", tt.Texts())
	}
	tt.Type("arch")
	if tt.HasText("MyGO!!!!!") || !tt.HasText("CRYCHIC") {
		t.Fatalf("search: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEnter)
	if cur != "crychic" || n != 1 || tt.HasText("Find a workspace") {
		t.Fatalf("enter: %q", cur)
	}
	_ = tt.Click("CRYCHIC")
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if cur != "mygo" || n != 2 {
		t.Fatalf("arrow: %q", cur)
	}
	_ = tt.Click("MyGO!!!!!")
	tt.Type("zzz")
	if !tt.HasText("No workspace matches") {
		t.Fatal("no empty state")
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("No workspace matches") || cur != "mygo" {
		t.Fatal("escape did not close")
	}
}

func TestWorkspaceSwitcherSnapshot(t *testing.T) {
	snapshot(t, "workspace_switcher", 300, 90, func(c *ui.Context) {
		cur := ui.Local(c.Root(), "cur", func() string { return "ave" })
		ui.Column(c).Padding(16).AlignItems(ui.Start).Children(func() { WorkspaceSwitcher(c, cur, switcherSpaces, WorkspaceSwitcherOptions{}) })
	})
}

func TestWorkspaceSwitcherLargeDataBuildsVisibleRows(t *testing.T) {
	items := make([]Workspace, 100000)
	for i := range items {
		items[i] = Workspace{ID: fmt.Sprint(i), Name: fmt.Sprintf("Workspace %d", i)}
	}
	cur := "0"
	tt := accountTester(500, 500, func(c *ui.Context) { WorkspaceSwitcher(c, &cur, items, WorkspaceSwitcherOptions{}) })
	if err := tt.Click("Workspace 0"); err != nil {
		t.Fatal(err)
	}
	if n := len(tt.Texts()); n > 100 {
		t.Fatalf("built %d text elements for a bounded viewport", n)
	}
	tt.Type("99999")
	tt.Key(0, ui.KeyEnter)
	if cur != "99999" {
		t.Fatalf("large-list search choice %q", cur)
	}
}

func TestWorkspaceSwitcherDisablingClosesAndInvalidCurrent(t *testing.T) {
	cur, disabled := "ave", false
	tt := accountTester(600, 500, func(c *ui.Context) {
		ui.Column(c).Disabled(disabled).Children(func() { WorkspaceSwitcher(c, &cur, switcherSpaces, WorkspaceSwitcherOptions{}) })
	})
	if err := tt.Click("Ave Mujica"); err != nil {
		t.Fatal(err)
	}
	disabled = true
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if tt.HasText("Find a workspace") || cur != "ave" {
		t.Fatal("disabled workspace popup stays active")
	}
	cur = "unknown"
	accountPanics(t, func(c *ui.Context) { WorkspaceSwitcher(c, &cur, switcherSpaces, WorkspaceSwitcherOptions{}) })
}
