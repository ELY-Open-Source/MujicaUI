# 490 Knob 旋钮

## 用途

圆形旋钮：上下拖动或方向键调整数值，显示刻度、数值弧和菱形指针。

## 最小示例

```go
gain := 35.0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Knob(c, &gain, mujica.KnobOptions{Label: "Gain", Ticks: 11})
}
```

签名：`func Knob(c *ui.Context, value *float64, opts KnobOptions) ChoiceResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Min`, `Max` | `float64` | 0–100 | 范围；`Max <= Min` panic。 |
| `Step` | `float64` | 连续 | 步长，按 NumberInput 网格吸附。 |
| `Ticks` | `int` | 0 | 刻度数；负数 panic。 |
| `Size` | `float32` | 56 | 直径。 |
| `Format` | `func(float64) string` | 自动 | 数值文字。 |
| `Disabled` | `bool` | `false` | 不可调。 |
| `Label` | `string` | “Knob” | 名称与标题。 |

## 状态

悬停与拖动时表盘变亮；焦点环为圆形。基于 MyGo SliderBase，读屏视为滑块。

## 事件

`Changed()`：值变化。

## 键盘操作

方向键按步调整（连续时每步 1%），Home/End 到两端；上拖增大，160 DIP 走完全程。

## 限制

无。
