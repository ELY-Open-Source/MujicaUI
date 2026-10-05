# 129 ViolinPlot 小提琴图

## 用途

比较多组样本的核密度轮廓，中位数以菱形标记，四分位以粗条表示。

## 最小示例

```go
groups := []chart.BoxGroup{{Name: "Choir", Values: []float64{1, 2, 3, 4, 7, 9}}, {Name: "Nave", Values: []float64{2, 4, 5, 6, 8}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.ViolinPlot(c, groups, chart.ViolinPlotOptions{Frame: chart.ChartFrame{Label: "Hymns"}, Bandwidth: 1})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`groups []BoxGroup`：唯一非空 Name 与样本 Values；NaN 为缺值，无穷值 panic。`Frame` 指定名称、高度、状态和占位文案。有限 `Bandwidth` 默认 0，按组自动估计；负值 panic。`Format` 自定义数值读数。

## 状态

默认、悬停、键盘读数、加载、错误、空数据、祖先禁用。单值或相同值样本保留密度轮廓。

## 事件

`ChartResult.Point` 为组索引；`Changed()` 表示读数组变化；`Retried()` 表示重试。

## 键盘操作

Tab 聚焦；方向键逐组读取中位数和四分位；Home/End 首尾；Escape 清除读数。

## 限制

每组轮廓使用 64 个采样点；大量组按像素位置抽样绘制，当前读数组保留。排序和核密度在样本切片身份或带宽变化时计算；更新数据需替换切片。
