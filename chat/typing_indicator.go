package chat

import (
	"fmt"
	"math"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// typingPeriod is one cycle of the dots, in milliseconds.
const typingPeriod = 1200

// TypingIndicator shows that name is typing: three fading dots, or static text under reduced motion.
func TypingIndicator(c *ui.Context, name string) *ui.Element {
	if name == "" {
		panic("mujica: TypingIndicator needs a name")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	label := fmt.Sprintf(core.Msg(c, "typing_indicator.label"), name)
	still := c.Preferences().ReduceMotion
	row := ui.Row(c).Gap(u*2).AlignItems(ui.Center).Padding(u, u*3).Radius(theme.CardRadius).
		Background(k.Surface).Border(theme.BorderWidth, k.Border).Role(ui.RoleStatus).Label(label).AlignSelf(ui.Start)
	row.Children(func() {
		if still {
			ui.Text(c, label).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine()
			return
		}
		ui.Box(c).Size(30, 12).Draw(func(p *ui.Painter, r ui.Rect) {
			// Repaints only while shown: the painter asks for the next frame.
			p.AnimationFrame()
			phase := float64(p.Now().UnixMilli()%typingPeriod) / typingPeriod
			for i := range 3 {
				a := 0.35 + 0.65*(0.5+0.5*math.Cos(2*math.Pi*(phase-float64(i)/6)))
				cx := internal.Snap(p, r.X+5+float32(i)*10)
				p.Fill(ui.Rect{X: cx - 3, Y: r.Y + r.H/2 - 3, W: 6, H: 6}, k.TextMuted.Alpha(float32(a)), 3)
			}
		})
	})
	return row
}
