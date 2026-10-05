package mujica

func init() {
	addMessages(map[string][2]string{
		"candlestick_chart.ohlc":   {"O %s H %s L %s C %s", "开 %s 高 %s 低 %s 收 %s"},
		"candlestick_chart.volume": {"Vol", "量"},
		"heikin_ashi.smoothed":     {"Heikin-Ashi", "平均 K 线"},
		"heikin_ashi.raw":          {"Candles", "原 K 线"},
		"heikin_ashi.switch":       {"Candle style", "K 线样式"},
		"compare_symbol.label":     {"Compared symbols", "对比品种"},
		"compare_symbol.remove":    {"Remove %s", "移除 %s"},
		"price_alert_line.tag":     {"Alert %s", "预警 %s"},
		"chart_sync.link":          {"Sync with other charts", "与其他图表同步"},
		"chart_sync.linked":        {"Synced", "已同步"},
		"chart_sync.unlinked":      {"Not synced", "未同步"},
		"drawing_tools.note":       {"Note", "标注"},
	})
}
