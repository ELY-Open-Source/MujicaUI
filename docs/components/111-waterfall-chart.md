# 111 WaterfallChart 瀑布图

## 用途

瀑布图：增减项累计，小计与总计柱显示累计值；增加为 Success 色并配 ▲，减少为 Danger 色并配 ▼，合计为中性色配 ■，柱端写出带符号的数值。

## 最小示例

```go
steps := []mujica.WaterfallStep{
	{Label: "Tithes", Value: 12},
	{Label: "Repairs", Value: -5},
	{Label: "Year", Kind: mujica.WaterfallTotal},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.WaterfallChart(c, steps, mujica.WaterfallChartOptions{Frame: mujica.ChartFrame{Label: "Ledger"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func WaterfallChart(c *ui.Context, steps []WaterfallStep, opts WaterfallChartOptions) ChartResult`；`WaterfallStep{Label, Value, Kind}`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `steps[i].Kind` | `WaterfallKind` | `WaterfallDelta` | `WaterfallDelta` 累加；`WaterfallSubtotal`、`WaterfallTotal` 显示当前累计，`Value` 须为 0。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 数值格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：逐柱移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。

## 限制

- 小计与总计都显示到该步为止的累计值，不清零。
