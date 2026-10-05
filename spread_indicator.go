package mujica

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"mujicaui/icons"
	"mujicaui/theme"
)

// SpreadIndicatorOptions configure a SpreadIndicator.
type SpreadIndicatorOptions struct {
	// Decimals of the spread; zero means 2.
	Decimals int
	// WarnPercent flags spreads this wide, in percent; zero never.
	WarnPercent float64
}

// SpreadIndicator shows ask minus bid and its percent.
func SpreadIndicator(c *ui.Context, bid, ask float64, opts SpreadIndicatorOptions) *ui.Element {
	if opts.WarnPercent < 0 {
		panic("mujica: SpreadIndicator with negative WarnPercent")
	}
	dec := priceDecimals("SpreadIndicator", opts.Decimals)
	k := Tokens(c)
	u := Density(c).Unit()
	spread := ask - bid
	pct := 0.0
	if mid := (ask + bid) / 2; mid > 0 {
		pct = spread / mid * 100
	}
	text := fmt.Sprintf(msg(c, "spread_indicator.value"), statisticFormat(spread, dec), statisticFormat(pct, 3))
	note, col, icon := "", k.TextMuted, ""
	switch {
	case spread < 0:
		note, col, icon = msg(c, "spread_indicator.crossed"), k.Danger, "circle-alert"
	case opts.WarnPercent > 0 && pct >= opts.WarnPercent:
		note, col, icon = msg(c, "spread_indicator.wide"), k.Warning, "triangle-alert"
	}
	e := ui.Row(c).Gap(u * 1.5).AlignItems(ui.Center).Shrink(0).Role(ui.RoleStatus).Label(text + " " + note)
	e.Children(func() {
		ui.Text(c, msg(c, "spread_indicator.title")).TextColor(k.TextMuted).FontSize(FontSize(c, theme.CaptionSize))
		ui.Text(c, text).FontFeatures("tnum").FontWeight(500).SingleLine()
		if note != "" {
			ui.Row(c).Gap(u / 2).AlignItems(ui.Center).TextColor(col).Children(func() {
				ui.Icon(c, icons.Must(icon)).Size(14, 14)
				ui.Text(c, note).FontSize(FontSize(c, theme.CaptionSize)).FontWeight(600).SingleLine()
			})
		}
	})
	return e
}
