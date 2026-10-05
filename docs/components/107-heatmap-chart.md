# 107 HeatmapChart 热力图

## 用途

二维矩阵热力图：按值分五级着色，右侧色阶图例标出最小与最大值；单元格读数；空值单元格为虚线框。

## 最小示例

```go
rows, cols := []string{"Mon", "Tue"}, []string{"Matins", "Lauds", "Vespers"}
values := [][]float64{{1, 4, 9}, {3, math.NaN(), 7}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.HeatmapChart(c, rows, cols, values, mujica.HeatmapChartOptions{Frame: mujica.ChartFrame{Label: "Hours"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func HeatmapChart(c *ui.Context, rows, cols []string, values [][]float64, opts HeatmapChartOptions) ChartResult`；`values[行][列]`，NaN 为空值。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 单元格与读数格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 单元格足够大时直接写出数值，不只依赖颜色。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试。

## 键盘操作

- Tab：聚焦热力图。
- ←/→：同一行左右移动；↑/↓：同一列上下移动；第一次按键进入左上角。
- Home/End：第一个/最后一个单元格；Escape：清除读数。

## 限制

- 颜色分五级，不做连续渐变。
