package mujica

func init() {
	addMessages(map[string][2]string{
		"dialog.close":                {"Close", "关闭"},
		"alert_dialog.cancel":         {"Cancel", "取消"},
		"alert_dialog.confirm":        {"Delete", "删除"},
		"alert_dialog.type_to_unlock": {"Type “%s” to confirm", "输入“%s”以确认"},
		"drawer.close":                {"Close", "关闭"},
		"popconfirm.confirm":          {"Confirm", "确定"},
		"popconfirm.cancel":           {"Cancel", "取消"},
	})
}
