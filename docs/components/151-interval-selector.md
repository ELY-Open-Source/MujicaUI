# 151 IntervalSelector

## 用途

1m 至 1M 周期；分组；收藏周期。

## 最小示例

```go
	value := chart.Interval("1h")
	favorites := []chart.Interval{"1m", "1h", "1D"}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.IntervalSelector(c, &value, &favorites, chart.IntervalSelectorOptions{})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

value 为 Intervals 中的周期；favorites 保存收藏列表；Disabled 禁用选择与收藏。

## 状态

浅深主题沿用窗口设置；交互组件遵循祖先 Disabled；调用方持有选择及展开状态。

## 事件

ChoiceResult.Changed() 表示选择变化。

## 键盘操作

Tab 进入收藏组，方向键移动焦点，Space 选择；Tab 遍历弹层周期和收藏按钮，Enter 激活，Escape 关闭。

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
