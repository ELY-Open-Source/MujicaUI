package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
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

// Button is a push button. Loading keeps its size and blocks clicks.
func Button(c *ui.Context, label string, opts ButtonOptions) *ui.Element {
	if label == "" && opts.Label == "" {
		panic("mujica: Button without a label needs Label to name it")
	}
	k := Tokens(c)
	d := Density(c)
	b := ui.ButtonBase(c).MinHeight(d.ControlHeight()).MinWidth(theme.MinHitSize).
		Padding(0, d.Unit()*4).Gap(d.Unit() * 2).Radius(theme.ControlRadius).
		Disabled(opts.Disabled || opts.Loading)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if opts.Label != "" {
		b.Label(opts.Label)
	}
	bg, fg, border := buttonColors(k, opts.Variant, b.Hovered() && !opts.Disabled, b.Pressed())
	b.Background(bg).TextColor(fg)
	if border.A > 0 {
		b.Border(theme.BorderWidth, border)
	}
	b.Children(func() {
		switch {
		case opts.Loading && opts.Icon != nil:
			// The spinner takes the icon's 16 DIP slot.
			Spinner(c, SpinnerOptions{Size: 16, Color: fg})
		case opts.Loading && label == "":
			Spinner(c, SpinnerOptions{Size: 16, Color: fg})
		case opts.Loading:
			// Over the hidden label, so the width stays.
			ui.Box(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Center().Children(func() {
				Spinner(c, SpinnerOptions{Size: 16, Color: fg})
			})
		case opts.Icon != nil:
			ui.Icon(c, opts.Icon).Size(16, 16)
		}
		if label != "" {
			t := ui.Text(c, label).SingleLine()
			if opts.Loading && opts.Icon == nil {
				t.TextColor(fg.Alpha(0))
			}
		}
	})
	return b
}

func buttonColors(k theme.Tokens, v ButtonVariant, hover, pressed bool) (bg, fg, border ui.Color) {
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
