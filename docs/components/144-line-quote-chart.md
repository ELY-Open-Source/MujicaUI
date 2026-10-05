# 144 LineQuoteChart

## 用途

分时图；昨收基线；均价线；盘前盘后分段。

## 最小示例

```go
	at := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
	points := []mujica.QuotePoint{{Time: at, Price: 101, Volume: 100}, {Time: at.Add(time.Hour), Price: 103, Volume: 200}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.LineQuoteChart(c, points, mujica.LineQuoteChartOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}, PrevClose: 100, RegularOpen: at.Add(30 * time.Minute), RegularClose: at.Add(45 * time.Minute)})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

`PrevClose` 为必填正数；RegularOpen、RegularClose 成对提供，零值展示同一时段。均价使用 Price × Volume 累计。

## 状态

图表通过 Frame 显示默认、加载、空数据、错误和重试；祖先 Disabled 阻止输入；浅深主题沿用窗口设置。

## 事件

ChartResult.Point 提供读数索引；Changed 表示读数变化；Retried 表示重试。

## 键盘操作

Tab 聚焦图表；方向键移动数据读数并朗读；Home/End 到首尾；Escape 清除读数。行情主图通过 +/- 缩放、Shift+左右方向键平移。 

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
