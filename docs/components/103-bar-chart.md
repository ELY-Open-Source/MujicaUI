# 103 BarChart 柱状图

## 用途

柱状图：纵向或横向，分组或堆叠，支持负值与柱端数值标签，键盘逐柱读数。

## 最小示例

```go
categories := []string{"Nave", "Choir", "Crypt"}
series := []mujica.ChartSeries{{Name: "Gold", Values: []float64{4, -2, 3}}, {Name: "Silver", Values: []float64{2, 5, 1}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.BarChart(c, categories, series, mujica.BarChartOptions{Frame: mujica.ChartFrame{Label: "Tithes"}, Labels: true})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func BarChart(c *ui.Context, categories []string, series []ChartSeries, opts BarChartOptions) ChartResult`；每条序列的值与类别一一对应，NaN 不画柱。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Horizontal` | `bool` | false | 类别在纵轴，柱横向延伸。 |
| `Stacked` | `bool` | false | 每个类别一根柱，正值向上、负值向下堆叠。 |
| `Labels` | `bool` | false | 在柱端写出数值（堆叠时不写）。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 数值格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 选中柱以描边标出，柱内有序列形状标记。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试；多序列时 `Toggled()` 报告图例切换。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：逐柱移动（按类别、再按序列），读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。

## 限制

- 极细的柱仍有 6 DIP 的悬停范围。
- 类别标签过密时会跳过重叠的标签，读数仍包含完整类别名。
