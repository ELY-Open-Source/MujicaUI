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
| `Format` | `func(float64) string` | 按 `Min`、`Step` 所需小数位（连续时至少 1 位） | 数值文字。 |
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

- 基于 `ui.SliderBase`，它向辅助技术报告 Slider 角色、范围与当前值（ui/base.go:128）。
- MyGo 只把方向键、Home、End 的按下交给 `SliderBase` 本身（ui/base.go:118，按键路由见 ui/input.go:495），以范围的 1% 移动并把结果舍入到 1e-9（ui/widgets.go:375），舍入可能让移动方向看起来相反。本组件让它操作副本，只用它判断“有按键按下”；方向与跳转取自松开的键（`HandleInput` 能收到 KeyUp），再按步长网格移动。因此键盘调整在松开按键时生效；按住不放时，松开后按 MyGo 收到的按下次数一次移动多步。
- 步长网格复用 NumberInput 的实现：值为 `Min + n·Step`，小数位取 `Min` 与 `Step` 所需，端点跳转（End/Home、另一端）先对齐网格再写值并触发事件；例如 `Max=10, Step=3` 时 End 得到 9。
- `SliderBase` 在构建时就处理指针（ui/base.go:111），禁用时外层先禁用，且只操作副本，保证值不变。
