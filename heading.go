package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// HeadingOptions configure a Heading. The zero value is a level 1 heading.
type HeadingOptions struct {
	// Level is 1, 2 or 3; zero means 1.
	Level int
	// Subtitle shows below the title in muted text.
	Subtitle string
	// Divider draws an ornament rule with a center diamond below.
	Divider bool
}

// Heading shows a serif title that follows the desktop's text scale.
func Heading(c *ui.Context, title string, opts HeadingOptions) *ui.Element {
	k := Tokens(c)
	size := headingSize(opts.Level)
	scale := c.Preferences().TextScale
	if scale <= 0 {
		scale = 1
	}
	u := Density(c).Unit()
	col := ui.Column(c).Gap(u).Shrink(0)
	col.Children(func() {
		ui.Text(c, title).Font(theme.SerifFont).FontSize(size * scale).FontWeight(600).
			LineHeight(1.25).TextColor(k.Text)
		if opts.Subtitle != "" {
			ui.Text(c, opts.Subtitle).FontSize(theme.BodySize * scale).LineHeight(theme.LineHeight).TextColor(k.TextMuted)
		}
		if opts.Divider {
			ui.Box(c).FillWidth().Height(9).Margin(u, 0, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
				headingOrnament(p, r, k.Ornament)
			})
		}
	})
	return col
}

func headingSize(level int) float32 {
	switch level {
	case 0, 1:
		return theme.H1Size
	case 2:
		return theme.H2Size
	case 3:
		return theme.H3Size
	}
	panic("mujica: Heading Level must be 1, 2 or 3")
}

// headingOrnament draws symmetric rules around a center diamond.
func headingOrnament(p *ui.Painter, r ui.Rect, c ui.Color) {
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	const d, gap = 7, 6
	internal.Rule(p, r.X, cx-gap, cy, c)
	internal.Rule(p, cx+gap, r.X+r.W, cy, c)
	p.FillPath(internal.Diamond(cx, cy, d), c)
}
