package input

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/navigation"
)

type menubarApp struct {
	wrap   bool
	zoom   string
	chosen []string
}

func (a *menubarApp) menus() []navigation.MenubarMenu {
	return []navigation.MenubarMenu{
		{Label: "File", Items: []navigation.MenubarItem{
			{ID: "new", Label: "New", Shortcut: navigation.MenubarShortcut{Mods: ui.Cmd, Key: ui.KeyN}},
			{ID: "open", Label: "Open…", Shortcut: navigation.MenubarShortcut{Mods: ui.Cmd, Key: ui.KeyO}},
			{Label: "Export", Items: []navigation.MenubarItem{
				{ID: "pdf", Label: "PDF"},
				{ID: "png", Label: "Image", Disabled: true},
				{ID: "svg", Label: "Vector"},
			}},
			{Kind: navigation.MenubarSeparator},
			{ID: "quit", Label: "Quit"},
		}},
		{Label: "View", Items: []navigation.MenubarItem{
			{ID: "wrap", Label: "Word wrap", Kind: navigation.MenubarCheck, Checked: a.wrap},
			{Kind: navigation.MenubarSeparator},
			{ID: "zoom-fit", Label: "Fit", Kind: navigation.MenubarRadio, Checked: a.zoom == "zoom-fit"},
			{ID: "zoom-100", Label: "Actual size", Kind: navigation.MenubarRadio, Checked: a.zoom == "zoom-100"},
		}},
		{Label: "Help", Items: []navigation.MenubarItem{{ID: "about", Label: "About"}}},
	}
}

func (a *menubarApp) view(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Column(c).Fill().Children(func() {
		if id, ok := navigation.Menubar(c, a.menus(), navigation.MenubarOptions{}).Chosen(); ok {
			a.chosen = append(a.chosen, id)
			switch id {
			case "wrap":
				a.wrap = !a.wrap
			case "zoom-fit", "zoom-100":
				a.zoom = id
			}
		}
	})
}

func TestMenubarPointer(t *testing.T) {
	a := &menubarApp{zoom: "zoom-fit"}
	tt := ui.NewTester(a.view, 700, 500)
	if err := tt.Click("View"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Word wrap"); err != nil {
		t.Fatal(err)
	}
	if !a.wrap || tt.HasText("Actual size") {
		t.Fatalf("check item: wrap %v, menu still open %v", a.wrap, tt.HasText("Actual size"))
	}
	_ = tt.Click("View")
	if _, ok := tt.Find("Word wrap, checked"); !ok {
		t.Fatal("checked state not named")
	}
	if err := tt.Click("Actual size"); err != nil {
		t.Fatal(err)
	}
	if a.zoom != "zoom-100" {
		t.Fatalf("radio chose %q", a.zoom)
	}
	_ = tt.Click("File")
	r, _ := tt.Find("Export")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Move(r.X+r.W/2+1, r.Y+r.H/2)
	if !tt.HasText("Vector") {
		t.Fatalf("hover did not open submenu: %q", tt.Texts())
	}
	_ = tt.Click("Image")
	if err := tt.Click("Vector"); err != nil {
		t.Fatal(err)
	}
	if got := a.chosen; len(got) != 3 || got[2] != "svg" {
		t.Fatalf("chosen %q", got)
	}
	_ = tt.Click("Help")
	tt.ClickAt(600, 400)
	if tt.HasText("About") {
		t.Fatal("click outside did not close")
	}
}

func TestMenubarHoverSwitchesMenus(t *testing.T) {
	a := &menubarApp{}
	tt := ui.NewTester(a.view, 700, 500)
	_ = tt.Click("File")
	r, _ := tt.Find("Help")
	tt.Move(r.X+2, r.Y+2)
	if tt.HasText("Quit") || !tt.HasText("About") {
		t.Fatalf("hover did not switch: %q", tt.Texts())
	}
}

func TestMenubarKeyboard(t *testing.T) {
	a := &menubarApp{}
	tt := ui.NewTester(a.view, 700, 500)
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("File") {
		t.Fatal("Tab did not reach the bar")
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyLeft)
	if !tt.Focused("File") {
		t.Fatal("arrows did not move between titles")
	}
	tt.Key(0, ui.KeyDown) // opens File on New
	tt.Key(0, ui.KeyDown) // Open
	tt.Key(0, ui.KeyDown) // Export
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("Vector") {
		t.Fatal("Right did not open the submenu")
	}
	tt.Key(0, ui.KeyDown) // skips disabled Image to Vector
	tt.Key(0, ui.KeyEnter)
	if len(a.chosen) != 1 || a.chosen[0] != "svg" {
		t.Fatalf("chosen %q", a.chosen)
	}
	if tt.HasText("Quit") || !tt.Focused("File") {
		t.Fatal("menu did not close back to its title")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyUp) // wraps past the separator to Quit
	tt.Key(0, ui.KeyRight)
	if !tt.HasText("Word wrap") {
		t.Fatal("Right on a plain item did not move to the next menu")
	}
	tt.Key(0, ui.KeySpace)
	if !a.wrap {
		t.Fatal("Space did not toggle the check item")
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Word wrap") {
		t.Fatal("Escape did not close")
	}
}

func TestMenubarKeysAfterClickInMenu(t *testing.T) {
	a := &menubarApp{}
	tt := ui.NewTester(a.view, 700, 500)
	_ = tt.Click("File")
	if err := tt.Click("Export"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("File") {
		t.Fatal("click into the menu left the title without focus")
	}
	tt.Key(0, ui.KeyDown) // PDF to Vector, skipping disabled Image
	tt.Key(0, ui.KeyEnter)
	if len(a.chosen) != 1 || a.chosen[0] != "svg" {
		t.Fatalf("arrows after a click chose %q", a.chosen)
	}
	_ = tt.Click("File")
	_ = tt.Click("Export")
	_ = tt.Click("Image") // disabled: the press still keeps keys working
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Vector") || !tt.HasText("Quit") {
		t.Fatalf("Escape did not close only the submenu: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Quit") {
		t.Fatal("second Escape did not close the menu")
	}
}

func TestMenubarTabOutClosesMenu(t *testing.T) {
	for _, back := range []bool{false, true} {
		a := &menubarApp{}
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{})
			ui.Column(c).Fill().Children(func() {
				Button(c, "Before", ButtonOptions{})
				navigation.Menubar(c, a.menus(), navigation.MenubarOptions{})
				Button(c, "After", ButtonOptions{})
			})
		}, 700, 500)
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyTab)
		tt.Key(0, ui.KeyDown)
		if !tt.HasText("Quit") {
			t.Fatal("menu did not open")
		}
		mods, want := ui.Modifiers(0), "After"
		if back {
			mods, want = ui.Shift, "Before"
		}
		tt.Key(mods, ui.KeyTab)
		if tt.HasText("Quit") || !tt.Focused(want) {
			t.Fatalf("shift=%v: menu open %v, focus on %s %v", back, tt.HasText("Quit"), want, tt.Focused(want))
		}
		tt.Key(0, ui.KeyEscape)
		if !tt.Focused(want) {
			t.Fatalf("shift=%v: Escape moved focus off %s", back, want)
		}
	}
}

func TestMenubarLongMenuScrolls(t *testing.T) {
	var items []navigation.MenubarItem
	for i := range 30 {
		items = append(items, navigation.MenubarItem{ID: fmt.Sprint(i), Label: fmt.Sprintf("Entry %02d", i)})
	}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		navigation.Menubar(c, []navigation.MenubarMenu{{Label: "Long", Items: items}}, navigation.MenubarOptions{})
	}, 400, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnd)
	tt.Frame()
	r, ok := tt.Find("Entry 29")
	if !ok || r.H == 0 || r.Y+r.H > 300 {
		t.Fatalf("last item outside the window: %+v", r)
	}
}

func TestMenubarAnnouncesHighlight(t *testing.T) {
	a := &menubarApp{wrap: true}
	tt := ui.NewTester(a.view, 700, 500)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown) // View opens on Word wrap
	tt.Key(0, ui.KeyDown) // Fit
	got := tt.Announcements()
	if len(got) < 2 || got[len(got)-2] != "Word wrap, checked" || got[len(got)-1] != "Fit" {
		t.Fatalf("announcements %q", got)
	}
}

func TestMenubarShortcutText(t *testing.T) {
	got := navigation.MenubarShortcutText(navigation.MenubarShortcut{Mods: ui.Ctrl | ui.Shift, Key: ui.KeyF5})
	want := "Ctrl+Shift+F5"
	if runtime.GOOS == "darwin" {
		want = "⌃⇧F5"
	}
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestMenubarInvalid(t *testing.T) {
	for name, menus := range map[string][]navigation.MenubarMenu{
		"empty":   nil,
		"title":   {{Items: []navigation.MenubarItem{{ID: "a", Label: "A"}}}},
		"id":      {{Label: "M", Items: []navigation.MenubarItem{{Label: "A"}}}},
		"submenu": {{Label: "M", Items: []navigation.MenubarItem{{Label: "A", Kind: navigation.MenubarCheck, Items: []navigation.MenubarItem{}}}}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("no panic")
				}
			}()
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				navigation.Menubar(c, menus, navigation.MenubarOptions{})
			}, 300, 100, 1)
		})
	}
}

func TestMenubarSnapshot(t *testing.T) {
	a := &menubarApp{wrap: true, zoom: "zoom-fit"}
	snapshot(t, "menubar", 480, 100, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			navigation.Menubar(c, a.menus(), navigation.MenubarOptions{})
		})
	})
}
