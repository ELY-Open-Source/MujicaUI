package chart_test

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
)

func ExampleMarketHeatmap() {
	var path []string
	root := chart.HeatNode{Name: "Market", Children: []chart.HeatNode{{Name: "Technology", Children: []chart.HeatNode{{Name: "AAPL", Cap: 3000, Change: 1.2}, {Name: "MSFT", Cap: 2800, Change: -0.6}}}, {Name: "Energy", Cap: 450, Change: -2.1}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		r := chart.MarketHeatmap(c, &path, root, chart.MarketHeatmapOptions{Frame: chart.ChartFrame{Label: "Market"}})
		if r.Changed() || r.ValueChanged() || r.Retried() {
			fmt.Println(path, r.Point)
		}
	}
	ui.Render(view, 960, 440, 1)
}
