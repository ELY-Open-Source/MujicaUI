# 194 EarningsCalendar

## 用途

财报日历。

## 最小示例

```go
var state data.DataTableState[string]
rows := []finance.EarningsEvent{{ID: "a", Symbol: "AAPL", Name: "Apple", Date: core.Date{Year: 2026, Month: 10, Day: 30}, Timing: finance.EarningsAfterClose, EPSEstimate: 1.6, EPSActual: 1.64}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.EarningsCalendar(c, &state, rows, finance.EarningsCalendarOptions{})
	if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() {
		fmt.Println(state.Selected())
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`s` 保存表格状态；`events` 使用稳定 ID、Date 和盘前盘后时段；EPSActual 的 NaN 表示待公布。`Version` 原地更新递增；`Status/Error/Disabled` 默认零值。

## 状态

默认、聚焦、禁用；异步结果支持加载、失败和空。

## 事件

if r.Changed() || r.Submitted() || r.SortChanged() || r.Retried() { 调用方回调(state.Selected()) }

## 键盘操作

Tab 聚焦表格；方向键及 Home/End 选择；Enter 提交；表头支持键盘排序。

## 限制

NaN 表示待公布；零预期仅显示金额差；财报数据由调用方提供。
