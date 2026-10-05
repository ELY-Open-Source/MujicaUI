package layout_test

import (
	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/layout"
)

func ExampleAppShell() {
	var collapsed bool
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		layout.AppShell(c, &collapsed, layout.AppShellOptions{
			TitleBar: func() { layout.TitleBar(c, "Court Records", layout.TitleBarOptions{}) },
			Sidebar:  func() { ui.Text(c, "Chronicles") },
			Content:  func() { ui.Text(c, "The Grand Hall") },
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTitleBar() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		layout.TitleBar(c, "Court Records", layout.TitleBarOptions{
			Trailing: func() { input.Button(c, "Share", input.ButtonOptions{Variant: input.Ghost}) },
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStatusBar() {
	opened := 0
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		bar := layout.StatusBar(c, layout.StatusBarOptions{
			Left:  []layout.StatusItem{{ID: "sync", Text: "Syncing", ShowProgress: true, Progress: 0.4}},
			Right: []layout.StatusItem{{ID: "bell", Label: "Notifications", Icon: icons.Must("bell"), Action: true}},
		})
		if bar.Clicked == "bell" {
			opened++
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleContainer() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		layout.Container(c, layout.ContainerOptions{MaxWidth: 320, Center: true}, func() {
			ui.Text(c, "A centered column of prose.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStack() {
	showDraft := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		layout.Stack(c, layout.StackOptions{Horizontal: true, Gap: layout.SpaceM}, func() {
			input.Button(c, "Save", input.ButtonOptions{})
			if showDraft {
				input.Button(c, "Draft", input.ButtonOptions{Variant: input.Secondary})
			}
			input.Button(c, "Cancel", input.ButtonOptions{Variant: input.Ghost})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleGrid() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		layout.Grid(c, layout.GridOptions{Columns: []ui.Track{ui.Fixed(96), ui.Fr(1)}}, func() {
			layout.GridCell(c, layout.GridCellOptions{ColumnSpan: -1}, func() { ui.Text(c, "Header") })
			ui.Text(c, "Name")
			ui.Text(c, "Ada")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleScrollArea() {
	var offset ui.ScrollState
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		layout.ScrollArea(c, &offset, layout.ScrollAreaOptions{Label: "Log"}, func() {
			for range 50 {
				ui.Text(c, "A line of the log")
			}
		}).Height(200)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSplitPane() {
	split := layout.SplitPaneState{Size: 160}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		layout.SplitPane(c, &split, layout.SplitPaneOptions{MinFirst: 120, Collapsible: layout.SplitPaneFirst},
			func() { ui.Text(c, "Files") },
			func() { ui.Text(c, "Editor") },
		).Element.Fill()
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAspectRatio() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Box(c).Size(320, 240).Children(func() {
			layout.AspectRatio(c, layout.AspectRatioOptions{Ratio: 16.0 / 9}, func() {
				ui.Text(c, "16 : 9")
			})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDivider() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			layout.Divider(c, layout.DividerOptions{Ornament: true})
			layout.Divider(c, layout.DividerOptions{Label: "Chapter II"})
		})
	}
	ui.Render(view, 400, 300, 1)
}
