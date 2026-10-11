package finance

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// PriceChangeBadgeOptions configure a PriceChangeBadge.
type PriceChangeBadgeOptions struct {
	// Decimals of the amount, zero meaning 2; percents show 2.
	Decimals int
	// HidePercent shows only the amount.
	HidePercent bool
}

// PriceChangeBadge shows a signed change and percent with an arrow.
func PriceChangeBadge(c *ui.Context, change, percent float64, opts PriceChangeBadgeOptions) *ui.Element {
	dec := core.PriceDecimals("PriceChangeBadge", opts.Decimals)
	sign := core.TrendSign(change)
	if sign == 0 {
		sign = core.TrendSign(percent)
	}
	col := core.TrendColor(c, sign)
	u := core.Density(c).Unit()
	text := data.StatisticSigned(change, dec)
	if !opts.HidePercent {
		text += " (" + data.StatisticSigned(percent, 2) + "%)"
	}
	e := ui.Row(c).Gap(u).AlignItems(ui.Center).Padding(u/2, u*2).Radius(theme.ControlRadius).Shrink(0).
		Background(col.Alpha(0.12)).TextColor(col).Role(ui.RoleText).Label(core.TrendWord(c, sign) + " " + text)
	e.Children(func() {
		ui.Icon(c, core.TrendIcon(sign)).Size(14, 14)
		ui.Text(c, text).FontSize(core.FontSize(c, theme.CaptionSize)).FontFeatures("tnum").FontWeight(600).SingleLine()
	})
	return e
}
