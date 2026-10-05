package main

var projectsFieldPages = []page{
	{ID: "369", Name: "BurndownChart", ZhName: "燃尽图", Category: projectsCategory, Demo: projectsDemoBurndownChart, Keys: "Tab 聚焦图例与图表；Space 切换图例；方向键逐日读取；Home/End 首尾；Escape 清除读数。", Code: `remaining := []float64{40, 32, 25, 18}
opts := project.BurndownChartOptions{Frame: chart.ChartFrame{Label: "Release burndown"}, Start: core.Date{Year: 2026, Month: 10, Day: 1}, End: core.Date{Year: 2026, Month: 10, Day: 14}, Total: 40}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := project.BurndownChart(c, remaining, opts)
	if r.Changed() {
		log.Println("point", r.Point)
	}
	if name, ok := r.Toggled(); ok {
		log.Println("series", name)
	}
	if r.Retried() {
		log.Println("retry")
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "370", Name: "PriorityIndicator", ZhName: "优先级指示", Category: projectsCategory, Demo: projectsDemoPriorityIndicator, Keys: "静态展示，通过无障碍名称读取完整等级。", Code: `priority := project.PriorityHigh
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	project.PriorityIndicator(c, priority, project.PriorityIndicatorOptions{})
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "371", Name: "StatusSelect", ZhName: "状态选择", Category: projectsCategory, Demo: projectsDemoStatusSelect, Keys: "Enter/Space 或上下键打开；上下、Home/End 定位；Enter 确认；Escape 关闭。", Code: `statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
value := "todo"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if project.StatusSelect(c, &value, statuses, project.StatusSelectOptions{}).Changed() {
		log.Println("status", value)
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "372", Name: "AssigneePicker", ZhName: "负责人选择", Category: projectsCategory, Demo: projectsDemoAssigneePicker, Keys: "Enter/Space 打开；输入搜索；下键进入选项；Space/Enter 勾选；Backspace 移除最后一项；Escape 关闭。", Code: `members := []project.ProjectMember{{ID: "ann", Name: "Ann"}, {ID: "ben", Name: "Ben"}}
ids := []string{"ann"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if project.AssigneePicker(c, &ids, members, project.AssigneePickerOptions{Max: 2}).Changed() {
		log.Println("assignees", ids)
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "373", Name: "DueDatePicker", ZhName: "截止日期选择", Category: projectsCategory, Demo: projectsDemoDueDatePicker, Keys: "Tab 到快捷项、日期输入与日历按钮；Enter 选择；日历方向键移动；Escape 关闭。", Code: `var due core.Date
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if project.DueDatePicker(c, &due, project.DueDatePickerOptions{}).Changed() {
		log.Println("due", due)
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "374", Name: "LabelManager", ZhName: "标签管理", Category: projectsCategory, Demo: projectsDemoLabelManager, Keys: "Tab 移动；Space 勾选；新增输入 Enter 创建；编辑颜色按钮打开编辑；Enter 保存名称；Escape 结束编辑。", Code: `labels := []project.TaskLabel{{ID: "ui", Name: "UI", Color: datetime.EventInfo}}
selected := []string{}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if project.LabelManager(c, &labels, &selected, project.LabelManagerOptions{}).Changed() {
		log.Println("labels changed", len(labels), selected)
	}
}
ui.Render(view, 1000, 700, 1)`},
}
