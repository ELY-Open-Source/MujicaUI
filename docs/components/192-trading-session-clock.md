# 192 TradingSessionClock

## 用途

交易时段时钟。

## 最小示例

```go
now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
markets := []finance.SessionMarket{{Name: "NYSE", Location: time.FixedZone("EDT", -4*3600), Sessions: []finance.TradingSession{{Phase: finance.MarketOpen, Start: 570, End: 960}}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.TradingSessionClock(c, markets, finance.TradingSessionClockOptions{Now: now})
}
ui.Render(view, 960, 440, 1)
```

## 参数

`markets` 每项包含 Name、Location 与有序时段；Start/End 是本地午夜后的分钟；`Holiday` 默认 false；`Now` 零值使用窗口时间。

## 状态

默认与数据边界。

## 事件

展示组件；数据由调用方更新。

## 键盘操作

展示组件。

## 限制

每天重复时段由调用方指定；节假日由 Holiday 指定；跨午夜时段拆分为两个日期的时段。
