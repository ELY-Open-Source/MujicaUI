# 368 SprintBoard · 迭代看板

## 用途

迭代看板用于项目管理界面。

## 最小示例

```go
statuses := []mujica.WorkflowStatus{{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted}, {ID: "done", Name: "Done", Group: mujica.StatusCompleted}}
tasks := []mujica.Task{{ID: "ENG-1", Title: "Plan the release", Status: "todo"}}
sprint := mujica.Sprint{Name: "Release sprint", Goal: "Ship", Start: mujica.Date{Year: 2026, Month: 10, Day: 1}, End: mujica.Date{Year: 2026, Month: 10, Day: 14}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.SprintBoard(c, sprint, &tasks, mujica.SprintBoardOptions{Board: mujica.KanbanBoardOptions{Columns: []mujica.KanbanColumn{{Status: "todo"}, {Status: "done", WIP: 2}}, Statuses: statuses}})
	if move, ok := r.Moved(); ok {
		log.Println("move", move.ID, move.To)
	}
	if id, ok := r.Opened(); ok {
		log.Println("open", id)
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`Sprint{Name, Goal, Start, End}`；日期必填且顺序有效。`Board` 为看板参数；`Today` 默认当前日期。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

与 `KanbanBoard` 相同；移动更新受控任务并刷新完成数。

## 键盘操作

看板卡片支持 Tab、Enter/Space、Alt+方向键与拖放。

## 限制

完成度按工作流状态统计，取消状态排除。日期为包含首尾的区间；外部原地修改递增 `Board.Version`。
