package input

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// PinInputOptions configure a PinInput. Length is 6 when zero.
type PinInputOptions struct {
	Length int
	// Label names the input; empty uses the built-in "Verification code".
	Label    string
	Disabled bool
	// Invalid draws the cells in the danger color.
	Invalid bool
}

// PinInputParts are the parts of a PinInput.
type PinInputParts struct {
	Element            *ui.Element
	changed, completed bool
}

// Changed reports whether the user changed the code this frame.
func (p PinInputParts) Changed() bool { return p.changed }

// Completed reports whether the user filled the last cell this frame.
func (p PinInputParts) Completed() bool { return p.completed }

// pinInputEvent is input taken between frames, applied in the next one.
type pinInputEvent struct {
	text      string
	backspace bool
	paste     bool
}

type pinInputState struct {
	queue   []pinInputEvent
	compose string
}

// PinInput edits a fixed-length digit code, one cell per digit, keeping leading zeros.
func PinInput(c *ui.Context, value *string, opts PinInputOptions) PinInputParts {
	n := opts.Length
	if n == 0 {
		n = 6
	}
	if n < 0 || len(*value) > n || strings.Trim(*value, "0123456789") != "" {
		panic(fmt.Sprintf("mujica: PinInput value %q does not fit %d digits", *value, n))
	}
	k := core.Tokens(c)
	d := core.Density(c)
	label := opts.Label
	if label == "" {
		label = core.Msg(c, "pin_input.label")
	}
	parts := PinInputParts{}
	e := ui.Row(c).Gap(d.Unit() * 2).Shrink(0).Role(ui.RoleTextField).Label(label).Disabled(opts.Disabled)
	parts.Element = e
	st := ui.Local(e, "pin", func() pinInputState { return pinInputState{} })
	put := func(s string) {
		var b strings.Builder
		for _, r := range s {
			if r >= '0' && r <= '9' {
				b.WriteRune(r)
			}
		}
		next := b.String()
		if len(next) > n {
			next = next[:n]
		}
		if next == *value {
			return
		}
		was := len(*value)
		*value = next
		parts.changed = true
		if len(next) == n && was < n {
			parts.completed = true
			slog.Debug("mujica: pin input completed", "length", n)
		}
	}
	disabled := e.IsDisabled()
	if disabled {
		st.queue, st.compose = st.queue[:0], ""
	}
	for _, ev := range st.queue {
		switch {
		case ev.paste:
			put(c.ReadClipboard())
		case ev.backspace && *value != "":
			put((*value)[:len(*value)-1])
		case ev.text != "":
			put(*value + ev.text)
		}
	}
	st.queue = st.queue[:0]
	if !disabled {
		e.Focusable().FocusRing(false).HandleInput(func(ev ui.InputEvent) bool {
			switch {
			case ev.Kind == ui.InputText:
				st.queue = append(st.queue, pinInputEvent{text: ev.Text})
				st.compose = ""
			case ev.Kind == ui.InputCompose:
				st.compose = ev.Text
			case ev.Kind == ui.InputCommand && ev.Text == "paste",
				ev.Kind == ui.InputKeyDown && ev.Mods == ui.Cmd && ev.Key == ui.KeyV:
				st.queue = append(st.queue, pinInputEvent{paste: true})
			case ev.Kind == ui.InputKeyDown && ev.Mods == 0 && ev.Key == ui.KeyBackspace:
				st.queue = append(st.queue, pinInputEvent{backspace: true})
			default:
				return false
			}
			return true
		})
	}
	size := d.ControlHeight() + d.Unit()*2
	active := min(len(*value), n-1)
	focused := e.Focused()
	compose := st.compose
	e.Description(fmt.Sprintf(core.Msg(c, "pin_input.filled"), len(*value), n))
	e.TextCaret(ui.Rect{X: float32(active) * (size + d.Unit()*2), Y: 0, W: size, H: size})
	e.Children(func() {
		for i := range n {
			cell := ui.Box(c).Size(size, size).Center().Radius(theme.ControlRadius).Background(k.Surface).Role(ui.RoleNone)
			border := k.ControlBorder
			switch {
			case opts.Invalid:
				border = k.Danger
			case focused && i == active:
				border = k.Focus
			}
			cell.Border(theme.BorderWidth, border)
			ch := ""
			if i < len(*value) {
				ch = (*value)[i : i+1]
			}
			showCompose := focused && i == active && compose != ""
			cell.Children(func() {
				switch {
				case showCompose:
					ui.Text(c, compose).FontSize(core.FontSize(c, theme.H3Size)).Underline().SingleLine()
				case ch != "":
					ui.Text(c, ch).FontSize(core.FontSize(c, theme.H3Size)).FontFeatures("tnum").SingleLine()
				}
			})
			caret := focused && i == active && len(*value) < n && !showCompose
			cell.DrawOver(func(p *ui.Painter, r ui.Rect) {
				if focused && i == active {
					internal.FocusRing(p, r, theme.ControlRadius, k.Focus)
				}
				if caret {
					h := r.H * 0.4
					p.Fill(ui.Rect{X: r.X + r.W/2 - 0.5, Y: r.Y + (r.H-h)/2, W: 1, H: h}, k.Text, 0)
				}
			})
		}
	})
	return parts
}
