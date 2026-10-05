# 117 Sunburst 旭日图

## 用途

层级环形图：中心是当前节点及总计，外圈依次是各层后代；点击或回车下钻，路径面包屑与返回按钮回到上层。

## 最小示例

```go
var path []string
root := chart.ChartNode{Name: "Abbey", Children: []chart.ChartNode{
	{Name: "Library", Children: []chart.ChartNode{{Name: "Psalters", Value: 30}, {Name: "Maps", Value: 10}}},
	{Name: "Kitchen", Value: 25},
}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.Sunburst(c, &path, root, chart.SunburstOptions{Frame: chart.ChartFrame{Label: "Rooms"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Sunburst(c *ui.Context, path *[]string, root ChartNode, opts SunburstOptions) ChartResult`；数据结构同 Treemap。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Depth` | `int` | `3` | 显示的层数；负数 panic。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 数值格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 外层颜色向表面色淡化，但保持 3:1 对比。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试；`ValueChanged()` 表示 `*path` 改变。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：按深度优先顺序在扇区间移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。
- Enter/Space：下钻到当前扇区；Backspace：返回上一层。

## 限制

- 弧长不足时不写标签，名称在读数中完整给出。
