# 364 TaskDetailPanel · 任务详情侧栏

## 用途

任务详情侧栏用于项目管理界面。

## 最小示例

```go
statuses := []project.WorkflowStatus{{ID: "todo", Name: "To do", Group: project.StatusUnstarted}, {ID: "done", Name: "Done", Group: project.StatusCompleted}}
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
ui.Render(view, 1000, 700, 1)
```

## 参数

`task *Task` 为受控任务；`Statuses`、`Members`、`Labels` 为关联数据；`Activity func()` 构建活动记录；`Today`、`Disabled` 控制日期与交互；`Version` 用于外部原地修改任务、成员或子任务。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 表示字段更改；`Closed()` 表示关闭按钮或 Escape。由调用方控制侧栏显隐。

## 键盘操作

Tab 在字段间移动；选择框使用方向键和 Enter；Escape 返回关闭事件；子任务 Enter 新增。

## 限制

组件是可嵌入的侧栏内容；容器负责布局与滚动。活动记录通过 `Activity` 插槽构建；外部原地更改后递增 `Version`。
