package navigation

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
)

func sidebarSample() []SidebarSection {
	return []SidebarSection{
		{Title: "Library", Items: []SidebarItem{
			{ID: "inbox", Label: "Inbox", Icon: icons.Must("inbox"), Count: 12},
			{ID: "starred", Label: "Starred", Icon: icons.Must("star")},
			{Label: "Projects", Items: []SidebarItem{
				{ID: "opera", Label: "Opera", Icon: icons.Must("folder")},
				{ID: "ballet", Label: "Ballet", Icon: icons.Must("folder")},
			}},
		}},
		{Title: "System", Items: []SidebarItem{
			{ID: "settings", Label: "Settings", Icon: icons.Must("settings")},
		}},
	}
}

func TestSidebarChooseAndKeys(t *testing.T) {
	sel, changes := "inbox", 0
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if Sidebar(c, &sel, sidebarSample(), SidebarOptions{}).Changed() {
			changes++
		}
	}, 400, 500)
	if err := tt.Click("Opera"); err != nil {
		t.Fatal(err)
	}
	if sel != "opera" || changes != 1 {
		t.Fatalf("click: %q %d", sel, changes)
	}
	tt.Key(0, ui.KeyDown)
	if sel != "ballet" {
		t.Fatalf("Down chose %q", sel)
	}
	tt.Key(0, ui.KeyEnd)
	if sel != "settings" {
		t.Fatalf("End chose %q", sel)
	}
	tt.Key(0, ui.KeyHome)
	if sel != "inbox" || changes != 4 {
		t.Fatalf("Home chose %q, changes %d", sel, changes)
	}
}

func TestSidebarCollapse(t *testing.T) {
	sel := "inbox"
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Sidebar(c, &sel, sidebarSample(), SidebarOptions{})
	}, 400, 500)
	if err := tt.Click("Projects"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Opera") {
		t.Fatal("nested group did not collapse")
	}
	if err := tt.Click("Library"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Starred") || !tt.HasText("Settings") {
		t.Fatalf("section did not collapse: %q", tt.Texts())
	}
	_ = tt.Click("Library")
	if !tt.HasText("Starred") || tt.HasText("Opera") {
		t.Fatalf("collapse state lost: %q", tt.Texts())
	}
}

func TestSidebarCompact(t *testing.T) {
	sel := "inbox"
	var w float32
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Row(c).FillHeight().Children(func() {
			w = Sidebar(c, &sel, sidebarSample(), SidebarOptions{Compact: true}).Element.Bounds().W
		})
	}, 400, 500)
	tt.Frame()
	if w != 56 {
		t.Fatalf("compact width %v", w)
	}
	if err := tt.Click("Ballet"); err != nil {
		t.Fatal(err)
	}
	if sel != "ballet" {
		t.Fatalf("compact click chose %q", sel)
	}
	if _, ok := tt.Find("Inbox (12)"); !ok {
		t.Fatalf("compact item lost its name: %q", tt.Texts())
	}
}

func TestSidebarInvalid(t *testing.T) {
	for name, tc := range map[string]struct {
		s    []SidebarSection
		opts SidebarOptions
	}{
		"title":   {[]SidebarSection{{Items: []SidebarItem{{ID: "a", Label: "A"}}}}, SidebarOptions{}},
		"id":      {[]SidebarSection{{Title: "T", Items: []SidebarItem{{Label: "A"}}}}, SidebarOptions{}},
		"compact": {[]SidebarSection{{Title: "T", Items: []SidebarItem{{ID: "a", Label: "A"}}}}, SidebarOptions{Compact: true}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("no panic")
				}
			}()
			sel := ""
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				Sidebar(c, &sel, tc.s, tc.opts)
			}, 300, 300, 1)
		})
	}
}

func TestSidebarGroupKeyboard(t *testing.T) {
	sel := "inbox"
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Sidebar(c, &sel, sidebarSample(), SidebarOptions{})
	}, 400, 500)
	// Tab stops: the item list, the Library title, the Projects group.
	for range 3 {
		tt.Key(0, ui.KeyTab)
	}
	tt.Key(0, ui.KeyEnter)
	if tt.HasText("Opera") {
		t.Fatal("Enter did not collapse the group")
	}
	tt.Key(0, ui.KeySpace)
	if !tt.HasText("Opera") {
		t.Fatal("Space did not expand the group")
	}
}

func TestSidebarSameNamedGroups(t *testing.T) {
	sel := "a1"
	sections := []SidebarSection{
		{Title: "Court", Items: []SidebarItem{{Label: "Drafts", Items: []SidebarItem{{ID: "a1", Label: "Edict"}}}}},
		{Title: "Chapel", Items: []SidebarItem{{Label: "Drafts", Items: []SidebarItem{{ID: "b1", Label: "Hymn"}}}}},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		Sidebar(c, &sel, sections, SidebarOptions{})
	}, 400, 500)
	if err := tt.Click("Drafts"); err != nil {
		t.Fatal(err)
	}
	if tt.HasText("Edict") || !tt.HasText("Hymn") {
		t.Fatalf("same-named groups share state: %q", tt.Texts())
	}
}

func TestSidebarSnapshot(t *testing.T) {
	snapshot(t, "sidebar", 420, 360, func(c *ui.Context) {
		a, b := "opera", "starred"
		ui.Row(c).Fill().Gap(24).Children(func() {
			Sidebar(c, &a, sidebarSample(), SidebarOptions{})
			Sidebar(c, &b, sidebarSample(), SidebarOptions{Compact: true})
		})
	})
}
