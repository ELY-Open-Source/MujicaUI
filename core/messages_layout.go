package core

func init() {
	addMessages(map[string][2]string{
		"app_shell.collapse":   {"Collapse sidebar", "收起侧栏"},
		"app_shell.expand":     {"Expand sidebar", "展开侧栏"},
		"status_bar.label":     {"Status bar", "状态栏"},
		"status_bar.hidden":    {"%d more items", "另有 %d 项"},
		"split_pane.divider":   {"Divider", "分隔条"},
		"split_pane.collapsed": {"Collapsed", "已折叠"},
	})
}
