# 174 OrderTable

委托表。

## 用途

委托列表：方向、类型、价格、成交进度与状态，按状态筛选；进行中的委托提供撤单。基于 DataTable。

## 最小示例

```go
var state mujica.OrderTableState
rows := []mujica.OrderRow{{ID: "o1", Symbol: "MUJC", Side: mujica.SideBuy, Price: 99.5, Quantity: 100, Filled: 40, Status: mujica.OrderPartial}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if id, ok := mujica.OrderTable(c, &state, rows, mujica.OrderTableOptions{}).Acted(); ok {
		fmt.Println("cancel", id)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `s` | `*OrderTableState` | 必填 | Filter 为状态筛选，Table 为表格状态。 |
| `rows` | `[]OrderRow` | 必填 | Filled 不得超过 Quantity，否则 panic。 |
| `Decimals / QuantityDecimals` | `int` | 2 / 0 | 小数位。 |
| `Version` | `int` | 0 | 原地修改行时递增。 |
| `Status / Error` | `DataStatus / string` | Ready | 加载与失败状态。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

全部、进行中、已成交、已关闭四种筛选；部分成交进度条；禁用。

## 事件

`TableActionResult.Acted()` 返回要撤销的委托 ID。

## 键盘操作

筛选分段用方向键切换；表格沿用 DataTable 键盘操作；Tab 到达撤单按钮。

## 限制

- 筛选只在数据或条件变化时执行。
