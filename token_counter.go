package mujica

import (
	"fmt"
	"strconv"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// tokenWarnShare is the share of the limit from which the count warns.
const tokenWarnShare = 0.9

// TokenCounter shows the caller's token estimate against the limit: muted, warning past 90%, danger over it.
func TokenCounter(c *ui.Context, used, limit int) *ui.Element {
	if used < 0 || limit <= 0 {
		panic("mujica: TokenCounter needs used >= 0 and limit > 0")
	}
	k := Tokens(c)
	u := Density(c).Unit()
	label := fmt.Sprintf(msg(c, "token_counter.label"), tokenThousands(used), tokenThousands(limit))
	color, icon := k.TextMuted, ""
	switch {
	case used > limit:
		color, icon = k.Danger, "circle-alert"
		label += ", " + fmt.Sprintf(msg(c, "token_counter.over"), tokenThousands(used-limit))
	case float64(used) >= tokenWarnShare*float64(limit):
		color, icon = k.Warning, "triangle-alert"
	}
	row := ui.Row(c).Gap(u).AlignItems(ui.Center).Shrink(0).TextColor(color).Role(ui.RoleStatus).Label(label)
	row.Children(func() {
		if icon != "" {
			ui.Icon(c, icons.Must(icon)).Size(12, 12)
		}
		ui.Text(c, "≈ "+tokenThousands(used)+" / "+tokenThousands(limit)).FontSize(FontSize(c, theme.CaptionSize)).FontFeatures("tnum").SingleLine()
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
