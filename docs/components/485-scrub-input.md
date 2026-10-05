# 485 ScrubInput 拖动数值输入

## 用途

按住标签左右拖动调整数值，每 4 DIP 一步；旁边的 NumberInput 可直接输入。范围、步长与精度规则与 NumberInput 相同。

## 最小示例

```go
opacity := 0.6
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ScrubInput(c, "Opacity", &opacity, mujica.ScrubInputOptions{Max: 1, Step: 0.1})
}
```

签名：`func ScrubInput(c *ui.Context, label string, value *float64, opts ScrubInputOptions) InputResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `label` | `string` | 必填 | 可拖动的标签，也是字段名称。 |
| `Min`, `Max` | `float64` | 无界 | 相等表示无界。 |
| `Step` | `float64` | 1 | 每一步的增量。 |
| `Precision` | `int` | 取自 Step | 保留小数位。 |
| `Disabled` | `bool` | `false` | 不可拖动也不可输入。 |

## 状态

标签悬停或拖动时变为强调色；拖到边界即停止。

## 事件

`Changed()`：拖动或输入改变了值；`Submitted()` 同 NumberInput。

## 键盘操作

标签不可聚焦；键盘在输入框内用 Up/Down 步进或直接输入。

## 限制

拖动只对指针可用。
