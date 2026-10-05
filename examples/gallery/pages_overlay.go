package main

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/overlay"
)

const overlayCategory = "Overlays"

var overlayPages = []page{
	{ID: "081", Name: "Dialog", ZhName: "对话框", Category: overlayCategory, Demo: dialogDemo, Code: dialogCode,
		Keys: "Tab / Shift+Tab 在对话框内循环；Escape 关闭（可关闭此规则）；关闭后焦点回到触发按钮。"},
	{ID: "082", Name: "AlertDialog", ZhName: "警示对话框", Category: overlayCategory, Demo: alertDialogDemo, Code: alertDialogCode,
		Keys: "打开时焦点在取消；Enter 只触发聚焦的按钮，输入框中的 Enter 不确认；Escape 等同取消；点击遮罩无效。"},
	{ID: "083", Name: "Drawer", ZhName: "抽屉", Category: overlayCategory, Demo: drawerDemo, Code: drawerCode,
		Keys: "模态：焦点限制在抽屉内并在关闭后恢复；Escape 关闭。非模态：窗口仍可操作，Escape 关闭。"},
	{ID: "084", Name: "Popover", ZhName: "弹出框", Category: overlayCategory, Demo: popoverDemo, Code: popoverCode,
		Keys: "Tab 从锚点进入弹层；Escape 或外部点击关闭并把焦点还给锚点。"},
	{ID: "085", Name: "Tooltip", ZhName: "文字提示", Category: overlayCategory, Demo: tooltipDemo, Code: tooltipCode,
		Keys: "键盘聚焦或悬停超过延迟后显示；Escape 隐藏；移开或失焦隐藏。"},
	{ID: "086", Name: "HoverCard", ZhName: "悬停卡片", Category: overlayCategory, Demo: hoverCardDemo, Code: hoverCardCode,
		Keys: "键盘聚焦锚点后显示，Tab 进入卡片；Escape 关闭。"},
	{ID: "087", Name: "DropdownMenu", ZhName: "下拉菜单", Category: overlayCategory, Demo: dropdownMenuDemo, Code: dropdownMenuCode,
		Keys: "Enter / Space / ↓ 打开；↑↓ Home End 移动；→ 打开子菜单，← 关闭；Enter 选择；Escape 关闭最上层；Tab 在任一层关闭整个菜单。"},
	{ID: "088", Name: "ContextMenu", ZhName: "右键菜单", Category: overlayCategory, Demo: contextMenuDemo, Code: contextMenuCode,
		Keys: "右键在指针处打开；焦点在区域内时按 Menu 键或 Shift+F10 在区域下方打开；其余按键同下拉菜单。"},
	{ID: "089", Name: "Popconfirm", ZhName: "气泡确认", Category: overlayCategory, Demo: popconfirmDemo, Code: popconfirmCode,
		Keys: "打开后焦点在确认（Danger 时在取消）；Enter 触发聚焦的按钮；Escape 或外部点击取消。"},
	{ID: "090", Name: "Overlay", ZhName: "覆盖层", Category: overlayCategory, Demo: overlayDemo, Code: overlayCode,
		Keys: "Escape 关闭；模态时焦点限制在覆盖层内，点击遮罩关闭。"},
}

// overlayControls lays out a demo's parameter row.
func overlayControls(c *ui.Context, fn func()) {
	ui.Row(c).Gap(16).Wrap().AlignItems(ui.Center).Padding(0, 0, 12, 0).Children(fn)
}

func overlayCheck(c *ui.Context, log func(string), v *bool, label string) {
	if ui.Checkbox(c, v, label).Changed() {
		log(fmt.Sprintf("%s = %v", label, *v))
	}
}

const dialogCode = `open := true
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
ui.Render(view, 400, 300, 1)`

type dialogDemoState struct {
	open, decorated, keepBackdrop, keepEscape, hideClose bool
	name                                                 string
}

func dialogDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "081", func() dialogDemoState { return dialogDemoState{decorated: true, name: "Ophelia"} })
	overlayControls(c, func() {
		overlayCheck(c, log, &s.decorated, "Decorated")
		overlayCheck(c, log, &s.keepBackdrop, "KeepOnBackdrop")
		overlayCheck(c, log, &s.keepEscape, "KeepOnEscape")
		overlayCheck(c, log, &s.hideClose, "HideClose")
	})
	if input.Button(c, "Rename…", input.ButtonOptions{}).Clicked() {
		s.open = true
		log("open")
	}
	was := s.open
	overlay.Dialog(c, &s.open, overlay.DialogOptions{
		Title: "Rename the chamber", Description: "The new name shows in every hall.", Decorated: s.decorated,
		KeepOnBackdrop: s.keepBackdrop, KeepOnEscape: s.keepEscape, HideClose: s.hideClose,
		Actions: func() {
			if input.Button(c, "Cancel", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				s.open = false
			}
			if input.Button(c, "Rename", input.ButtonOptions{}).Clicked() {
				log("renamed to " + s.name)
				s.open = false
			}
		},
	}, func() { ui.TextInput(c, &s.name) })
	if was && !s.open {
		log("closed")
	}
}

const alertDialogCode = `open := true
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
ui.Render(view, 400, 300, 1)`

type alertDialogDemoState struct{ open, typed bool }

func alertDialogDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "082", func() alertDialogDemoState { return alertDialogDemoState{} })
	overlayControls(c, func() { overlayCheck(c, log, &s.typed, "ConfirmText") })
	if input.Button(c, "Delete archive", input.ButtonOptions{Variant: input.Danger}).Clicked() {
		s.open = true
		log("open")
	}
	opts := overlay.AlertDialogOptions{Title: "Delete the archive?", Description: "Twelve letters will be destroyed. This cannot be undone."}
	if s.typed {
		opts.ConfirmText = "archive"
	}
	r := overlay.AlertDialog(c, &s.open, opts)
	if r.Confirmed {
		log("confirmed")
	}
	if r.Canceled {
		log("canceled")
	}
}

const drawerCode = `open := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if input.Button(c, "Filters", input.ButtonOptions{}).Clicked() {
		open = true
	}
	overlay.Drawer(c, &open, overlay.DrawerOptions{Title: "Filters", Side: overlay.DrawerRight, Size: 320}, func() {
		ui.Text(c, "Guests are placed by rank.")
	})
}
ui.Render(view, 400, 300, 1)`

type drawerDemoState struct {
	open, nonModal bool
	side           string
	size           float64
}

func drawerDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "083", func() drawerDemoState { return drawerDemoState{side: "Right", size: 320} })
	sides := map[string]overlay.DrawerSide{"Right": overlay.DrawerRight, "Left": overlay.DrawerLeft, "Top": overlay.DrawerTop, "Bottom": overlay.DrawerBottom}
	overlayControls(c, func() {
		if ui.Select(c, &s.side, []string{"Right", "Left", "Top", "Bottom"}).Changed() {
			log("side = " + s.side)
		}
		overlayCheck(c, log, &s.nonModal, "NonModal")
		overlaySlider(c, log, &s.size, "Size", "DIP", 120, 900, 20)
	})
	if input.Button(c, "Open drawer", input.ButtonOptions{}).Clicked() {
		s.open = true
		log("open")
	}
	was := s.open
	overlay.Drawer(c, &s.open, overlay.DrawerOptions{Title: "Seating chart", Side: sides[s.side], Size: float32(s.size), NonModal: s.nonModal,
		Actions: func() {
			if input.Button(c, "Apply", input.ButtonOptions{}).Clicked() {
				log("applied")
				s.open = false
			}
		}}, func() {
		ui.Text(c, "Guests are placed by rank, then by arrival.").TextColor(core.Tokens(c).TextMuted)
	})
	if was && !s.open {
		log("closed")
	}
}

const popoverCode = `open := true
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
ui.Render(view, 400, 300, 1)`

type popoverDemoState struct {
	open, titled bool
	rsvp         bool
}

func popoverDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "084", func() popoverDemoState { return popoverDemoState{titled: true} })
	overlayControls(c, func() { overlayCheck(c, log, &s.titled, "Title") })
	b := input.Button(c, "Guests", input.ButtonOptions{Variant: input.Secondary})
	if b.Clicked() {
		s.open = !s.open
		log(fmt.Sprintf("open = %v", s.open))
	}
	opts := overlay.PopoverOptions{Label: "Guest list"}
	if s.titled {
		opts.Title = "Guest list"
	}
	was := s.open
	overlay.Popover(c, b, &s.open, opts, func() {
		ui.Text(c, "Twelve guests confirmed; three await a reply.")
		if ui.Checkbox(c, &s.rsvp, "Remind the others").Changed() {
			log(fmt.Sprintf("remind = %v", s.rsvp))
		}
	})
	if was && !s.open {
		log("closed")
	}
}

const tooltipCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	b := input.Button(c, "Seal", input.ButtonOptions{})
	overlay.Tooltip(c, b, "Seal the letter with wax", overlay.TooltipOptions{})
}
ui.Render(view, 400, 300, 1)`

func tooltipDemo(c *ui.Context, log func(string)) {
	delay := ui.Local(c.Root(), "085", func() float64 { return 500 })
	overlayControls(c, func() {
		overlaySlider(c, log, delay, "Delay", "ms", 100, 1500, 100)
	})
	opts := overlay.TooltipOptions{Delay: time.Duration(*delay) * time.Millisecond}
	ui.Row(c).Gap(12).Children(func() {
		b := input.Button(c, "Seal", input.ButtonOptions{})
		if b.Clicked() {
			log("seal clicked")
		}
		overlay.Tooltip(c, b, "Seal the letter with wax", opts)
		d := input.Button(c, "Burn", input.ButtonOptions{Variant: input.Danger})
		if d.Clicked() {
			log("burn clicked")
		}
		overlay.Tooltip(c, d, "Burn every copy of the letter; nothing is kept", opts)
	})
}

const hoverCardCode = `view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	b := input.Button(c, "Lady Sakiko", input.ButtonOptions{Variant: input.Ghost})
	overlay.HoverCard(c, b, overlay.HoverCardOptions{}, func() {
		ui.Text(c, "Keyboardist of the court.")
	})
}
ui.Render(view, 400, 300, 1)`

type hoverCardDemoState struct {
	shown                        bool
	openDelay, closeDelay, width float64
}

func hoverCardDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "086", func() hoverCardDemoState {
		return hoverCardDemoState{openDelay: 500, closeDelay: 300, width: 280}
	})
	overlayControls(c, func() {
		overlaySlider(c, log, &s.openDelay, "OpenDelay", "ms", 100, 1500, 100)
		overlaySlider(c, log, &s.closeDelay, "CloseDelay", "ms", 100, 1500, 100)
		overlaySlider(c, log, &s.width, "Width", "DIP", 200, 420, 20)
	})
	b := input.Button(c, "Lady Sakiko", input.ButtonOptions{Variant: input.Ghost})
	card := overlay.HoverCard(c, b, overlay.HoverCardOptions{
		OpenDelay:  time.Duration(s.openDelay) * time.Millisecond,
		CloseDelay: time.Duration(s.closeDelay) * time.Millisecond,
		Width:      float32(s.width),
	}, func() {
		ui.Text(c, "Lady Sakiko").FontWeight(600)
		ui.Text(c, "Keyboardist of the court; writes the masquerade's score.").TextColor(core.Tokens(c).TextMuted)
		if input.Button(c, "Follow", input.ButtonOptions{Variant: input.Outline}).Clicked() {
			log("follow")
		}
	})
	if (card != nil) != s.shown {
		s.shown = card != nil
		log(fmt.Sprintf("card shown = %v", s.shown))
	}
}

// overlaySlider is a labeled step slider logging its value.
func overlaySlider(c *ui.Context, log func(string), v *float64, name, unit string, lo, hi, step float64) {
	ui.Text(c, fmt.Sprintf("%s %.0f %s", name, *v, unit))
	if ui.StepSlider(c, v, lo, hi, step).Width(140).Changed() {
		log(fmt.Sprintf("%s = %.0f %s", name, *v, unit))
	}
}

const dropdownMenuCode = `pinned := false
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
ui.Render(view, 400, 300, 1)`

type menuDemoState struct{ pinned, sealed, locked, disabled bool }

// overlayMenu fills the shared demo menu, logging every choice.
func overlayMenu(s *menuDemoState, log func(string)) func(m *overlay.PopupMenu) {
	return func(m *overlay.PopupMenu) {
		m.Group("Letter", func() {
			if m.Item("Rename", overlay.PopupMenuItemOptions{}) {
				log("rename")
			}
			if m.Check("Pinned", &s.pinned, overlay.PopupMenuItemOptions{KeepOpen: true}) {
				log(fmt.Sprintf("pinned = %v", s.pinned))
			}
			if m.Check("Sealed", &s.sealed, overlay.PopupMenuItemOptions{}) {
				log(fmt.Sprintf("sealed = %v", s.sealed))
			}
		})
		m.Group("Move", func() {
			m.Submenu("Move to", overlay.PopupMenuItemOptions{}, func(m *overlay.PopupMenu) {
				for _, f := range []string{"Archive", "Drafts", "Vault"} {
					if m.Item(f, overlay.PopupMenuItemOptions{}) {
						log("move to " + f)
					}
				}
			})
		})
		m.Separator()
		if m.Item("Burn", overlay.PopupMenuItemOptions{Disabled: s.locked}) {
			log("burn")
		}
	}
}

func dropdownMenuDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "087", func() menuDemoState { return menuDemoState{locked: true} })
	overlayControls(c, func() {
		overlayCheck(c, log, &s.locked, "Burn disabled")
		overlayCheck(c, log, &s.disabled, "Disabled")
	})
	overlay.DropdownMenu(c, "Actions", overlay.DropdownMenuOptions{Disabled: s.disabled}, overlayMenu(s, log))
}

const contextMenuCode = `view := func(c *ui.Context) {
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
ui.Render(view, 400, 300, 1)`

func contextMenuDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "088", func() menuDemoState { return menuDemoState{} })
	overlayControls(c, func() {
		overlayCheck(c, log, &s.locked, "Burn disabled")
		overlayCheck(c, log, &s.disabled, "Disabled")
	})
	k := core.Tokens(c)
	row := ui.Column(c).Height(160).FillWidth().Padding(16).Gap(8).Background(k.Surface).Border(1, k.Border).Radius(4)
	row.Children(func() {
		ui.Text(c, "Right-click anywhere here, or focus the button and press Shift+F10.").TextColor(k.TextMuted)
		if input.Button(c, "Letter to the Countess", input.ButtonOptions{Variant: input.Ghost}).Clicked() {
			log("letter clicked")
		}
	})
	overlay.ContextMenu(c, row, overlay.ContextMenuOptions{Disabled: s.disabled}, overlayMenu(s, log))
}

const popconfirmCode = `open := false
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
ui.Render(view, 400, 300, 1)`

type popconfirmDemoState struct{ open, danger bool }

func popconfirmDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "089", func() popconfirmDemoState { return popconfirmDemoState{} })
	overlayControls(c, func() { overlayCheck(c, log, &s.danger, "Danger") })
	b := input.Button(c, "Archive", input.ButtonOptions{Variant: input.Secondary})
	if b.Clicked() {
		s.open = true
		log("open")
	}
	r := overlay.Popconfirm(c, b, &s.open, overlay.PopconfirmOptions{Title: "Archive this letter?", Description: "It leaves the inbox.", Danger: s.danger})
	if r.Confirmed {
		log("confirmed")
	}
	if r.Canceled {
		log("canceled")
	}
}

const overlayCode = `open := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	overlay.Overlay(c, &open, overlay.OverlayOptions{Place: overlay.OverlayBottomRight, Label: "Now playing"}, func() {
		ui.Text(c, "Symbol I: △")
	})
}
ui.Render(view, 400, 300, 1)`

type overlayDemoState struct {
	open, modal, dim bool
	place            string
}

func overlayDemo(c *ui.Context, log func(string)) {
	s := ui.Local(c.Root(), "090", func() overlayDemoState { return overlayDemoState{place: "BottomRight"} })
	names := []string{"Center", "Top", "Bottom", "Left", "Right", "TopLeft", "TopRight", "BottomLeft", "BottomRight"}
	place := overlay.OverlayCenter
	for i, n := range names {
		if n == s.place {
			place = overlay.OverlayPlace(i)
		}
	}
	overlayControls(c, func() {
		if ui.Select(c, &s.place, names).Changed() {
			log("place = " + s.place)
		}
		overlayCheck(c, log, &s.modal, "Modal")
		if s.modal {
			overlayCheck(c, log, &s.dim, "Dim")
		}
	})
	if input.Button(c, "Show overlay", input.ButtonOptions{}).Clicked() {
		s.open = true
		log("open")
	}
	was := s.open
	overlay.Overlay(c, &s.open, overlay.OverlayOptions{Place: place, Modal: s.modal, Dim: s.modal && s.dim, Label: "Now playing"}, func() {
		ui.Text(c, "Now playing").FontWeight(600)
		ui.Text(c, "Symbol I: △").TextColor(core.Tokens(c).TextMuted)
		if input.Button(c, "Close", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
			s.open = false
		}
	})
	if was && !s.open {
		log("closed")
	}
}
