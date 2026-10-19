package data

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

var descItems = []DescriptionItem{
	{Label: "Archivist", Value: "Sister Agnès", Copyable: true},
	{Label: "Shelf", Value: "North vault, row 7"},
	{Label: "Seal", Value: "Wine-red wax, cracked", Span: -1},
}

func TestDescriptionListColumnsAndCopy(t *testing.T) {
	width := float32(560)
	copied := -1
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(width).Children(func() {
			if i, ok := DescriptionList(c, descItems, DescriptionListOptions{}).Copied(); ok {
				copied = i
			}
		})
	}, 700, 400)
	a, _ := tt.Find("Archivist")
	b, _ := tt.Find("Shelf")
	seal, _ := tt.Find("Seal")
	if a.Y != b.Y || b.X <= a.X || seal.Y <= a.Y {
		t.Fatalf("two columns expected: %v %v %v", a, b, seal)
	}
	if err := tt.Click("Copy Archivist"); err != nil {
		t.Fatal(err)
	}
	if tt.Clipboard() != "Sister Agnès" || copied != 0 {
		t.Fatalf("clipboard %q copied %d", tt.Clipboard(), copied)
	}
	if _, ok := tt.Find("Copied"); !ok {
		t.Fatal("copied mark missing")
	}
	width = 300
	tt.Frame()
	a, _ = tt.Find("Archivist")
	b, _ = tt.Find("Shelf")
	if b.Y <= a.Y || b.X != a.X {
		t.Fatalf("narrow width kept columns: %v %v", a, b)
	}
}

func TestDescriptionListSpan(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Width(900).Children(func() {
			DescriptionList(c, []DescriptionItem{{Label: "A", Value: "a", Span: 2}, {Label: "B", Value: "b"}, {Label: "C", Value: "c"}}, DescriptionListOptions{Columns: 3})
		})
	}, 950, 300)
	a, _ := tt.Find("a")
	b, _ := tt.Find("b")
	cc, _ := tt.Find("c")
	if b.Y != a.Y || cc.Y <= a.Y {
		t.Fatalf("span 2 of 3: a %v b %v c %v", a, b, cc)
	}
}

func TestDescriptionListBadOptions(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		DescriptionList(c, nil, DescriptionListOptions{Columns: -1})
	}, 100, 100)
}

func TestDescriptionListSnapshot(t *testing.T) {
	snapshot(t, "description_list", 560, 200, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() { DescriptionList(c, descItems, DescriptionListOptions{}) })
	})
}
