# 146 DepthChart

## 用途

买卖深度累计面积；中间价；悬停读数。

## 最小示例

```go
	bids := []mujica.DepthLevel{{Price: 99, Size: 10}, {Price: 98, Size: 20}}
	asks := []mujica.DepthLevel{{Price: 101, Size: 15}, {Price: 102, Size: 25}}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.DepthChart(c, bids, asks, mujica.DepthChartOptions{Frame: mujica.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

买盘价格递减，卖盘递增，Size 为正数；最优买价小于最优卖价。Decimals 默认两位。

## 状态

图表通过 Frame 显示默认、加载、空数据、错误和重试；祖先 Disabled 阻止输入；浅深主题沿用窗口设置。

## 事件

ChartResult.Point 提供读数索引；Changed 表示读数变化；Retried 表示重试。

## 键盘操作

Tab 聚焦图表；方向键移动数据读数并朗读；Home/End 到首尾；Escape 清除读数。行情主图通过 +/- 缩放、Shift+左右方向键平移。 

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
