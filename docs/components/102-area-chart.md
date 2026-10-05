# 102 AreaChart 面积图

## 用途

面积图：单序列填充到基线，或多序列堆叠；基线可配置；悬停与键盘读数显示各序列原值（不是堆叠后的值）。

## 最小示例

```go
series := []chart.ChartSeries{
	{Name: "Matins", Values: []float64{3, 4, 6, 5}},
	{Name: "Vespers", Values: []float64{2, 3, 2, 4}},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.AreaChart(c, series, chart.AreaChartOptions{Frame: chart.ChartFrame{Label: "Attendance"}, Stacked: true})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func AreaChart(c *ui.Context, series []ChartSeries, opts AreaChartOptions) ChartResult`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `X`, `Times` | | nil | 同 LineChart。 |
| `Stacked` | `bool` | false | 可见序列依次堆叠；隐藏序列不参与堆叠。 |
| `Baseline` | `float64` | `0` | 面积的起始值；坐标轴总包含基线。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 读数格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试；`Toggled()` 同 LineChart。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：在数据点之间移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。

## 限制

- 缺值处面积断开；堆叠时缺值按 0 计入后续序列。
- 填充为半透明的序列色，边线保持实色以满足 3:1 对比。
