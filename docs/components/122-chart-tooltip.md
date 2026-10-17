# 122 ChartTooltip 读数浮层

## 用途

在指针或键盘数据点旁显示多序列读数，自动避让窗口边缘。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	k := core.Tokens(c)
	chart.ChartTooltip(c, 40, 40, []chart.ChartTooltipLine{{Text: "Lauds"}, {Text: "Pews 5", Color: k.Info}, {Text: "Candles 2", Color: k.Success, Shape: 1}})
}
ui.Render(view, 400, 200, 1)
```

## 参数

`x, y float32`：窗口 DIP 坐标。`lines []ChartTooltipLine` 必须非空；`Text` 为读数，`Color` 为标记色，透明色隐藏标记；`Shape` 为标记形状。

## 状态

160 ms 进入过渡；减少动画直接出现。宽度和高度受窗口限制，长文字换行，内容可滚动。

## 事件

返回浮层 `ui.Element`；读数事件由所属 Chart 提供。

## 键盘操作

通过 `Chart` 使用时，方向键改变锚点，Escape 清除浮层；浮层内容沿用滚动容器操作。

## 限制

调用方传入当前读数和锚点，负责显示时机；所属图表应同时显示并朗读读数。
