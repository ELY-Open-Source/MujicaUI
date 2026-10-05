# 283 CostBreakdown 费用分解

## 用途

按模型、工具或步骤汇总调用次数、token 与费用，表格可排序，底部显示合计。

## 最小示例

```go
var state mujica.CostBreakdownState
items := []mujica.CostItem{
	{Model: "opus", Tool: "search", InputTokens: 1000, OutputTokens: 200, Cost: 0.30},
	{Model: "haiku", Tool: "search", InputTokens: 4000, OutputTokens: 500, Cost: 0.05},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.CostBreakdown(c, &state, items, mujica.CostBreakdownOptions{GroupBy: mujica.CostByModel}).Element.Height(260)
}
ui.Render(view, 400, 300, 1)
```

签名：`func CostBreakdown(c *ui.Context, s *CostBreakdownState, items []CostItem, opts CostBreakdownOptions) DataTableView`

## 参数

`CostItem`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Model, Tool, Step | string | "" | 分组依据；为空归入“（未命名）”。 |
| InputTokens, OutputTokens | int | 0 | token 数；为负 panic。 |
| Cost | float64 | 0 | 费用；为负或非有限值 panic。 |

`CostBreakdownOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| GroupBy | CostGroup | CostByModel | CostByModel、CostByTool、CostByStep。 |
| Currency | string | "$" | 金额前缀。 |
| Decimals | int | 0（即 2） | 金额小数位；为负 panic。 |
| Version | int | 0 | 原地修改数据时递增；新切片或长度变化自动刷新。 |

`CostBreakdownState.Table` 保存排序与选择。汇总只在 items 或分组方式变化时计算。给返回的 `Element` 一个高度。

## 状态

- 费用列含占比条与百分比，数字为等宽数字。
- 合计行为衬线“合计”、token 总数与加粗金额，朗读完整合计。

## 事件

- 同 DataTable：`Changed()`、`SortChanged()`。

## 键盘操作

- 同 DataTable：↑/↓ 选择行，列标题可用键盘排序。

## 限制

- 金额以浮点汇总并按小数位显示，不做货币舍入规则。
