# 120 RealtimeChart 实时图表

## 用途

实时折线图：调用方向 `RealtimeBuffer` 追加数据后刷新，图表显示截至最新数据的固定时间窗口；可暂停视图；组件自身不设计时器，隐藏或空闲时不重绘。

## 最小示例

```go
buf := chart.NewRealtimeBuffer(600, "Bell")
start := time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)
for i := range 60 {
	buf.Append(start.Add(time.Duration(i)*time.Second), float64(i%7))
}
var state chart.RealtimeChartState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.RealtimeChart(c, &state, buf, chart.RealtimeChartOptions{Frame: chart.ChartFrame{Label: "Sound"}, Window: 30 * time.Second})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func RealtimeChart(c *ui.Context, s *RealtimeChartState, buf *RealtimeBuffer, opts RealtimeChartOptions) ChartResult`；`NewRealtimeBuffer(capacity, names...)` 保存最近 `capacity` 个样本，`Append(t, values...)` 每个序列一个值且时间不倒退。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Window` | `time.Duration` | 必填 | 显示的时间窗口，止于最新样本（暂停时止于暂停那一刻）。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 读数格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 暂停时按钮显示“继续”，新数据继续进入缓冲区但不改变视图。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试；`ValueChanged()` 表示暂停或继续；`Toggled()` 报告图例切换。

## 键盘操作

- Tab：暂停按钮、图例、图表。
- 图表读数键同 LineChart；暂停按钮：Enter/Space。

## 限制

- 数据追加在后台时，请通过 `Window.Update` 回到界面线程追加并刷新。
- 缓冲区按块丢弃旧样本，最多占用两倍容量。
