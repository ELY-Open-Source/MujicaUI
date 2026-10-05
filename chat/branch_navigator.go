package chat

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
	"github.com/ZacharyZhang-NY/MujicaUI/internal"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

// BranchNavigator switches among count branches of a message, showing "< 2 / 3 >"; Left and Right step.
func BranchNavigator(c *ui.Context, index *int, count int) core.ChoiceResult {
	if count < 1 || *index < 0 || *index >= count {
		panic("mujica: BranchNavigator needs 0 <= index < count")
	}
	k := core.Tokens(c)
	label := fmt.Sprintf(core.Msg(c, "branch_navigator.label"), *index+1, count)
	var r core.ChoiceResult
	r.Element = ui.Row(c).Focusable().AlignItems(ui.Center).Radius(theme.ControlRadius).Role(ui.RoleGroup).Label(label)
	internal.Ring(r.Element, k.Focus, theme.ControlRadius)
	disabled := r.Element.IsDisabled()
	step := func(d int) {
		if i := *index + d; i >= 0 && i < count && !disabled {
			*index, r = i, core.NewChoiceResult(r.Element, true)
			c.Announce(fmt.Sprintf(core.Msg(c, "branch_navigator.label"), i+1, count))
			slog.Debug("mujica: branch_navigator change", "index", i)
		}
	}
	switch {
	case r.Element.Shortcut(0, ui.KeyLeft):
		step(-1)
	case r.Element.Shortcut(0, ui.KeyRight):
		step(1)
	case r.Element.Shortcut(0, ui.KeyHome):
		step(-*index)
	case r.Element.Shortcut(0, ui.KeyEnd):
		step(count - 1 - *index)
	}
	r.Element.Children(func() {
		if input.NumberInputButton(c, "chevron-left", core.Msg(c, "branch_navigator.previous"), *index > 0) {
			step(-1)
		}
		ui.Text(c, fmt.Sprintf("%d / %d", *index+1, count)).FontFeatures("tnum").FontSize(core.FontSize(c, theme.CaptionSize)).
			TextColor(k.TextMuted).SingleLine()
		if input.NumberInputButton(c, "chevron-right", core.Msg(c, "branch_navigator.next"), *index < count-1) {
			step(1)
		}
	})
	return r
}
