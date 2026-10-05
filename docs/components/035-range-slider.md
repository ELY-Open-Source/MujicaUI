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
| `Format` | `func(float64) string` | 按 `Min`、`Step` 所需小数位（连续时至少 1 位） | 数值文字。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | "" | 区间名称；两个滑块分别名为“`Label` 最小值”“`Label` 最大值”，数值与范围由 Slider 角色报告。 |

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

- MyGo 自带的 `ui.RangeSlider` 使用玻璃质感滑块（ui/widgets.go:419），不符合本库风格；MyGo 也没有双滑块 Base，故每个滑块各是一个 `SliderBase`（报告 Slider 角色、范围与当前值，ui/base.go:128），拖动由外层的 `Dragged`、`PointerPosition` 处理。
- 两个滑块的范围都报告为整条轨道 `Min`–`Max`；“不能越过另一端”这一动态上下限无法告诉辅助技术。
- 按键处理与 Slider 相同：距端点不足 1% 时等到按键松开再移动（见 034 限制）。
