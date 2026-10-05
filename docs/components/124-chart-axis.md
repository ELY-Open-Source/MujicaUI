# 124 ChartAxis 坐标轴

## 用途

提供数值、时间、类别和对数刻度，以及值和像素之间的映射。

## 最小示例

```go
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	axis := mujica.ChartAxis(c, mujica.ChartAxisOptions{Side: mujica.SideLeft, Kind: mujica.AxisLog, Min: 1, Max: 1000})
	ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
		axis.Draw(p, ui.Rect{X: r.X + 60, Y: r.Y + 20, W: r.W - 80, H: r.H - 60})
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`Kind` 默认 `AxisLinear`；`Side` 默认 `SideBottom`。数值与对数轴需要有限 `Min < Max`，对数轴 `Min > 0`。时间轴需要 `From < To`。类别轴需要非空 `Categories`。`Ticks` 为递增且位于范围内的指定刻度，nil 自动生成。`Format func(float64) string` 自定义标签。

## 状态

展示部件，继承当前主题、文本缩放和字体。

## 事件

`Map`、`Invert`、`Ticks`、`Label` 返回映射与标签；组件自身没有业务事件。

## 键盘操作

所属 Chart 负责聚焦、方向键和读数。坐标轴作为绘图部件随图表展示。

## 限制

为标签预留边距。时间值采用 Unix 秒，默认格式使用 UTC；需要本地时区时通过 Format 格式化。数据和刻度切片在使用期间保持稳定。
