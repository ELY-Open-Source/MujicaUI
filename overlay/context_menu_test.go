package overlay

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

func contextMenuTester(a *menuApp, disabled bool) *ui.Tester {
	return ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(20).Gap(8).Children(func() {
			row := ui.Row(c).Height(80).FillWidth().Padding(12).Background(core.Tokens(c).Surface)
			row.Children(func() {
				input.Button(c, "Letter row", input.ButtonOptions{Variant: input.Ghost})
			})
			ContextMenu(c, row, ContextMenuOptions{Disabled: disabled}, a.build)
		})
	}, 500, 320)
}

func TestContextMenuOpensAtPointer(t *testing.T) {
	a := &menuApp{}
	tt := contextMenuTester(a, false)
	tt.RightClickAt(200, 60)
	tt.Frame()
	r, ok := tt.Find("Rename")
	// The label sits past the menu padding and the check column.
	if !ok || r.X < 230 || r.X > 250 || r.Y < 60 || r.Y > 110 {
		t.Fatalf("menu item %+v not at the pointer (200, 60)", r)
	}
	if tt.Menu() != nil {
		t.Fatal("the system menu opened")
	}
	_ = tt.Click("Rename")
	if len(a.chosen) != 1 || tt.HasText("Rename") {
		t.Fatalf("chosen %q", a.chosen)
	}
}

func TestContextMenuOnInnerButtonAndEdge(t *testing.T) {
	a := &menuApp{}
	tt := contextMenuTester(a, false)
	b, _ := tt.Find("Letter row")
	tt.RightClickAt(b.X+4, b.Y+4)
	if !tt.HasText("Rename") {
		t.Fatal("right click on an inner button did not open the menu")
	}
	tt.Key(0, ui.KeyEscape)
	tt.RightClickAt(480, 95)
	tt.Frame()
	last, _ := tt.Find("Burn")
	first, _ := tt.Find("Rename")
	if last.Y+last.H > 320 || first.Y > 95 || last.X+last.W > 500 {
		t.Fatalf("menu not kept in the window: first %+v last %+v", first, last)
	}
}

func TestContextMenuKeyboard(t *testing.T) {
	for _, key := range []struct {
		mods ui.Modifiers
		key  ui.Key
	}{{0, ui.KeyContextMenu}, {ui.Shift, ui.KeyF10}} {
		a := &menuApp{}
		tt := contextMenuTester(a, false)
		tt.Key(0, ui.KeyTab)
		tt.Key(key.mods, key.key)
		tt.Frame()
		if !tt.Focused("Rename") {
			t.Fatalf("%v: menu did not open with focus", key)
		}
		item, _ := tt.Find("Rename")
		if item.Y < 100 {
			t.Fatalf("%v: keyboard menu %+v not below the focused target", key, item)
		}
		tt.Key(0, ui.KeyEscape)
		if tt.HasText("Rename") || !tt.Focused("Letter row") {
			t.Fatalf("%v: Escape did not close and restore", key)
		}
	}
}

func TestContextMenuDisabled(t *testing.T) {
	a := &menuApp{}
	tt := contextMenuTester(a, true)
	tt.RightClickAt(200, 60)
	tt.Key(0, ui.KeyTab)
	tt.Key(ui.Shift, ui.KeyF10)
	if tt.HasText("Rename") {
		t.Fatal("disabled context menu opened")
	}
}

func TestContextMenuSnapshot(t *testing.T) {
	a := &menuApp{}
	snapshotAfter(t, "context_menu", 460, 320, func(c *ui.Context) {
		ui.Column(c).Padding(20).Children(func() {
			row := ui.Row(c).Height(36).FillWidth().Padding(0, 12).AlignItems(ui.Center).Background(core.Tokens(c).Selection)
			row.Children(func() { ui.Text(c, "Letter to the Countess") })
			ContextMenu(c, row, ContextMenuOptions{}, a.build)
		})
	}, func(tt *ui.Tester) {
		tt.RightClickAt(160, 40)
		r, ok := tt.Find("Pinned")
		if !ok {
			t.Fatal("menu did not open")
		}
		tt.Move(r.X+4, r.Y+4) // hover highlight on Pinned
	})
}
