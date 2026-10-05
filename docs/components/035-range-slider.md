# 035 RangeSlider · 范围滑块

## 用途

用两个菱形滑块选择一个区间，如价格范围。两端不会交叉（允许相等）。

## 最小示例

```go
low, high := 20.0, 80.0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.RangeSlider(c, &low, &high, mujica.RangeSliderOptions{Step: 10, Label: "Price"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func RangeSlider(c *ui.Context, low, high *float64, opts RangeSliderOptions) ChoiceResult`。`*low > *high` 时两者交换；超出范围的值被夹回并对齐步长。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Min`, `Max` | `float64` | 0, 0 → 0–100 | 范围。 |
| `Step` | `float64` | 0 | 步长，0 为连续。 |
| `Marks` | `[]SliderMark` | nil | 刻度与文字。 |
| `ShowValue` | `bool` | false | 悬停、拖动或聚焦时显示两端数值；低端提示向左、高端向右，互不遮挡。 |
| `Format` | `func(float64) string` | 按步长精度 | 数值文字。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | "" | 区间名称；两个滑块分别读作“`Label` 最小值, 20”“`Label` 最大值, 80”。 |

## 状态

默认、悬停、拖动、聚焦（各滑块独立焦点环）、禁用。

## 事件

`Changed()` 在任一端变化的那一帧返回 true。

## 键盘操作

- Tab：依次聚焦低端、高端滑块。
- ←/↓、→/↑：一个步长；不能越过另一端。
- Home / End：到范围端点或另一端所在位置。
- 指针按下轨道任意处：移动较近的一端，并把焦点交给它。

## 限制

- MyGo 自带的 `ui.RangeSlider` 使用玻璃质感滑块（ui/widgets.go:419），不符合本库风格；MyGo 也没有双滑块 Base，故本组件用公开的 `Dragged`、`PointerPosition`、`Shortcut` 组合实现。
- 滑块的无障碍角色为 Slider，但数值范围属性只能由 `SliderBase` 设置（ui/base.go:128），因此当前值写在名称里（如“Price 最小值, 20”），屏幕阅读器读名称而非数值属性。
