# 171 OrderEntry

下单面板。

## 用途

限价、市价、止损与条件限价单的录入：方向、类型、数量、价格、触发价与有效期；同步与调用方异步校验；预估金额；提交防重入。建立在 Form 之上。

## 最小示例

```go
var state finance.OrderEntryState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if order, ok := finance.OrderEntry(c, &state, finance.OrderEntryOptions{Symbol: "MUJC", Last: 100.2}).Submitted(); ok {
		fmt.Println("place", order.Quantity)
		state.EndSubmit("")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `s` | `*OrderEntryState` | 必填 | Draft 为当前订单，Form 为表单状态。 |
| `Symbol` | `string` | 必填 | 代码，空值 panic。 |
| `Last` | `float64` | 0 | 最新价，用于市价预估和止损方向检查。 |
| `TickSize / LotSize` | `float64` | 0.01 / 1 | 价格与数量步长。 |
| `MaxQuantity` | `float64` | 0 | 数量上限，0 不限。 |
| `Initial` | `OrderDraft` | 买入限价 | 重置恢复的订单。 |
| `Validate` | `func(OrderDraft) string` | nil | 调用方同步规则。 |
| `Check` | `func(OrderDraft, uint64)` | nil | 异步检查（如可用资金），以 `s.Resolve(version, msg)` 回传；旧版本结果被丢弃。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认、字段错误（失焦或提交后）、异步检查中、提交中（按钮加载）、表单错误、禁用。

## 事件

`OrderEntryResult.Submitted()` 返回通过全部校验的订单；处理完毕后调用 `s.EndSubmit(err)`。

## 键盘操作

Tab 依次经过方向、类型、各字段与按钮；数量与价格字段 ↑/↓ 按步长调整；提交失败聚焦第一个错误字段。

## 限制

- 止损单检查：买入触发价须高于最新价，卖出须低于。
- 数量字段的异步检查覆盖整张订单：任一字段变化都会让旧结果失效。
