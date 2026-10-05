# 026 NumberInput 数字输入框

## 用途

输入与步进数值，支持范围、步长、精度与空值。合并 MyGo 的 Number input（输入模式）与 Stepper（`StepperOnly` 仅步进模式）。

## 最小示例

```go
if mujica.NumberInput(c, &copies, mujica.NumberInputOptions{Min: 1, Max: 99, Label: "Copies"}).Changed() {
	fmt.Println("copies", copies)
}
```

签名：`func NumberInput(c *ui.Context, value *float64, opts NumberInputOptions) InputResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*float64` | 必填 | 数值；空值用 `math.NaN()` 表示，仅在 `AllowEmpty` 时允许。 |
| `Min` / `Max` | `float64` | 都为 0 | 范围；相等（含零值）表示不限范围；`Min > Max` 时 panic。 |
| `Step` | `float64` | `1` | 方向键与按钮的步长；负数 panic。 |
| `Precision` | `int` | 取 `Step` 的小数位 | 保留与显示的小数位。 |
| `AllowEmpty` | `bool` | `false` | 允许清空，值变为 NaN；从空值步进时从 0（夹在范围内）开始。 |
| `StepperOnly` | `bool` | `false` | 只显示数值与上下箭头，不可输入；必须给出 `Min < Max`，且不能与 `AllowEmpty` 同用。 |
| `Placeholder` | `string` | `""` | 空值时的占位。 |
| `Label` | `string` | `""` | 无障碍名称。 |
| `Disabled` | `bool` | `false` | 禁用。 |
| `Error` | `string` | `""` | 调用方给出的错误。 |

## 状态

输入模式同 TextInput；输入非数字或超出范围时显示内置错误（“请输入数字”“请输入 a 到 b 之间的数字”），值保持上一次有效值；失去焦点后文本恢复为当前值并清除该错误。到达边界时对应按钮禁用。

## 事件

`Changed()`：值因输入、方向键、按钮或步进器改变。

## 键盘操作

输入模式：Up/Down 按步长增减。仅步进模式：Up/Right 增加、Down/Left 减少、Home 到最小、End 到最大；按住箭头持续步进。

## 限制

- 输入模式下的加减按钮不进入 Tab 序列，键盘用 Up/Down。
- 仅步进模式复用 MyGo `ui.Stepper`，箭头高度来自主题间距，不随密度的控件高度变化。
- MyGo `ui.Stepper` 的箭头不检查禁用（`indicators.go:216`），禁用时组件向它传入副本，保证值不被修改。
- MyGo 的步进器吸附依赖有限边界（`widgets.go` 中 `snap`），所以仅步进模式要求有界。
