# 113 BulletChart 子弹图

## 用途

子弹图：实际值的细条、目标值的垂直线与定性区间带，可横向或纵向。

## 最小示例

```go
target := 250.0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.BulletChart(c, 270, chart.BulletChartOptions{Frame: chart.ChartFrame{Label: "Alms"}, Target: &target, Bands: []float64{150, 225, 300}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func BulletChart(c *ui.Context, value float64, opts BulletChartOptions) ChartResult`；值须 ≥ 0。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `72` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Target` | `*float64` | nil | 目标值，画为与刻度垂直的粗线。 |
| `Bands` | `[]float64` | nil | 各区间的上界，须为正且递增。 |
| `Max` | `float64` | 自动取整 | 刻度终点；0 按值、目标与区间取整。 |
| `Vertical` | `bool` | false | 纵向摆放。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 数值格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 区间由浅到深的中性色带表示，读数写出区间范围。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：依次读实际值、目标值与各区间，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。

## 限制

- 超出 `Max` 的值画到刻度终点，读数仍为原值。
