# 112 FunnelChart 漏斗图

## 用途

漏斗图：阶段宽度与值成比例，右侧写出值与相对上一阶段的转化率；读数同时给出相对首阶段的比例。

## 最小示例

```go
stages := []chart.ChartSlice{{Label: "Visitors", Value: 1000}, {Label: "Novices", Value: 400}, {Label: "Vows", Value: 100}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.FunnelChart(c, stages, chart.FunnelChartOptions{Frame: chart.ChartFrame{Label: "Order"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func FunnelChart(c *ui.Context, stages []ChartSlice, opts FunnelChartOptions) ChartResult`；值须为有限且 ≥ 0。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 数值格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 全部为 0 时显示无数据。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：逐阶段移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。

## 限制

- 上一阶段为 0 时转化率显示“—”。
