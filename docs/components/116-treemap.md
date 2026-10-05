# 116 Treemap 矩形树图

## 用途

层级矩形树图：按大小平铺当前节点的子项（squarify 布局），点击或回车下钻有子项的矩形，返回按钮与面包屑回到上层；标签按面积显示并截断。

## 最小示例

```go
var path []string
root := chart.ChartNode{Name: "Abbey", Children: []chart.ChartNode{
	{Name: "Library", Children: []chart.ChartNode{{Name: "Psalters", Value: 30}, {Name: "Maps", Value: 10}}},
	{Name: "Kitchen", Value: 25},
}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.Treemap(c, &path, root, chart.TreemapOptions{Frame: chart.ChartFrame{Label: "Rooms"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Treemap(c *ui.Context, path *[]string, root ChartNode, opts TreemapOptions) ChartResult`；`ChartNode{Name, Value, Children}`，叶子的 `Value` 为大小，分支为子项之和；`*path` 是从根的子项起的名称路径。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 数值格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 当前矩形内描边。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试；`ValueChanged()` 表示下钻、返回或点击面包屑改变了 `*path`。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：按大小顺序在矩形间移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。
- Enter/Space：下钻；Backspace：返回上一层；返回按钮与面包屑可用 Tab 到达。

## 限制

- 名称为空、兄弟同名或值为负时 panic；`*path` 指向不存在的节点时 panic。
- 只显示当前节点的一层子项。
