package input

import (
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// ScrubInputOptions configure a ScrubInput; the range and step rules are NumberInput's.
type ScrubInputOptions struct {
	Min, Max  float64
	Step      float64
	Precision int
	Disabled  bool
}

// scrubInputPixels is how far a drag moves for one step.
const scrubInputPixels float32 = 4

// ScrubInput edits *value by dragging its label sideways, one Step per 4 DIPs, or by typing in its field.
func ScrubInput(c *ui.Context, label string, value *float64, opts ScrubInputOptions) InputResult {
	if label == "" || value == nil {
		panic("mujica: ScrubInput needs a label and a value")
	}
	nopts := NumberInputOptions{Min: opts.Min, Max: opts.Max, Step: opts.Step, Precision: opts.Precision, Label: label, Disabled: opts.Disabled}
	s := NumberInputResolve(nopts)
	k := core.Tokens(c)
	var r InputResult
	row := ui.Row(c).AlignItems(ui.Center).Gap(core.Density(c).Unit() * 2).Disabled(opts.Disabled)
	row.Children(func() {
		h := ui.Box(c).Shrink(0).MinHeight(theme.MinHitSize).MinWidth(theme.MinHitSize).Center().Cursor(ui.CursorResizeEW).Role(ui.RoleNone)
		drag := ui.Local(h, "scrub_input.drag", func() float64 { return math.NaN() })
		dx, _, held := h.Dragged()
		scrubbed := false
		switch {
		case !held || h.IsDisabled():
			*drag = math.NaN()
		case math.IsNaN(*drag):
			*drag = 0
			fallthrough
		default:
			start := *value
			*drag += float64(dx)
			if n := int(*drag / float64(scrubInputPixels)); n != 0 {
				*drag -= float64(n) * float64(scrubInputPixels)
				// One grid step at a time, so a long drag stops at the bound.
				v := start
				for range max(n, -n) {
					v = s.StepBy(v, max(-1, min(1, n)))
				}
				if v != *value {
					*value, scrubbed = v, true
				}
			}
		}
		fg := k.TextMuted
		if held || h.Hovered() {
			fg = k.AccentText
		}
		h.Children(func() {
			ui.Text(c, label).TextColor(fg).Underline().DecorationColor(k.Border).SingleLine()
		})
		r = NumberInput(c, value, nopts)
		r.Element.Width(core.Density(c).Unit() * 36)
		r.changed = r.changed || scrubbed
		if scrubbed {
			slog.Debug("mujica: scrub_input scrubbed", "label", label, "value", *value)
		}
	})
	r.Element = row
	return r
}
