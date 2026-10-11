package navigation

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TestNavigationPopupsFadeIn checks each popup fades in over PopupDuration instead of appearing at once.
func TestNavigationPopupsFadeIn(t *testing.T) {
	for _, tc := range []struct {
		name, open, sample string
		view               func(c *ui.Context)
	}{
		{"breadcrumb", "Show hidden path", "Alpha", func(c *ui.Context) {
			Breadcrumb(c, breadcrumbPath("Root", "Alpha", "Beta", "Gamma", "Omega", "Leaf"), BreadcrumbOptions{})
		}},
		{"pagination", "10 / page", "50 / page", func(c *ui.Context) {
			page, size := 1, 10
			Pagination(c, &page, &size, PaginationOptions{Total: 100, PageSizes: []int{10, 50}})
		}},
		{"navigation menu", "Library", "Letters", func(c *ui.Context) {
			NavigationMenu(c, ui.NewRouter("/"), navigationMenuSample(), NavigationMenuOptions{})
		}},
		{"menubar", "File", "Quit", func(c *ui.Context) {
			Menubar(c, (&menubarApp{}).menus(), MenubarOptions{})
		}},
		{"command palette", "", "Esc close", func(c *ui.Context) {
			open := true
			CommandPalette(c, &open, commandPaletteSample(), CommandPaletteOptions{})
		}},
	} {
		opened := time.Now()
		tt := ui.NewTester(func(c *ui.Context) {
			core.Use(c, core.Settings{Mode: core.Light})
			ui.Column(c).Fill().Padding(16).AlignItems(ui.Start).Children(func() { tc.view(c) })
		}, 700, 500)
		if tc.open != "" {
			if err := tt.Click(tc.open); err != nil {
				t.Fatal(err)
			}
		}
		r, ok := tt.Find(tc.sample)
		if !ok {
			t.Fatalf("%s: no %q", tc.name, tc.sample)
		}
		x, y := int(r.X-3), int(r.Y+r.H/2)
		early := tt.Image().RGBAAt(x, y)
		mid := time.Since(opened) < theme.PopupDuration/2
		time.Sleep(200 * time.Millisecond)
		tt.Frame()
		// A slow run may finish the fade before the first sample.
		if late := tt.Image().RGBAAt(x, y); mid && early == late {
			t.Errorf("%s: popup appeared at once (%v)", tc.name, late)
		}
	}
}

// TestNavigationMenuLongPanelScrolls checks a tall group stays in the window and follows the focus.
func TestNavigationMenuLongPanelScrolls(t *testing.T) {
	var links []NavigationLink
	for i := range 30 {
		links = append(links, NavigationLink{Label: fmt.Sprintf("Link %02d", i), Path: fmt.Sprintf("/l/%d", i)})
	}
	tt := ui.NewTester(func(c *ui.Context) {
		core.Use(c, core.Settings{})
		NavigationMenu(c, ui.NewRouter("/"), []NavigationLink{{Label: "Group", Items: links}}, NavigationMenuOptions{})
	}, 400, 300)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnd)
	tt.Frame()
	if !tt.Focused("Link 29") {
		t.Fatal("End did not focus the last link")
	}
	r, _ := tt.Find("Link 29")
	if r.H == 0 || r.Y+r.H > 300 {
		t.Fatalf("focused link off screen at %+v", r)
	}
}

// TestNavigationCaptionsFollowTextScale checks small text grows with the desktop text scale.
func TestNavigationCaptionsFollowTextScale(t *testing.T) {
	for _, tc := range []struct {
		text string
		view func(c *ui.Context)
	}{
		{"Name and seal", func(c *ui.Context) {
			cur := 0
			Steps(c, &cur, []Step{{Title: "Account", Description: "Name and seal"}}, StepsOptions{})
		}},
		{"12", func(c *ui.Context) {
			sel := "a"
			Sidebar(c, &sel, []SidebarSection{{Title: "Mail", Items: []SidebarItem{{ID: "a", Label: "Inbox", Icon: icons.Must("inbox"), Count: 12}}}}, SidebarOptions{})
		}},
		{"Verse of the court", func(c *ui.Context) {
			navigationMenuItem(c, NavigationLink{Label: "Poems", Path: "/p", Description: "Verse of the court"}, false, false)
		}},
		{"Esc close", func(c *ui.Context) {
			open := true
			CommandPalette(c, &open, []Command{{ID: "a", Label: "A"}}, CommandPaletteOptions{})
		}},
	} {
		height := func(scale float32) float32 {
			tt := ui.NewTester(func(c *ui.Context) {
				core.Use(c, core.Settings{})
				ui.Column(c).Fill().AlignItems(ui.Start).Children(func() { tc.view(c) })
			}, 900, 700)
			tt.SetPreferences(ui.Preferences{TextScale: scale})
			tt.Frame()
			r, ok := tt.Find(tc.text)
			if !ok {
				t.Fatalf("no %q", tc.text)
			}
			return r.H
		}
		if h1, h2 := height(1), height(2); h2 < h1*1.8 {
			t.Errorf("%q: height %v at 2x text, %v at 1x", tc.text, h2, h1)
		}
	}
}
