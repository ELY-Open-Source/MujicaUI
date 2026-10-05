package chat

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/icons"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// ScrollToBottomButton shows while away is true, with the unread count; click it to return to the end.
// While hidden it returns an empty element whose Clicked is false.
func ScrollToBottomButton(c *ui.Context, away bool, unread int) *ui.Element {
	if unread < 0 {
		panic("mujica: ScrollToBottomButton unread must not be negative")
	}
	if !away {
		return ui.Box(c).Absolute().Invisible()
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	label := core.Msg(c, "scroll_to_bottom_button.label")
	if unread > 0 {
		label = fmt.Sprintf(core.Msg(c, "scroll_to_bottom_button.unread"), unread)
	}
	wrap := ui.Row(c).Absolute().Left(0).Right(0).Bottom(u * 3).Justify(ui.Center).PassThrough()
	var b *ui.Element
	wrap.Children(func() {
		b = ui.ButtonBase(c).MinHeight(theme.MinHitSize).MinWidth(theme.MinHitSize).Padding(0, u*3).Gap(u*2).
			Radius(theme.MinHitSize/2).Label(label).Background(k.Surface).Border(theme.BorderWidth, k.ControlBorder).
			TextColor(k.Text).Shadow(0, 2, 8, 0, ui.RGBA(0, 0, 0, 0.14)).Transition(core.ColorFade)
		b.Opacity(core.OverlayEnter(c, b))
		internal.Ring(b, k.Focus, theme.MinHitSize/2)
		if b.Hovered() {
			b.Background(k.SurfaceHover)
		}
		b.Children(func() {
			ui.Icon(c, icons.Must("arrow-down")).Size(16, 16)
			if unread > 0 {
				ui.Text(c, fmt.Sprint(unread)).FontFeatures("tnum").FontWeight(600).TextColor(k.AccentText).SingleLine()
			}
		})
	})
	return b
}
