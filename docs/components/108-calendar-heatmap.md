# 108 CalendarHeatmap 日历热力图

## 用途

按日的贡献图：一列一周，按当年最大值分四级着色；年份切换；周起始可配置；逐日读数。

## 最小示例

```go
year := 2026
values := map[core.Date]float64{{Year: 2026, Month: time.March, Day: 4}: 5}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.CalendarHeatmap(c, &year, values, chart.CalendarHeatmapOptions{Frame: chart.ChartFrame{Label: "Prayers"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func CalendarHeatmap(c *ui.Context, year *int, values map[Date]float64, opts CalendarHeatmapOptions) ChartResult`；`*year` 须在 1–9999。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `WeekStart` | `WeekStart` | 跟随语言 | 每列第一行是星期几。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 读数格式。 |

## 状态

默认、无记录日（浅色格）、四级色阶、悬停或键盘读数（描边）、加载、失败。下方图例从“少”到“多”。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试；`ValueChanged()` 表示上一年/下一年按钮改变了 `*year`。

## 键盘操作

- Tab：年份按钮，再到日历。
- ←/→：前后一周；↑/↓：前后一天；Home/End：全年第一天与最后一天。
- 上一年/下一年按钮：Enter/Space。

## 限制

- 方格较小，只作读数目标，不是独立操作，因此不要求 28 DIP 命中区。
- 值 ≤ 0 视为无记录。
