# 190 IndexCard

## 用途

指数卡片。

## 最小示例

```go
quote := mujica.IndexQuote{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.IndexCard(c, quote, mujica.IndexCardOptions{Clickable: true})
	if r.Changed() {
		fmt.Println(quote.Name)
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`q.Name` 必填；Price、Change、ChangePercent 和 Series 来自调用方。`Decimals` 默认 2；`Clickable/Disabled` 默认 false。

## 状态

默认；交互区域支持聚焦与选中，Disabled 参数控制禁用。

## 事件

if r.Changed() { 调用方回调(quote.Name) }

## 键盘操作

Tab 聚焦卡片或迷你图；Enter/Space 激活卡片，图表方向键读数。

## 限制

迷你图复用 Sparkline；行情获取由调用方负责。
