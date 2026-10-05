# 370 PriorityIndicator · 优先级指示

## 用途

优先级指示用于项目管理界面。

## 最小示例

```go
priority := mujica.PriorityHigh
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.PriorityIndicator(c, priority, mujica.PriorityIndicatorOptions{})
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`p Priority` 决定等级；`IconOnly` 默认 false，开启后保留无障碍名称。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

静态展示。

## 键盘操作

静态展示，通过无障碍名称读取完整等级。

## 限制

等级仅接受 None、Low、Medium、High、Urgent。
