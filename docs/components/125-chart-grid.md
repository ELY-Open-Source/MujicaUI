# 125 ChartGrid 图表网格

## 用途

按坐标轴主刻度和次刻度绘制对齐设备像素的网格线。

## 最小示例

```go
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	x := mujica.ChartAxis(c, mujica.ChartAxisOptions{Max: 10})
	y := mujica.ChartAxis(c, mujica.ChartAxisOptions{Side: mujica.SideLeft, Max: 100})
	grid := mujica.ChartGrid(c, mujica.ChartGridOptions{X: &x, Y: &y, Minor: 2})
	ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
		plot := ui.Rect{X: r.X + 60, Y: r.Y + 20, W: r.W - 80, H: r.H - 60}
		grid.Draw(p, plot)
		x.Draw(p, plot)
		y.Draw(p, plot)
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`X, Y *ChartAxisPart` 至少一条轴；`Minor` 默认 0，正整数表示每个主刻度间隔的等分数。类别轴的 Minor 必须为 0。负值 panic。

## 状态

展示部件；主网格使用 Border，次网格使用较淡的 Border。

## 事件

返回 `ChartGridPart`，由 `Draw(p, plot)` 绘制。

## 键盘操作

所属图表的键盘操作保持生效；网格随图表展示。

## 限制

对数轴次刻度按对数距离等分。大范围刻度按绘图像素密度抽样绘制。
