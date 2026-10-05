# 189 MarketOverview

## 用途

市场概览。

## 最小示例

```go
indices := []mujica.IndexQuote{{Name: "S&P 500", Price: 5751.07, Change: 51.13, ChangePercent: 0.9, Series: []float64{5700, 5730, 5720, 5751}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.MarketOverview(c, indices, mujica.MarketOverviewOptions{Breadth: mujica.MarketBreadth{Advancing: 2100, Declining: 1400, Unchanged: 300}})
	if r.Chosen >= 0 {
		fmt.Println(indices[r.Chosen].Name)
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`indices` 为指数报价组；`Breadth` 是非负涨跌家数；`AsOf` 零值隐藏时间戳；`Decimals` 默认 2；`Disabled` 默认 false。

## 状态

默认；交互区域支持聚焦与选中，Disabled 参数控制禁用。

## 事件

if r.Chosen >= 0 { 调用方回调(indices[r.Chosen].Name) }

## 键盘操作

Tab 到达指数卡片；Enter/Space 选择。

## 限制

提供展示与点击事件；行情获取由调用方负责。
