# 379 TimeTracker · 工时计时

## 用途

工时计时用于项目管理界面。

## 最小示例

```go
var state mujica.TimeTrackerState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.TimeTracker(c, &state, mujica.TimeTrackerOptions{}).Changed() {
		log.Println("tracked", state.Total(c.Now()))
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`state *TimeTrackerState`；`Now` 固定时钟供测试；`Disabled` 禁用操作；`Label` 为读数名称。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 表示开始、停止或手动修正；`Entries` 和 `Adjust` 由调用方持有。

## 键盘操作

Tab 移动；Enter/Space 开始停止；停止后打开修正，编辑时长，Apply 提交。

## 限制

只有最后一项可处于运行态；时段按顺序且互不重叠；累计时长非负。隐藏时保留计时状态。
