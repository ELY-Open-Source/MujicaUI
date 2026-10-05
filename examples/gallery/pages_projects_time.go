package main

var projectsTimePages = []page{
	{ID: "375", Name: "IssueCard", ZhName: "任务卡片", Category: projectsCategory, Demo: projectsDemoIssueCard, Keys: "Tab 聚焦；Enter/Space 激活。", Code: `statuses := []mujica.WorkflowStatus{{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted}, {ID: "done", Name: "Done", Group: mujica.StatusCompleted}}
task := mujica.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: mujica.PriorityHigh}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.IssueCard(c, task, mujica.IssueCardOptions{Statuses: statuses}).Clicked() {
		log.Println("open", task.ID)
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "376", Name: "IssueIdBadge", ZhName: "任务编号", Category: projectsCategory, Demo: projectsDemoIssueIdBadge, Keys: "Copyable 开启时 Tab 聚焦，Enter/Space 复制。", Code: `id := "ENG-123"
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.IssueIdBadge(c, id, mujica.IssueIdBadgeOptions{Copyable: true}).Changed() {
		log.Println("copied", id)
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "377", Name: "MilestoneProgress", ZhName: "里程碑进度", Category: projectsCategory, Demo: projectsDemoMilestoneProgress, Keys: "静态展示，进度条通过无障碍名称与百分比读取。", Code: `milestone := mujica.Milestone{Title: "Release", Due: mujica.Date{Year: 2026, Month: 10, Day: 14}, Done: 3, Total: 8}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.MilestoneProgress(c, milestone, mujica.MilestoneProgressOptions{})
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "378", Name: "WorkloadView", ZhName: "成员负载", Category: projectsCategory, Demo: projectsDemoWorkloadView, Keys: "静态行通过无障碍名称读取负载与超出量；滚动查看更多成员。", Code: `loads := []mujica.Workload{{Member: mujica.ProjectMember{ID: "ann", Name: "Ann"}, Assigned: 46, Capacity: 40}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.WorkloadView(c, loads, mujica.WorkloadViewOptions{})
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "379", Name: "TimeTracker", ZhName: "工时计时", Category: projectsCategory, Demo: projectsDemoTimeTracker, Keys: "Tab 移动；Enter/Space 开始停止；停止后打开修正，编辑时长，Apply 提交。", Code: `var state mujica.TimeTrackerState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.TimeTracker(c, &state, mujica.TimeTrackerOptions{}).Changed() {
		log.Println("tracked", state.Total(c.Now()))
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "380", Name: "PomodoroTimer", ZhName: "番茄钟", Category: projectsCategory, Demo: projectsDemoPomodoroTimer, Keys: "Tab 移动；Enter/Space 开始、暂停、跳过与重置。", Code: `var state mujica.PomodoroState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.PomodoroTimer(c, &state, mujica.PomodoroTimerOptions{Focus: 25 * time.Minute})
	if r.Changed() {
		log.Println("timer", state.Phase(), state.Running())
	}
	if phase, ok := r.Finished(); ok {
		log.Println("finished", phase)
	}
}
ui.Render(view, 1000, 700, 1)`},
}
