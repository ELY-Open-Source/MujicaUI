# 173 PositionTable

持仓表。

## 用途

持仓列表：数量、成本、现价、市值、浮动盈亏与比例，每行提供平仓操作。基于 DataTable，以稳定行 ID 保存选择与排序。

## 最小示例

```go
var state mujica.DataTableState[string]
rows := []mujica.PositionRow{{ID: "p1", Symbol: "MUJC", Quantity: 100, AvgCost: 400, Price: 412.08}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if id, ok := mujica.PositionTable(c, &state, rows, mujica.PositionTableOptions{}).Acted(); ok {
		fmt.Println("close", id)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `s` | `*DataTableState[string]` | 必填 | 表格状态。 |
| `rows` | `[]PositionRow` | 必填 | ID 必须稳定。 |
| `Decimals` | `int` | 2 | 小数位。 |
| `Version` | `int` | 0 | 原地修改行时递增。 |
| `Status / Error` | `DataStatus / string` | Ready | 加载与失败状态。 |
| `Disabled` | `bool` | false | 禁用平仓。 |

## 状态

默认、排序、选中、加载、失败、空、禁用。

## 事件

`TableActionResult.Acted()` 返回要平仓的行 ID；`Table` 为 DataTable 事件。

## 键盘操作

沿用 DataTable 的键盘操作；Tab 到达各行的平仓按钮。

## 限制

- 只构建可见行。
- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
