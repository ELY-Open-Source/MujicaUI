# 442 StatCard · 统计卡片

## 用途

展示调用方统计值、变化率和迷你趋势图。

## 最小示例

```go
value := "12,480"
opts := devtools.StatCardOptions{Label: "Requests", Change: 0.12, HasChange: true, Trend: []float64{3, 5, 4, 8}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if devtools.StatCard(c, value, opts).Changed() {
		fmt.Println("retry")
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

Label 必填；Change 是相对变化；HasChange 开启变化说明；LowerIsBetter 反转好坏；Trend 为趋势值；Status、Error 控制加载与错误。

## 状态

值与变化率、趋势、加载、失败、父级禁用。

## 事件

Changed() 表示重试。

## 键盘操作

Tab 聚焦趋势或重试；方向键读取数据；Enter 激活重试。

## 限制

卡片最小高度 116 DIP，随桌面文本缩放增长；统计值由调用方计算。
