package mujica

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// TextOptions configure a Text. The zero value wraps every line.
type TextOptions struct {
	// MaxLines cuts the text with an ellipsis after that many lines; zero never cuts.
	MaxLines int
	// Selectable lets the user select and copy the text.
	Selectable bool
}

// Text shows body text that wraps at the width it gets.
func Text(c *ui.Context, s string, opts TextOptions) *ui.Element {
	if opts.MaxLines < 0 {
		panic("mujica: Text MaxLines must not be negative")
	}
	e := ui.Text(c, s).LineHeight(theme.LineHeight)
	switch {
	case opts.MaxLines == 1:
		e.SingleLine()
	case opts.MaxLines > 1:
		e.MaxLines(opts.MaxLines)
	}
	if opts.Selectable {
		e.Selectable()
	}
	return e
}

// EditableTextOptions configure the editable mode of Text.
type EditableTextOptions struct {
	// Saving shows the caller's save in progress and blocks editing.
	Saving bool
	// Error is the caller's save failure, shown below the text.
	Error string
	// Label names the text for assistive technology.
	Label string
}

// EditableText is Text's editable mode. Double click or Enter edits; Enter
// or leaving confirms, Escape cancels. Changed reports a confirmed new value.
func EditableText(c *ui.Context, value *string, opts EditableTextOptions) *ui.Element {
	k := Tokens(c)
	u := Density(c).Unit()
	line := k.ControlBorder
	if opts.Error != "" {
		line = k.Danger
	}
	var box *ui.Element
	ui.Column(c).Gap(u).AlignItems(ui.Start).Children(func() {
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Disabled(opts.Saving).Children(func() {
			box = ui.EditableText(c, value)
			pad := u / 2
			box.Padding(pad+1, pad*2+1).MinWidth(theme.MinHitSize*2).MinHeight(theme.MinHitSize).Radius(theme.ControlRadius).
				BorderWidth(0, 0, 1, 0).BorderColor(line).BorderStyle(ui.BorderDashed)
			internal.Ring(box, k.Focus, theme.ControlRadius)
			if box.Hovered() {
				box.Background(k.SurfaceHover)
			}
			internal.Fade(box, Motion(c, theme.StateDuration))
			if opts.Label != "" {
				box.Label(opts.Label)
			}
			if opts.Saving {
				Spinner(c, SpinnerOptions{Size: 12, Label: msg(c, "editable_text.saving")})
			}
		})
		if opts.Error != "" {
			box.Error(opts.Error)
			ui.Text(c, opts.Error).TextColor(k.Danger).FontSize(theme.CaptionSize)
		}
	})
	if box.Changed() {
		slog.Debug("mujica: editabletext changed", "len", len(*value))
	}
	return box
}
