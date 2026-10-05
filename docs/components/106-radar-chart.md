# 106 RadarChart 雷达图

## 用途

雷达图：三个及以上维度，多序列多边形，刻度环，按维度悬停或键盘读数。

## 最小示例

```go
axes := []string{"Voice", "Latin", "Chant", "Organ"}
series := []mujica.ChartSeries{{Name: "Agnès", Values: []float64{8, 6, 9, 4}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.RadarChart(c, axes, series, mujica.RadarChartOptions{Frame: mujica.ChartFrame{Label: "Skills"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func RadarChart(c *ui.Context, axes []string, series []ChartSeries, opts RadarChartOptions) ChartResult`；每条序列每个维度一个值，值须 ≥ 0。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Max` | `float64` | 自动取整 | 外环代表的值；0 按数据向上取整。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 读数格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 当前维度的轴线加粗、标签变为正文色。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试；多序列时 `Toggled()` 报告图例切换。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：按维度移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。

## 限制

- 缺值按 0 画在中心，读数显示“无数值”。
