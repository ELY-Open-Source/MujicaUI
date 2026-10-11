package overlay

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
)

type menuApp struct {
	chosen      []string
	pinned      bool
	locked      bool
	pinDisabled bool
}

func (a *menuApp) build(m *PopupMenu) {
	m.Group("Letter", func() {
		if m.Item("Rename", PopupMenuItemOptions{}) {
			a.chosen = append(a.chosen, "Rename")
		}
		if m.Check("Pinned", &a.pinned, PopupMenuItemOptions{KeepOpen: true, Disabled: a.pinDisabled}) {
			a.chosen = append(a.chosen, "Pinned")
		}
	})
	m.Group("Move", func() {
		m.Submenu("Move to", PopupMenuItemOptions{}, func(m *PopupMenu) {
			for _, f := range []string{"Archive", "Drafts"} {
				if m.Item(f, PopupMenuItemOptions{}) {
					a.chosen = append(a.chosen, f)
				}
			}
		})
	})
	m.Separator()
	if m.Item("Burn", PopupMenuItemOptions{Disabled: !a.locked}) {
		a.chosen = append(a.chosen, "Burn")
	}
}

func (a *menuApp) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Column(c).Padding(20).Gap(8).AlignItems(ui.Start).Children(func() {
		DropdownMenu(c, "Actions", DropdownMenuOptions{}, a.build)
		input.Button(c, "After", input.ButtonOptions{})
	})
}

func TestDropdownMenuChooseClosesAndRestores(t *testing.T) {
	a := &menuApp{}
	tt := ui.NewTester(a.view, 500, 400)
	if err := tt.Click("Actions"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Letter") || !tt.Focused("Rename") {
		t.Fatalf("menu not open with focus on the first item: %q", tt.Texts())
	}
	_ = tt.Click("Rename")
	if len(a.chosen) != 1 || a.chosen[0] != "Rename" || tt.HasText("Letter") {
		t.Fatalf("chosen %q, still open %v", a.chosen, tt.HasText("Letter"))
	}
	if !tt.Focused("Actions") {
		t.Fatal("focus not back on the trigger")
	}
}

func TestDropdownMenuCheckKeepsOpenAndDisabledBlocks(t *testing.T) {
	a := &menuApp{}
	tt := ui.NewTester(a.view, 500, 400)
	_ = tt.Click("Actions")
	_ = tt.Click("Pinned")
	if !a.pinned || !tt.HasText("Letter") {
		t.Fatalf("pinned %v, open %v", a.pinned, tt.HasText("Letter"))
	}
	_ = tt.Click("Burn")
	if len(a.chosen) != 1 || !tt.HasText("Letter") {
		t.Fatalf("disabled item acted: %q", a.chosen)
	}
}

func TestDropdownMenuKeyboard(t *testing.T) {
	a := &menuApp{}
	tt := ui.NewTester(a.view, 500, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Rename") {
		t.Fatal("Down did not open the menu on the first item")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Move to") {
		t.Fatal("arrows did not reach Move to")
	}
	tt.Key(0, ui.KeyDown) // Burn is disabled: wraps to Rename
	if !tt.Focused("Rename") {
		t.Fatal("arrow did not skip the disabled item")
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Archive") {
		t.Fatalf("Right did not enter the submenu: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyLeft)
	if !tt.Focused("Move to") || tt.HasText("Archive") {
		t.Fatal("Left did not close the submenu")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if len(a.chosen) != 1 || a.chosen[0] != "Drafts" || tt.HasText("Letter") {
		t.Fatalf("submenu choice %q, open %v", a.chosen, tt.HasText("Letter"))
	}
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Letter") || !tt.Focused("Actions") {
		t.Fatal("Escape did not close and restore")
	}
}

func TestDropdownMenuSubmenuFlips(t *testing.T) {
	a := &menuApp{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Fill().Justify(ui.End).Padding(20).Children(func() {
			DropdownMenu(c, "Actions", DropdownMenuOptions{}, a.build)
		})
	}, 420, 400)
	_ = tt.Click("Actions")
	r, _ := tt.Find("Move to")
	tt.Move(r.X+5, r.Y+5)
	tt.Frame()
	s, ok := tt.Find("Archive")
	if !ok || s.X+s.W > 420 || s.X > r.X {
		t.Fatalf("submenu %+v not flipped left of %+v", s, r)
	}
}

func TestDropdownMenuCheckDisabledThisFrameKeepsValue(t *testing.T) {
	a := &menuApp{}
	tt := ui.NewTester(a.view, 500, 400)
	_ = tt.Click("Actions")
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Pinned") {
		t.Fatal("Pinned not focused")
	}
	a.pinDisabled = true // takes effect in the frame that handles the key
	tt.Key(0, ui.KeySpace)
	if a.pinned || len(a.chosen) != 0 {
		t.Fatalf("disabled check flipped: pinned %v chosen %q", a.pinned, a.chosen)
	}
}

func TestDropdownMenuWideSubmenuStaysInWindow(t *testing.T) {
	wide := strings.Repeat("Wide ", 30)
	var sub ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).Fill().Padding(20, 0, 0, 200).Children(func() {
			DropdownMenu(c, "Actions", DropdownMenuOptions{}, func(m *PopupMenu) {
				m.Item(wide, PopupMenuItemOptions{})
				m.Submenu("Move to", PopupMenuItemOptions{}, func(m *PopupMenu) {
					m.Item(wide+"archive", PopupMenuItemOptions{})
				})
			})
		})
	}, 800, 400)
	_ = tt.Click("Actions")
	r, _ := tt.Find("Move to")
	tt.Move(r.X+5, r.Y+5)
	tt.Frame()
	s, ok := tt.Find(wide + "archive")
	if !ok {
		t.Fatal("submenu did not open")
	}
	sub = s
	if sub.X < 4 || sub.X+sub.W > 796 {
		t.Fatalf("submenu item %+v outside the 800 DIP window", sub)
	}
}

func TestDropdownMenuDisabled(t *testing.T) {
	a := &menuApp{}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		DropdownMenu(c, "Actions", DropdownMenuOptions{Disabled: true}, a.build)
	}, 300, 300)
	_ = tt.Click("Actions")
	if tt.HasText("Letter") {
		t.Fatal("disabled trigger opened")
	}
}

func TestDropdownMenuTabInSubmenuClosesMenu(t *testing.T) {
	a := &menuApp{}
	tt := ui.NewTester(a.view, 500, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Archive") {
		t.Fatal("submenu not entered")
	}
	tt.Key(0, ui.KeyTab)
	if tt.HasText("Archive") || tt.HasText("Rename") {
		t.Fatal("Tab in a submenu left the menu open")
	}
	if !tt.Focused("Actions") {
		t.Fatal("focus did not return to the trigger")
	}
}

func TestDropdownMenuSnapshot(t *testing.T) {
	a := &menuApp{pinned: true}
	snapshotAfter(t, "dropdown_menu", 520, 320, func(c *ui.Context) {
		ui.Column(c).Padding(20).AlignItems(ui.Start).Children(func() {
			DropdownMenu(c, "Actions", DropdownMenuOptions{}, a.build)
		})
	}, func(tt *ui.Tester) {
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyDown) // opens with the first item focused
		tt.Key(0, ui.KeyEnd)
		tt.Key(0, ui.KeyRight) // opens Move to, focus ring on Archive
		if !tt.Focused("Archive") {
			t.Fatal("submenu did not open")
		}
	})
}
