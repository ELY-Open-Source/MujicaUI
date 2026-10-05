package mujica

func init() {
	addMessages(map[string][2]string{
		"violin_plot.readout":           {"%s · median %s · quartiles %s–%s · %d values", "%s · 中位数 %s · 四分位 %s–%s · %d 个值"},
		"parallel_coordinates.filtered": {"filtered out", "已筛除"},
		"parallel_coordinates.brush":    {"%s filter %s–%s", "%s 筛选 %s–%s"},
		"chord_diagram.out":             {"out", "流出"},
		"chord_diagram.in":              {"in", "流入"},
		"chart_annotation.readout":      {"Annotation %d · %s", "标注 %d · %s"},
		"chart_annotation.x":            {"x %s", "x %s"},
		"chart_annotation.y":            {"y %s", "y %s"},
		"spark_bar.empty":               {"No data", "暂无数据"},
		"chart_empty_state.empty":       {"No data", "暂无数据"},
		"chart_legend.more":             {"+%d more", "另有 %d 项"},
		"chart_legend.fewer":            {"Show fewer", "收起"},
		"chart_brush.window":            {"%s: %s to %s", "%s：%s 至 %s"},
		"chart_brush.start":             {"Range start %s", "区间起点 %s"},
		"chart_brush.end":               {"Range end %s", "区间终点 %s"},
		"trend_indicator.up":            {"up", "上升"},
		"trend_indicator.down":          {"down", "下降"},
		"trend_indicator.flat":          {"unchanged", "持平"},
	})
}
