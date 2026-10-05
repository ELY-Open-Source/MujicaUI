# 119 NetworkGraph 网络图

## 用途

网络图：节点与边，使用调用方给定的坐标或简单的力导布局；拖动节点、拖动空白处平移、滚轮或 +/- 缩放。

## 最小示例

```go
var state mujica.NetworkGraphState
nodes := []mujica.GraphNode{{ID: "a", Label: "Abbot"}, {ID: "b", Label: "Prior"}, {ID: "c", Label: "Novice"}}
edges := []mujica.GraphEdge{{From: "a", To: "b"}, {From: "b", To: "c"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.NetworkGraph(c, &state, nodes, edges, mujica.NetworkGraphOptions{Frame: mujica.ChartFrame{Label: "House"}, Force: true})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func NetworkGraph(c *ui.Context, s *NetworkGraphState, nodes []GraphNode, edges []GraphEdge, opts NetworkGraphOptions) ChartResult`；`GraphNode{ID, Label, X, Y}`，`GraphEdge{From, To}`。`NetworkGraphState` 保存拖动后的位置与视图：`Position(id)`、`Zoom()`、`ResetView()`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Force` | `bool` | false | 用力导模拟布局（固定 300 步，结果确定）；最多 500 个节点。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 当前节点外加圆环，相连的边加粗。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：在节点间移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。
- Shift+方向键：移动当前节点；+/-：缩放；0：复位视图。

## 限制

- 布局只在节点或边的切片变化时重新计算，拖动不会重新适配视图。
- 不支持边标签与有向箭头。
