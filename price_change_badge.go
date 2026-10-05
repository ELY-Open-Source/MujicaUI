package mujica

import (
	"github.com/egoist/mygo/ui"

	"mujicaui/theme"
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
	dec := priceDecimals("PriceChangeBadge", opts.Decimals)
	sign := trendSign(change)
	if sign == 0 {
		sign = trendSign(percent)
	}
	col := trendColor(c, sign)
	u := Density(c).Unit()
	text := statisticSigned(change, dec)
	if !opts.HidePercent {
		text += " (" + statisticSigned(percent, 2) + "%)"
	}
	e := ui.Row(c).Gap(u).AlignItems(ui.Center).Padding(u/2, u*2).Radius(theme.ControlRadius).Shrink(0).
		Background(col.Alpha(0.12)).TextColor(col).Role(ui.RoleText).Label(trendWord(c, sign) + " " + text)
	e.Children(func() {
		ui.Icon(c, trendIcon(sign)).Size(14, 14)
		ui.Text(c, text).FontSize(FontSize(c, theme.CaptionSize)).FontFeatures("tnum").FontWeight(600).SingleLine()
	})
	return e
}
