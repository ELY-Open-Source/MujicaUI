# 034 Slider · 滑块

## 用途

在连续或分步的数值范围内取值，如音量。滑块为菱形，命中区域是整个矩形控件。

## 最小示例

```go
volume := 40.0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Slider(c, &volume, mujica.SliderOptions{Step: 5, ShowValue: true, Label: "Volume"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Slider(c *ui.Context, value *float64, opts SliderOptions) ChoiceResult`。范围外的值在构建时被夹到范围内并对齐步长（不算一次用户变更）。`Min >= Max` 或 `Step < 0` 时 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Min`, `Max` | `float64` | 0, 0 → 0–100 | 范围；二者都为 0 时取 0–100。 |
| `Step` | `float64` | 0 | 步长；值为 `Min + n·Step`。0 为连续。 |
| `Marks` | `[]SliderMark` | nil | 刻度：`Value` 处画刻线，`Label` 非空时在下方显示文字。 |
| `ShowValue` | `bool` | false | 悬停、拖动或聚焦时在滑块上方显示数值提示（预留高度，不改变布局）。 |
| `Format` | `func(float64) string` | 按步长精度 | 数值文字。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | "" | 无障碍名称。 |

## 状态

默认、悬停（菱形放大 2 DIP）、拖动、聚焦（菱形外菱形焦点环）、禁用。已选段用 `AccentText`，深色主题对比度满足 3:1。

## 事件

`Changed()` 在值变化的那一帧返回 true。

## 键盘操作

- ←/↓、→/↑：减少或增加一个步长（连续模式为范围的 1%）。
- Home / End：最小值 / 最大值（对齐到步长）。

## 限制

- 基于 `ui.SliderBase`。MyGo 的方向键步长固定为范围的 1%（ui/base.go:116），本组件把每次按键结果换算为一个 `Step`；当 `Step` 小于范围的 1% 且滑块距端点不足 1% 时，方向键会直接到端点。
- `SliderBase` 在构建时就处理指针（ui/base.go:111），禁用时组件改为让它操作副本，保证值不变。
