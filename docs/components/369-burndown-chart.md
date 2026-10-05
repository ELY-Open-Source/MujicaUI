# 369 BurndownChart · 燃尽图

## 用途

燃尽图用于项目管理界面。

## 最小示例

```go
remaining := []float64{40, 32, 25, 18}
opts := mujica.BurndownChartOptions{Frame: mujica.ChartFrame{Label: "Release burndown"}, Start: mujica.Date{Year: 2026, Month: 10, Day: 1}, End: mujica.Date{Year: 2026, Month: 10, Day: 14}, Total: 40}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.BurndownChart(c, remaining, opts)
	if r.Changed() {
		log.Println("point", r.Point)
	}
	if name, ok := r.Toggled(); ok {
		log.Println("series", name)
	}
	if r.Retried() {
		log.Println("retry")
	}
}
ui.Render(view, 1000, 700, 1)
```

## 参数

`Frame` 提供必填 `Label`、默认 240 DIP `Height`、`Status/Error/Empty`；`Start < End`，`Total > 0`；`Format` 自定义读数；`Version` 刷新原地修改。

## 状态

业务状态由调用方持有；主题、密度与文本缩放遵循窗口设置。交互控件提供焦点、悬停和禁用状态；展示组件保留语义名称。

## 事件

`Changed()` 与 `Point` 返回活跃日期索引；`Toggled()` 返回图例名称；`Retried()` 请求重试。

## 键盘操作

Tab 聚焦图例与图表；Space 切换图例；方向键逐日读取；Home/End 首尾；Escape 清除读数。

## 限制

实际值以 `NaN` 表示缺失；数组从 Start 连续对应日期。数据加载与失败由调用方提供；数值须有限且非负。
