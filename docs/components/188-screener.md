# 188 Screener

## 用途

选股器。

## 最小示例

```go
var state mujica.ScreenerState
fields := []mujica.ScreenerField{{ID: "pe", Label: "P/E", Decimals: 1}}
rows := []mujica.ScreenerRow{{Symbol: "AAPL", Values: map[string]float64{"pe": 34}}, {Symbol: "XOM", Values: map[string]float64{"pe": 13}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.Screener(c, &state, fields, rows, mujica.ScreenerOptions{})
	if r.Changed || r.Save || r.Table.Changed() || r.Table.Submitted() || r.Table.SortChanged() || r.Table.Retried() {
		fmt.Println(state.Conditions)
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`s.Conditions` 保存条件；`fields` 至少一项，ID 唯一；`rows.Values` 为字段值。`Version` 原地更新递增；`Status/Error/Disabled` 默认零值。

## 状态

默认、聚焦、禁用；异步结果支持加载、失败和空。

## 事件

if r.Changed || r.Save || r.Table.Changed() || r.Table.Submitted() || r.Table.SortChanged() || r.Table.Retried() { 调用方回调(state.Conditions) }

## 键盘操作

Tab 到达字段、比较、数值和保存；Select 方向键选择；NumberInput 方向键步进；表格方向键选择。

## 限制

全部条件取交集；保存由调用方处理；缺失值无法通过条件。
