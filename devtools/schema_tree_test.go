package devtools

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
)

var schemaSample = []SchemaNode{{Name: "court", Kind: SchemaDatabase, Children: []SchemaNode{
	{Name: "guests", Kind: SchemaTable, Children: []SchemaNode{
		{Name: "id", Kind: SchemaColumn, Detail: "bigint primary key"},
		{Name: "name", Kind: SchemaColumn, Detail: "text not null"},
		{Name: "guests_name_idx", Kind: SchemaIndex, Detail: "(name)"},
	}},
	{Name: "seating", Kind: SchemaView},
}}}

func schemaTester(s *SchemaTreeState, roots []SchemaNode, changes *int) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if SchemaTree(c, s, roots, SchemaTreeOptions{}).Changed() {
			*changes++
		}
	}, 520, 420)
}

func TestSchemaTreeKeyboardOpenAndChoose(t *testing.T) {
	var s SchemaTreeState
	n := 0
	tt := schemaTester(&s, schemaSample, &n)
	if _, ok := s.Selected(); ok || tt.HasText("guests") {
		t.Fatalf("start: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown) // court
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown) // guests
	tt.Key(0, ui.KeyRight)
	for _, s := range []string{"guests", "seating", "bigint primary key", "(name)"} {
		if !tt.HasText(s) {
			t.Fatalf("missing %q in %q", s, tt.Texts())
		}
	}
	tt.Key(0, ui.KeyDown)
	if p, _ := s.Selected(); p != "court/guests/id" || n == 0 {
		t.Fatalf("selected %q changes %d", p, n)
	}
}

func TestSchemaTreeEmptyRulesMany(t *testing.T) {
	var s SchemaTreeState
	n := 0
	tt := schemaTester(&s, nil, &n)
	if !tt.HasText("No databases") {
		t.Fatal("empty text")
	}
	cols := make([]SchemaNode, 100000)
	for i := range cols {
		cols[i] = SchemaNode{Name: fmt.Sprint("c", i), Kind: SchemaColumn}
	}
	big := []SchemaNode{{Name: "wide", Kind: SchemaTable, Children: cols}}
	var bs SchemaTreeState
	tt = schemaTester(&bs, big, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyRight)
	if got := len(tt.Texts()); got > 400 || !tt.HasText("c0") {
		t.Fatalf("%d texts for 100000 columns", got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("duplicate path did not panic")
		}
	}()
	var ds SchemaTreeState
	schemaTester(&ds, []SchemaNode{{Name: "a"}, {Name: "a"}}, &n)
}

func TestSchemaTreeSnapshot(t *testing.T) {
	var s SchemaTreeState
	snapshotAfter(t, "schema_tree", 500, 300, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { SchemaTree(c, &s, schemaSample, SchemaTreeOptions{Height: 260}) })
	}, func(tt *ui.Tester) {
		s = SchemaTreeState{}
		for _, k := range []ui.Key{ui.KeyTab, ui.KeyDown, ui.KeyRight, ui.KeyDown, ui.KeyRight, ui.KeyDown} {
			tt.Key(0, k)
		}
	})
}

func TestSchemaTreeAncestorDisabled(t *testing.T) {
	var s SchemaTreeState
	disabled := false
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Disabled(disabled).Children(func() {
			if SchemaTree(c, &s, schemaSample, SchemaTreeOptions{}).Changed() {
				changes++
			}
		})
	}, 520, 420)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	before, _ := s.Selected()
	old := changes
	disabled = true
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	if after, _ := s.Selected(); after != before || changes != old || tt.HasText("guests") {
		t.Fatalf("disabled schema: before=%q after=%q events=%d old=%d texts=%q", before, after, changes, old, tt.Texts())
	}
}

func TestSchemaTreeDisclosureHitArea(t *testing.T) {
	var state SchemaTreeState
	changes := 0
	tt := schemaTester(&state, schemaSample, &changes)
	row, ok := tt.Find("court")
	if !ok {
		t.Fatal("database row missing")
	}
	tt.ClickAt(row.X+31, row.Y+8)
	if !tt.HasText("guests") {
		t.Fatal("disclosure missed 28 DIP hit area")
	}
}

func TestSchemaTreeClearsSelectionWhenEmpty(t *testing.T) {
	var state SchemaTreeState
	roots := schemaSample
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SchemaTree(c, &state, roots, SchemaTreeOptions{})
	}, 520, 420)
	if err := tt.Click("court"); err != nil {
		t.Fatal(err)
	}
	if _, ok := state.Selected(); !ok {
		t.Fatal("database selection missing")
	}
	roots = nil
	tt.Frame()
	if _, ok := state.Selected(); ok || !tt.HasText("No databases") {
		t.Fatal("empty schema retained selection")
	}
}
