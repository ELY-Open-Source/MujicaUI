# 168 DepthLadder

价格阶梯。

## 用途

按最小变动价位排列的价格阶梯（DOM）：买量、价格、卖量三列，最新价居中并加框；点击买量或卖量单元格回调下单。

## 最小示例

```go
var state mujica.DepthLadderState
book := mujica.OrderBookData{Bids: []mujica.BookLevel{{Price: 100, Size: 5}}, Asks: []mujica.BookLevel{{Price: 100.01, Size: 2}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if price, side, ok := mujica.DepthLadder(c, &state, 100, book, mujica.DepthLadderOptions{TickSize: 0.01}).Clicked(); ok {
		fmt.Println("place", side, "at", price)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `s` | `*DepthLadderState` | 必填 | Offset 为相对最新价滚动的档数。 |
| `last` | `float64` | 必填 | 最新价。 |
| `book` | `OrderBookData` | 必填 | 盘口数据。 |
| `TickSize` | `float64` | 必填 | 最小变动价位，须为价格精度的整数倍，否则 panic。 |
| `Decimals / SizeDecimals` | `int` | 2 / 0 | 小数位。 |
| `Rows` | `int` | 15 | 显示的行数。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认、悬停、聚焦、最新价行（底色、框线与菱形）、已滚动（显示回到最新价按钮）、禁用。

## 事件

`DepthLadderResult.Clicked()` 返回价格与方向（买量列为买入，卖量列为卖出）。

## 键盘操作

Tab 进入阶梯，方向键在单元格间移动，Enter/Space 下单；PageUp/PageDown 滚动半屏；滚轮逐档滚动。

## 限制

- 只构建可见行；按价位建立的数量索引只在数据变化时重建。
- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
