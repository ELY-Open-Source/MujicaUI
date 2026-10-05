# 152 ChartTypeSwitcher

## 用途

K 线、美国线、折线、面积、平均 K 线切换。

## 最小示例

```go
	value := chart.ChartCandles
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.ChartTypeSwitcher(c, &value, chart.ChartTypeSwitcherOptions{})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

ChartType 默认 ChartCandles；支持 ChartOHLC、ChartLine、ChartArea、ChartHeikinAshi；Disabled 禁用操作。

## 状态

浅深主题沿用窗口设置；交互组件遵循祖先 Disabled；调用方持有选择及展开状态。

## 事件

ChoiceResult.Changed() 表示选择变化。

## 键盘操作

Tab 进入分组，左右方向键移动焦点，Space 选择。

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
