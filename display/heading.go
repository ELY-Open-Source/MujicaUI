package display

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/internal"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// HeadingOptions configure a Heading; zero is level 1.
type HeadingOptions struct {
	// Level is 1, 2 or 3; zero means 1.
	Level int
	// Subtitle shows below the title in muted text.
	Subtitle string
	// Divider draws an ornament rule with a center diamond below.
	Divider bool
}

// Heading shows a serif title scaled with desktop text.
func Heading(c *ui.Context, title string, opts HeadingOptions) *ui.Element {
	k := core.Tokens(c)
	size := headingSize(opts.Level)
	u := core.Density(c).Unit()
	col := ui.Column(c).Gap(u).Shrink(0)
	col.Children(func() {
		ui.Text(c, title).Font(theme.SerifFont).FontSize(core.FontSize(c, size)).FontWeight(600).
			LineHeight(1.25).TextColor(k.Text)
		if opts.Subtitle != "" {
			ui.Text(c, opts.Subtitle).FontSize(core.FontSize(c, theme.BodySize)).LineHeight(theme.LineHeight).TextColor(k.TextMuted)
		}
		if opts.Divider {
			ui.Box(c).FillWidth().Height(9).Margin(u, 0, 0, 0).Draw(func(p *ui.Painter, r ui.Rect) {
				internal.RuleDiamond(p, r, 6, 7, k.Ornament)
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
