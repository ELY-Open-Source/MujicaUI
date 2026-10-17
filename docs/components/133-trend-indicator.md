# 133 TrendIndicator 趋势指示

## 用途

以方向箭头、带符号数值和无障碍方向文字展示指标变化。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.TrendIndicator(c, 12.5, chart.TrendIndicatorOptions{Precision: 1, Unit: "%", Note: "vs Lent"})
}
ui.Render(view, 400, 100, 1)
```

## 参数

`change float64` 必须有限；`Precision` 默认 0，负值 panic；`Unit` 和 `Note` 默认空；`Invert` 默认 false，开启后下降采用成功语义；`Format` 替换数值格式。

## 状态

上升、下降、持平；正负号与箭头提供颜色之外的区分；方向词支持中英文和消息覆盖。

## 事件

返回 `ui.Element`，展示调用方提供的变化值。

## 键盘操作

展示部件；辅助技术读取方向和完整变化值。

## 限制

变化值由调用方计算；组件提供显示和语义颜色。
