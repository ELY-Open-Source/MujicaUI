package mujica

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// ChartTypeSwitcherOptions configure a ChartTypeSwitcher.
type ChartTypeSwitcherOptions struct {
	Disabled bool
}

// ChartTypeSwitcher chooses how CandlestickChart draws: candles, OHLC bars, line, area or Heikin-Ashi.
func ChartTypeSwitcher(c *ui.Context, value *ChartType, opts ChartTypeSwitcherOptions) ChoiceResult {
	if value == nil {
		panic("mujica: ChartTypeSwitcher needs a value")
	}
	if *value < ChartCandles || *value > ChartHeikinAshi {
		panic(fmt.Sprintf("mujica: unknown ChartType %d", *value))
	}
	var items []ToggleItem[ChartType]
	for t := ChartCandles; t <= ChartHeikinAshi; t++ {
		items = append(items, ToggleItem[ChartType]{Value: t, Label: msg(c, fmt.Sprintf("chart_type_switcher.type%d", t))})
	}
	sel := []ChartType{*value}
	var r ChoiceResult
	ui.Row(c).Disabled(opts.Disabled).Children(func() {
		r = ToggleGroup(c, &sel, items, ToggleGroupOptions{Label: msg(c, "chart_type_switcher.label")})
	})
	if r.Changed() {
		*value = sel[0]
	}
	return r
}
