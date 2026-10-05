package mujica

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
)

var treeSelectTestNodes = []TreeSelectNode{
	{ID: "court", Label: "Court", Children: []TreeSelectNode{
		{ID: "king", Label: "King"},
		{ID: "queen", Label: "Queen"},
		{ID: "jester", Label: "Jester", Disabled: true},
	}},
	{ID: "church", Label: "Church", Children: []TreeSelectNode{
		{ID: "bishop", Label: "Bishop"},
	}},
	{ID: "guard", Label: "Guard"},
}

func treeSelectTester(sel *[]string, opts TreeSelectOptions, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Padding(10).Children(func() {
			if TreeSelect(c, sel, treeSelectTestNodes, opts).Changed() {
				*changes++
			}
		})
	}, 500, 500)
}

func TestTreeSelectSingle(t *testing.T) {
	var sel []string
	changes := 0
	tt := treeSelectTester(&sel, TreeSelectOptions{}, &changes)
	_ = tt.Click("Select…")
	if tt.HasText("King") {
		t.Fatal("collapsed children shown")
	}
	if err := tt.Click("Expand Court"); err != nil {
		t.Fatal(err)
	}
	_ = tt.Click("Queen")
	if !slices.Equal(sel, []string{"queen"}) || changes != 1 || tt.HasText("King") {
		t.Fatalf("sel %v changes %d open %v", sel, changes, tt.HasText("King"))
	}
}

func TestTreeSelectCascade(t *testing.T) {
	var sel []string
	changes := 0
	tt := treeSelectTester(&sel, TreeSelectOptions{Multiple: true}, &changes)
	_ = tt.Click("Select…")
	_ = tt.Click("Expand Court")
	_ = tt.Click("Court")
	// The disabled Jester stays out.
	if !slices.Equal(sel, []string{"king", "queen"}) {
		t.Fatalf("cascade checked %v", sel)
	}
	// The chosen labels also show as tags: reach King with the keys.
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeySpace)
	if !slices.Equal(sel, []string{"queen"}) {
		t.Fatalf("after unchecking King: %v", sel)
	}
	if got := treeSelectCheck(&treeSelectTestNodes[0], sel, TreeSelectOptions{Multiple: true}); got != 2 {
		t.Fatalf("Court state %v, want mixed", got)
	}
	_ = tt.Click("Court")
	if !slices.Equal(sel, []string{"king", "queen"}) || changes != 3 {
		t.Fatalf("mixed parent click: %v changes %d", sel, changes)
	}
}

func TestTreeSelectIndependent(t *testing.T) {
	var sel []string
	changes := 0
	tt := treeSelectTester(&sel, TreeSelectOptions{Multiple: true, Rule: TreeSelectIndependent}, &changes)
	_ = tt.Click("Select…")
	_ = tt.Click("Court")
	if !slices.Equal(sel, []string{"court"}) {
		t.Fatalf("independent %v", sel)
	}
}

func TestTreeSelectSearch(t *testing.T) {
	var sel []string
	changes := 0
	tt := treeSelectTester(&sel, TreeSelectOptions{Multiple: true}, &changes)
	_ = tt.Click("Select…")
	tt.Type("bish")
	if !tt.HasText("Church") || !tt.HasText("Bishop") || tt.HasText("Court") || tt.HasText("Guard") {
		t.Fatalf("search shows %v", tt.Texts())
	}
	tt.Type("zz")
	if !tt.HasText("No results") {
		t.Fatal("no-results missing")
	}
}

func TestTreeSelectKeyboard(t *testing.T) {
	var sel []string
	changes := 0
	tt := treeSelectTester(&sel, TreeSelectOptions{Multiple: true, Label: "Roles"}, &changes)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Court") {
		t.Fatal("Down from search did not reach the tree")
	}
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("King") {
		t.Fatal("Right did not expand")
	}
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("King") {
		t.Fatal("second Right did not enter the children")
	}
	tt.Key(0, ui.KeySpace)
	tt.Key(0, ui.KeyLeft)
	if !tt.Focused("Court") {
		t.Fatal("Left did not go to the parent")
	}
	tt.Key(0, ui.KeyLeft)
	if tt.HasText("Queen") {
		t.Fatal("Left did not collapse")
	}
	if !slices.Equal(sel, []string{"king"}) {
		t.Fatalf("sel %v", sel)
	}
}

func TestTreeSelectSnapshot(t *testing.T) {
	snapshot(t, "tree_select", 460, 110, func(c *ui.Context) {
		single, multi := []string{"bishop"}, []string{"king", "queen", "bishop", "guard"}
		ui.Row(c).Padding(20).Gap(16).Children(func() {
			TreeSelect(c, &single, treeSelectTestNodes, TreeSelectOptions{})
			TreeSelect(c, &multi, treeSelectTestNodes, TreeSelectOptions{Multiple: true, MaxTags: 2})
		})
	})
}
