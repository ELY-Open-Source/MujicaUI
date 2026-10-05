package mujica

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
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
		Use(c, Settings{})
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
