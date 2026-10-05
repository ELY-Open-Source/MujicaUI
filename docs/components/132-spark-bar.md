# 132 SparkBar 迷你柱

## 用途

行内展示正负数值柱，零线区分方向，末项以强调颜色和菱形标记。

## 最小示例

```go
values := []float64{3, -2, 4, 1, -3}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.SparkBar(c, values, chart.SparkBarOptions{Frame: chart.ChartFrame{Label: "Alms"}, BarWidth: 10})
}
ui.Render(view, 400, 100, 1)
```

## 参数

`values []float64`：有限数值，NaN 表示缺值。`Frame.Label` 必填；`Frame.Height` 默认 28 DIP。`Frame.Status`、`Error`、`Empty` 控制占位。`BarWidth` 默认 6 DIP，负值 panic；整体宽度受窗口限制。`Format` 自定义读数。

## 状态

默认显示末值；悬停或聚焦显示当前点。加载、错误与空数据保持绘图区尺寸；祖先禁用阻止操作。

## 事件

返回 `ChartResult`，`Point` 为当前索引，`Changed()` 表示读数变化，`Retried()` 表示重试。

## 键盘操作

Tab 聚焦；左右箭头逐点读数并朗读；Home/End 到首尾；Escape 恢复末值。

## 限制

大数据按像素列保留首、极值和末项，绘制数量受宽度约束。数据范围与抽样缓存按切片身份更新；更换数据时替换切片。
