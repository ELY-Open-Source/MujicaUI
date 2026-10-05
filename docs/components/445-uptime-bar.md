# 445 UptimeBar · 可用性条

## 用途

按日展示可用性，以条高与文字区分状态。

## 最小示例

```go
days := []mujica.UptimeDay{{Date: mujica.Date{Year: 2026, Month: 10, Day: 1}, Uptime: 1}, {Date: mujica.Date{Year: 2026, Month: 10, Day: 2}, Uptime: 0.91, Note: "Recovery"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.UptimeBar(c, days, mujica.UptimeBarOptions{Frame: mujica.ChartFrame{Label: "API uptime"}})
	if r.Changed() {
		fmt.Println("day", r.Point)
	}
	if r.Retried() {
		fmt.Println("retry")
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

Frame.Label 必填；Frame.Height 默认 36 DIP；Frame.Status、Error 控制状态；Degraded 默认 0.999，Down 默认 0.95；Uptime 为 0 到 1，NaN 表示缺少数据。

## 状态

正常、降级、中断、缺测、空数据、加载、失败、父级禁用。

## 事件

Changed() 报告读取点变化；Point 为日期索引；Retried() 报告重试。

## 键盘操作

Tab 聚焦图；方向键逐日移动；Home、End 到边界；Escape 清除读取。

## 限制

大数据按像素列汇总；更新数据时提交新切片。
