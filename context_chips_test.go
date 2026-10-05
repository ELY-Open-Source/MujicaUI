package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestContextChipsRemoveAndFold(t *testing.T) {
	items := append([]ContextItem{}, mentionSample...)
	items = append(items, ContextItem{ID: "w2", Label: "Hall of Mirrors", Kind: ContextWeb})
	off := false
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Disabled(off).Children(func() {
			if id, ok := ContextChips(c, items, ContextChipsOptions{Max: 3}).Removed(); ok {
				for i, it := range items {
					if it.ID == id {
						items = append(items[:i], items[i+1:]...)
						break
					}
				}
			}
		})
	}, 600, 200)
	if tt.HasText("seating.xlsx") || !tt.HasText("+2 more") {
		t.Fatalf("fold: %q", tt.Texts())
	}
	_ = tt.Click("+2 more")
	if !tt.HasText("Hall of Mirrors") || !tt.HasText("Show less") {
		t.Fatalf("expand: %q", tt.Texts())
	}
	_ = tt.Click("Remove guests.csv")
	if len(items) != 4 || tt.HasText("guests.csv") {
		t.Fatalf("remove left %d items", len(items))
	}
	off = true
	_ = tt.Click("Remove Protocol")
	if len(items) != 4 {
		t.Fatal("disabled chips removed an item")
	}
}

func TestContextChipsNegativeMax(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic for a negative Max")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ContextChips(c, nil, ContextChipsOptions{Max: -1})
	}, 100, 100)
}

func TestContextChipsSnapshot(t *testing.T) {
	snapshot(t, "context_chips", 560, 90, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() { ContextChips(c, mentionSample, ContextChipsOptions{Max: 3}) })
	})
}

func TestContextChipsKeyboardAndEmpty(t *testing.T) {
	items, removed := mentionSample, ""
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		if id, ok := ContextChips(c, items, ContextChipsOptions{Max: 1}).Removed(); ok {
			removed = id
		}
	}, 400, 160)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if removed != "w1" {
		t.Fatalf("removed %q, want w1", removed)
	}
	items = nil
	tt.Frame()
	if tt.HasText("+3 more") || tt.HasText("Versailles") {
		t.Fatal("empty chips retained content")
	}
}
