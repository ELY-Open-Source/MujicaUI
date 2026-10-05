package main

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/navigation"
)

var navigationPages = []page{
	{ID: "059", Name: "Tabs", ZhName: "标签页", Category: "Navigation", Demo: navigationTabsDemo, Code: navigationTabsCode,
		Keys: "Tab enters the tab list at the chosen tab. Left/Right (Up/Down when vertical) choose the next enabled tab, skipping disabled ones and wrapping; Home and End choose the first and last."},
	{ID: "060", Name: "SegmentedControl", ZhName: "分段控件", Category: "Navigation", Demo: navigationSegmentedDemo, Code: navigationSegmentedCode,
		Keys: "One Tab stop. Arrow keys choose the previous or next segment and wrap; Home and End choose the ends."},
	{ID: "061", Name: "Breadcrumb", ZhName: "面包屑", Category: "Navigation", Demo: navigationBreadcrumbDemo, Code: navigationBreadcrumbCode,
		Keys: "Tab moves through the links and the fold button; Enter or Space activates. In the fold menu, Tab moves between items and Escape closes it, returning focus to the button."},
	{ID: "062", Name: "Pagination", ZhName: "分页", Category: "Navigation", Demo: navigationPaginationDemo, Code: navigationPaginationCode,
		Keys: "The page buttons are one Tab stop: Left/Right move, Home/End jump to the ends, Enter or Space goes to the page. The size select opens with Up/Down; the jump field goes on Enter."},
	{ID: "063", Name: "Steps", ZhName: "步骤条", Category: "Navigation", Demo: navigationStepsDemo, Code: navigationStepsCode,
		Keys: "Clickable steps are one Tab stop; arrow keys move between them and Enter or Space goes to the step."},
	{ID: "064", Name: "Sidebar", ZhName: "侧边栏", Category: "Navigation", Demo: navigationSidebarDemo, Code: navigationSidebarCode,
		Keys: "One Tab stop. Up/Down choose the previous or next item, Home/End the ends, typing the first letters of a label jumps to it. Section and group titles open and close with a click."},
	{ID: "065", Name: "NavigationMenu", ZhName: "导航菜单", Category: "Navigation", Demo: navigationMenuDemo, Code: navigationMenuCode,
		Keys: "Left/Right move between top links; Down or Enter opens a group with focus on its first link, Up/Down move inside, Escape closes. Back/Forward also answer Cmd+[ / Cmd+] (Alt+Left/Right off macOS)."},
	{ID: "066", Name: "Menubar", ZhName: "菜单栏", Category: "Navigation", Demo: navigationMenubarDemo, Code: navigationMenubarCode,
		Keys: "Tab reaches the bar; Left/Right move between menus, Down opens one. Up/Down/Home/End move through items, Right opens a submenu (or the next menu), Left closes it, Enter/Space chooses, Escape closes one level."},
	{ID: "067", Name: "Toolbar", ZhName: "工具栏", Category: "Navigation", Demo: navigationToolbarDemo, Code: navigationToolbarCode,
		Keys: "One Tab stop. Left/Right move between actions, Home/End to the ends, Enter or Space runs one. Actions that do not fit move to the » menu at the end."},
	{ID: "068", Name: "CommandPalette", ZhName: "命令面板", Category: "Navigation", Demo: navigationPaletteDemo, Code: navigationPaletteCode,
		Keys: "Cmd+K (Ctrl+K off macOS) opens it with focus in the search. Type to filter, Up/Down choose (wrapping), Enter runs and closes, Escape closes without running."},
}

// navigationCheck is a check box that logs its changes.
func navigationCheck(c *ui.Context, v *bool, label string, log func(string)) {
	was := *v
	ui.Checkbox(c, v, label)
	if *v != was {
		log(fmt.Sprintf("%s: %v", label, *v))
	}
}

// navigationSlider is a labeled whole-number slider that logs its changes.
func navigationSlider(c *ui.Context, v *float64, label string, lo, hi float64, log func(string)) {
	was := math.Round(*v)
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		ui.Text(c, fmt.Sprintf("%s %.0f", label, math.Round(*v)))
		ui.StepSlider(c, v, lo, hi, 1).Width(160)
	})
	if math.Round(*v) != was {
		log(fmt.Sprintf("%s: %.0f", label, math.Round(*v)))
	}
}

func navigationControls(c *ui.Context, fn func()) { ui.Row(c).Wrap().Gap(16).Children(fn) }

func navigationDemo(c *ui.Context, controls, body func()) {
	ui.Column(c).Gap(20).AlignItems(ui.Stretch).Children(func() {
		navigationControls(c, controls)
		body()
	})
}

const navigationTabsCode = `navigation.Tabs(c, &tab, []navigation.Tab{
	{Label: "General", Panel: func() { ui.Text(c, "General settings") }},
	{Label: "Fonts", Icon: icons.Must("italic"), Panel: func() { ui.Text(c, "Fonts") }},
	{Label: "Archive", Disabled: true},
}, navigation.TabsOptions{Label: "Settings"})`

func navigationTabsDemo(c *ui.Context, log func(string)) {
	type state struct {
		tab                     int
		vertical, archive, many bool
	}
	s := ui.Local(c.Root(), "059", func() state { return state{archive: true} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.vertical, "Vertical", log)
		navigationCheck(c, &s.archive, "Disable Archive", log)
		navigationCheck(c, &s.many, "Many tabs (overflow)", log)
	}, func() {
		names := []string{"General", "Fonts", "Archive", "Advanced"}
		if s.many {
			names = append(names, "Seals", "Heraldry", "Ledgers", "Envoys", "Banquets", "Chapel")
		}
		var tabs []navigation.Tab
		for i, n := range names {
			tabs = append(tabs, navigation.Tab{Label: n, Disabled: n == "Archive" && s.archive, Panel: func() {
				ui.Text(c, fmt.Sprintf("The %s panel (tab %d).", n, i+1))
			}})
		}
		tabs[1].Icon = icons.Must("italic")
		s.tab = min(s.tab, len(tabs)-1)
		if tabs[s.tab].Disabled {
			s.tab = 0
		}
		ui.Column(c).Width(420).Children(func() {
			if navigation.Tabs(c, &s.tab, tabs, navigation.TabsOptions{Vertical: s.vertical, Label: "Settings"}).Changed() {
				log("Changed: tab " + names[s.tab])
			}
		})
	})
}

const navigationSegmentedCode = `navigation.SegmentedControl(c, &view, []navigation.Segment{
	{Label: "List", Icon: icons.Must("list"), HideLabel: true},
	{Label: "Grid", Icon: icons.Must("layout-grid"), HideLabel: true},
}, navigation.SegmentedControlOptions{Label: "View"})`

func navigationSegmentedDemo(c *ui.Context, log func(string)) {
	type state struct {
		text, icon             int
		equal, iconOnly, disab bool
	}
	s := ui.Local(c.Root(), "060", func() state { return state{} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.equal, "Equal widths", log)
		navigationCheck(c, &s.iconOnly, "Icons only", log)
		navigationCheck(c, &s.disab, "Disabled", log)
	}, func() {
		opts := navigation.SegmentedControlOptions{Equal: s.equal, Disabled: s.disab, Label: "Range"}
		ui.Row(c).Gap(24).AlignItems(ui.Center).Children(func() {
			if navigation.SegmentedControl(c, &s.text, []navigation.Segment{{Label: "Day"}, {Label: "Week"}, {Label: "Fortnight"}}, opts).Changed() {
				log(fmt.Sprintf("Changed: range %d", s.text))
			}
			opts.Label = "View"
			if navigation.SegmentedControl(c, &s.icon, []navigation.Segment{
				{Label: "List", Icon: icons.Must("list"), HideLabel: s.iconOnly},
				{Label: "Grid", Icon: icons.Must("layout-grid"), HideLabel: s.iconOnly},
			}, opts).Changed() {
				log(fmt.Sprintf("Changed: view %d", s.icon))
			}
		})
	})
}

const navigationBreadcrumbCode = `if i, ok := navigation.Breadcrumb(c, path, navigation.BreadcrumbOptions{}).Chosen(); ok {
	path = path[:i+1]
}`

var navigationBreadcrumbNames = []string{"Home", "Library", "Manuscripts", "Seventeenth century", "Correspondence of the court", "Letters to the Duchess", "Draft of the winter decree"}

func navigationBreadcrumbDemo(c *ui.Context, log func(string)) {
	type state struct{ depth, max float64 }
	s := ui.Local(c.Root(), "061", func() state { return state{depth: 6, max: 4} })
	navigationDemo(c, func() {
		navigationSlider(c, &s.depth, "Depth", 1, 7, log)
		navigationSlider(c, &s.max, "MaxItems", 3, 7, log)
	}, func() {
		var path []navigation.BreadcrumbItem
		for i, n := range navigationBreadcrumbNames[:int(s.depth)] {
			it := navigation.BreadcrumbItem{Label: n}
			if i == 0 {
				it.Icon = icons.Must("house")
			}
			path = append(path, it)
		}
		if i, ok := navigation.Breadcrumb(c, path, navigation.BreadcrumbOptions{MaxItems: int(s.max)}).Chosen(); ok {
			log("Chosen: " + path[i].Label)
			s.depth = float64(i + 1)
		}
	})
}

const navigationPaginationCode = `if navigation.Pagination(c, &page, &size, navigation.PaginationOptions{Total: 420, PageSizes: []int{10, 20, 50}, Jump: true}).Changed() {
	load(page, size)
}`

func navigationPaginationDemo(c *ui.Context, log func(string)) {
	type state struct {
		page, size         int
		total              float64
		sizes, jump, disab bool
	}
	s := ui.Local(c.Root(), "062", func() state { return state{page: 1, size: 20, total: 420, sizes: true, jump: true} })
	navigationDemo(c, func() {
		navigationSlider(c, &s.total, "Total", 0, 1000, log)
		navigationCheck(c, &s.sizes, "Page sizes", log)
		navigationCheck(c, &s.jump, "Jump", log)
		navigationCheck(c, &s.disab, "Disabled", log)
	}, func() {
		opts := navigation.PaginationOptions{Total: int(s.total), Jump: s.jump, Disabled: s.disab}
		if s.sizes {
			opts.PageSizes = []int{10, 20, 50}
		} else {
			s.size = 20
		}
		if navigation.Pagination(c, &s.page, &s.size, opts).Changed() {
			log(fmt.Sprintf("Changed: page %d, size %d", s.page, s.size))
		}
	})
}

const navigationStepsCode = `navigation.Steps(c, &current, []navigation.Step{
	{Title: "Account", Clickable: true},
	{Title: "Profile", Description: "Name and seal"},
	{Title: "Confirm"},
}, navigation.StepsOptions{Label: "Sign-up"})`

func navigationStepsDemo(c *ui.Context, log func(string)) {
	type state struct {
		current                   int
		vertical, failed, clickOK bool
	}
	s := ui.Local(c.Root(), "063", func() state { return state{current: 1, clickOK: true} })
	steps := []navigation.Step{
		{Title: "Account", Description: "Name and seal"},
		{Title: "Profile", Description: "Titles and house"},
		{Title: "Payment", Description: "Treasury"},
		{Title: "Confirm"},
	}
	navigationDemo(c, func() {
		navigationCheck(c, &s.vertical, "Vertical", log)
		navigationCheck(c, &s.failed, "Payment failed", log)
		navigationCheck(c, &s.clickOK, "Done steps clickable", log)
		if input.Button(c, "Back", input.ButtonOptions{Variant: input.Secondary, Disabled: s.current == 0}).Clicked() {
			s.current--
			log(fmt.Sprintf("Back: step %d", s.current+1))
		}
		if input.Button(c, "Next", input.ButtonOptions{Disabled: s.current == len(steps)}).Clicked() {
			s.current++
			log(fmt.Sprintf("Next: step %d", s.current+1))
		}
	}, func() {
		for i := range steps {
			steps[i].Clickable = s.clickOK && i < s.current
		}
		steps[2].Error = s.failed
		if navigation.Steps(c, &s.current, steps, navigation.StepsOptions{Vertical: s.vertical, Label: "Sign-up"}).Changed() {
			log("Changed: step " + steps[s.current].Title)
		}
	})
}

const navigationSidebarCode = `navigation.Sidebar(c, &selected, []navigation.SidebarSection{
	{Title: "Mail", Items: []navigation.SidebarItem{
		{ID: "inbox", Label: "Inbox", Icon: icons.Must("inbox"), Count: 3},
		{Label: "Folders", Items: []navigation.SidebarItem{{ID: "work", Label: "Work", Icon: icons.Must("folder")}}},
	}},
}, navigation.SidebarOptions{})`

func navigationSidebarDemo(c *ui.Context, log func(string)) {
	type state struct {
		selected string
		compact  bool
		unread   float64
	}
	s := ui.Local(c.Root(), "064", func() state { return state{selected: "inbox", unread: 12} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.compact, "Compact", log)
		navigationSlider(c, &s.unread, "Unread", 0, 1200, log)
	}, func() {
		sections := []navigation.SidebarSection{
			{Title: "Library", Items: []navigation.SidebarItem{
				{ID: "inbox", Label: "Inbox", Icon: icons.Must("inbox"), Count: int(s.unread)},
				{ID: "starred", Label: "Starred", Icon: icons.Must("star")},
				{Label: "Projects", Items: []navigation.SidebarItem{
					{ID: "opera", Label: "Opera", Icon: icons.Must("folder")},
					{ID: "ballet", Label: "Ballet", Icon: icons.Must("folder")},
				}},
			}},
			{Title: "System", Items: []navigation.SidebarItem{
				{ID: "settings", Label: "Settings", Icon: icons.Must("settings")},
				{ID: "trash", Label: "Trash", Icon: icons.Must("trash-2")},
			}},
		}
		ui.Row(c).Height(360).Border(1, core.Tokens(c).Border).Radius(4).Clip().Children(func() {
			if navigation.Sidebar(c, &s.selected, sections, navigation.SidebarOptions{Compact: s.compact}).Changed() {
				log("Changed: " + s.selected)
			}
			ui.Text(c, "Showing "+s.selected).Padding(24)
		})
	})
}

const navigationMenuCode = `navigation.NavigationMenu(c, router, []navigation.NavigationLink{
	{Label: "Home", Path: "/"},
	{Label: "Library", Items: []navigation.NavigationLink{
		{Label: "Poems", Path: "/library/poems", Description: "Verse of the court"},
	}},
}, navigation.NavigationMenuOptions{History: true})
router.View(c, func(r *ui.Route) { ui.Text(c, r.Path()) })`

func navigationMenuDemo(c *ui.Context, log func(string)) {
	type state struct {
		router  *ui.Router
		history bool
		last    string
	}
	s := ui.Local(c.Root(), "065", func() state { return state{router: ui.NewRouter("/"), history: true, last: "/"} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.history, "Back and forward", log)
	}, func() {
		navigation.NavigationMenu(c, s.router, []navigation.NavigationLink{
			{Label: "Home", Path: "/", Icon: icons.Must("house")},
			{Label: "Library", Items: []navigation.NavigationLink{
				{Label: "Poems", Path: "/library/poems", Description: "Verse of the court"},
				{Label: "Letters", Path: "/library/letters", Description: "Correspondence"},
				{Label: "Archive", Items: []navigation.NavigationLink{
					{Label: "Older letters", Path: "/archive/letters"},
				}},
			}},
			{Label: "About", Path: "/about"},
		}, navigation.NavigationMenuOptions{History: s.history})
		if p := s.router.Path(); p != s.last {
			log(fmt.Sprintf("Path: %s → %s", s.last, p))
			s.last = p
		}
		ui.Column(c).Height(120).Children(func() {
			s.router.View(c, func(r *ui.Route) {
				r.Title(r.Path())
				ui.Text(c, "Page "+r.Path()).Padding(16)
			})
		})
	})
}

const navigationMenubarCode = `res := navigation.Menubar(c, []navigation.MenubarMenu{
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
}`

func navigationMenubarDemo(c *ui.Context, log func(string)) {
	type state struct {
		wrap, ruler, locked bool
		zoom                string
	}
	s := ui.Local(c.Root(), "066", func() state { return state{wrap: true, zoom: "fit"} })
	navigationDemo(c, func() {
		navigationCheck(c, &s.locked, "Lock editing (disables Edit items)", log)
	}, func() {
		menus := []navigation.MenubarMenu{
			{Label: "File", Items: []navigation.MenubarItem{
				{ID: "new", Label: "New", Shortcut: navigation.MenubarShortcut{Mods: ui.Cmd, Key: ui.KeyN}},
				{ID: "open", Label: "Open…", Shortcut: navigation.MenubarShortcut{Mods: ui.Cmd, Key: ui.KeyO}},
				{Label: "Export", Items: []navigation.MenubarItem{
					{ID: "pdf", Label: "PDF"},
					{ID: "png", Label: "Image"},
					{Label: "Archive", Items: []navigation.MenubarItem{{ID: "zip", Label: "Zip"}, {ID: "tar", Label: "Tarball"}}},
				}},
				{Kind: navigation.MenubarSeparator},
				{ID: "quit", Label: "Quit", Shortcut: navigation.MenubarShortcut{Mods: ui.Cmd, Key: ui.KeyQ}},
			}},
			{Label: "Edit", Items: []navigation.MenubarItem{
				{ID: "undo", Label: "Undo", Disabled: s.locked, Shortcut: navigation.MenubarShortcut{Mods: ui.Cmd, Key: ui.KeyZ}},
				{ID: "redo", Label: "Redo", Disabled: s.locked, Shortcut: navigation.MenubarShortcut{Mods: ui.Cmd | ui.Shift, Key: ui.KeyZ}},
			}},
			{Label: "View", Items: []navigation.MenubarItem{
				{ID: "wrap", Label: "Word wrap", Kind: navigation.MenubarCheck, Checked: s.wrap},
				{ID: "ruler", Label: "Ruler", Kind: navigation.MenubarCheck, Checked: s.ruler},
				{Kind: navigation.MenubarSeparator},
				{ID: "fit", Label: "Fit to window", Kind: navigation.MenubarRadio, Checked: s.zoom == "fit"},
				{ID: "actual", Label: "Actual size", Kind: navigation.MenubarRadio, Checked: s.zoom == "actual"},
			}},
		}
		ui.Column(c).Height(260).Border(1, core.Tokens(c).Border).Radius(4).Clip().Children(func() {
			if id, ok := navigation.Menubar(c, menus, navigation.MenubarOptions{}).Chosen(); ok {
				switch id {
				case "wrap":
					s.wrap = !s.wrap
				case "ruler":
					s.ruler = !s.ruler
				case "fit", "actual":
					s.zoom = id
				}
				log("Chosen: " + id)
			}
			ui.Text(c, fmt.Sprintf("Word wrap %v · Ruler %v · Zoom %s", s.wrap, s.ruler, s.zoom)).Padding(16)
		})
	})
}

const navigationToolbarCode = `navigation.Toolbar(c, navigation.ToolbarOptions{Label: "Format"}, func() {
	navigation.ToolbarGroup(c, func() {
		navigation.ToolbarButton(c, "Bold", navigation.ToolbarButtonOptions{Icon: icons.Must("bold")})
		navigation.ToolbarButton(c, "Italic", navigation.ToolbarButtonOptions{Icon: icons.Must("italic")})
	})
	navigation.ToolbarSeparator(c)
	if navigation.ToolbarButton(c, "Copy", navigation.ToolbarButtonOptions{Icon: icons.Must("copy"), ShowLabel: true}).Clicked() {
		copySelection()
	}
})`

func navigationToolbarDemo(c *ui.Context, log func(string)) {
	type state struct {
		width     float64
		labels    bool
		noDeletes bool
	}
	s := ui.Local(c.Root(), "067", func() state { return state{width: 520, noDeletes: true} })
	navigationDemo(c, func() {
		navigationSlider(c, &s.width, "Width", 120, 640, log)
		navigationCheck(c, &s.labels, "Show labels", log)
		navigationCheck(c, &s.noDeletes, "Disable Delete", log)
	}, func() {
		ui.Column(c).Width(float32(s.width)).Children(func() {
			navigation.Toolbar(c, navigation.ToolbarOptions{Label: "Format"}, func() {
				button := func(label, icon string, disabled bool) {
					if navigation.ToolbarButton(c, label, navigation.ToolbarButtonOptions{Icon: icons.Must(icon), ShowLabel: s.labels, Disabled: disabled}).Clicked() {
						log("Clicked: " + label)
					}
				}
				navigation.ToolbarGroup(c, func() {
					button("Bold", "bold", false)
					button("Italic", "italic", false)
					button("Underline", "underline", false)
				})
				navigation.ToolbarSeparator(c)
				button("Copy", "copy", false)
				button("Image", "image", false)
				button("Delete", "trash-2", s.noDeletes)
				navigation.ToolbarSeparator(c)
				button("Settings", "settings", false)
			})
		})
	})
}

const navigationPaletteCode = `res := navigation.CommandPalette(c, &open, []navigation.Command{
	{ID: "new", Label: "New document", Group: "File"},
	{ID: "theme", Label: "Toggle dark mode", Group: "View", Keywords: []string{"night"}},
}, navigation.CommandPaletteOptions{})
if id, ok := res.Chosen(); ok {
	run(id)
}`

func navigationPaletteDemo(c *ui.Context, log func(string)) {
	type state struct{ open, wasOpen bool }
	s := ui.Local(c.Root(), "068", func() state { return state{} })
	navigationDemo(c, func() {
		if input.Button(c, "Open command palette", input.ButtonOptions{Icon: icons.Must("command")}).Clicked() {
			s.open = true
		}
		ui.Text(c, "or press Cmd+K / Ctrl+K").TextColor(core.Tokens(c).TextMuted)
	}, func() {
		res := navigation.CommandPalette(c, &s.open, []navigation.Command{
			{ID: "new", Label: "New document", Group: "File", Icon: icons.Must("file")},
			{ID: "open", Label: "Open recent", Group: "File", Icon: icons.Must("folder")},
			{ID: "save", Label: "Save all", Group: "File", Icon: icons.Must("file-text")},
			{ID: "theme", Label: "Toggle dark mode", Group: "View", Keywords: []string{"theme", "night"}, Icon: icons.Must("moon")},
			{ID: "zoom-in", Label: "Zoom in", Group: "View", Icon: icons.Must("zoom-in")},
			{ID: "zoom-out", Label: "Zoom out", Group: "View", Icon: icons.Must("zoom-out")},
			{ID: "settings", Label: "Open settings", Group: "Application", Keywords: []string{"preferences"}, Icon: icons.Must("settings")},
		}, navigation.CommandPaletteOptions{})
		if s.open != s.wasOpen {
			log(fmt.Sprintf("Open: %v", s.open))
			s.wasOpen = s.open
		}
		if id, ok := res.Chosen(); ok {
			log("Run: " + id)
		}
	})
}
