package input

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
)

func TestBulkActionBarCountActionsClear(t *testing.T) {
	selected, deletes := 0, 0
	dis := false
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if Button(c, "Select", ButtonOptions{}).Clicked() {
			selected += 3
		}
		r := BulkActionBar(c, selected, BulkActionBarOptions{Disabled: dis}, func() {
			if Button(c, "Delete", ButtonOptions{Variant: Danger}).Clicked() {
				deletes++
			}
		})
		if r.Changed() {
			selected = 0
		}
	}, 600, 300)
	if tt.HasText("selected") {
		t.Fatal("bar shown with nothing selected")
	}
	_ = tt.Click("Select")
	if !tt.HasText("3 selected") || !slices.Contains(tt.Announcements(), "3 selected") {
		t.Fatalf("texts %q", tt.Texts())
	}
	_ = tt.Click("Delete")
	if deletes != 1 {
		t.Fatal("action slot did not run")
	}
	dis = true
	tt.Frame()
	_ = tt.Click("Clear selection")
	if selected != 3 {
		t.Fatal("disabled bar cleared")
	}
	dis = false
	tt.Frame()
	_ = tt.Click("Delete") // focus inside the bar
	tt.Key(0, ui.KeyEscape)
	if selected != 0 || tt.HasText("selected") {
		t.Fatalf("Escape left %d selected", selected)
	}
	_ = tt.Click("Select")
	if err := tt.Click("Clear selection"); err != nil || selected != 0 {
		t.Fatalf("clear button: %v %d", err, selected)
	}
}

func TestBulkActionBarRejectsNegative(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("no panic")
		}
	}()
	ui.NewTester(func(c *ui.Context) { core.Use(c, core.Settings{}); BulkActionBar(c, -1, BulkActionBarOptions{}, nil) }, 200, 100)
}

func TestBulkActionBarSnapshot(t *testing.T) {
	snapshot(t, "bulk_action_bar", 600, 90, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			BulkActionBar(c, 12, BulkActionBarOptions{}, func() {
				Button(c, "Archive", ButtonOptions{Variant: Secondary})
				Button(c, "Delete", ButtonOptions{Variant: Danger})
			})
		})
	})
}
