# 118 SankeyChart 桑基图

## 用途

桑基图：节点按最长路径分列，列内顺序与传入顺序一致（稳定）；流量带宽度与值成比例；悬停或键盘读到节点或流量时，高亮经过它的全部路径。

## 最小示例

```go
nodes := []chart.SankeyNode{{ID: "tithe", Label: "Tithes"}, {ID: "fund", Label: "Fund"}, {ID: "roof", Label: "Roof"}}
links := []chart.SankeyLink{{From: "tithe", To: "fund", Value: 60}, {From: "fund", To: "roof", Value: 50}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.SankeyChart(c, nodes, links, chart.SankeyChartOptions{Frame: chart.ChartFrame{Label: "Money"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func SankeyChart(c *ui.Context, nodes []SankeyNode, links []SankeyLink, opts SankeyChartOptions) ChartResult`；`SankeyNode{ID, Label}`，`SankeyLink{From, To, Value}`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 数值格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 高亮时其余流量淡化。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：先按节点、再按流量移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。

## 限制

- 链接成环、端点未知或值不为正时 panic。
- 布局只在节点或链接切片变化时计算。
