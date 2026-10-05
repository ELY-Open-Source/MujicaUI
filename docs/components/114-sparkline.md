# 114 Sparkline 迷你图

## 用途

行内迷你图：折线、柱形或面积，无坐标轴；末值以菱形标记并在右侧写出数值；聚焦后可逐点读数。

## 最小示例

```go
values := []float64{3, 5, 4, 7, 6, 9}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.Sparkline(c, values, chart.SparklineOptions{Frame: chart.ChartFrame{Label: "Bells"}, Kind: chart.SparklineArea})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Sparkline(c *ui.Context, values []float64, opts SparklineOptions) ChartResult`；NaN 为缺值。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `28` | 迷你图高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Kind` | `SparklineKind` | `SparklineLine` | `SparklineLine`、`SparklineBar` 或 `SparklineArea`。 |
| `Width` | `float32` | `96` | 迷你图宽度（DIP）。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 右侧数值格式。 |

## 状态

默认（显示末值）、悬停或键盘读数（右侧数值变为当前点，点外加圆环）、加载（小型加载标记）、失败（错误文字与重试图标）、无数据。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试。

## 键盘操作

- Tab：聚焦迷你图。
- ←/→：逐点移动；Home/End：首尾；Escape：恢复显示末值。

## 限制

- 不显示图例与提示行，以便放在表格行或文字中。
