package mujica

func init() {
	addMessages(map[string][2]string{
		"key_value_input.key":    {"Key", "键"},
		"key_value_input.value":  {"Value", "值"},
		"key_value_input.row":    {"Row %d", "第 %d 行"},
		"key_value_input.remove": {"Remove row %d", "删除第 %d 行"},
		"key_value_input.add":    {"Add row", "添加一行"},

		"stat_card.change": {"Change %s", "变化 %s"},
		"devtools.ago":     {"%s ago", "%s 前"},

		"alert_list.info":     {"Info", "提示"},
		"alert_list.warning":  {"Warning", "警告"},
		"alert_list.critical": {"Critical", "严重"},
		"alert_list.ack":      {"Acknowledge", "确认"},
		"alert_list.acked":    {"Acknowledged", "已确认"},
		"alert_list.silence":  {"Silence 1 h", "静默 1 小时"},
		"alert_list.silenced": {"Silenced until %s", "静默至 %s"},
		"alert_list.empty":    {"No alerts", "没有告警"},
	})
}
