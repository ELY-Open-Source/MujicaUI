# 148 MultiChartLayout

## 用途

1、2、4、6 宫格多图；独立品种与周期；布局切换不丢状态。

## 最小示例

```go
	state := chart.MultiChartState{Count: 2}
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.MultiChartLayout(c, &state, chart.MultiChartLayoutOptions{Height: 200}, func(i int, cell *chart.ChartCell) {
			input.TextInput(c, &cell.Symbol, input.TextInputOptions{Label: "Symbol"})
		})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

Count 支持 1、2、4、6，零值为 1；Cells 保存六个独立品种、周期与窗口；Active 指向可见单元。Height 默认 480。

## 状态

浅深主题沿用窗口设置；交互组件遵循祖先 Disabled；调用方持有选择及展开状态。

## 事件

ChoiceResult.Changed() 表示选择变化。

## 键盘操作

Tab 进入布局组，方向键移动焦点，Space 切换；Tab 进入各图表独立操作。

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
