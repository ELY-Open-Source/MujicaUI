package main

const projectsCategory = "Projects"

var projectsPages = append(append(projectsTaskPages, projectsFieldPages...), projectsTimePages...)

var projectsTaskPages = []page{
	{ID: "361", Name: "KanbanBoard", ZhName: "任务看板", Category: projectsCategory, Demo: projectsDemoKanbanBoard, Keys: "Tab 聚焦卡片；Enter/Space 打开；Alt+左右跨列；Alt+上下排序；可拖放到卡片之前或列末。", Code: `statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
tasks := []project.Task{{ID: "ENG-1", Title: "Plan the release", Status: "todo"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := project.KanbanBoard(c, &tasks, project.KanbanBoardOptions{Columns: []project.KanbanColumn{{Status: "todo"}, {Status: "done", WIP: 2}}, Statuses: statuses})
	if move, ok := r.Moved(); ok {
		log.Println("move", move.ID, move.To, move.Index)
	}
	if id, ok := r.Opened(); ok {
		log.Println("open", id)
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "362", Name: "TaskList", ZhName: "任务列表", Category: projectsCategory, Demo: projectsDemoTaskList, Keys: "Tab 在筛选器和任务控件间移动；空格切换完成；Enter 打开；F2 编辑标题。", Code: `statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
tasks := []project.Task{{ID: "ENG-1", Title: "Plan the release", Status: "todo"}}
filter := project.TaskShowAll
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := project.TaskList(c, tasks, &filter, project.TaskListOptions{GroupBy: project.TaskGroupStatus, Statuses: statuses})
	if id, ok := r.Changed(); ok {
		log.Println("edit", id)
	}
	if id, ok := r.Opened(); ok {
		log.Println("open", id)
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "363", Name: "TaskItem", ZhName: "任务行", Category: projectsCategory, Demo: projectsDemoTaskItem, Keys: "Tab 在复选框和标题间移动；空格切换完成；Enter 打开；F2 或双击开始编辑；Enter 保存；Escape 取消。", Code: `task := project.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: project.PriorityHigh}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := project.TaskItem(c, &task, project.TaskItemOptions{})
	if id, ok := r.Changed(); ok {
		log.Println("edit", id)
	}
	if id, ok := r.Opened(); ok {
		log.Println("open", id)
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "364", Name: "TaskDetailPanel", ZhName: "任务详情侧栏", Category: projectsCategory, Demo: projectsDemoTaskDetailPanel, Keys: "Tab 在字段间移动；选择框使用方向键和 Enter；Escape 返回关闭事件；子任务 Enter 新增。", Code: `statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
task := project.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: project.PriorityHigh}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := project.TaskDetailPanel(c, &task, project.TaskDetailPanelOptions{Statuses: statuses, Activity: func() { ui.Text(c, "Release planned") }})
	if r.Changed() {
		log.Println("task edited", task.ID)
	}
	if r.Closed() {
		log.Println("close", task.ID)
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "365", Name: "SubtaskList", ZhName: "子任务列表", Category: projectsCategory, Demo: projectsDemoSubtaskList, Keys: "Tab 移动；Space 切换完成；新增输入框 Enter 提交；Alt+上下排序；拖动握柄排序。", Code: `subtasks := []project.Subtask{{ID: "one", Title: "Review", Done: true}, {ID: "two", Title: "Ship"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if project.SubtaskList(c, &subtasks, project.SubtaskListOptions{}).Changed() {
		log.Println("subtasks changed", len(subtasks))
	}
}
ui.Render(view, 1000, 700, 1)`},
	{ID: "368", Name: "SprintBoard", ZhName: "迭代看板", Category: projectsCategory, Demo: projectsDemoSprintBoard, Keys: "看板卡片支持 Tab、Enter/Space、Alt+方向键与拖放。", Code: `statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
tasks := []project.Task{{ID: "ENG-1", Title: "Plan the release", Status: "todo"}}
sprint := project.Sprint{Name: "Release sprint", Goal: "Ship", Start: core.Date{Year: 2026, Month: 10, Day: 1}, End: core.Date{Year: 2026, Month: 10, Day: 14}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := project.SprintBoard(c, sprint, &tasks, project.SprintBoardOptions{Board: project.KanbanBoardOptions{Columns: []project.KanbanColumn{{Status: "todo"}, {Status: "done", WIP: 2}}, Statuses: statuses}})
	if move, ok := r.Moved(); ok {
		log.Println("move", move.ID, move.To)
	}
	if id, ok := r.Opened(); ok {
		log.Println("open", id)
	}
}
ui.Render(view, 1000, 700, 1)`},
}
