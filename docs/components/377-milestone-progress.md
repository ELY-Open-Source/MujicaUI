# 377 MilestoneProgress · 里程碑进度

## 用途

里程碑进度用于项目管理界面。

## 最小示例

```go
milestone := project.Milestone{Title: "Release", Due: core.Date{Year: 2026, Month: 10, Day: 14}, Done: 3, Total: 8}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	project.MilestoneProgress(c, milestone, project.MilestoneProgressOptions{})
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`Milestone{Title, Due, Done, Total}`；`Today` 默认当前日期；零 Due 表示无日期。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

静态展示。

## 键盘操作

静态展示，进度条通过无障碍名称与百分比读取。

## 限制

要求 0 ≤ Done ≤ Total；零总量显示 0%；完成后显示达成标记。
