package overlay

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func splitButtonTester(dis *bool, runs *int, chosen *string) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Padding(20).Children(func() {
			r := SplitButton(c, "Save", SplitButtonOptions{Disabled: *dis}, func(m *PopupMenu) {
				for _, s := range []string{"Save as…", "Save all"} {
					if m.Item(s, PopupMenuItemOptions{}) {
						*chosen = s
					}
				}
			})
			if r.Changed() {
				*runs++
			}
		})
	}, 500, 300)
}

func TestSplitButtonMainAndMenu(t *testing.T) {
	dis, runs, chosen := false, 0, ""
	tt := splitButtonTester(&dis, &runs, &chosen)
	if err := tt.Click("Save"); err != nil {
		t.Fatal(err)
	}
	if runs != 1 || tt.HasText("Save all") {
		t.Fatalf("main click: runs %d, menu open %v", runs, tt.HasText("Save all"))
	}
	if err := tt.Click("More options"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Save all"); err != nil {
		t.Fatal(err)
	}
	if chosen != "Save all" || runs != 1 || tt.HasText("Save as…") {
		t.Fatalf("menu choice %q runs %d", chosen, runs)
	}
}

func TestSplitButtonKeyboardHalves(t *testing.T) {
	dis, runs, chosen := false, 0, ""
	tt := splitButtonTester(&dis, &runs, &chosen)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Save") {
		t.Fatal("first Tab misses the main half")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("More options") {
		t.Fatal("second Tab misses the menu half")
	}
	tt.Key(0, ui.KeyDown)
	if runs != 1 || !tt.HasText("Save as…") {
		t.Fatalf("runs %d, menu %v", runs, tt.HasText("Save as…"))
	}
	tt.Key(0, ui.KeyEnter)
	if chosen != "Save as…" {
		t.Fatalf("Enter chose %q", chosen)
	}
}

func TestSplitButtonDisabled(t *testing.T) {
	dis, runs, chosen := true, 0, ""
	tt := splitButtonTester(&dis, &runs, &chosen)
	_ = tt.Click("Save")
	_ = tt.Click("More options")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	if runs != 0 || tt.HasText("Save all") {
		t.Fatalf("disabled: runs %d menu %v", runs, tt.HasText("Save all"))
	}
	defer func() {
		if recover() == nil {
			t.Fatal("Danger variant accepted")
		}
	}()
	ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		SplitButton(c, "x", SplitButtonOptions{Variant: input.Danger}, func(*PopupMenu) {})
	}, 200, 100)
}

func TestSplitButtonSnapshot(t *testing.T) {
	snapshot(t, "split_button", 460, 80, func(c *ui.Context) {
		ui.Row(c).Padding(16).Gap(12).Children(func() {
			menu := func(m *PopupMenu) { m.Item("Other", PopupMenuItemOptions{}) }
			SplitButton(c, "Publish", SplitButtonOptions{}, menu)
			SplitButton(c, "Export", SplitButtonOptions{Variant: input.Secondary}, menu)
			SplitButton(c, "Merge", SplitButtonOptions{Disabled: true}, menu)
		})
	})
}
