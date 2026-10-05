# 138 DensityPlot 密度图

## 用途

将多组样本的高斯核密度曲线叠加在同一坐标域，以形状和颜色区分序列。

## 最小示例

```go
groups := []mujica.BoxGroup{{Name: "Matins", Values: []float64{1, 2, 2, 3, 7}}, {Name: "Vespers", Values: []float64{3, 4, 5, 6, 6}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.DensityPlot(c, groups, mujica.DensityPlotOptions{Frame: mujica.ChartFrame{Label: "Ages"}, Bandwidth: 1, Points: 128})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`groups []BoxGroup` 为唯一非空组名和样本；NaN 为缺值。有限 `Bandwidth` 默认 0，按组自动估计；负值 panic。`Points` 默认 128，显式值须至少 2。`Frame` 和 `Format` 控制图表与密度读数。

## 状态

默认、序列显隐、悬停、键盘读数、加载、错误、空数据、祖先禁用。

## 事件

`Changed()` 表示采样位置变化；`Toggled()` 表示序列显隐；`Retried()` 表示重试。

## 键盘操作

Tab 先到图例，再到图表；图例 Space 切换显隐；图表方向键读取共享 X 位置的各组密度；Home/End 首尾；Escape 清除。

## 限制

核密度计算在数据、带宽或 Points 变化时执行；绘制曲线按像素列抽样。100,000 个样本默认生成每组 128 个密度点。更新数据需替换样本切片。
