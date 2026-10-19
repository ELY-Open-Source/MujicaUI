package display

import (
	"runtime"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// Kbd shows keys as caps; "Mod" is Command or Ctrl.
func Kbd(c *ui.Context, keys ...string) ui.Element {
	if len(keys) == 0 {
		panic("mujica: Kbd needs at least one key")
	}
	k := core.Tokens(c)
	mac := runtime.GOOS == "darwin"
	spoken := make([]string, len(keys))
	for i, key := range keys {
		if key == "" {
			panic("mujica: Kbd key must not be empty")
		}
		_, spoken[i] = KbdKey(key, mac)
	}
	row := ui.Row(c).Gap(2).AlignItems(ui.Center).Shrink(0).Role(ui.RoleText).Label(strings.Join(spoken, "+"))
	row.Children(func() {
		for i, key := range keys {
			if i > 0 && !mac {
				ui.Text(c, "+").FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
			}
			face, _ := KbdKey(key, mac)
			ui.Box(c.Key(i)).Row().Center().MinWidth(core.FontSize(c, 20)).MinHeight(core.FontSize(c, 20)).Padding(0, 5).
				Radius(theme.ControlRadius).Background(k.Surface).BorderColor(k.ControlBorder).
				BorderWidth(1, 1, 2, 1).Children(func() {
				ui.Text(c, face).Font(theme.MonoFont).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.Text).SingleLine()
			})
		}
	})
	return row
}

// kbdNames maps keys to macOS cap, word, other cap.
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

// KbdKey returns how key shows and reads on the platform.
func KbdKey(key string, mac bool) (shown, spoken string) {
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
