package mujica

func init() {
	addMessages(map[string][2]string{
		"loading": {"Loading", "加载中"},

		"alert.close":   {"Dismiss", "关闭"},
		"alert.info":    {"Information", "信息"},
		"alert.success": {"Success", "成功"},
		"alert.warning": {"Warning", "警告"},
		"alert.error":   {"Error", "错误"},

		"toast.close":  {"Dismiss notification", "关闭通知"},
		"toast.region": {"Notifications", "通知"},

		"progress.indeterminate": {"In progress", "进行中"},
		"meter.normal":           {"Normal", "正常"},
		"meter.warning":          {"High", "偏高"},
		"meter.critical":         {"Critical", "严重"},
		"meter.low_warning":      {"Low", "偏低"},

		"empty_state.no_data.title":          {"Nothing here yet", "暂无内容"},
		"empty_state.no_data.description":    {"Items you add will appear here.", "添加的内容会显示在这里。"},
		"empty_state.no_results.title":       {"No results", "没有结果"},
		"empty_state.no_results.description": {"Try other words or clear the filters.", "请换个关键词，或清除筛选条件。"},
		"empty_state.no_results.query":       {"Nothing matches “%s”.", "没有与“%s”匹配的内容。"},

		"result.success.title":   {"Done", "操作成功"},
		"result.failure.title":   {"Something went wrong", "操作失败"},
		"result.forbidden.title": {"No access", "无访问权限"},

		"status_indicator.success": {"Healthy", "正常"},
		"status_indicator.warning": {"Degraded", "警告"},
		"status_indicator.error":   {"Failing", "错误"},
		"status_indicator.offline": {"Offline", "离线"},
		"status_indicator.busy":    {"Working", "处理中"},

		"notification_center.title":         {"Notifications", "通知"},
		"notification_center.unread":        {"%d unread", "%d 条未读"},
		"notification_center.unread_mark":   {"Unread", "未读"},
		"notification_center.mark_read":     {"Mark read", "标为已读"},
		"notification_center.mark_unread":   {"Mark unread", "标为未读"},
		"notification_center.mark_all_read": {"Mark all read", "全部已读"},
		"notification_center.clear":         {"Clear", "清除"},
		"notification_center.clear_all":     {"Clear all", "全部清除"},
		"notification_center.empty":         {"You're all caught up", "没有新通知"},
		"notification_center.empty_hint":    {"New notifications will appear here.", "新的通知会显示在这里。"},

		"loading_overlay.label": {"Loading…", "正在加载…"},
	})
}
