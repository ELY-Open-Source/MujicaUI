package layout

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/MujicaUI/core"
	"github.com/ELY-Open-Source/MujicaUI/internal"
	"github.com/ELY-Open-Source/MujicaUI/theme"
)

// dividerDiamond is the ornament's size in DIPs.
const dividerDiamond float32 = 7

// DividerOptions configure a Divider.
type DividerOptions struct {
	// Vertical draws a line down a row instead of across a column.
	Vertical bool
	// Label sits in the middle of a horizontal divider.
	Label string
	// Ornament draws the line in the ornament color with a single diamond in the middle; a Label replaces the diamond.
	Ornament bool
}

// Divider is a thin separating line. It never takes focus and its ornament is not read out.
func Divider(c *ui.Context, opts DividerOptions) ui.Element {
	if opts.Vertical && opts.Label != "" {
		panic("mujica: a vertical Divider cannot have a Label")
	}
	k := core.Tokens(c)
	line := k.Border
	if opts.Ornament {
		line = k.Ornament
	}
	if opts.Label == "" && !opts.Ornament {
		e := ui.Box(c).Background(line).Shrink(0).AlignSelf(ui.Stretch)
		if opts.Vertical {
			return e.Width(theme.BorderWidth)
		}
		return e.Height(theme.BorderWidth)
	}
	if opts.Label == "" && !opts.Vertical {
		gap := core.Density(c).Unit()*2 + dividerDiamond/2
		return ui.Box(c).Height(dividerDiamond).Shrink(0).AlignSelf(ui.Stretch).Draw(func(p *ui.Painter, r ui.Rect) {
			internal.RuleDiamond(p, r, gap, dividerDiamond, line)
		})
	}
	e := ui.Row(c).Gap(core.Density(c).Unit() * 2).Shrink(0).AlignSelf(ui.Stretch)
	if opts.Vertical {
		e.Column().AlignItems(ui.Center)
	}
	segment := func() {
		s := ui.Box(c).Grow(1).Background(line)
		if opts.Vertical {
			s.Width(theme.BorderWidth)
		} else {
			s.Height(theme.BorderWidth)
		}
	}
	return e.Children(func() {
		segment()
		if opts.Label != "" {
			ui.Text(c, opts.Label).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.CaptionSize)).
				TextColor(k.TextMuted).SingleLine().Shrink(0)
		} else {
			ui.Box(c).Size(dividerDiamond, dividerDiamond).Shrink(0).Role(ui.RoleNone).
				Draw(func(p *ui.Painter, r ui.Rect) {
					p.FillPath(internal.Diamond(r.X+r.W/2, r.Y+r.H/2, r.W), k.Ornament)
				})
		}
		segment()
	})
}
