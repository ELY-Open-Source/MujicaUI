# 130 ParallelCoordinates 平行坐标

## 用途

在多条平行轴间展示每行数值；轴上刷选筛除行，当前读数行用粗线高亮。

## 最小示例

```go
filters := map[string][2]float64{}
axes := []string{"Age", "Hymns", "Vigils"}
rows := [][]float64{{20, 3, 9}, {35, 8, 2}, {50, 5, 5}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ParallelCoordinates(c, &filters, axes, rows, mujica.ParallelCoordinatesOptions{Frame: mujica.ChartFrame{Label: "Monks"}, Names: []string{"Agnes", "Clare", "Hild"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`filters *map[string][2]float64` 必填，键为轴名，值为有限递增范围。`axes []string` 至少两个唯一非空名称；每行长度与轴数一致。NaN 为缺值。`Names` 为逐行名称，nil 显示序号。`Frame` 与 `Format` 同其他图表。

## 状态

默认、高亮、筛除、加载、错误、空数据、祖先禁用；筛除行的读数附带本地化文字。

## 事件

`Changed()` 表示读数行变化；`ValueChanged()` 表示刷选范围变化；`Retried()` 表示重试。

## 键盘操作

Tab 聚焦；方向键逐行读数；Home/End 首尾。Alt+PageUp/PageDown 选择轴；Alt+上下调整下端；Alt+Shift+上下调整上端；Alt+Delete 清除该轴刷选。

## 限制

指针拖动轴设置范围，轴上点击清除范围。100,000 行绘制数量受图表高度约束，键盘可读取全部行。范围与过滤仅在数据或筛选变化时计算；更新数据需替换外层 rows 切片。
