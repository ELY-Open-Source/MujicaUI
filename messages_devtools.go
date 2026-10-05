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

		"alert_list.info":             {"Info", "提示"},
		"alert_list.warning":          {"Warning", "警告"},
		"alert_list.critical":         {"Critical", "严重"},
		"alert_list.ack":              {"Acknowledge", "确认"},
		"alert_list.acked":            {"Acknowledged", "已确认"},
		"alert_list.silence":          {"Silence 1 h", "静默 1 小时"},
		"alert_list.silenced":         {"Silenced until %s", "静默至 %s"},
		"incident_card.investigating": {"Investigating", "调查中"},
		"incident_card.identified":    {"Identified", "已定位"},
		"incident_card.monitoring":    {"Monitoring", "观察中"},
		"incident_card.resolved":      {"Resolved", "已解决"},
		"incident_card.none":          {"none", "无"},
		"incident_card.minor":         {"minor", "轻微"},
		"incident_card.major":         {"major", "严重"},
		"incident_card.critical":      {"critical", "致命"},
		"incident_card.impact":        {"Impact: %s", "影响：%s"},
		"incident_card.started":       {"Started %s", "开始于 %s"},

		"service_status.operational":   {"Operational", "正常"},
		"service_status.degraded":      {"Degraded", "降级"},
		"service_status.outage":        {"Outage", "中断"},
		"service_status.all_ok":        {"All systems operational", "所有服务正常"},
		"service_status.some_degraded": {"Some services degraded", "部分服务降级"},
		"service_status.some_down":     {"Some services down", "部分服务中断"},
		"service_status.counts":        {"%d up · %d degraded · %d down", "%d 正常 · %d 降级 · %d 中断"},
		"service_status.empty":         {"No services", "没有服务"},

		"alert_list.empty": {"No alerts", "没有告警"},
	})
}
