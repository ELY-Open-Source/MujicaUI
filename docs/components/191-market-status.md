# 191 MarketStatus

## 用途

市场状态。

## 最小示例

```go
now := time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC)
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	finance.MarketStatus(c, finance.MarketOpen, finance.MarketStatusOptions{Market: "NYSE", Now: now, Next: now.Add(5 * time.Hour), NextPhase: finance.MarketAfterHours})
}
ui.Render(view, 960, 440, 1)
```

## 参数

`phase` 支持五个 MarketPhase；`Market` 默认空；`Next` 零值隐藏；`NextPhase` 默认 Closed；`Now` 零值使用窗口时间。

## 状态

默认与数据边界。

## 事件

展示组件；数据由调用方更新。

## 键盘操作

展示组件。

## 限制

使用窗口时间时每分钟刷新；传入 Now 时由调用方推进时间。
