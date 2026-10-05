# 373 DueDatePicker · 截止日期选择

## 用途

截止日期选择用于项目管理界面。

## 最小示例

```go
var due core.Date
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if project.DueDatePicker(c, &due, project.DueDatePickerOptions{}).Changed() {
		log.Println("due", due)
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`value *Date`；`Today` 默认当前日期；`WeekStart` 默认本地周起始；`Done`、`Disabled`、`Label` 控制展示与交互。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 表示日期更改；返回 `Input` 为日期编辑器。

## 键盘操作

Tab 到快捷项、日期输入与日历按钮；Enter 选择；日历方向键移动；Escape 关闭。

## 限制

零日期表示清空；Next week 选择下一个周起始日。`Done` 让已完成任务显示普通日期。
