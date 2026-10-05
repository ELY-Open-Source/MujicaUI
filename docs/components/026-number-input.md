# 026 NumberInput 数字输入框

## 用途

输入与步进数值，支持范围、步长、精度与空值。合并 MyGo 的 Number input（输入模式）与 Stepper（`StepperOnly` 仅步进模式）；两种模式的步进计算都由组件完成。

## 最小示例

```go
if input.NumberInput(c, &copies, input.NumberInputOptions{Min: 1, Max: 99, Label: "Copies"}).Changed() {
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
| `Precision` | `int` | 取 `Step` 的小数位 | 保留与显示的小数位。范围边界向内取到该精度上（如 [0.11, 0.29] 精度 1 时有效值只有 0.2）；范围内没有任何该精度的值时 panic。 |
| `AllowEmpty` | `bool` | `false` | 允许清空，值变为 NaN；从空值步进时从 0（夹在范围内）开始。 |
| `StepperOnly` | `bool` | `false` | 只显示数值与上下箭头，不可输入；必须给出 `Min < Max`，且不能与 `AllowEmpty` 同用。 |
| `Placeholder` | `string` | `""` | 空值时的占位。 |
| `Label` | `string` | `""` | 无障碍名称。 |
| `Disabled` | `bool` | `false` | 禁用；任一祖先容器禁用时同样生效。 |
| `Error` | `string` | `""` | 调用方给出的错误。 |

## 状态

输入模式同 TextInput；输入非数字或超出范围时显示内置错误（“请输入数字”“请输入 a 到 b 之间的数字”），值保持上一次有效值；失去焦点后错误文本与提示保留，直到用户改正，或调用方从外部改变了值。值在聚焦期间被外部修改时，输入框立即显示新值。到达边界时对应按钮禁用。

## 事件

`Changed()`：值因输入、方向键、按钮或步进器改变。

## 键盘操作

步进沿以 `Min`（无界时为 0）为起点、`Step` 为间隔的网格移动，结果按精度取整，并停在范围内（范围先向内取到精度上，如 `Max` 0.29999999、步长 0.1 时最大值为 0.2）。不在网格上的值先走到该方向上最近的网格点。网格序号按整数计算，结果与原始 `Min`/`Max` 精确比较，越界则退回一格，不使用容差。

输入模式：Up/Down 按步长增减。仅步进模式：Up/Right 增加、Down/Left 减少、Home 到 `Min`、End 到 `Max`（如 [0, 10] 步长 3 时 End 为 10，再 Down 为 9）。

## 限制

- 输入模式下的加减按钮不进入 Tab 序列，键盘用 Up/Down。
- 仅步进模式不使用 MyGo `ui.Stepper`：它的箭头不检查禁用（`ui/indicators.go:216`），吸附固定按 9 位小数取整（`ui/widgets.go` 中 `snap`），步长 1e-10 时无法前进。组件自绘箭头并自行计算；按住箭头不会连续步进，需重复点击或使用方向键。
- 仅步进模式要求有界，以便 Home/End 有目标值。
- 仅步进模式的值通过无障碍描述朗读，MyGo 未公开为自定义元素设置数值范围的接口。
