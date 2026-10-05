package mujica

func init() {
	addMessages(map[string][2]string{
		"chart_legend.more":  {"+%d more", "另有 %d 项"},
		"chart_legend.fewer": {"Show fewer", "收起"},
		"chart_brush.window": {"%s: %s to %s", "%s：%s 至 %s"},
		"chart_brush.start":  {"Range start %s", "区间起点 %s"},
		"chart_brush.end":    {"Range end %s", "区间终点 %s"},
	})
}
