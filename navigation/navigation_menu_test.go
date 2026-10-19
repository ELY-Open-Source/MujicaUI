package navigation

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
)

func navigationMenuSample() []NavigationLink {
	return []NavigationLink{
		{Label: "Home", Path: "/", Icon: icons.Must("house")},
		{Label: "Library", Items: []NavigationLink{
			{Label: "Poems", Path: "/library/poems", Description: "Verse of the court"},
			{Label: "Letters", Path: "/library/letters"},
			{Label: "Archive", Items: []NavigationLink{
				{Label: "Older letters", Path: "/archive/letters"},
			}},
		}},
		{Label: "About", Path: "/about"},
	}
}

func TestNavigationMenuPush(t *testing.T) {
	r := ui.NewRouter("/")
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		NavigationMenu(c, r, navigationMenuSample(), NavigationMenuOptions{History: true})
	}, 700, 400)
	if err := tt.Click("About"); err != nil {
		t.Fatal(err)
	}
	if r.Path() != "/about" {
		t.Fatalf("path %q", r.Path())
	}
	if err := tt.Click("Library"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Verse of the court") || !tt.HasText("Older letters") {
		t.Fatalf("group did not open: %q", tt.Texts())
	}
	if err := tt.Click("Older letters"); err != nil {
		t.Fatal(err)
	}
	if r.Path() != "/archive/letters" || tt.HasText("Verse of the court") {
		t.Fatalf("path %q, panel open %v", r.Path(), tt.HasText("Verse of the court"))
	}
	if err := tt.Click("Back"); err != nil {
		t.Fatal(err)
	}
	if r.Path() != "/about" {
		t.Fatalf("Back went to %q", r.Path())
	}
	if err := tt.Click("Forward"); err != nil {
		t.Fatal(err)
	}
	if r.Path() != "/archive/letters" {
		t.Fatalf("Forward went to %q", r.Path())
	}
}

func TestNavigationMenuKeyboard(t *testing.T) {
	r := ui.NewRouter("/")
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		NavigationMenu(c, r, navigationMenuSample(), NavigationMenuOptions{})
	}, 700, 400)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	if !tt.Focused("Library") {
		t.Fatal("Right did not reach Library")
	}
	tt.Key(0, ui.KeyDown)
	if !tt.Focused("Poems") {
		t.Fatalf("Down did not open into the panel: %q", tt.Texts())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if r.Path() != "/library/letters" {
		t.Fatalf("keyboard went to %q", r.Path())
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEscape)
	if tt.HasText("Poems") || !tt.Focused("Library") {
		t.Fatalf("Escape: open %v, focus on Library %v", tt.HasText("Poems"), tt.Focused("Library"))
	}
}

func TestNavigationMenuActive(t *testing.T) {
	l := navigationMenuSample()
	for _, tc := range []struct {
		path string
		want [3]bool
	}{
		{"/", [3]bool{true, false, false}},
		{"/library/poems", [3]bool{false, true, false}},
		{"/archive/letters/7", [3]bool{false, true, false}},
		{"/about", [3]bool{false, false, true}},
		{"/aboutx", [3]bool{false, false, false}},
	} {
		for i := range l {
			got := navigationMenuActive(l[i], tc.path)
			if got != tc.want[i] {
				t.Errorf("%s: link %s active %v", tc.path, l[i].Label, got)
			}
		}
	}
}

func TestNavigationMenuInvalid(t *testing.T) {
	for name, links := range map[string][]NavigationLink{
		"empty": nil,
		"path":  {{Label: "X"}},
		"label": {{Path: "/x"}},
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("no panic")
				}
			}()
			ui.Render(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				NavigationMenu(c, ui.NewRouter("/"), links, NavigationMenuOptions{})
			}, 300, 100, 1)
		})
	}
}

func TestNavigationMenuSnapshot(t *testing.T) {
	r := ui.NewRouter("/")
	r.Push("/library/poems")
	snapshot(t, "navigation_menu", 620, 120, func(c *ui.Context) {
		ui.Column(c).Padding(16).Children(func() {
			NavigationMenu(c, r, navigationMenuSample(), NavigationMenuOptions{History: true})
		})
	})
}
