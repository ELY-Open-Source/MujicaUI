# 104 PieChart 饼图

## 用途

饼图与环形图：扇区按占比绘制并标注百分比；过窄的扇区可合并为“其他”；环形图中心显示总计。

## 最小示例

```go
slices := []chart.ChartSlice{{Label: "Wax", Value: 50}, {Label: "Oil", Value: 30}, {Label: "Salt", Value: 4}, {Label: "Ash", Value: 3}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.PieChart(c, slices, chart.PieChartOptions{Frame: chart.ChartFrame{Label: "Stores"}, Donut: true, MinAngle: 20})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func PieChart(c *ui.Context, slices []ChartSlice, opts PieChartOptions) ChartResult`；`ChartSlice{Label, Value}`，值须为有限且 ≥ 0。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Donut` | `bool` | false | 环形图，中心写总计。 |
| `MinAngle` | `float64` | `0` | 小于该角度（度）的扇区达到两个及以上时合并为“其他”；须在 [0, 360)。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 数值格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 当前扇区向外凸出。总计为 0 时显示无数据。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：在数据点之间移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。

## 限制

- 只有一个窄扇区时不合并，以免“其他”只代表一项。
- 扇区内的形状标记与百分比文字按对比度自动选用正文色或表面色。
