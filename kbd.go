package mujica

import (
	"runtime"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
)

// Kbd shows a key or key combination as key caps. It registers nothing.
// "Mod" is Command on macOS and Ctrl elsewhere; Ctrl, Alt, Shift, Enter,
// Esc, Tab, Backspace, Delete, Up, Down, Left and Right follow the platform.
func Kbd(c *ui.Context, keys ...string) *ui.Element {
	if len(keys) == 0 {
		panic("mujica: Kbd needs at least one key")
	}
	k := Tokens(c)
	mac := runtime.GOOS == "darwin"
	spoken := make([]string, len(keys))
	for i, key := range keys {
		if key == "" {
			panic("mujica: Kbd key must not be empty")
		}
		_, spoken[i] = kbdKey(key, mac)
	}
	row := ui.Row(c).Gap(2).AlignItems(ui.Center).Shrink(0).Role(ui.RoleText).Label(strings.Join(spoken, "+"))
	row.Children(func() {
		for i, key := range keys {
			if i > 0 && !mac {
				ui.Text(c, "+").FontSize(theme.CaptionSize).TextColor(k.TextMuted)
			}
			face, _ := kbdKey(key, mac)
			ui.Box(c).Key(i).Row().Center().MinWidth(20).Height(20).Padding(0, 5).
				Radius(theme.ControlRadius).Background(k.Surface).BorderColor(k.ControlBorder).
				BorderWidth(1, 1, 2, 1).Children(func() {
				ui.Text(c, face).Font(theme.MonoFont).FontSize(theme.CaptionSize).TextColor(k.Text).SingleLine()
			})
		}
	})
	return row
}

// kbdNames maps a key to its macOS cap, macOS word and other platforms' cap.
var kbdNames = map[string]struct{ macCap, macWord, other string }{
	"mod":       {"⌘", "Command", "Ctrl"},
	"ctrl":      {"⌃", "Control", "Ctrl"},
	"alt":       {"⌥", "Option", "Alt"},
	"shift":     {"⇧", "Shift", "Shift"},
	"enter":     {"↩", "Return", "Enter"},
	"esc":       {"esc", "Escape", "Esc"},
	"tab":       {"⇥", "Tab", "Tab"},
	"backspace": {"⌫", "Delete", "Backspace"},
	"delete":    {"⌦", "Forward Delete", "Delete"},
	"up":        {"↑", "Up", "↑"},
	"down":      {"↓", "Down", "↓"},
	"left":      {"←", "Left", "←"},
	"right":     {"→", "Right", "→"},
}

// kbdKey returns how key shows and how it is read on the platform.
func kbdKey(key string, mac bool) (shown, spoken string) {
	n, ok := kbdNames[strings.ToLower(key)]
	if !ok {
		if len([]rune(key)) == 1 {
			key = strings.ToUpper(key)
		}
		return key, key
	}
	if mac {
		return n.macCap, n.macWord
	}
	word := n.other
	switch strings.ToLower(key) {
	case "up", "down", "left", "right":
		word = n.macWord
	}
	return n.other, word
}
