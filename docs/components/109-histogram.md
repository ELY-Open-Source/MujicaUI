# 109 Histogram 直方图

## 用途

分箱计数直方图：箱数或箱宽可配置，可切换为累计模式；分箱只在数据或参数变化时计算一次。

## 最小示例

```go
values := []float64{1, 2, 2, 3, 4, 5, 9, 10}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.Histogram(c, values, chart.HistogramOptions{Frame: chart.ChartFrame{Label: "Ages"}, BinWidth: 5})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Histogram(c *ui.Context, values []float64, opts HistogramOptions) ChartResult`；NaN 与无穷值被忽略。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Bins` | `int` | Sturges 规则 | 箱数；负数 panic。 |
| `BinWidth` | `float64` | 0 | 箱宽；设定时优先于 `Bins`，箱边取箱宽的整数倍。超过 10000 箱 panic。 |
| `Cumulative` | `bool` | false | 每箱计入不超过其上界的全部值。 |
| `Format` | `func(float64) string` | 按箱宽取精度 | 箱边格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：逐箱移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。

## 限制

- 按箱数分箱时最后一箱包含上界，读数用 `]` 表示。
