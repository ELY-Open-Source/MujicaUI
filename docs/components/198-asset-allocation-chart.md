# 198 AssetAllocationChart

## 用途

资产配置图。

## 最小示例

```go
items := []finance.AllocationItem{{Name: "Stocks", Value: 60, Target: 50}, {Name: "Bonds", Value: 30, Target: 40}, {Name: "Cash", Value: 10, Target: 10}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.AssetAllocationChart(c, items, finance.AssetAllocationChartOptions{Frame: chart.ChartFrame{Label: "Allocation"}})
	if r.Changed() || r.Retried() {
		fmt.Println(r.Point)
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`Version` 在原地修改数据时递增。环形图合并小于 1 度的多个扇区；列表虚拟化保留每项的占比、目标与偏离。

`items` 包含 Name、Value、Target；目标总和为 100 或全部 0。`Frame.Label` 必填；`Tolerance` 默认 5 个百分点；`Format` 可覆盖数值。

## 状态

默认与数据边界；图表的 Frame 支持加载、空、失败占位。

## 事件

if r.Changed() || r.Retried() { 调用方回调(r.Point) }

## 键盘操作

Tab 聚焦环形图；方向键逐项读数；Home/End 到首尾。

## 限制

列表同时提供占比、目标、偏离；负资产值触发 panic。
