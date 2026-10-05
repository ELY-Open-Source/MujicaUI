# 371 StatusSelect · 状态选择

## 用途

状态选择用于项目管理界面。

## 最小示例

```go
statuses := []mujica.WorkflowStatus{{ID: "todo", Name: "To do", Group: mujica.StatusUnstarted}, {ID: "done", Name: "Done", Group: mujica.StatusCompleted}}
value := "todo"
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.StatusSelect(c, &value, statuses, mujica.StatusSelectOptions{}).Changed() {
		log.Println("status", value)
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`value *string`、`statuses []WorkflowStatus`；状态 `Group` 为 Backlog、Unstarted、Started、Completed、Canceled；`Disabled` 和 `Label` 控制交互与名称。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 表示受控状态 ID 更改。

## 键盘操作

Enter/Space 或上下键打开；上下、Home/End 定位；Enter 确认；Escape 关闭。

## 限制

状态 ID 唯一且非空；当前值须属于状态列表。状态按五个工作流阶段分组。
