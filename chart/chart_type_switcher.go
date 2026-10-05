package chart

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/input"
)

// ChartTypeSwitcherOptions configure a ChartTypeSwitcher.
type ChartTypeSwitcherOptions struct {
	Disabled bool
}

// ChartTypeSwitcher chooses candles, OHLC, line, area, or Heikin-Ashi.
func ChartTypeSwitcher(c *ui.Context, value *ChartType, opts ChartTypeSwitcherOptions) core.ChoiceResult {
	if value == nil {
		panic("mujica: ChartTypeSwitcher needs a value")
	}
	if *value < ChartCandles || *value > ChartHeikinAshi {
		panic(fmt.Sprintf("mujica: unknown ChartType %d", *value))
	}
	var items []input.ToggleItem[ChartType]
	for t := ChartCandles; t <= ChartHeikinAshi; t++ {
		items = append(items, input.ToggleItem[ChartType]{Value: t, Label: core.Msg(c, fmt.Sprintf("chart_type_switcher.type%d", t))})
	}
	sel := []ChartType{*value}
	var r core.ChoiceResult
	ui.Row(c).Disabled(opts.Disabled).Children(func() {
		r = input.ToggleGroup(c, &sel, items, input.ToggleGroupOptions{Label: core.Msg(c, "chart_type_switcher.label")})
	})
	if r.Changed() {
		*value = sel[0]
		slog.Debug("mujica: chart_type_switcher changed", "type", *value)
	}
	return r
}
