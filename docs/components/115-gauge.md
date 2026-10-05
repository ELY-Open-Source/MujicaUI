# 115 Gauge 仪表盘

## 用途

仪表盘：半圆或整圆表盘，按阈值分段着色，指针指向当前值，下方写出数值、单位和所在分段的名称，不只靠颜色表达区间。

## 最小示例

```go
segments := []mujica.GaugeSegment{{To: 60, Label: "Calm", Tone: mujica.GaugeGood}, {To: 100, Label: "Crowded", Tone: mujica.GaugeDanger}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Gauge(c, 72, mujica.GaugeOptions{Frame: mujica.ChartFrame{Label: "Nave"}, Max: 100, Unit: "%", Segments: segments})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Gauge(c *ui.Context, value float64, opts GaugeOptions) ChartResult`；值须在 `Min`–`Max` 之间。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `200` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Min`, `Max` | `float64` | 必须 `Min < Max` | 刻度范围。 |
| `Segments` | `[]GaugeSegment` | nil | 分段 `{To, Label, Tone}`，上界递增且最后一段止于 `Max`；`Tone` 为 `GaugeNeutral`、`GaugeGood`、`GaugeWarning`、`GaugeDanger`。 |
| `Unit` | `string` | 空 | 数值后的单位。 |
| `Full` | `bool` | false | 整圆表盘，从底部顺时针一周。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 数值格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 当前分段加宽显示。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：依次读当前值与各分段，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。

## 限制

- 不做指针动画，值变化时直接切换。
