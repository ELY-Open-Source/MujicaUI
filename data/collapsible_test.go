package data

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

func TestCollapsibleToggleAndFocus(t *testing.T) {
	open, changes, clicks := false, 0, 0
	disabled := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).Children(func() {
			if Collapsible(c, &open, CollapsibleOptions{Title: "Details", Disabled: disabled}, func() {
				if input.Button(c, "Inner", input.ButtonOptions{}).Clicked() {
					clicks++
				}
			}).Changed() {
				changes++
			}
			input.Button(c, "After", input.ButtonOptions{})
		})
	}, 300, 300)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
	if tt.HasText("Inner") {
		t.Fatal("closed content built")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("Tab entered hidden content")
	}
	if err := tt.Click("Details"); err != nil {
		t.Fatal(err)
	}
	if !open || changes != 1 || !tt.HasText("Inner") {
		t.Fatalf("open %v changes %d", open, changes)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Inner") {
		t.Fatal("Tab skipped open content")
	}
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if open || changes != 2 {
		t.Fatalf("Space did not close: open %v changes %d", open, changes)
	}
	open = true // controlled from outside
	tt.Frame()
	if !tt.HasText("Inner") {
		t.Fatal("external open ignored")
	}
	disabled = true
	tt.Frame()
	_ = tt.Click("Details")
	if !open || changes != 2 {
		t.Fatal("disabled trigger toggled")
	}
}

func TestCollapsibleNeedsTitle(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	open := false
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Collapsible(c, &open, CollapsibleOptions{}, func() {})
	}, 100, 100)
}

func TestCollapsibleSnapshot(t *testing.T) {
	snapshot(t, "collapsible", 400, 200, func(c *ui.Context) {
		a := ui.Local(c.Root(), "a", func() *bool { b := true; return &b })
		b := ui.Local(c.Root(), "b", func() *bool { b := false; return &b })
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			Collapsible(c, *a, CollapsibleOptions{Title: "Rites of the evening"}, func() {
				ui.Text(c, "Candles are lit at the seventh bell.")
			})
			Collapsible(c, *b, CollapsibleOptions{Title: "Sealed records", Disabled: true}, func() {})
		})
	})
}

func TestCollapsibleClosingContentLeavesTabOrder(t *testing.T) {
	open := true
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(10).Children(func() {
			Collapsible(c, &open, CollapsibleOptions{Title: "Details"}, func() {
				input.Button(c, "Inner", input.ButtonOptions{})
			})
			input.Button(c, "After", input.ButtonOptions{})
		})
	}, 300, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace) // closes; the content shrinks for 120 ms
	if open || !tt.HasText("Inner") {
		t.Fatalf("expected a closing animation: open %v", open)
	}
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("After") {
		t.Fatal("Tab entered closing content")
	}
}
