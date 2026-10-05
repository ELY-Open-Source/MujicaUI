package mujica

func init() {
	addMessages(map[string][2]string{
		"candlestick_chart.ohlc":    {"O %s H %s L %s C %s", "开 %s 高 %s 低 %s 收 %s"},
		"candlestick_chart.volume":  {"Vol", "量"},
		"heikin_ashi.smoothed":      {"Heikin-Ashi", "平均 K 线"},
		"heikin_ashi.raw":           {"Candles", "原 K 线"},
		"heikin_ashi.switch":        {"Candle style", "K 线样式"},
		"compare_symbol.label":      {"Compared symbols", "对比品种"},
		"compare_symbol.remove":     {"Remove %s", "移除 %s"},
		"price_alert_line.tag":      {"Alert %s", "预警 %s"},
		"chart_sync.link":           {"Sync with other charts", "与其他图表同步"},
		"chart_sync.linked":         {"Synced", "已同步"},
		"chart_sync.unlinked":       {"Not synced", "未同步"},
		"indicator_pane.remove":     {"Remove %s", "移除 %s"},
		"line_quote_chart.price":    {"Price", "价格"},
		"line_quote_chart.average":  {"Average", "均价"},
		"line_quote_chart.pre":      {"Pre-market", "盘前"},
		"line_quote_chart.after":    {"After hours", "盘后"},
		"depth_chart.bids":          {"Bids", "买盘"},
		"depth_chart.asks":          {"Asks", "卖盘"},
		"depth_chart.bid":           {"Bid", "买"},
		"depth_chart.ask":           {"Ask", "卖"},
		"depth_chart.total":         {"cumulative %s", "累计 %s"},
		"depth_chart.mid":           {"Mid", "中间价"},
		"volume_profile.value_area": {"Value area", "价值区"},
		"volume_profile.other":      {"Outside", "区外"},
		"volume_profile.poc":        {"POC", "控制点"},
		"volume_profile.vah":        {"VAH", "价值区上沿"},
		"volume_profile.val":        {"VAL", "价值区下沿"},
		"drawing_tools.note":        {"Note", "标注"},
	})
}
