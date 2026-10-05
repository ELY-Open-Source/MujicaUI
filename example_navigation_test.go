package mujica_test

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
)

func ExampleTabs() {
	tab := 0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.Tabs(c, &tab, []mujica.Tab{
			{Label: "General", Panel: func() { ui.Text(c, "General settings") }},
			{Label: "Fonts", Icon: icons.Must("italic"), Panel: func() { ui.Text(c, "Fonts") }},
			{Label: "Archive", Disabled: true},
		}, mujica.TabsOptions{Label: "Settings"}).Changed() {
			fmt.Println("tab", tab)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSegmentedControl() {
	view := 0
	build := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.SegmentedControl(c, &view, []mujica.Segment{
			{Label: "List", Icon: icons.Must("list"), HideLabel: true},
			{Label: "Grid", Icon: icons.Must("layout-grid"), HideLabel: true},
		}, mujica.SegmentedControlOptions{Label: "View"}).Changed() {
			fmt.Println("view", view)
		}
	}
	ui.Render(build, 400, 300, 1)
}

func ExampleBreadcrumb() {
	path := []mujica.BreadcrumbItem{{Label: "Home", Icon: icons.Must("house")}, {Label: "Library"}, {Label: "Poems"}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if i, ok := mujica.Breadcrumb(c, path, mujica.BreadcrumbOptions{}).Chosen(); ok {
			path = path[:i+1]
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePagination() {
	page, size := 1, 20
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.Pagination(c, &page, &size, mujica.PaginationOptions{Total: 420, PageSizes: []int{10, 20, 50}, Jump: true}).Changed() {
			fmt.Println("load page", page, "of size", size)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSteps() {
	current := 1
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Steps(c, &current, []mujica.Step{
			{Title: "Account", Clickable: true},
			{Title: "Profile", Description: "Name and seal"},
			{Title: "Confirm"},
		}, mujica.StepsOptions{Label: "Sign-up"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSidebar() {
	selected := "inbox"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Sidebar(c, &selected, []mujica.SidebarSection{
			{Title: "Mail", Items: []mujica.SidebarItem{
				{ID: "inbox", Label: "Inbox", Icon: icons.Must("inbox"), Count: 3},
				{Label: "Folders", Items: []mujica.SidebarItem{{ID: "work", Label: "Work", Icon: icons.Must("folder")}}},
			}},
		}, mujica.SidebarOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleNavigationMenu() {
	router := ui.NewRouter("/")
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.NavigationMenu(c, router, []mujica.NavigationLink{
			{Label: "Home", Path: "/"},
			{Label: "Library", Items: []mujica.NavigationLink{
				{Label: "Poems", Path: "/library/poems", Description: "Verse of the court"},
			}},
		}, mujica.NavigationMenuOptions{History: true})
		router.View(c, func(r *ui.Route) { ui.Text(c, r.Path()) })
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMenubar() {
	wrap := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		res := mujica.Menubar(c, []mujica.MenubarMenu{
			{Label: "File", Items: []mujica.MenubarItem{
				{ID: "new", Label: "New", Shortcut: mujica.MenubarShortcut{Mods: ui.Cmd, Key: ui.KeyN}},
				{Label: "Export", Items: []mujica.MenubarItem{{ID: "pdf", Label: "PDF"}}},
			}},
			{Label: "View", Items: []mujica.MenubarItem{
				{ID: "wrap", Label: "Word wrap", Kind: mujica.MenubarCheck, Checked: wrap},
			}},
		}, mujica.MenubarOptions{})
		if id, ok := res.Chosen(); ok && id == "wrap" {
			wrap = !wrap
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToolbar() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Toolbar(c, mujica.ToolbarOptions{Label: "Format"}, func() {
			mujica.ToolbarGroup(c, func() {
				mujica.ToolbarButton(c, "Bold", mujica.ToolbarButtonOptions{Icon: icons.Must("bold")})
				mujica.ToolbarButton(c, "Italic", mujica.ToolbarButtonOptions{Icon: icons.Must("italic")})
			})
			mujica.ToolbarSeparator(c)
			if mujica.ToolbarButton(c, "Copy", mujica.ToolbarButtonOptions{Icon: icons.Must("copy"), ShowLabel: true}).Clicked() {
				fmt.Println("copy")
			}
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCommandPalette() {
	open := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		res := mujica.CommandPalette(c, &open, []mujica.Command{
			{ID: "new", Label: "New document", Group: "File"},
			{ID: "theme", Label: "Toggle dark mode", Group: "View", Keywords: []string{"night"}},
		}, mujica.CommandPaletteOptions{})
		if id, ok := res.Chosen(); ok {
			fmt.Println("run", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}
