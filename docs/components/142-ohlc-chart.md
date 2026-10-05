# 142 OHLCChart

## 用途

美国线；与 K 线共用数据接口。

## 最小示例

```go
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.OHLCChart(c, candles, &window, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

与 CandlestickChart 共用 Candle、MarketView 和 CandlestickOptions；函数将 Type 设置为 ChartOHLC。

## 状态

图表通过 Frame 显示默认、加载、空数据、错误和重试；祖先 Disabled 阻止输入；浅深主题沿用窗口设置。

## 事件

ChartResult.Point 提供读数索引；Changed 表示读数变化；Retried 表示重试。

## 键盘操作

Tab 聚焦图表；方向键移动数据读数并朗读；Home/End 到首尾；Escape 清除读数。行情主图通过 +/- 缩放、Shift+左右方向键平移。 

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
