package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// ToolbarOptions configure a Toolbar.
type ToolbarOptions struct {
	// Label names the toolbar for assistive technology.
	Label string
}

// Toolbar is one Tab stop of actions; those that do not fit go to an overflow menu.
func Toolbar(c *ui.Context, opts ToolbarOptions, fn func()) *ui.Element {
	k := Tokens(c)
	tb := ui.Toolbar(c, fn).Background(k.Surface).Border(theme.BorderWidth, k.Border).Radius(theme.CardRadius).
		MinHeight(Density(c).ControlHeight() + 8)
	if opts.Label != "" {
		tb.Label(opts.Label)
	}
	return tb
}

// ToolbarButtonOptions configure a ToolbarButton.
type ToolbarButtonOptions struct {
	Icon *ui.SVG
	// ShowLabel shows the label beside the icon; without an icon it always shows.
	ShowLabel bool
	Disabled  bool
}

// ToolbarButton is an action of a Toolbar; label names it and its overflow menu item.
func ToolbarButton(c *ui.Context, label string, opts ToolbarButtonOptions) *ui.Element {
	if label == "" {
		panic("mujica: ToolbarButton needs a label")
	}
	k := Tokens(c)
	d := Density(c)
	u := d.Unit()
	text := opts.ShowLabel || opts.Icon == nil
	b := ui.ButtonBase(c).MinHeight(d.ControlHeight()).MinWidth(d.ControlHeight()).Gap(u * 2).
		Radius(theme.ControlRadius).Label(label).Disabled(opts.Disabled).TextColor(k.Text)
	if text {
		b.Padding(0, u*3)
	}
	if !text {
		b.Tooltip(label)
	}
	internal.Ring(b, k.Focus, theme.ControlRadius)
	switch {
	case b.Pressed():
		b.Background(k.Selection)
	case b.Hovered():
		b.Background(k.SurfaceHover)
	}
	internal.Ease(b, Motion(c, theme.StateDuration), k.Text)
	b.Children(func() {
		if opts.Icon != nil {
			ui.Icon(c, opts.Icon).Size(16, 16)
		}
		if text {
			ui.Text(c, label).SingleLine()
		}
	})
	return b
}

// ToolbarGroup keeps related actions together; the overflow menu sets groups apart.
func ToolbarGroup(c *ui.Context, fn func()) *ui.Element {
	return ui.Row(c).AlignItems(ui.Center).Gap(Density(c).Unit() * 0.5).Shrink(0).Role(ui.RoleGroup).Children(fn)
}

// ToolbarSeparator draws a thin rule between actions.
func ToolbarSeparator(c *ui.Context) {
	ui.Box(c).Size(1, Density(c).ControlHeight()-12).Margin(0, Density(c).Unit()).Shrink(0).
		Background(Tokens(c).Ornament).Role(ui.RoleNone)
}
