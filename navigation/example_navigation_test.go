package navigation_test

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/navigation"
)

func ExampleTabs() {
	tab := 0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if navigation.Tabs(c, &tab, []navigation.Tab{
			{Label: "General", Panel: func() { ui.Text(c, "General settings") }},
			{Label: "Fonts", Icon: icons.Must("italic"), Panel: func() { ui.Text(c, "Fonts") }},
			{Label: "Archive", Disabled: true},
		}, navigation.TabsOptions{Label: "Settings"}).Changed() {
			fmt.Println("tab", tab)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSegmentedControl() {
	view := 0
	build := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if navigation.SegmentedControl(c, &view, []navigation.Segment{
			{Label: "List", Icon: icons.Must("list"), HideLabel: true},
			{Label: "Grid", Icon: icons.Must("layout-grid"), HideLabel: true},
		}, navigation.SegmentedControlOptions{Label: "View"}).Changed() {
			fmt.Println("view", view)
		}
	}
	ui.Render(build, 400, 300, 1)
}

func ExampleBreadcrumb() {
	path := []navigation.BreadcrumbItem{{Label: "Home", Icon: icons.Must("house")}, {Label: "Library"}, {Label: "Poems"}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if i, ok := navigation.Breadcrumb(c, path, navigation.BreadcrumbOptions{}).Chosen(); ok {
			path = path[:i+1]
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePagination() {
	page, size := 1, 20
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if navigation.Pagination(c, &page, &size, navigation.PaginationOptions{Total: 420, PageSizes: []int{10, 20, 50}, Jump: true}).Changed() {
			fmt.Println("load page", page, "of size", size)
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSteps() {
	current := 1
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		navigation.Steps(c, &current, []navigation.Step{
			{Title: "Account", Clickable: true},
			{Title: "Profile", Description: "Name and seal"},
			{Title: "Confirm"},
		}, navigation.StepsOptions{Label: "Sign-up"})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSidebar() {
	selected := "inbox"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		navigation.Sidebar(c, &selected, []navigation.SidebarSection{
			{Title: "Mail", Items: []navigation.SidebarItem{
				{ID: "inbox", Label: "Inbox", Icon: icons.Must("inbox"), Count: 3},
				{Label: "Folders", Items: []navigation.SidebarItem{{ID: "work", Label: "Work", Icon: icons.Must("folder")}}},
			}},
		}, navigation.SidebarOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleNavigationMenu() {
	router := ui.NewRouter("/")
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		navigation.NavigationMenu(c, router, []navigation.NavigationLink{
			{Label: "Home", Path: "/"},
			{Label: "Library", Items: []navigation.NavigationLink{
				{Label: "Poems", Path: "/library/poems", Description: "Verse of the court"},
			}},
		}, navigation.NavigationMenuOptions{History: true})
		router.View(c, func(r *ui.Route) { ui.Text(c, r.Path()) })
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleMenubar() {
	wrap := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res := navigation.Menubar(c, []navigation.MenubarMenu{
			{Label: "File", Items: []navigation.MenubarItem{
				{ID: "new", Label: "New", Shortcut: navigation.MenubarShortcut{Mods: ui.Cmd, Key: ui.KeyN}},
				{Label: "Export", Items: []navigation.MenubarItem{{ID: "pdf", Label: "PDF"}}},
			}},
			{Label: "View", Items: []navigation.MenubarItem{
				{ID: "wrap", Label: "Word wrap", Kind: navigation.MenubarCheck, Checked: wrap},
			}},
		}, navigation.MenubarOptions{})
		if id, ok := res.Chosen(); ok && id == "wrap" {
			wrap = !wrap
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleToolbar() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		navigation.Toolbar(c, navigation.ToolbarOptions{Label: "Format"}, func() {
			navigation.ToolbarGroup(c, func() {
				navigation.ToolbarButton(c, "Bold", navigation.ToolbarButtonOptions{Icon: icons.Must("bold")})
				navigation.ToolbarButton(c, "Italic", navigation.ToolbarButtonOptions{Icon: icons.Must("italic")})
			})
			navigation.ToolbarSeparator(c)
			if navigation.ToolbarButton(c, "Copy", navigation.ToolbarButtonOptions{Icon: icons.Must("copy"), ShowLabel: true}).Clicked() {
				fmt.Println("copy")
			}
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleCommandPalette() {
	open := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		res := navigation.CommandPalette(c, &open, []navigation.Command{
			{ID: "new", Label: "New document", Group: "File"},
			{ID: "theme", Label: "Toggle dark mode", Group: "View", Keywords: []string{"night"}},
		}, navigation.CommandPaletteOptions{})
		if id, ok := res.Chosen(); ok {
			fmt.Println("run", id)
		}
	}
	ui.Render(view, 400, 300, 1)
}
