# 167 TimeAndSales

逐笔成交。

## 用途

逐笔成交列表，最新在上，只构建可见行；大单以菱形与粗体标记；可暂停跟随。

## 最小示例

```go
var state finance.TimeAndSalesState
trades := []finance.TradePrint{{ID: "1", Time: time.Now(), Price: 100.01, Size: 300, Side: finance.SideBuy}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.TimeAndSales(c, &state, trades, finance.TimeAndSalesOptions{LargeSize: 1000})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `s` | `*TimeAndSalesState` | 必填 | Paused 暂停跟随。 |
| `trades` | `[]TradePrint` | 必填 | 按时间先后追加；ID 唯一。 |
| `Decimals / SizeDecimals` | `int` | 2 / 0 | 价格与数量小数位。 |
| `LargeSize` | `float64` | 0 | 数量达到即标记为大单。 |
| `Height` | `float32` | 320 | 高度（DIP）。 |
| `Status / Error` | `DataStatus / string` | Ready | 加载、失败等状态。 |

## 状态

跟随、暂停（显示新成交数量）、加载、失败、空列表。

## 事件

`ChoiceResult.Changed()` 在暂停或继续时为 true。

## 键盘操作

Tab 到达暂停按钮（Enter/Space 切换）与列表；列表内 ↑/↓ 移动。

## 限制

- 基于 List 虚拟化，100,000 条成交只构建可见行。
- 暂停时冻结显示的条数，新成交不打断阅读。
- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
