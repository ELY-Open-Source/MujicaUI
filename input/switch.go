package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// SwitchOptions configure a Switch.
type SwitchOptions struct {
	Disabled bool
	// Loading shows a pending result: the switch keeps its place and takes no input.
	Loading bool
	// Description is read out after the label.
	Description string
}

const (
	switchW float32 = 36
	switchH float32 = 20
	switchK float32 = 14
)

// Switch turns *on over by click or Space; Changed reports it, the caller applies it.
func Switch(c *ui.Context, on *bool, label string, opts SwitchOptions) ui.Element {
	k := core.Tokens(c)
	d := core.Density(c)
	blocked := opts.Disabled || opts.Loading
	var sw ui.Element
	// Disabled on a wrapper blocks the base before it takes input.
	ui.Row(c).Shrink(0).Disabled(blocked).Children(func() {
		sw = ui.SwitchBase(c, on).Gap(d.Unit() * 2).AlignItems(ui.Center).
			MinHeight(d.ControlHeight()).FocusRing(false)
	})
	if opts.Description != "" {
		sw.Description(opts.Description)
	}
	if sw.Changed() {
		slog.Debug("mujica: switch changed", "label", label, "on", *on)
	}
	hover := sw.Hovered()
	pos := sw.Animate("on", switchFloat(*on), core.Motion(c, theme.StateDuration))
	sw.Children(func() {
		track := ui.Box(c).Size(switchW, switchH).Shrink(0).Radius(theme.CardRadius)
		bg, border := k.Surface, k.ControlBorder
		thumb := k.ControlBorder
		if *on {
			bg, border, thumb = k.Accent, k.Accent, k.OnAccent
			if hover {
				bg, border = k.AccentHover, k.AccentHover
			}
		} else if hover {
			bg, border, thumb = k.SurfaceHover, k.Text, k.Text
		}
		track.Background(bg).Border(theme.BorderWidth, border)
		track.Draw(func(p *ui.Painter, r ui.Rect) {
			if *on {
				// A short upright bar marks "on" without relying on color.
				p.Fill(ui.Rect{X: r.X + 9, Y: r.Y + 6, W: 1.5, H: r.H - 12}, k.OnAccent, 0)
			}
		})
		track.DrawOver(func(p *ui.Painter, r ui.Rect) {
			if sw.FocusVisible() {
				internal.FocusRing(p, r, theme.CardRadius, k.Focus)
			}
		})
		inset := (switchH - switchK) / 2
		x := inset + pos*(switchW-switchK-2*inset)
		track.Children(func() {
			t := ui.Box(c).Absolute().Left(x).Top(inset).Size(switchK, switchK).Radius(theme.ControlRadius).Background(thumb).Center()
			if opts.Loading {
				t.Children(func() { core.Spinner(c, core.SpinnerOptions{Size: 10, Color: bg}) })
			}
		})
		if label != "" {
			ui.Text(c, label).SingleLine()
		}
	})
	return sw
}

func switchFloat(b bool) float32 {
	if b {
		return 1
	}
	return 0
}
