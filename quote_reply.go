package mujica

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"mujicaui/internal"
	"mujicaui/theme"
)

// QuoteReply shows the author and a one-line excerpt of a quoted message; Clicked jumps to it.
func QuoteReply(c *ui.Context, author, excerpt string) *ui.Element {
	if author == "" {
		panic("mujica: QuoteReply needs an author")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	line := strings.Join(strings.Fields(excerpt), " ")
	b := ui.ButtonBase(c).MinHeight(theme.MinHitSize).MaxWidth(420).Gap(u*2).Padding(u, u*3, u, u*2).Justify(ui.Start).
		Radius(theme.ControlRadius).Background(k.SurfaceHover).BorderWidth(0, 0, 0, 3).BorderColor(k.AccentText).
		Label(fmt.Sprintf(msg(c, "quote_reply.label"), author, line)).Transition(colorFade)
	internal.Ring(b, k.Focus, theme.ControlRadius)
	if b.Hovered() {
		b.Background(k.Selection)
	}
	b.Children(func() {
		ui.Text(c, author).FontWeight(600).FontSize(FontSize(c, theme.CaptionSize)).SingleLine().Shrink(0)
		ui.Text(c, line).FontSize(FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted).SingleLine().MinWidth(0).Shrink(1)
	})
	return b
}
