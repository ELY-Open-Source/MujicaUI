# 193 EconomicCalendar

## 用途

经济日历。

## 最小示例

```go
var state finance.EconomicCalendarState
rows := []finance.EconomicEvent{{ID: "jobs", Country: "US", Title: "Payrolls", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Importance: 3, Forecast: "150K", Previous: "142K", Actual: "254K"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.EconomicCalendar(c, &state, rows, finance.EconomicCalendarOptions{})
	if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
		fmt.Println(state.Table.Selected())
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`s` 保存国家、重要性和表格状态；`events` 的 ID 稳定且重要性 1—3。`Version` 原地更新递增；`Status/Error/Disabled` 默认零值。

## 状态

默认、聚焦、禁用；异步结果支持加载、失败和空。

## 事件

if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() { 调用方回调(state.Table.Selected()) }

## 键盘操作

Tab 到达国家过滤、重要性和表格；方向键选择；Enter 提交行。

## 限制

事件数据与发布结果由调用方提供；过滤与排序缓存直到数据或条件变更。
