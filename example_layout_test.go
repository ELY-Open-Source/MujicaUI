package mujica_test

import (
	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
	"mujicaui/icons"
)

func ExampleAppShell() {
	var collapsed bool
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.AppShell(c, &collapsed, mujica.AppShellOptions{
			TitleBar: func() { mujica.TitleBar(c, "Court Records", mujica.TitleBarOptions{}) },
			Sidebar:  func() { ui.Text(c, "Chronicles") },
			Content:  func() { ui.Text(c, "The Grand Hall") },
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTitleBar() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TitleBar(c, "Court Records", mujica.TitleBarOptions{
			Trailing: func() { mujica.Button(c, "Share", mujica.ButtonOptions{Variant: mujica.Ghost}) },
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStatusBar() {
	opened := 0
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		bar := mujica.StatusBar(c, mujica.StatusBarOptions{
			Left:  []mujica.StatusItem{{ID: "sync", Text: "Syncing", ShowProgress: true, Progress: 0.4}},
			Right: []mujica.StatusItem{{ID: "bell", Label: "Notifications", Icon: icons.Must("bell"), Action: true}},
		})
		if bar.Clicked == "bell" {
			opened++
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleContainer() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Container(c, mujica.ContainerOptions{MaxWidth: 320, Center: true}, func() {
			ui.Text(c, "A centered column of prose.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleStack() {
	showDraft := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Stack(c, mujica.StackOptions{Horizontal: true, Gap: mujica.SpaceM}, func() {
			mujica.Button(c, "Save", mujica.ButtonOptions{})
			mujica.StackItem(!showDraft, func() { mujica.Button(c, "Draft", mujica.ButtonOptions{Variant: mujica.Secondary}) })
			mujica.Button(c, "Cancel", mujica.ButtonOptions{Variant: mujica.Ghost})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleGrid() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Grid(c, mujica.GridOptions{Columns: []ui.Track{ui.Fixed(96), ui.Fr(1)}}, func() {
			mujica.GridCell(c, mujica.GridCellOptions{ColumnSpan: -1}, func() { ui.Text(c, "Header") })
			ui.Text(c, "Name")
			ui.Text(c, "Ada")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleScrollArea() {
	var offset ui.ScrollState
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.ScrollArea(c, &offset, mujica.ScrollAreaOptions{Label: "Log"}, func() {
			for range 50 {
				ui.Text(c, "A line of the log")
			}
		}).Height(200)
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleSplitPane() {
	split := mujica.SplitPaneState{Size: 160}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.SplitPane(c, &split, mujica.SplitPaneOptions{MinFirst: 120, Collapsible: mujica.SplitPaneFirst},
			func() { ui.Text(c, "Files") },
			func() { ui.Text(c, "Editor") },
		).Element.Fill()
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAspectRatio() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		ui.Box(c).Size(320, 240).Children(func() {
			mujica.AspectRatio(c, mujica.AspectRatioOptions{Ratio: 16.0 / 9}, func() {
				ui.Text(c, "16 : 9")
			})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDivider() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			mujica.Divider(c, mujica.DividerOptions{Ornament: true})
			mujica.Divider(c, mujica.DividerOptions{Label: "Chapter II"})
		})
	}
	ui.Render(view, 400, 300, 1)
}
