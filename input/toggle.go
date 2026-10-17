package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ToggleOptions configure a Toggle.
type ToggleOptions struct {
	// Icon goes before the label.
	Icon     *ui.SVG
	Disabled bool
	// Label names an icon-only toggle for assistive technology.
	Label string
}

// Toggle is a button that stays pressed while *on. Click or Space flips it; Changed reports it.
func Toggle(c *ui.Context, on *bool, label string, opts ToggleOptions) ui.Element {
	if label == "" && opts.Label == "" {
		panic("mujica: Toggle needs a label or Label")
	}
	var b ui.Element
	ui.Row(c).Shrink(0).Disabled(opts.Disabled).Children(func() {
		b = ui.ToggleBase(c, on)
	})
	if b.Changed() {
		slog.Debug("mujica: toggle changed", "label", label+opts.Label, "on", *on)
	}
	toggleLook(c, b, *on, label, opts.Icon)
	if opts.Label != "" {
		b.Label(opts.Label)
	}
	return b
}

// toggleLook styles a pressed-or-not button and builds its content.
func toggleLook(c *ui.Context, b ui.Element, on bool, label string, icon *ui.SVG) {
	k := core.Tokens(c)
	d := core.Density(c)
	b.MinHeight(d.ControlHeight()).MinWidth(d.ControlHeight()).Padding(0, d.Unit()*3).Gap(d.Unit() * 2).
		Radius(theme.ControlRadius)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	bg, fg, border := k.Surface, k.Text, k.ControlBorder
	switch {
	case on:
		bg, fg, border = k.Selection, k.AccentText, k.AccentText
	case b.Pressed():
		bg = k.Selection
	case b.Hovered():
		bg = k.SurfaceHover
	}
	b.Background(bg).TextColor(fg).Border(theme.BorderWidth, border)
	if on {
		toggleMark(b, k)
	}
	b.Children(func() {
		if icon != nil {
			ui.Icon(c, icon).Size(16, 16)
		}
		if label != "" {
			ui.Text(c, label).SingleLine()
		}
	})
}

// toggleMark draws a bar under the content: the pressed state beside the color.
func toggleMark(b ui.Element, k theme.Tokens) {
	b.Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X + r.W*0.25, Y: r.Y + r.H - 4, W: r.W * 0.5, H: 2}, k.AccentText, 1)
	})
}
