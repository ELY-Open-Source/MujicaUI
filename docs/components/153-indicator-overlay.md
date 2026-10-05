# 153 IndicatorOverlay

## 用途

主图叠加 MA、EMA、BOLL、VWAP；参数可配；计算由调用方或内置辅助函数。

## 最小示例

```go
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []chart.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window chart.MarketView
	overlays := []chart.IndicatorOverlay{{Kind: chart.IndicatorMA, Params: []float64{2}}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.CandlestickChart(c, candles, &window, chart.CandlestickOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}, Overlays: overlays})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

Kind 支持 MA、EMA、BOLL、VWAP；Params 为 nil 使用默认参数；Values 可覆盖内置计算，长度等于 candles，NaN 表示缺测。

## 状态

图表通过 Frame 显示默认、加载、空数据、错误和重试；祖先 Disabled 阻止输入；浅深主题沿用窗口设置。

## 事件

ChartResult.Point 提供读数索引；Changed 表示读数变化；Retried 表示重试。

## 键盘操作

Tab 聚焦图表；方向键移动数据读数并朗读；Home/End 到首尾；Escape 清除读数。行情主图通过 +/- 缩放、Shift+左右方向键平移。 

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
