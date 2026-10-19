package input

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// ButtonVariant is a button's emphasis. The zero value is Primary.
type ButtonVariant int

const (
	Primary ButtonVariant = iota
	Secondary
	Outline
	Ghost
	Danger
)

// ButtonOptions configure a Button.
type ButtonOptions struct {
	Variant ButtonVariant
	// Icon goes before the label.
	Icon     *ui.SVG
	Loading  bool
	Disabled bool
	// Label names an icon-only button for assistive technology; required without a label.
	Label string
}

// buttonSpinnerSize is the leading icon slot, which the spinner takes while loading.
const buttonSpinnerSize float32 = 16

// Button is a push button. Loading keeps its size and blocks clicks.
func Button(c *ui.Context, label string, opts ButtonOptions) ui.Element {
	if label == "" && opts.Label == "" {
		panic("mujica: Button without a label needs Label to name it")
	}
	k := core.Tokens(c)
	d := core.Density(c)
	pad, gap := d.Unit()*4, d.Unit()*2
	if opts.Loading && opts.Icon == nil && label != "" {
		// The spinner slot comes out of the padding, so the width holds.
		pad -= (buttonSpinnerSize + gap) / 2
	}
	b := ui.ButtonBase(c).MinHeight(d.ControlHeight()).MinWidth(theme.MinHitSize).
		Padding(0, pad).Gap(gap).Radius(theme.ControlRadius).
		Disabled(opts.Disabled || opts.Loading)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if opts.Label != "" {
		b.Label(opts.Label)
	}
	bg, fg, border := ButtonColors(k, opts.Variant, b.Hovered() && !opts.Disabled, b.Pressed())
	b.Background(bg).TextColor(fg)
	if border.A > 0 {
		b.Border(theme.BorderWidth, border)
	}
	b.Children(func() {
		if opts.Loading {
			core.Spinner(c, core.SpinnerOptions{Size: buttonSpinnerSize, Color: fg})
		} else if opts.Icon != nil {
			ui.Icon(c, opts.Icon).Size(buttonSpinnerSize, buttonSpinnerSize)
		}
		if label != "" {
			ui.Text(c, label).SingleLine()
		}
	})
	return b
}

func ButtonColors(k theme.Tokens, v ButtonVariant, hover, pressed bool) (bg, fg, border ui.Color) {
	switch v {
	case Secondary:
		bg, fg, border = k.Surface, k.Text, k.ControlBorder
	case Outline:
		bg, fg, border = ui.Transparent, k.AccentText, k.AccentText
	case Ghost:
		bg, fg = ui.Transparent, k.AccentText
	case Danger:
		bg, fg = k.DangerFill, k.OnDanger
		if pressed {
			bg = k.DangerFill.Mix(ui.RGB(0, 0, 0), 0.25)
		} else if hover {
			bg = k.DangerFill.Mix(ui.RGB(0, 0, 0), 0.12)
		}
		return
	default:
		bg, fg = k.Accent, k.OnAccent
		if pressed {
			bg = k.AccentPressed
		} else if hover {
			bg = k.AccentHover
		}
		return
	}
	if pressed {
		bg = k.Selection
	} else if hover {
		bg = k.SurfaceHover
	}
	return
}
