package messaging

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// UnreadDivider marks where unread messages begin; count zero shows "New".
func UnreadDivider(c *ui.Context, count int) *ui.Element {
	if count < 0 {
		panic("mujica: UnreadDivider count must not be negative")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	text := core.Msg(c, "unread_divider.new")
	if count > 0 {
		text = fmt.Sprintf(core.Msg(c, "unread_divider.count"), count)
	}
	row := ui.Row(c).FillWidth().Gap(u * 2).AlignItems(ui.Center).MinHeight(theme.MinHitSize).Role(ui.RoleStatus).Label(text)
	row.Children(func() {
		line := func() { ui.Box(c).Height(1).Grow(1).Background(k.Danger) }
		line()
		ui.Text(c, text).FontSize(core.FontSize(c, theme.CaptionSize)).FontWeight(600).TextColor(k.Danger).SingleLine().
			Padding(0, u*2).Radius(theme.ControlRadius).Border(theme.BorderWidth, k.Danger)
		line()
	})
	return row
}
