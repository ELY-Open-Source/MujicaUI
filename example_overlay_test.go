package mujica_test

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	mujica "mujicaui"
)

func ExampleDialog() {
	open := true
	name := "Ophelia"
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		if mujica.Button(c, "Rename…", mujica.ButtonOptions{}).Clicked() {
			open = true
		}
		mujica.Dialog(c, &open, mujica.DialogOptions{
			Title:       "Rename the chamber",
			Description: "The new name shows in every hall.",
			Decorated:   true,
			Actions: func() {
				if mujica.Button(c, "Cancel", mujica.ButtonOptions{Variant: mujica.Secondary}).Clicked() {
					open = false
				}
				if mujica.Button(c, "Rename", mujica.ButtonOptions{}).Clicked() {
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
		mujica.Use(c, mujica.Settings{})
		if mujica.Button(c, "Delete archive", mujica.ButtonOptions{Variant: mujica.Danger}).Clicked() {
			open = true
		}
		r := mujica.AlertDialog(c, &open, mujica.AlertDialogOptions{
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
		mujica.Use(c, mujica.Settings{})
		if mujica.Button(c, "Filters", mujica.ButtonOptions{}).Clicked() {
			open = true
		}
		mujica.Drawer(c, &open, mujica.DrawerOptions{Title: "Filters", Side: mujica.DrawerRight, Size: 320}, func() {
			ui.Text(c, "Guests are placed by rank.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePopover() {
	open := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		b := mujica.Button(c, "Guests", mujica.ButtonOptions{Variant: mujica.Secondary})
		if b.Clicked() {
			open = !open
		}
		mujica.Popover(c, b, &open, mujica.PopoverOptions{Title: "Guest list"}, func() {
			ui.Text(c, "Twelve guests confirmed.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleTooltip() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		b := mujica.Button(c, "Seal", mujica.ButtonOptions{})
		mujica.Tooltip(c, b, "Seal the letter with wax", mujica.TooltipOptions{})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleHoverCard() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		b := mujica.Button(c, "Lady Sakiko", mujica.ButtonOptions{Variant: mujica.Ghost})
		mujica.HoverCard(c, b, mujica.HoverCardOptions{}, func() {
			ui.Text(c, "Keyboardist of the court.")
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleDropdownMenu() {
	pinned := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DropdownMenu(c, "Actions", mujica.DropdownMenuOptions{}, func(m *mujica.PopupMenu) {
			m.Group("Letter", func() {
				if m.Item("Rename", mujica.PopupMenuItemOptions{}) {
					fmt.Println("rename")
				}
				m.Check("Pinned", &pinned, mujica.PopupMenuItemOptions{KeepOpen: true})
			})
			m.Submenu("Move to", mujica.PopupMenuItemOptions{}, func(m *mujica.PopupMenu) {
				m.Item("Archive", mujica.PopupMenuItemOptions{})
			})
			m.Separator()
			m.Item("Burn", mujica.PopupMenuItemOptions{Disabled: true})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleContextMenu() {
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		row := ui.Row(c).Padding(12).Children(func() {
			ui.Text(c, "Letter to the Countess")
		})
		mujica.ContextMenu(c, row, mujica.ContextMenuOptions{}, func(m *mujica.PopupMenu) {
			if m.Item("Reply", mujica.PopupMenuItemOptions{}) {
				fmt.Println("reply")
			}
			m.Item("Delete", mujica.PopupMenuItemOptions{})
		})
	}
	ui.Render(view, 400, 300, 1)
}

func ExamplePopconfirm() {
	open := false
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		b := mujica.Button(c, "Archive", mujica.ButtonOptions{Variant: mujica.Secondary})
		if b.Clicked() {
			open = true
		}
		if mujica.Popconfirm(c, b, &open, mujica.PopconfirmOptions{Title: "Archive this letter?"}).Confirmed {
			fmt.Println("archived")
		}
	}
	ui.Render(view, 400, 300, 1)
}

func ExampleOverlay() {
	open := true
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.Overlay(c, &open, mujica.OverlayOptions{Place: mujica.OverlayBottomRight, Label: "Now playing"}, func() {
			ui.Text(c, "Symbol I: △")
		})
	}
	ui.Render(view, 400, 300, 1)
}
