# 161 TickerTape

行情滚动条。

## 用途

横向循环滚动的行情条，显示代码、价格与涨跌幅。指针悬停或键盘聚焦时暂停；减少动画时改为静态分页。

## 最小示例

```go
items := []finance.TickerItem{{Symbol: "MUJC", Price: 412.08, ChangePercent: -1.35}, {Symbol: "AVE", Price: 88.2, ChangePercent: 2.1}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.TickerTape(c, items, finance.TickerTapeOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `items` | `[]TickerItem` | 必填 | 代码、价格、涨跌幅（%）。 |
| `Decimals` | `int` | 2 | 价格小数位，0 表示 2，负数 panic。 |
| `Speed` | `float32` | 40 | 每秒滚动的 DIP。 |
| `PerPage` | `int` | 4 | 减少动画时每页条数。 |
| `Label` | `string` | 内置“Market ticker” | 无障碍名称。 |

## 状态

滚动、暂停（悬停或聚焦）、静态分页（减少动画）、空列表。

## 事件

无业务事件。

## 键盘操作

Tab 聚焦滚动条即暂停；减少动画时 ←/→ 翻页，也可用上一页、下一页按钮。

## 限制

- 滚动以四分之一 DIP 取整，暂停时停止请求重绘。
- 第二份副本用于无缝循环，屏幕阅读器可能读到重复内容。
- 涨跌颜色由 `UseTrend(c, GreenUp|RedUp)` 按窗口设置，默认绿涨红跌；方向同时用箭头与正负号表达，不只靠颜色。
