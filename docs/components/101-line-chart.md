# 101 LineChart 折线图

## 用途

多序列折线图：数值或时间横轴，缺值（NaN）处断开，悬停与键盘读数，图例勾选框控制序列显隐。大数据按像素列抽稀绘制，10 万点保持流畅。

## 最小示例

```go
series := []mujica.ChartSeries{
	{Name: "Pews", Values: []float64{3, 5, math.NaN(), 4, 8}},
	{Name: "Candles", Values: []float64{5, 2, 6, 7, 1}},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.LineChart(c, series, mujica.LineChartOptions{Frame: mujica.ChartFrame{Label: "Chapel"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func LineChart(c *ui.Context, series []ChartSeries, opts LineChartOptions) ChartResult`；`ChartSeries{Name, Values}`，各序列等长且名称唯一。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `X` | `[]float64` | nil | 点 i 的横坐标，须递增；nil 时为下标。 |
| `Times` | `[]time.Time` | nil | 时间轴；与 `X` 二选一，须递增。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 读数中的数值格式。 |

## 状态

默认、悬停或键盘读数（竖线 + 各序列放大标记）、隐藏序列（图例划线、空心标记）、加载、失败、无数据（全部为缺值时）。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试；`Toggled()` 返回本帧被切换显隐的序列名。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：在数据点之间移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。
- 图例勾选框：Space 切换序列显隐。

## 限制

- 数据切片以首元素地址与长度识别；原地修改数据后请传入新切片，缓存的坐标范围才会更新。
- 每条序列除颜色外还有不同形状的标记（圆、方、菱形、三角），密集时按间隔绘制。
