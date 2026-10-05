# 147 VolumeProfile

## 用途

价格分布成交量；价值区与控制点。

## 最小示例

```go
	bins := []chart.ProfileBin{{Low: 99, High: 100, Volume: 30}, {Low: 100, High: 101, Volume: 80}, {Low: 101, High: 102, Volume: 40}}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.VolumeProfile(c, bins, chart.VolumeProfileOptions{Frame: chart.ChartFrame{Label: "BTC", Height: 200}})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

ProfileBin 按价格递增排列；Volume 非负。ValueArea 默认 0.7；VolumeProfileOf 可按 OHLC 区间生成分布。

## 状态

图表通过 Frame 显示默认、加载、空数据、错误和重试；祖先 Disabled 阻止输入；浅深主题沿用窗口设置。

## 事件

ChartResult.Point 提供读数索引；Changed 表示读数变化；Retried 表示重试。

## 键盘操作

Tab 聚焦图表；方向键移动数据读数并朗读；Home/End 到首尾；Escape 清除读数。行情主图通过 +/- 缩放、Shift+左右方向键平移。 

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
