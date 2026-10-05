# 492 ColorPalette 预设色板

## 用途

从预设色块中选择颜色，可附带自定义颜色入口（复用 ColorPicker）。选中色块有内框与对勾，标记颜色按对比度自动取深或浅。

## 最小示例

```go
color := ui.Hex("#7D2034")
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.ColorPalette(c, &color, []input.ColorSwatch{
		{Color: ui.Hex("#7D2034"), Name: "Wine"},
		{Color: ui.Hex("#AD9164"), Name: "Old gold"},
		{Color: ui.Hex("#355E80"), Name: "Slate blue"},
	}, input.ColorPaletteOptions{Custom: true, Label: "Tag color"})
}
```

签名：`func ColorPalette(c *ui.Context, value *ui.Color, swatches []ColorSwatch, opts ColorPaletteOptions) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `swatches` | `[]ColorSwatch` | 必填 | 颜色与名称；缺名称 panic。 |
| `Columns` | `int` | 8 | 每行色块数；负数 panic。 |
| `Custom` | `bool` | `false` | 显示自定义颜色入口。 |
| `Disabled` | `bool` | `false` | 不可选。 |
| `Label` | `string` | `""` | 组名称。 |

## 状态

悬停加深边框；选中有内框与对勾；每个色块提示其名称。

## 事件

`Changed()`：颜色变化（含自定义）。

## 键盘操作

色板是一个 Tab 停靠点；左右键逐个移动，上下键跨行，Enter/Space 选择。

## 限制

无。
