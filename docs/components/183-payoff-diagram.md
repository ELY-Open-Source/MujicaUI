# 183 PayoffDiagram

## 用途

到期损益图。

## 最小示例

```go
legs := []finance.OptionLeg{{Kind: finance.LegCall, Strike: 100, Premium: 3, Quantity: 1}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.PayoffDiagram(c, legs, finance.PayoffDiagramOptions{Frame: chart.ChartFrame{Label: "Call payoff"}, Spot: 103})
	if r.Changed() || r.Retried() {
		fmt.Println(r.Point)
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`legs` 为认购、认沽或标的组合；负数量表示卖出。`Frame.Label` 必填，`Height` 默认 240；`Spot` 默认 0；`Low/High` 默认行权价上下 30%；`Points` 默认 101；`Format` 可覆盖读数。

## 状态

默认与数据边界；图表的 Frame 支持加载、空、失败占位。

## 事件

if r.Changed() || r.Retried() { 调用方回调(r.Point) }

## 键盘操作

Tab 聚焦；方向键、Home/End 逐点读数；Escape 清除。

## 限制

计算到期损益，合约乘数由 Quantity 表达；价格采样数量由 Points 指定。
