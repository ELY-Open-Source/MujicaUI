package chart

import (
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/icons"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// TrendIndicatorOptions configure a TrendIndicator.
type TrendIndicatorOptions struct {
	// Precision is the decimals shown, 0 or more.
	Precision int
	// Unit follows the number, as "%".
	Unit string
	// Note explains the change, as "vs last week".
	Note string
	// Invert marks falls as good, as for costs or errors.
	Invert bool
	// Format replaces the built-in signed number.
	Format func(v float64) string
}

// TrendIndicator shows signed change, direction, and semantic color.
func TrendIndicator(c *ui.Context, change float64, opts TrendIndicatorOptions) ui.Element {
	if opts.Precision < 0 || math.IsNaN(change) || math.IsInf(change, 0) {
		panic("mujica: TrendIndicator needs a number and Precision ≥ 0")
	}
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	word, text := trendWords(c, change, opts)
	col, icon := k.TextMuted, "minus"
	switch {
	case change > 0:
		col, icon = k.Success, "arrow-up"
	case change < 0:
		col, icon = k.Danger, "arrow-down"
	}
	if opts.Invert && change != 0 {
		col = map[bool]ui.Color{true: k.Danger, false: k.Success}[change > 0]
	}
	e := ui.Row(c).Gap(u).AlignItems(ui.Center).Wrap().Role(ui.RoleText).Label(word + " " + text)
	e.Children(func() {
		ui.Row(c).Gap(u / 2).AlignItems(ui.Center).TextColor(col).Children(func() {
			ui.Icon(c, icons.Must(icon)).Size(14, 14)
			ui.Text(c, text).FontWeight(600).SingleLine()
		})
		if opts.Note != "" {
			ui.Text(c, opts.Note).FontSize(core.FontSize(c, theme.CaptionSize)).TextColor(k.TextMuted)
		}
	})
	return e
}

// trendWords formats direction and signed change.
func trendWords(c *ui.Context, change float64, opts TrendIndicatorOptions) (word, text string) {
	word = core.Msg(c, "trend_indicator.flat")
	switch {
	case change > 0:
		word = core.Msg(c, "trend_indicator.up")
	case change < 0:
		word = core.Msg(c, "trend_indicator.down")
	}
	if opts.Format != nil {
		return word, opts.Format(change)
	}
	return word, data.StatisticSigned(change, opts.Precision) + opts.Unit
}
