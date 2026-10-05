# 375 IssueCard · 任务卡片

## 用途

任务卡片用于项目管理界面。

## 最小示例

```go
statuses := []mujica.WorkflowStatus{{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted}, {ID: "done", Name: "Done", Group: mujica.StatusCompleted}}
task := mujica.Task{ID: "ENG-1", Title: "Plan the release", Status: "todo", Priority: mujica.PriorityHigh}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.IssueCard(c, task, mujica.IssueCardOptions{Statuses: statuses}).Clicked() {
		log.Println("open", task.ID)
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`task Task`；`Statuses` 解析状态；`Members`、`Labels` 解析关联数据；`Selected` 显示选中底色与边标；`Disabled` 禁用激活。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

返回 `Element.Clicked()` 为卡片激活。

## 键盘操作

Tab 聚焦；Enter/Space 激活。

## 限制

长标题最多两行；头像组折叠到三个；字段与关联数据由调用方提供。
