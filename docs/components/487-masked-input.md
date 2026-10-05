# 487 MaskedInput 掩码输入

## 用途

按格式掩码输入电话、卡号等：`9` 数字、`a` 字母、`*` 二者皆可，其余字符原样显示。值只保存填入的字符。

## 最小示例

```go
phone := ""
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.MaskedInput(c, &phone, mujica.MaskedInputOptions{Mask: "(999) 999-9999", Label: "Phone"})
}
```

签名：`func MaskedInput(c *ui.Context, value *string, opts MaskedInputOptions) InputResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Mask` | `string` | 必填 | 掩码；没有可填位 panic。 |
| `value` | `*string` | 必填 | 已填字符；超长或不符掩码 panic。 |
| `Label` | `string` | `""` | 无障碍名称。 |
| `Disabled` | `bool` | `false` | 不可编辑。 |
| `Error` | `string` | `""` | 错误信息。 |

## 状态

未填位置显示下划线；输入法组合期间显示带下划线的组合文字，不提交、不跳格，只有提交的文字进入掩码。

## 事件

`Changed()`：值变化；`Submitted()`：Enter。另有 `MaskFormat(mask, raw)` 生成带格式文字。

## 键盘操作

左右方向键、Home/End 移动光标；Backspace/Delete 删除；粘贴只取符合掩码的字符；Enter 提交。

## 限制

自绘编辑器：没有文字选择与撤销，复制得到未格式化的值。
