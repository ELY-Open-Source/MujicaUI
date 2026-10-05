# 486 CurrencyInput 货币输入

## 用途

以最小单位（分）的整数保存金额，避免浮点误差；显示币种符号，未聚焦时按千分位分组，聚焦后显示纯数字便于编辑。

## 最小示例

```go
cents := int64(123450)
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.CurrencyInput(c, &cents, mujica.CurrencyInputOptions{Symbol: "$", Label: "Price"})
}
```

签名：`func CurrencyInput(c *ui.Context, value *int64, opts CurrencyInputOptions) InputResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*int64` | 必填 | 金额的最小单位数。 |
| `Symbol` | `string` | `""` | 前缀符号，如 `$`、`¥`。 |
| `Decimals` | `int` | 2 | 小数位，0–6，越界 panic。 |
| `NoFraction` | `bool` | `false` | 无小数（如日元）；与 `Decimals` 同设 panic。 |
| `Label` | `string` | `""` | 无障碍名称。 |
| `Disabled` | `bool` | `false` | 不可编辑。 |
| `Error` | `string` | `""` | 调用方错误信息。 |

## 状态

非法输入或小数位过多时保留输入并显示具体错误，值不变。

## 事件

`Changed()`：合法输入改变了值。

## 键盘操作

Tab 聚焦；文本编辑、选择、剪贴板与撤销由 MyGo 编辑器提供。

## 限制

千分位与小数点固定为 `,` 与 `.`。
