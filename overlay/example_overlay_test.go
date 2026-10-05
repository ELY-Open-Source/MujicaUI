package overlay_test

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/overlay"
)

func ExampleDialog() {
	open := true
	name := "Ophelia"
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Button(c, "Rename…", input.ButtonOptions{}).Clicked() {
			open = true
		}
		overlay.Dialog(c, &open, overlay.DialogOptions{
			Title:       "Rename the chamber",
			Description: "The new name shows in every hall.",
			Decorated:   true,
			Actions: func() {
				if input.Button(c, "Cancel", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
					open = false
				}
				if input.Button(c, "Rename", input.ButtonOptions{}).Clicked() {
					open = false
				}
			},
		}, func() {
			ui.TextInput(c, &name)
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleAlertDialog() {
	open := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Button(c, "Delete archive", input.ButtonOptions{Variant: input.Danger}).Clicked() {
			open = true
		}
		r := overlay.AlertDialog(c, &open, overlay.AlertDialogOptions{
			Title:       "Delete the archive?",
			Description: "Twelve letters will be destroyed.",
			ConfirmText: "archive",
		})
		if r.Confirmed {
			fmt.Println("deleted")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDrawer() {
	open := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		if input.Button(c, "Filters", input.ButtonOptions{}).Clicked() {
			open = true
		}
		overlay.Drawer(c, &open, overlay.DrawerOptions{Title: "Filters", Side: overlay.DrawerRight, Size: 320}, func() {
			ui.Text(c, "Guests are placed by rank.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePopover() {
	open := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		b := input.Button(c, "Guests", input.ButtonOptions{Variant: input.Secondary})
		if b.Clicked() {
			open = !open
		}
		overlay.Popover(c, b, &open, overlay.PopoverOptions{Title: "Guest list"}, func() {
			ui.Text(c, "Twelve guests confirmed.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTooltip() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		b := input.Button(c, "Seal", input.ButtonOptions{})
		overlay.Tooltip(c, b, "Seal the letter with wax", overlay.TooltipOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleHoverCard() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		b := input.Button(c, "Lady Sakiko", input.ButtonOptions{Variant: input.Ghost})
		overlay.HoverCard(c, b, overlay.HoverCardOptions{}, func() {
			ui.Text(c, "Keyboardist of the court.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDropdownMenu() {
	pinned := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		overlay.DropdownMenu(c, "Actions", overlay.DropdownMenuOptions{}, func(m *overlay.PopupMenu) {
			m.Group("Letter", func() {
				if m.Item("Rename", overlay.PopupMenuItemOptions{}) {
					fmt.Println("rename")
				}
				m.Check("Pinned", &pinned, overlay.PopupMenuItemOptions{KeepOpen: true})
			})
			m.Submenu("Move to", overlay.PopupMenuItemOptions{}, func(m *overlay.PopupMenu) {
				m.Item("Archive", overlay.PopupMenuItemOptions{})
			})
			m.Separator()
			m.Item("Burn", overlay.PopupMenuItemOptions{Disabled: true})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleContextMenu() {
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		row := ui.Row(c).Padding(12).Children(func() {
			ui.Text(c, "Letter to the Countess")
		})
		overlay.ContextMenu(c, row, overlay.ContextMenuOptions{}, func(m *overlay.PopupMenu) {
			if m.Item("Reply", overlay.PopupMenuItemOptions{}) {
				fmt.Println("reply")
			}
			m.Item("Delete", overlay.PopupMenuItemOptions{})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePopconfirm() {
	open := false
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		b := input.Button(c, "Archive", input.ButtonOptions{Variant: input.Secondary})
		if b.Clicked() {
			open = true
		}
		if overlay.Popconfirm(c, b, &open, overlay.PopconfirmOptions{Title: "Archive this letter?"}).Confirmed {
			fmt.Println("archived")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOverlay() {
	open := true
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		overlay.Overlay(c, &open, overlay.OverlayOptions{Place: overlay.OverlayBottomRight, Label: "Now playing"}, func() {
			ui.Text(c, "Symbol I: △")
		})
	}
	ui.Render(view, 400, 300, 1)
}
