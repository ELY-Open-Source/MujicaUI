# 155 IndicatorSelector

## 用途

指标搜索、添加、移除、参数编辑。

## 最小示例

```go
	var overlays []chart.IndicatorOverlay
	var panes chart.IndicatorPanes
	view := func(c *ui.Context) {
		core.Use(c, core.Settings{})
		chart.IndicatorSelector(c, &overlays, &panes, chart.IndicatorSelectorOptions{})
	}
	ui.Render(view, 640, 400, 1)
```

## 参数

overlays 与 panes 均由调用方持有；添加默认参数，可搜索、移除和编辑。Disabled 关闭面板并阻止编辑。

## 状态

浅深主题沿用窗口设置；交互组件遵循祖先 Disabled；调用方持有选择及展开状态。

## 事件

ChoiceResult.Changed() 表示选择变化。

## 键盘操作

Tab 聚焦按钮，Enter 打开；Tab 到搜索和参数字段；方向键步进参数；Escape 关闭。

## 限制

数据由调用方提供；更新数据时替换切片，缓存按切片身份失效。行情读取、交易及持久化由调用方负责。
