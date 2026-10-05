# 126 ChartBrush 区间刷选

## 用途

选择图表可见区间，支持移动窗口、调整两端和拖动新建范围。

## 最小示例

```go
value := [2]float64{20, 60}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ChartBrush(c, &value, mujica.ChartBrushOptions{Label: "Months", Max: 100, Step: 1, Values: []float64{2, 4, 3, 8, 5}})
}
ui.Render(view, 400, 160, 1)
```

## 参数

`value *[2]float64` 为域内严格递增区间；`Label` 必填；有限 `Min < Max`。`Step` 默认跨度的 1/100，负值 panic。`Values` 为预览折线。`Height` 默认 48 DIP，显式高度至少 28 DIP。`Format` 自定义值文本。`Disabled` 默认 false。

## 状态

显示选区边框与两端手柄；自身或祖先禁用时保持区间。拖动窗口到边缘保持跨度。

## 事件

返回 `ChoiceResult`；`Changed()` 表示区间变化；新范围会被朗读，调用方据此调整图表范围。

## 键盘操作

Tab 按起点、窗口、终点顺序聚焦；方向键移动一个步长；手柄 PageUp/PageDown 移动十个步长；Home/End 到域边界。

## 限制

范围至少保留一个 Step。预览切片按身份缓存，更新数据时替换切片。调用方负责图表缩放和查询。
