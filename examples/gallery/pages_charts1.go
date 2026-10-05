package main

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ZacharyZhang-NY/MujicaUI/chart"
	"github.com/ZacharyZhang-NY/MujicaUI/core"
	"github.com/ZacharyZhang-NY/MujicaUI/data"
	"github.com/ZacharyZhang-NY/MujicaUI/theme"
)

const c1Category = "Charts I"

const c1Keys = "Tab focuses the chart (legend check boxes first); arrows move between points and the readout is announced; Home/End jump; Escape clears."

var charts1Pages = []page{
	{ID: "101", Name: "LineChart", ZhName: "折线图", Category: c1Category, Demo: c1LineDemo, Code: c1LineCode, Keys: c1Keys + " Space on a legend entry hides or shows its series."},
	{ID: "102", Name: "AreaChart", ZhName: "面积图", Category: c1Category, Demo: c1AreaDemo, Code: c1AreaCode, Keys: c1Keys},
	{ID: "103", Name: "BarChart", ZhName: "柱状图", Category: c1Category, Demo: c1BarDemo, Code: c1BarCode, Keys: c1Keys + " Arrows step bar by bar."},
	{ID: "104", Name: "PieChart", ZhName: "饼图", Category: c1Category, Demo: c1PieDemo, Code: c1PieCode, Keys: c1Keys},
	{ID: "105", Name: "ScatterChart", ZhName: "散点图", Category: c1Category, Demo: c1ScatterDemo, Code: c1ScatterCode, Keys: c1Keys + " Enter or Space selects the point; drag a box to select."},
	{ID: "106", Name: "RadarChart", ZhName: "雷达图", Category: c1Category, Demo: c1RadarDemo, Code: c1RadarCode, Keys: c1Keys},
	{ID: "107", Name: "HeatmapChart", ZhName: "热力图", Category: c1Category, Demo: c1HeatDemo, Code: c1HeatCode, Keys: "Tab focuses the map; Left/Right move along a row, Up/Down along a column; Home/End jump; Escape clears."},
	{ID: "108", Name: "CalendarHeatmap", ZhName: "日历热力图", Category: c1Category, Demo: c1CalendarDemo, Code: c1CalendarCode, Keys: "Left/Right move a week, Up/Down a day; Home/End jump to the year's ends; the year buttons take Enter or Space."},
	{ID: "109", Name: "Histogram", ZhName: "直方图", Category: c1Category, Demo: c1HistogramDemo, Code: c1HistogramCode, Keys: c1Keys + " Arrows step bin by bin."},
	{ID: "110", Name: "BoxPlot", ZhName: "箱线图", Category: c1Category, Demo: c1BoxDemo, Code: c1BoxCode, Keys: c1Keys + " Arrows step group by group."},
	{ID: "111", Name: "WaterfallChart", ZhName: "瀑布图", Category: c1Category, Demo: c1WaterfallDemo, Code: c1WaterfallCode, Keys: c1Keys},
	{ID: "112", Name: "FunnelChart", ZhName: "漏斗图", Category: c1Category, Demo: c1FunnelDemo, Code: c1FunnelCode, Keys: c1Keys + " Arrows step stage by stage."},
	{ID: "113", Name: "BulletChart", ZhName: "子弹图", Category: c1Category, Demo: c1BulletDemo, Code: c1BulletCode, Keys: c1Keys + " Arrows read the value, the target, then each range."},
	{ID: "114", Name: "Sparkline", ZhName: "迷你图", Category: c1Category, Demo: c1SparkDemo, Code: c1SparkCode, Keys: "Tab focuses the sparkline; Left/Right read points; Home/End jump; Escape shows the last value again."},
	{ID: "115", Name: "Gauge", ZhName: "仪表盘", Category: c1Category, Demo: c1GaugeDemo, Code: c1GaugeCode, Keys: c1Keys + " Arrows read the value, then each segment."},
	{ID: "116", Name: "Treemap", ZhName: "矩形树图", Category: c1Category, Demo: c1TreemapDemo, Code: c1TreemapCode, Keys: c1Keys + " Enter or Space drills into a branch; Backspace goes back; Back and the breadcrumb take Tab."},
	{ID: "117", Name: "Sunburst", ZhName: "旭日图", Category: c1Category, Demo: c1SunburstDemo, Code: c1SunburstCode, Keys: c1Keys + " Enter or Space drills in; Backspace goes back."},
	{ID: "118", Name: "SankeyChart", ZhName: "桑基图", Category: c1Category, Demo: c1SankeyDemo, Code: c1SankeyCode, Keys: c1Keys + " Reading a node or link highlights every path through it."},
	{ID: "119", Name: "NetworkGraph", ZhName: "网络图", Category: c1Category, Demo: c1NetworkDemo, Code: c1NetworkCode, Keys: c1Keys + " Shift+arrows move the node; +/- zoom; 0 resets the view; drag nodes or the background."},
	{ID: "120", Name: "RealtimeChart", ZhName: "实时图表", Category: c1Category, Demo: c1RealtimeDemo, Code: c1RealtimeCode, Keys: "Tab reaches Pause, the legend and the chart; chart keys as in LineChart."},
}

var c1Statuses = []string{"Ready", "Loading", "Failed", "Empty"}

// c1Frame shows the status picker and returns the frame.
func c1Frame(c *ui.Context, label string, status *string) (chart.ChartFrame, bool) {
	ui.Select(c, status, c1Statuses).Width(160).AlignSelf(ui.Start)
	f := chart.ChartFrame{Label: label}
	switch *status {
	case "Loading":
		f.Status = data.DataLoading
	case "Failed":
		f.Status, f.Error = data.DataFailed, "The archive did not answer"
	}
	return f, *status != "Empty"
}

// c1Section lays out a title, controls, then the chart.
func c1Section(c *ui.Context, title string, controls, body func()) {
	ui.Column(c).Gap(12).Children(func() {
		ui.Text(c, title).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size))
		ui.Row(c).Gap(16).Wrap().AlignItems(ui.Center).Children(controls)
		ui.Column(c).Width(560).Children(body)
	})
}

// c1Log records a chart's events in the gallery log.
func c1Log(log func(string), name string, r chart.ChartResult) {
	if r.Changed() {
		log(fmt.Sprintf("%s: point %d", name, r.Point))
	}
	if s, ok := r.Toggled(); ok {
		log(name + ": toggled " + s)
	}
	if r.ValueChanged() {
		log(name + ": value changed")
	}
	if r.Retried() {
		log(name + ": retry")
	}
}
