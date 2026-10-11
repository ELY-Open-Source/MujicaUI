package chat

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// tokenWarnShare sets the warning threshold.
const tokenWarnShare = 0.9

// TokenCounter displays estimates, limits and threshold warnings.
func TokenCounter(c *ui.Context, used, limit int) *ui.Element {
	if used < 0 || limit <= 0 {
		panic("mujica: TokenCounter needs used >= 0 and limit > 0")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	label := fmt.Sprintf(core.Msg(c, "token_counter.label"), tokenThousands(used), tokenThousands(limit))
	color, icon := k.TextMuted, ""
	switch {
	case used > limit:
		color, icon = k.Danger, "circle-alert"
		label += ", " + fmt.Sprintf(core.Msg(c, "token_counter.over"), tokenThousands(used-limit))
	case float64(used) >= tokenWarnShare*float64(limit):
		color, icon = k.Warning, "triangle-alert"
	}
	row := ui.Row(c).Gap(u).AlignItems(ui.Center).Shrink(0).TextColor(color).Role(ui.RoleStatus).Label(label)
	row.Children(func() {
		if icon != "" {
			ui.Icon(c, icons.Must(icon)).Size(12, 12)
		}
		ui.Text(c, "≈ "+tokenThousands(used)+" / "+tokenThousands(limit)).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum").SingleLine()
	})
	return row
}

// tokenThousands writes n with comma thousands separators.
func tokenThousands(n int) string {
	s := strconv.Itoa(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}
