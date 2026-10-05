# 165 OrderBook

盘口。

## 用途

按价格显示卖档（上）与买档（下），带累计深度条与价差；可按价格步长合并档位；点击档位回传价格。

## 最小示例

```go
var state finance.OrderBookState
book := finance.OrderBookData{
	Bids: []finance.BookLevel{{Price: 100.00, Size: 5}, {Price: 99.99, Size: 3}},
	Asks: []finance.BookLevel{{Price: 100.01, Size: 1}, {Price: 100.02, Size: 2}},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if price, side, ok := finance.OrderBook(c, &state, book, finance.OrderBookOptions{Groups: []float64{0.01, 0.05}}).Clicked(); ok {
		fmt.Println("order at", price, side)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `s` | `*OrderBookState` | 必填 | Group 为合并步长，0 取 Groups 第一项；不在 Groups 中 panic。 |
| `book` | `OrderBookData` | 必填 | Bids、Asks 与 Version（原地修改时递增）。 |
| `Decimals / SizeDecimals` | `int` | 2 / 0 | 价格与数量小数位。 |
| `Depth` | `int` | 8 | 每侧档数。 |
| `Groups` | `[]float64` | nil | 可选合并步长，须为价格精度的整数倍，否则 panic。 |
| `Disabled` | `bool` | false | 禁用点击与合并选择。 |

## 状态

默认、悬停、聚焦、禁用；空盘口时只显示表头。

## 事件

`OrderBookResult.Clicked()` 返回被点击档位的价格与方向。

## 键盘操作

Tab 进入档位列表，↑/↓ 在档位间移动，Enter/Space 选中价格。

## 限制

- 合并以整数价格单位计算：买档向下取整、卖档向上取整。只在数据或步长变化时重新聚合。
- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
