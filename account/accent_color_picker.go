package account

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/icons"
	"github.com/ELY-Open-Source/MujicaUI/input"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// AccentColorPickerOptions configure an AccentColorPicker.
type AccentColorPickerOptions struct {
	// Candidates defaults to the built-in court palette.
	Candidates []ui.Color
	Disabled   bool
}

// accentPickerColors offers readable accents and marked low-contrast examples.
var accentPickerColors = []ui.Color{
	ui.Hex("#7D2034"), ui.Hex("#355E80"), ui.Hex("#326247"), ui.Hex("#4B3A5A"),
	ui.Hex("#805714"), ui.Hex("#211C1D"), ui.Hex("#AD9164"), ui.Hex("#D9A5B3"),
}

// accentPickerReadable reports whether accent carries OnAccent text at 4.5:1.
func accentPickerReadable(k theme.Tokens, accent ui.Color) bool {
	return theme.Contrast(k.OnAccent, accent) >= 4.5
}

// AccentColorPicker selects an accent and marks insufficient contrast.
func AccentColorPicker(c *ui.Context, color *ui.Color, opts AccentColorPickerOptions) core.ChoiceResult {
	cands := opts.Candidates
	if cands == nil {
		cands = accentPickerColors
	}
	if len(cands) == 0 {
		panic("mujica: AccentColorPicker needs candidates")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	var r core.ChoiceResult
	r.Element = ui.Column(c).Gap(d.Unit() * 2)
	r.Element.Children(func() {
		low := false
		ui.Row(c).Gap(d.Unit() * 2).Wrap().Role(ui.RoleRadioGroup).Label(core.Msg(c, "accent_color_picker.label")).
			FocusGroup(ui.Horizontal).Disabled(opts.Disabled).Children(func() {
			selected := color
			if r.Element.IsDisabled() || opts.Disabled {
				copy := *color
				selected = &copy
			}
			for _, col := range cands {
				ok := accentPickerReadable(k, col)
				low = low || !ok
				hex := strings.ToUpper(input.ColorPickerHex(col))
				sw := ui.RadioBase(c, selected, col).Size(36, 36).Center().Radius(theme.CardRadius).Label(hex)
				if !ok {
					sw.Description(core.Msg(c, "accent_color_picker.low"))
				}
				if sw.Changed() && !sw.IsDisabled() {
					r = core.NewChoiceResult(r.Element, true)
					slog.Debug("mujica: accent_color_picker chose", "color", hex)
				}
				internal.Ring(sw, k.Focus, theme.CardRadius)
				on := *color == col
				sw.Draw(func(p *ui.Painter, rect ui.Rect) {
					p.Fill(rect, col, theme.CardRadius)
					p.Stroke(rect, k.Border, theme.CardRadius, 1)
					if on {
						p.Stroke(ui.Rect{X: rect.X - 3, Y: rect.Y - 3, W: rect.W + 6, H: rect.H + 6}, k.Text, theme.CardRadius+3, 1.5)
					}
				})
				sw.Children(func() {
					switch {
					case on:
						ui.Icon(c, icons.Must("check")).Size(16, 16).TextColor(accentPickerMark(col))
					case !ok:
						ui.Icon(c, icons.Must("triangle-alert")).Size(14, 14).TextColor(accentPickerMark(col))
					}
				})
			}
		})
		if low {
			ui.Row(c).Gap(d.Unit()).AlignItems(ui.Center).TextColor(k.TextMuted).Children(func() {
				ui.Icon(c, icons.Must("triangle-alert")).Size(12, 12)
				ui.Text(c, core.Msg(c, "accent_color_picker.low_note")).FontSize(core.FontSize(c, theme.CaptionSize))
			})
		}
		if !accentPickerReadable(k, *color) {
			ui.Text(c, fmt.Sprintf(core.Msg(c, "accent_color_picker.chosen_low"), strings.ToUpper(input.ColorPickerHex(*color)))).
				TextColor(k.Warning).FontSize(core.FontSize(c, theme.CaptionSize)).Role(ui.RoleStatus)
		}
	})
	return r
}

// accentPickerMark is black or white, whichever reads on col.
func accentPickerMark(col ui.Color) ui.Color {
	if theme.Contrast(col, ui.RGB(255, 255, 255)) >= theme.Contrast(col, ui.RGB(0, 0, 0)) {
		return ui.RGB(255, 255, 255)
	}
	return ui.RGB(0, 0, 0)
}
