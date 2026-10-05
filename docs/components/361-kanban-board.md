# 361 KanbanBoard · 任务看板

## 用途

任务看板用于项目管理界面。

## 最小示例

```go
statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
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
ui.Render(view, 1000, 700, 1)
```

## 参数

`tasks *[]Task` 为受控数据。`Columns` 至少一列；`Statuses` 提供其状态；`Members`、`Labels` 解析关联 ID。`Version` 标识原地数据更改；`Disabled` 禁用移动与打开；`Label` 为无障碍名称。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Moved()` 返回已应用的 `KanbanMove{ID, From, To, Index}`；`Opened()` 返回卡片 ID。

## 键盘操作

Tab 聚焦卡片；Enter/Space 打开；Alt+左右跨列；Alt+上下排序；可拖放到卡片之前或列末。

## 限制

修改原切片任务时递增 `Version`。列按 `Columns` 顺序排列；`WIP=0` 表示自由容量。网络与持久化由调用方处理。
