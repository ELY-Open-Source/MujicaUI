# 149 ChartSync

## 用途

多图十字光标与时间范围同步；可逐图解除。

## 最小示例

```go
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	var sync mujica.ChartSync
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		sync.Link(c, "price")
		mujica.CandlestickChart(c, candles, &window, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}, Sync: &sync, SyncID: "price"})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

Count、Start、End 保存共享时间窗口；不同周期按时间边界映射；Cursor 保存时间光标；Unlinked 按稳定图表 ID 解联。通过 CandlestickOptions 或 VolumeChartOptions 的 Sync 接入。

## 状态

浅深主题沿用窗口设置；交互组件遵循祖先 Disabled；调用方持有选择及展开状态。

## 事件

ChoiceResult.Changed() 表示选择变化；ChartSync.Link 返回 ChoiceResult。

## 键盘操作

Tab 聚焦 Link 开关，Space 切换关联；图表方向键移动同步光标。

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
