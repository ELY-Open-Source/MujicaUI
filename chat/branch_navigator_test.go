package chat

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestBranchNavigatorButtonsAndKeys(t *testing.T) {
	i, changes, off := 1, 0, false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(off).Children(func() {
			if BranchNavigator(c, &i, 3).Changed() {
				changes++
			}
		})
	}, 300, 100)
	if !tt.HasText("2 / 3") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Next branch")
	_ = tt.Click("Next branch")
	if i != 2 || changes != 1 {
		t.Fatalf("index %d changes %d at the end", i, changes)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyHome)
	tt.Key(0, ui.KeyLeft)
	if i != 0 || changes != 3 {
		t.Fatalf("keys: index %d changes %d", i, changes)
	}
	tt.Key(0, ui.KeyEnd)
	if i != 2 {
		t.Fatalf("End gave %d", i)
	}
	off = true
	tt.Key(0, ui.KeyLeft)
	_ = tt.Click("Previous branch")
	if i != 2 {
		t.Fatalf("disabled navigator moved to %d", i)
	}
}

func TestBranchNavigatorBounds(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for an index past count")
		}
	}()
	i := 3
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		BranchNavigator(c, &i, 3)
	}, 100, 100)
}

func TestBranchNavigatorSnapshot(t *testing.T) {
	i := 1
	snapshot(t, "branch_navigator", 200, 60, func(c *ui.Context) {
		ui.Column(c).Padding(12).Children(func() { BranchNavigator(c, &i, 3) })
	})
}
