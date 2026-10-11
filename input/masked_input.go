package input

import (
	"log/slog"
	"strings"
	"unicode"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// MaskedInputOptions configure a MaskedInput.
type MaskedInputOptions struct {
	// Mask is required: 9 takes a digit, a a letter, * either; anything else shows as is.
	Mask     string
	Label    string
	Disabled bool
	// Error marks the value invalid and shows the message below.
	Error string
}

// maskedInputState is the caret over the typed characters and any input method composition.
type maskedInputState struct {
	caret              int
	compose            string
	changed, submitted bool
}

// maskedInputFits reports whether r may fill the mask slot m.
func maskedInputFits(m, r rune) bool {
	switch m {
	case '9':
		return unicode.IsDigit(r)
	case 'a':
		return unicode.IsLetter(r)
	case '*':
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}
	return false
}

func maskedInputSlot(m rune) bool { return m == '9' || m == 'a' || m == '*' }

// MaskFormat lays raw into mask; unfilled slots show as underscores.
func MaskFormat(mask, raw string) string {
	rs := []rune(raw)
	var b strings.Builder
	i := 0
	for _, m := range mask {
		switch {
		case !maskedInputSlot(m):
			b.WriteRune(m)
		case i < len(rs):
			b.WriteRune(rs[i])
			i++
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}

// maskedInputAt is the mask index of slot n, or the mask's length past the last.
func maskedInputAt(mask []rune, n int) int {
	for i, m := range mask {
		if maskedInputSlot(m) {
			if n == 0 {
				return i
			}
			n--
		}
	}
	return len(mask)
}

// MaskedInput edits *value, the characters typed into Mask's slots; an input method's composition commits only when done.
func MaskedInput(c *ui.Context, value *string, opts MaskedInputOptions) InputResult {
	if value == nil {
		panic("mujica: MaskedInput needs a value")
	}
	mask := []rune(opts.Mask)
	var slots []rune
	for _, m := range mask {
		if maskedInputSlot(m) {
			slots = append(slots, m)
		}
	}
	if len(slots) == 0 {
		panic("mujica: MaskedInput Mask has no slots")
	}
	raw := []rune(*value)
	if len(raw) > len(slots) {
		panic("mujica: MaskedInput value is longer than its Mask")
	}
	for i, r := range raw {
		if !maskedInputFits(slots[i], r) {
			panic("mujica: MaskedInput value does not fit its Mask")
		}
	}
	k := core.Tokens(c)
	var r InputResult
	look := &textInputLook{disabled: opts.Disabled, err: opts.Error}
	r.Element, r.Input = textInputFrame(c, look, func(frame *ui.Element) *ui.Element {
		in := ui.Row(c).Grow(1).Focusable().Role(ui.RoleTextField).Label(opts.Label).Description(MaskFormat(opts.Mask, *value))
		st := ui.Local(in, "masked_input", func() *maskedInputState { return &maskedInputState{caret: len(raw)} })
		s := *st
		s.caret = min(s.caret, len(raw))
		insert := func(text string) {
			for _, ch := range text {
				if s.caret < len(slots) && len(raw) < len(slots) && maskedInputFits(slots[s.caret], ch) {
					raw = append(raw[:s.caret], append([]rune{ch}, raw[s.caret:]...)...)
					s.caret++
					s.changed = true
				}
			}
		}
		in.HandleInput(func(ev ui.InputEvent) bool {
			if in.IsDisabled() {
				return false
			}
			raw = []rune(*value)
			switch ev.Kind {
			case ui.InputText:
				s.compose = ""
				insert(ev.Text)
			case ui.InputCompose:
				s.compose = ev.Text
			case ui.InputCommand:
				switch ev.Text {
				case "paste":
					insert(c.ReadClipboard())
				case "copy":
					c.WriteClipboard(*value)
				default:
					return false
				}
			case ui.InputKeyDown:
				if s.compose != "" {
					return false // the input method owns the keys while composing
				}
				switch ev.Key {
				case ui.KeyBackspace:
					if s.caret > 0 {
						raw = append(raw[:s.caret-1], raw[s.caret:]...)
						s.caret--
						s.changed = true
					}
				case ui.KeyDelete:
					if s.caret < len(raw) {
						raw = append(raw[:s.caret], raw[s.caret+1:]...)
						s.changed = true
					}
				case ui.KeyLeft:
					s.caret = max(0, s.caret-1)
				case ui.KeyRight:
					s.caret = min(len(raw), s.caret+1)
				case ui.KeyHome:
					s.caret = 0
				case ui.KeyEnd:
					s.caret = len(raw)
				case ui.KeyEnter:
					s.submitted = true
				default:
					return false
				}
			default:
				return false
			}
			*value = string(raw)
			return true
		})
		if s.changed {
			r.changed, s.changed = true, false
		}
		if s.submitted {
			r.submitted, s.submitted = true, false
		}
		shown := []rune(MaskFormat(opts.Mask, *value))
		at := maskedInputAt(mask, s.caret)
		size := core.FontSize(c, theme.BodySize)
		w, h := c.MeasureText(0, ui.Span{Text: string(shown[:at]), Size: size})
		if h == 0 {
			_, h = c.MeasureText(0, ui.Span{Text: "0", Size: size})
		}
		in.MinHeight(h).TextCaret(ui.Rect{X: w, Y: 0, W: 1, H: h}).Cursor(ui.CursorText)
		filled := maskedInputAt(mask, len([]rune(*value)))
		in.Children(func() {
			ui.Text(c, string(shown[:at])).TextColor(k.Text).FontSize(size).SingleLine()
			if s.compose != "" {
				ui.Text(c, s.compose).TextColor(k.Text).FontSize(size).Underline().SingleLine()
			}
			if at < filled {
				ui.Text(c, string(shown[at:filled])).TextColor(k.Text).FontSize(size).SingleLine()
			}
			ui.Text(c, string(shown[max(at, filled):])).TextColor(k.TextMuted).FontSize(size).SingleLine()
		})
		in.DrawOver(func(p *ui.Painter, rr ui.Rect) {
			if in.Focused() && !in.IsDisabled() && s.compose == "" {
				p.Fill(ui.Rect{X: rr.X + w, Y: rr.Y + (rr.H-h)/2, W: 1, H: h}, k.Text, 0)
			}
		})
		return in
	}, nil)
	if r.changed {
		slog.Debug("mujica: masked_input change", "filled", len([]rune(*value)), "slots", len(slots))
	}
	return r
}
