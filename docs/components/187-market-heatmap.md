# 187 MarketHeatmap

## 用途

市场热力图。

## 最小示例

```go
var path []string
root := mujica.HeatNode{Name: "Market", Children: []mujica.HeatNode{{Name: "Technology", Children: []mujica.HeatNode{{Name: "AAPL", Cap: 3000, Change: 1.2}, {Name: "MSFT", Cap: 2800, Change: -0.6}}}, {Name: "Energy", Cap: 450, Change: -2.1}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.MarketHeatmap(c, &path, root, mujica.MarketHeatmapOptions{Frame: mujica.ChartFrame{Label: "Market"}})
	if r.Changed() || r.ValueChanged() || r.Retried() {
		fmt.Println(path, r.Point)
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`Version` 在原地修改数据时递增。小于 4 DIP 的瓦片按屏幕网格汇总绘制，键盘仍逐项读取原始数据。

`path` 保存下钻路径；`root` 名称必填，子项名称唯一；`Cap` 为面积，`Change` 为百分比。`Frame.Label` 必填；`Scale` 默认 3。

## 状态

默认与数据边界；图表的 Frame 支持加载、空、失败占位。

## 事件

if r.Changed() || r.ValueChanged() || r.Retried() { 调用方回调(path, r.Point) }

## 键盘操作

Tab 聚焦图；方向键逐块读数；Enter 下钻；Backspace 返回；路径按钮可用 Tab。

## 限制

市场数据由调用方提供；原地修改数据后更换 Children 切片或递增 Version。
