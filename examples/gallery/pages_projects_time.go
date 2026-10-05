package main

var projectsTimePages = []page{
	{ID: "375", Name: "IssueCard", ZhName: "任务卡片", Category: projectsCategory, Demo: projectsDemoIssueCard, Keys: "Tab 聚焦；Enter/Space 激活。", Code: `statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
task := project.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: project.PriorityHigh}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if project.IssueCard(c, task, project.IssueCardOptions{Statuses: statuses}).Clicked() {
		log.Println("open", task.ID)
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "376", Name: "IssueIdBadge", ZhName: "任务编号", Category: projectsCategory, Demo: projectsDemoIssueIdBadge, Keys: "Copyable 开启时 Tab 聚焦，Enter/Space 复制。", Code: `id := "ENG-123"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if project.IssueIdBadge(c, id, project.IssueIdBadgeOptions{Copyable: true}).Changed() {
		log.Println("copied", id)
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "377", Name: "MilestoneProgress", ZhName: "里程碑进度", Category: projectsCategory, Demo: projectsDemoMilestoneProgress, Keys: "静态展示，进度条通过无障碍名称与百分比读取。", Code: `milestone := project.Milestone{Title: "Release", Due: core.Date{Year: 2026, Month: 10, Day: 14}, Done: 3, Total: 8}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	project.MilestoneProgress(c, milestone, project.MilestoneProgressOptions{})
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "378", Name: "WorkloadView", ZhName: "成员负载", Category: projectsCategory, Demo: projectsDemoWorkloadView, Keys: "静态行通过无障碍名称读取负载与超出量；滚动查看更多成员。", Code: `loads := []project.Workload{{Member: project.ProjectMember{ID: "ann", Name: "Ann"}, Assigned: 46, Capacity: 40}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	project.WorkloadView(c, loads, project.WorkloadViewOptions{})
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "379", Name: "TimeTracker", ZhName: "工时计时", Category: projectsCategory, Demo: projectsDemoTimeTracker, Keys: "Tab 移动；Enter/Space 开始停止；停止后打开修正，编辑时长，Apply 提交。", Code: `var state project.TimeTrackerState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if project.TimeTracker(c, &state, project.TimeTrackerOptions{}).Changed() {
		log.Println("tracked", state.Total(c.Now()))
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "380", Name: "PomodoroTimer", ZhName: "番茄钟", Category: projectsCategory, Demo: projectsDemoPomodoroTimer, Keys: "Tab 移动；Enter/Space 开始、暂停、跳过与重置。", Code: `var state project.PomodoroState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := project.PomodoroTimer(c, &state, project.PomodoroTimerOptions{Focus: 25 * time.Minute})
	if r.Changed() {
		log.Println("timer", state.Phase(), state.Running())
	}
	if phase, ok := r.Finished(); ok {
		log.Println("finished", phase)
	}
}
ui.Render(view, 1000, 700, 1)`},
}
