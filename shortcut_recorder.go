package mujica

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/internal"
	"mujicaui/theme"
)

// Shortcut combines modifiers and a key; zero means unassigned.
type Shortcut struct {
	Mods ui.Modifiers
	Key  ui.Key
}

// IsZero reports an empty shortcut.
func (s Shortcut) IsZero() bool { return s == Shortcut{} }

// keys returns platform-ordered key names.
func (s Shortcut) keys() []string {
	var out []string
	mac := runtime.GOOS == "darwin"
	add := func(m ui.Modifiers, name string) {
		if s.Mods&m != 0 {
			out = append(out, name)
		}
	}
	if mac {
		add(ui.Ctrl, "ctrl")
		add(ui.Alt, "alt")
		add(ui.Shift, "shift")
		add(ui.Super, "mod")
	} else {
		add(ui.Ctrl, "ctrl")
		add(ui.Super, "Win")
		add(ui.Alt, "alt")
		add(ui.Shift, "shift")
	}
	return append(out, shortcutKeyName(s.Key))
}

// String spells s as keys joined by "+", as "ctrl+shift+K".
func (s Shortcut) String() string {
	if s.IsZero() {
		return ""
	}
	return strings.Join(s.keys(), "+")
}

// shortcutKeyName is the Kbd name of key.
func shortcutKeyName(k ui.Key) string {
	switch {
	case k >= ui.KeyA && k <= ui.KeyZ:
		return string(rune('A' + k - ui.KeyA))
	case k >= ui.Key0 && k <= ui.Key9:
		return string(rune('0' + k - ui.Key0))
	case k >= ui.KeyF1 && k <= ui.KeyF12:
		return fmt.Sprintf("F%d", k-ui.KeyF1+1)
	}
	if n, ok := shortcutNames[k]; ok {
		return n
	}
	panic(fmt.Sprintf("mujica: Shortcut key %d has no name", k))
}

var shortcutNames = map[ui.Key]string{
	ui.KeyEnter: "enter", ui.KeyEscape: "esc", ui.KeyBackspace: "backspace", ui.KeyTab: "tab", ui.KeySpace: "Space",
	ui.KeyDelete: "delete", ui.KeyHome: "Home", ui.KeyEnd: "End", ui.KeyPageUp: "PgUp", ui.KeyPageDown: "PgDn",
	ui.KeyLeft: "left", ui.KeyRight: "right", ui.KeyUp: "up", ui.KeyDown: "down", ui.KeyMinus: "-", ui.KeyEqual: "=",
	ui.KeyComma: ",", ui.KeyPeriod: ".", ui.KeySlash: "/", ui.KeySemicolon: ";", ui.KeyQuote: "'",
	ui.KeyBracketLeft: "[", ui.KeyBracketRight: "]", ui.KeyBackslash: "\\", ui.KeyBackquote: "`",
	ui.KeyInsert: "Insert", ui.KeyContextMenu: "Menu", ui.KeyBack: "Back", ui.KeyForward: "Forward",
}

// shortcutKbd displays keycaps or a muted unassigned mark.
func shortcutKbd(c *ui.Context, s Shortcut) {
	if s.IsZero() {
		ui.Text(c, "—").TextColor(Tokens(c).TextMuted)
		return
	}
	Kbd(c, s.keys()...)
}

// ShortcutBinding names an action that already owns a shortcut.
type ShortcutBinding struct {
	Action   string
	Shortcut Shortcut
}

// ShortcutRecorderOptions configure a ShortcutRecorder.
type ShortcutRecorderOptions struct {
	// Taken lists occupied shortcuts; conflicts preserve the previous value.
	Taken []ShortcutBinding
	// Label names the recorder, as the action it sets.
	Label    string
	Disabled bool
}

type shortcutRecorderState struct {
	recording, needMod bool
	caught             *Shortcut
	conflict           string
}

// ShortcutRecorder records key combinations; Escape cancels.
func ShortcutRecorder(c *ui.Context, value *Shortcut, opts ShortcutRecorderOptions) ChoiceResult {
	if opts.Label == "" {
		panic("mujica: ShortcutRecorder needs a Label")
	}
	d := Density(c)
	var r ChoiceResult
	r.Element = ui.Column(c).Gap(d.Unit()).AlignItems(ui.Start).Disabled(opts.Disabled)
	r.Element.Children(func() {
		b := ui.ButtonBase(c).MinHeight(d.ControlHeight()).MinWidth(160).Padding(0, d.Unit()*3).Gap(d.Unit() * 2).
			Radius(theme.ControlRadius).Justify(ui.SpaceBetween)
		st := ui.Local(b, "rec", func() shortcutRecorderState { return shortcutRecorderState{} })
		disabled := b.IsDisabled()
		if disabled || !b.Focused() {
			st.recording, st.caught = false, nil
		}
		if st.caught != nil {
			s := *st.caught
			st.caught, st.recording, st.conflict = nil, false, ""
			for _, t := range opts.Taken {
				if t.Shortcut == s && s != *value {
					st.conflict = fmt.Sprintf(msg(c, "shortcut_recorder.conflict"), s.String(), t.Action)
				}
			}
			if st.conflict == "" && s != *value {
				*value, r.changed = s, true
				slog.Debug("mujica: shortcut_recorder set", "shortcut", s.String())
			}
		}
		if b.Clicked() && !disabled {
			st.recording, st.conflict, st.needMod = true, "", false
			slog.Debug("mujica: shortcut_recorder recording")
		}
		if st.recording {
			b.HandleInput(func(ev ui.InputEvent) bool {
				if b.IsDisabled() {
					return false
				}
				return shortcutRecorderKey(st, ev)
			})
		}
		shortcutRecorderLook(c, b, st, *value, opts.Label)
		shortcutRecorderNote(c, st)
	})
	return r
}

// shortcutRecorderKey captures keys during recording.
func shortcutRecorderKey(st *shortcutRecorderState, ev ui.InputEvent) bool {
	switch {
	case ev.Kind == ui.InputText:
		return true
	case ev.Kind != ui.InputKeyDown || ev.Key == ui.KeyUnknown:
		return false
	case ev.Key == ui.KeyTab && ev.Mods&^ui.Shift == 0:
		return false
	case ev.Key == ui.KeyEscape && ev.Mods == 0:
		st.recording = false
	case ev.Mods&^ui.Shift == 0 && !(ev.Key >= ui.KeyF1 && ev.Key <= ui.KeyF12):
		st.needMod = true
	default:
		s := Shortcut{ev.Mods, ev.Key}
		st.caught = &s
	}
	return true
}

// shortcutRecorderLook displays the shortcut or recording prompt.
func shortcutRecorderLook(c *ui.Context, b *ui.Element, st *shortcutRecorderState, value Shortcut, name string) {
	k := Tokens(c)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	label := name + ": " + value.String()
	if value.IsZero() {
		label = name + ": " + msg(c, "shortcut_recorder.none")
	}
	b.Label(label)
	bg, border := k.Surface, k.ControlBorder
	if st.recording {
		bg, border = k.Selection, k.AccentText
	} else if b.Hovered() {
		bg = k.SurfaceHover
	}
	b.Background(bg).Border(theme.BorderWidth, border)
	b.Children(func() {
		if st.recording {
			ui.Text(c, msg(c, "shortcut_recorder.press")).TextColor(k.AccentText).SingleLine()
		} else {
			shortcutKbd(c, value)
		}
		ui.Icon(c, icons.Must("command")).Size(14, 14).TextColor(k.TextMuted)
	})
}

// shortcutRecorderNote displays conflict and modifier guidance.
func shortcutRecorderNote(c *ui.Context, st *shortcutRecorderState) {
	text := st.conflict
	if st.needMod && st.recording {
		text = msg(c, "shortcut_recorder.needs_modifier")
	}
	if text == "" {
		return
	}
	ui.Row(c).Gap(4).AlignItems(ui.Center).TextColor(Tokens(c).Warning).Role(ui.RoleStatus).Children(func() {
		ui.Icon(c, icons.Must("triangle-alert")).Size(14, 14)
		ui.Text(c, text).FontSize(FontSize(c, theme.CaptionSize))
	})
}
