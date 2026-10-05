# 159 PriceAlertLine

## 用途

价格预警线；拖动调整；触发回调；不自行通知。

## 最小示例

```go
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	alerts := []mujica.PriceAlertLine{{ID: "take-profit", Price: 105}}
	triggered := ""
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CandlestickChart(c, candles, &window, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}, Alerts: &alerts, OnAlert: func(alert mujica.PriceAlertLine, last mujica.Candle) {
			triggered = alert.ID
		}})
		ui.Text(c, triggered)
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

ID 为唯一非空字符串；Price 为有限价格；Label 可选。CandlestickOptions.Alerts 接入主图；OnAlert 接收预警和最新 Candle。

## 状态

图表通过 Frame 显示默认、加载、空数据、错误和重试；祖先 Disabled 阻止输入；浅深主题沿用窗口设置。

## 事件

ChartResult.Point 提供读数索引；Changed 表示读数变化；Retried 表示重试。 OnAlert 在最新收盘价格跨线时触发一次；拖动直接更新 Alerts。

## 键盘操作

聚焦图表后方向键查看价格；PageUp/PageDown 选择预警线；Alt+上下方向键按价格精度步进；拖动预警线调整价格。

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。 百分比对比模式展示归一后的品种曲线，预警使用价格模式。
