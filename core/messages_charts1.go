package core

func init() {
	addMessages(map[string][2]string{
		"chart.empty":             {"No data", "暂无数据"},
		"chart.hint":              {"Focus the chart and use the arrow keys to read values", "聚焦图表后用方向键读取数值"},
		"chart.missing":           {"no value", "无数值"},
		"pie_chart.other":         {"Other", "其他"},
		"pie_chart.total":         {"Total", "总计"},
		"scatter_chart.size":      {"size", "大小"},
		"scatter_chart.selected":  {"selected", "已选"},
		"calendar_heatmap.prev":   {"Previous year", "上一年"},
		"calendar_heatmap.next":   {"Next year", "下一年"},
		"calendar_heatmap.less":   {"Less", "少"},
		"calendar_heatmap.more":   {"More", "多"},
		"box_plot.readout":        {"%s · median %s · quartiles %s–%s · whiskers %s–%s · %d outliers", "%s · 中位数 %s · 四分位 %s–%s · 须线 %s–%s · 离群点 %d 个"},
		"waterfall_chart.running": {"running total", "累计"},
		"funnel_chart.readout":    {"%s · %s · %s of the stage before · %s of the first", "%s · %s · 较上一阶段 %s · 占首阶段 %s"},
		"bullet_chart.actual":     {"Actual", "实际"},
		"bullet_chart.target":     {"Target", "目标"},
		"bullet_chart.band":       {"Range %d: %s–%s", "区间 %d：%s–%s"},
		"chart.back":              {"Back", "返回"},
		"chart.opens":             {"Enter opens", "回车展开"},
		"sankey_chart.in":         {"in", "流入"},
		"sankey_chart.out":        {"out", "流出"},
		"network_graph.readout":   {"%s · %d links: %s", "%s · %d 条连线：%s"},
		"realtime_chart.pause":    {"Pause", "暂停"},
		"realtime_chart.resume":   {"Resume", "继续"},
	})
}
