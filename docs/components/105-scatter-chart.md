# 105 ScatterChart 散点图

## 用途

散点图：按分组着色并配不同形状，点大小可映射为气泡，拖动框选或用键盘选中点。

## 最小示例

```go
var selected []chart.ScatterKey
groups := []chart.ScatterGroup{{Name: "Monks", Points: []chart.ScatterPoint{{X: 1, Y: 2, Size: 4}, {X: 3, Y: 5, Size: 9}}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.ScatterChart(c, &selected, groups, chart.ScatterChartOptions{Frame: chart.ChartFrame{Label: "Choir"}, Bubbles: true})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ScatterChart(c *ui.Context, selected *[]ScatterKey, groups []ScatterGroup, opts ScatterChartOptions) ChartResult`；`ScatterGroup{Name, Points []ScatterPoint{X, Y, Size}}`，`ScatterKey{Group, Point}` 标识一个点。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Bubbles` | `bool` | false | 按 `Size` 映射点面积（直径 4–18 DIP）。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 读数格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 选中的点外加圆环，读数注明“已选”。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试；`ValueChanged()` 表示 `*selected` 被框选或键盘改变；多组时 `Toggled()` 报告图例切换。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：在数据点之间移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。
- Enter/Space：选中或取消当前点。
- 在绘图区拖动：框选，替换现有选择。

## 限制

- 框选替换选择，不支持按住修饰键追加。
- 点数较多时逐点绘制，没有抽稀。
