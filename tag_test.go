package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

type tagCounts struct{ clicked, closed int }

func tagTester(opts TagOptions, n *tagCounts) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Row(c).AlignItems(ui.Start).Children(func() {
			v := Tag(c, "Velvet", opts)
			if v.Clicked {
				n.clicked++
			}
			if v.Closed {
				n.closed++
			}
		})
	}, 300, 100)
}

func TestTagClickAndCloseAreSeparate(t *testing.T) {
	var n tagCounts
	tt := tagTester(TagOptions{Clickable: true, Closable: true, Icon: "star"}, &n)
	if err := tt.Click("Remove Velvet"); err != nil {
		t.Fatal(err)
	}
	if n != (tagCounts{0, 1}) {
		t.Fatalf("after close click %+v", n)
	}
	if err := tt.Click("Velvet"); err != nil {
		t.Fatal(err)
	}
	if n != (tagCounts{1, 1}) {
		t.Fatalf("after label click %+v", n)
	}
}

func TestTagKeyboard(t *testing.T) {
	var n tagCounts
	tt := tagTester(TagOptions{Clickable: true, Closable: true}, &n)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyDelete)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Remove Velvet") {
		t.Fatal("remove button not second focus stop")
	}
	tt.Key(0, ui.KeySpace)
	if n != (tagCounts{1, 2}) {
		t.Fatalf("counts %+v", n)
	}
}

func TestTagPlainTakesNoFocus(t *testing.T) {
	var n tagCounts
	tt := tagTester(TagOptions{}, &n)
	_ = tt.Click("Velvet")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if n != (tagCounts{}) || tt.Focused("Velvet") {
		t.Fatalf("plain tag reacted: %+v", n)
	}
}

func TestTagDisabled(t *testing.T) {
	var n tagCounts
	tt := tagTester(TagOptions{Clickable: true, Closable: true, Disabled: true}, &n)
	_ = tt.Click("Velvet")
	_ = tt.Click("Remove Velvet")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyDelete)
	if n != (tagCounts{}) {
		t.Fatalf("disabled tag reacted: %+v", n)
	}
}

func TestTagSnapshot(t *testing.T) {
	snapshot(t, "tag", 480, 80, func(c *ui.Context) {
		ui.Row(c).Padding(24).Gap(8).AlignItems(ui.Center).Children(func() {
			Tag(c, "Baroque", TagOptions{})
			Tag(c, "Favorite", TagOptions{Icon: "star", Clickable: true})
			Tag(c, "Archive", TagOptions{Closable: true})
			Tag(c, "Locked", TagOptions{Icon: "lock", Closable: true, Disabled: true})
		})
	})
}
