# 150 TimeRangeSelector

## 用途

1D、5D、1M、3M、YTD、1Y、5Y、自定义；键盘导航。

## 最小示例

```go
	value := mujica.MarketRange{Preset: mujica.Range1M}
	view := func(c *ui.Context) {
		mujica.Use(c, mujica.Settings{})
		mujica.TimeRangeSelector(c, &value, mujica.TimeRangeSelectorOptions{})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

MarketRange.Preset 默认 Range1D；Custom 保存 DateRange。Today 零值跟随系统日期，Bounds 使用调用方时间和时区。

## 状态

浅深主题沿用窗口设置；交互组件遵循祖先 Disabled；调用方持有选择及展开状态。

## 事件

ChoiceResult.Changed() 表示选择变化。

## 键盘操作

Tab 进入预设组，左右方向键移动焦点，Space 选择；自定义日期通过 DateRangePicker 键盘操作。

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
