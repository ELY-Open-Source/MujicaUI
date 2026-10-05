# 156 DrawingTools

## 用途

趋势线、水平线、斐波那契、矩形、文字标注；选中、移动、删除、撤销。

## 最小示例

```go
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	candles := []mujica.Candle{
		{Time: at, Open: 100, High: 103, Low: 99, Close: 102, Volume: 1200},
		{Time: at.Add(time.Hour), Open: 102, High: 104, Low: 100, Close: 101, Volume: 1600},
		{Time: at.Add(2 * time.Hour), Open: 101, High: 105, Low: 100, Close: 104, Volume: 900},
	}
	var window mujica.MarketView
	drawings := mujica.DrawingTools{Tool: mujica.ToolTrendLine}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.CandlestickChart(c, candles, &window, mujica.CandlestickOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}, Drawings: &drawings})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

Items 保存时间和价格锚点；Tool 默认选择工具；Text 为文字标注内容；Locked 阻止编辑；Selected、Delete、Undo、Clear 管理画线。

## 状态

图表通过 Frame 显示默认、加载、空数据、错误和重试；祖先 Disabled 阻止输入；浅深主题沿用窗口设置。

## 事件

ChartResult.Point 提供读数索引；Changed 表示读数变化；Retried 表示重试。

## 键盘操作

聚焦图表后方向键移动读数；Enter 在当前读数放置锚点，双点工具按两次 Enter；Alt+Enter 循环选择画线；点击放置锚点；Alt+方向键移动选中画线；Delete/Backspace 删除；Command/Ctrl+Z 撤销；Escape 取消。

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。 图形锚点使用已有 Candle 时间。
