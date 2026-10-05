package mujica

import (
	"runtime"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// overlay draws the open menu and its submenus above the window.
func (m *menubarRun) overlay(titles []ui.Rect) {
	c, st := m.c, m.st
	ui.Overlay(c, func() {
		back := ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0)
		if back.Clicked() {
			m.close()
			return
		}
		if back.Pressed() {
			m.focus(st.open)
		}
		// Title hotspots above the backdrop switch menus as the pointer moves.
		for i, r := range titles {
			hot := ui.Box(c).Absolute().Left(r.X).Top(r.Y).Size(r.W, r.H).Cursor(ui.CursorPointer)
			if hot.Clicked() {
				if i == st.open {
					m.close()
				} else {
					m.openMenu(i, false)
				}
				return
			}
			if hot.Pressed() {
				m.focus(st.open)
			}
			if hot.Hovered() && i != st.open {
				m.openMenu(i, false)
			}
		}
		if st.open < 0 {
			return
		}
		r := titles[st.open]
		x, y := r.X, r.Y+r.H+2
		var prev ui.Rect
		for l := 0; l < len(st.hl) && st.open >= 0; l++ {
			panel, anchor := m.panel(l, x, y, prev)
			if anchor == nil || l+1 >= len(st.hl) {
				break
			}
			a := anchor.Bounds()
			prev = panel.Bounds()
			x, y = prev.X+prev.W-2, a.Y-Density(c).Unit()
		}
	})
}

// panel builds level l at (x, y) and returns it with its highlighted row.
func (m *menubarRun) panel(l int, x, y float32, parent ui.Rect) (*ui.Element, *ui.Element) {
	c, st := m.c, m.st
	k := Tokens(c)
	u := Density(c).Unit()
	items := m.items(l)
	ww, wh := c.Size()
	// Long menus scroll within the window instead of running off it.
	panel := ui.Scroll(c).Absolute().Padding(u).MinWidth(200).MaxHeight(wh-8).AlignItems(ui.Stretch).Role(ui.RolePopup).
		Radius(theme.PopupRadius).Background(k.Surface).Border(theme.BorderWidth, k.Border).
		Shadow(0, 6, 16, 0, ui.RGBA(0, 0, 0, 0.18))
	internal.PopIn(c, panel, Motion(c, theme.PopupDuration))
	pb := panel.Bounds()
	if x+pb.W > ww-4 {
		x = ww - 4 - pb.W
		if l > 0 {
			x = parent.X - pb.W + 2
		}
	}
	if y+pb.H > wh-4 {
		y = max(4, wh-4-pb.H)
	}
	panel.Left(max(4, x)).Top(y)
	if panel.Clicked() || panel.Pressed() {
		m.focus(st.open)
	}
	var anchor *ui.Element
	panel.Children(func() {
		for j, it := range items {
			if it.Kind == MenubarSeparator {
				ui.Box(c).Height(1).Margin(u, u*2).Background(k.Border).Role(ui.RoleNone)
				continue
			}
			row := m.item(it, st.hl[l] == j)
			if st.hl[l] == j {
				anchor = row
				if st.reveal && l == len(st.hl)-1 {
					row.ScrollIntoView()
					st.reveal = false
				}
			}
			if row.Clicked() {
				st.hl = append(st.hl[:l], j)
				m.activate(l, j)
				return
			}
			if row.Pressed() {
				// A press takes the focus off the title, which drives the keys.
				m.focus(st.open)
			}
			px, py, over := row.PointerPosition()
			b := row.Bounds()
			at := [2]float32{b.X + px, b.Y + py}
			if over && at != st.ptr && !it.Disabled {
				st.ptr = at
				if st.hl[l] != j || it.Items != nil && len(st.hl) == l+1 {
					st.hl = append(st.hl[:l], j)
					if it.Items != nil {
						st.hl = append(st.hl, -1)
					}
					c.Invalidate()
				}
			}
		}
	})
	return panel, anchor
}

func (m *menubarRun) item(it MenubarItem, hl bool) *ui.Element {
	c := m.c
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	name := menubarName(c, it)
	row := ui.Row(c).AlignItems(ui.Center).Gap(u*2).MinHeight(d.ControlHeight()-4).Padding(0, u*2).
		Radius(theme.ControlRadius).Label(name).TextColor(k.Text).Disabled(it.Disabled).Role(ui.RoleButton)
	if hl {
		row.Background(k.Selection)
	}
	internal.Ease(row, Motion(c, theme.StateDuration), k.Text)
	row.Children(func() {
		mark := ui.Box(c).Size(16, 16).Shrink(0).Center().Role(ui.RoleNone)
		if it.Checked && it.Kind == MenubarCheck {
			mark.TextColor(k.AccentText).Children(func() { ui.Icon(c, icons.Must("check")).Size(14, 14) })
		}
		if it.Checked && it.Kind == MenubarRadio {
			mark.Draw(func(p *ui.Painter, r ui.Rect) {
				p.FillPath(internal.Diamond(r.X+r.W/2, r.Y+r.H/2, 8), k.AccentText)
			})
		}
		ui.Text(c, it.Label).SingleLine()
		ui.Box(c).Grow(1).MinWidth(u * 6)
		if it.Shortcut.Key != 0 {
			ui.Text(c, menubarShortcutText(it.Shortcut)).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
		}
		if it.Items != nil {
			ui.Icon(c, icons.Must("chevron-right")).Size(14, 14).TextColor(k.TextMuted)
		}
	})
	return row
}

// menubarName is an item's label, with its checked state for assistive technology.
func menubarName(c *ui.Context, it MenubarItem) string {
	if (it.Kind == MenubarCheck || it.Kind == MenubarRadio) && it.Checked {
		return it.Label + ", " + msg(c, "menubar.checked")
	}
	return it.Label
}

var menubarKeyNames = map[ui.Key]string{
	ui.KeyEnter: "Enter", ui.KeyEscape: "Esc", ui.KeyBackspace: "Backspace", ui.KeyTab: "Tab", ui.KeySpace: "Space",
	ui.KeyDelete: "Del", ui.KeyHome: "Home", ui.KeyEnd: "End", ui.KeyPageUp: "PgUp", ui.KeyPageDown: "PgDn",
	ui.KeyLeft: "←", ui.KeyRight: "→", ui.KeyUp: "↑", ui.KeyDown: "↓", ui.KeyMinus: "-", ui.KeyEqual: "=",
	ui.KeyComma: ",", ui.KeyPeriod: ".", ui.KeySlash: "/", ui.KeySemicolon: ";", ui.KeyQuote: "'",
	ui.KeyF1: "F1", ui.KeyF2: "F2", ui.KeyF3: "F3", ui.KeyF4: "F4", ui.KeyF5: "F5", ui.KeyF6: "F6",
	ui.KeyF7: "F7", ui.KeyF8: "F8", ui.KeyF9: "F9", ui.KeyF10: "F10", ui.KeyF11: "F11", ui.KeyF12: "F12",
	ui.KeyBracketLeft: "[", ui.KeyBracketRight: "]", ui.KeyBackslash: "\\", ui.KeyBackquote: "`",
}

// menubarShortcutText spells a shortcut the platform's way: ⌘S on macOS, Ctrl+S elsewhere.
func menubarShortcutText(s MenubarShortcut) string {
	key, ok := menubarKeyNames[s.Key]
	switch {
	case ok:
	case s.Key >= ui.KeyA && s.Key <= ui.KeyZ:
		key = string(rune('A' + s.Key - ui.KeyA))
	case s.Key >= ui.Key0 && s.Key <= ui.Key9:
		key = string(rune('0' + s.Key - ui.Key0))
	default:
		panic("mujica: MenubarShortcut key has no name")
	}
	mac := runtime.GOOS == "darwin"
	var parts []string
	for _, m := range []struct {
		mod       ui.Modifiers
		mac, name string
	}{{ui.Ctrl, "⌃", "Ctrl"}, {ui.Alt, "⌥", "Alt"}, {ui.Shift, "⇧", "Shift"}, {ui.Super, "⌘", "Win"}} {
		if s.Mods&m.mod != 0 {
			if mac {
				parts = append(parts, m.mac)
			} else {
				parts = append(parts, m.name)
			}
		}
	}
	if mac {
		return strings.Join(parts, "") + key
	}
	return strings.Join(append(parts, key), "+")
}
