# 044 PinInput

验证码输入。

## 用途

输入固定位数的数字验证码。值是字符串，保留前导零。

## 最小示例

```go
code := ""
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.PinInput(c, &code, mujica.PinInputOptions{Length: 6}).Completed() {
		fmt.Println("verify", code)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*string` | 必填 | 当前数字串；含非数字或超出位数时 panic。 |
| `Length` | `int` | 6 | 位数。 |
| `Label` | `string` | 内置“Verification code/验证码” | 无障碍名称。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Invalid` | `bool` | false | 格子边框显示危险色。 |

## 状态

默认、聚焦（当前格强调边框、焦点环与光标）、输入法组合中（当前格显示带下划线的组合文字）、禁用、错误。

## 事件

`PinInputParts.Changed()` 值改变；`Completed()` 在用户填满最后一格的那一帧为 true。

## 键盘操作

数字键依次填格；非数字被忽略；Backspace 清除最后一位；Cmd/Ctrl+V 或编辑菜单粘贴整段，取其中的数字填满各格。

## 限制

- 只支持数字。编辑按顺序进行：只能在末尾追加或删除，不能单独改中间某一格。
- 输入法组合期间不提交、不跳格；组合结束后才写入值。
- 整个组件是一个焦点停靠点，屏幕阅读器以文本框读出名称与“已输入 n / N 位”的说明。
