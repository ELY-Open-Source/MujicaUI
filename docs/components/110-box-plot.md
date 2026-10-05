# 110 BoxPlot 箱线图

## 用途

箱线图：四分位箱、中位线、1.5 倍四分距内的须线与离群点，多组对比。统计只在数据变化时计算。

## 最小示例

```go
groups := []mujica.BoxGroup{{Name: "Choir", Values: []float64{1, 2, 3, 4, 5, 40}}, {Name: "Nave", Values: []float64{4, 5, 6, 7}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.BoxPlot(c, groups, mujica.BoxPlotOptions{Frame: mujica.ChartFrame{Label: "Hymns"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func BoxPlot(c *ui.Context, groups []BoxGroup, opts BoxPlotOptions) ChartResult`；`BoxGroup{Name, Values}`，名称必填。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Frame.Label` | `string` | 必填 | 图表的无障碍名称；为空 panic。 |
| `Frame.Height` | `float32` | `240` | 绘图区高度（DIP）；负数 panic。 |
| `Frame.Status` | `DataStatus` | `DataReady` | `DataLoading` 显示加载标记，`DataFailed` 显示错误与“重试”；占位与图表同尺寸。 |
| `Frame.Error` | `string` | 内置“加载失败” | 失败时的具体说明。 |
| `Frame.Empty` | `string` | 内置“暂无数据” | 无数据时的说明。 |
| `Format` | `func(float64) string` | 6 位有效数字 | 读数格式。 |

## 状态

默认、悬停读数、键盘读数（焦点环 + 当前点加粗或描边）、加载、失败（可重试）、无数据。图表区始终保持同样尺寸。 离群点画为空心圆。

## 事件

返回 `ChartResult`：`Point` 为当前读数点（无则 -1）；`Changed()` 读数点变化；`Retried()` 点击了重试。

## 键盘操作

- Tab：聚焦绘图区（图例勾选框在其前）。
- ←/→（以及 ↑/↓）：按组移动，读数显示在图表下方并向读屏软件播报；第一次按键进入第一个点。
- Home/End：第一个/最后一个点；Escape：清除读数。

## 限制

- 四分位数用线性插值计算。
