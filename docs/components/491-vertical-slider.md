# 491 VerticalSlider 纵向滑块

## 用途

纵向放置的 Slider：上大下小，规则与 Slider 相同（步长网格、刻度、数值气泡）。

## 最小示例

```go
level := 40.0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.VerticalSlider(c, &level, input.VerticalSliderOptions{Step: 10, Label: "Level"})
}
```

签名：`func VerticalSlider(c *ui.Context, value *float64, opts VerticalSliderOptions) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Min`, `Max` | `float64` | 0–100 | 范围。 |
| `Step` | `float64` | 连续 | 步长。 |
| `Marks` | `[]SliderMark` | `nil` | 刻度与标签。 |
| `Height` | `float32` | 160 | 轨道长度；负数 panic。 |
| `ShowValue` | `bool` | `false` | 悬停、拖动或聚焦时显示数值。 |
| `Format` | `func(float64) string` | 自动 | 数值文字。 |
| `Disabled` | `bool` | `false` | 不可调。 |
| `Label` | `string` | `""` | 名称。 |

## 状态

菱形滑块，悬停与拖动放大；填充段为 AccentText。

## 事件

`Changed()`：值变化。

## 键盘操作

Up/Right 增大，Down/Left 减小，Home/End 到两端；点击或拖动设置。

## 限制

无。
