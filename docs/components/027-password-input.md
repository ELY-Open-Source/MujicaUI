# 027 PasswordInput 密码输入框

## 用途

遮蔽显示的密码输入，带显示/隐藏切换。允许粘贴；遮蔽时不允许复制与剪切。基于 `ui.TextInputBase(...).Password()`，辅助技术读不到遮蔽中的值。组件从不记录密码内容。

## 最小示例

```go
if input.PasswordInput(c, &password, input.PasswordInputOptions{Label: "Password"}).Submitted() {
	fmt.Println("signing in")
}
```

签名：`func PasswordInput(c *ui.Context, value *string, opts PasswordInputOptions) InputResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*string` | 必填 | 密码。 |
| `Placeholder` | `string` | `""` | 占位文字。 |
| `Label` | `string` | `""` | 无障碍名称。 |
| `Disabled` | `bool` | `false` | 禁用（切换按钮同时禁用）。 |
| `Error` | `string` | `""` | 错误文字。 |

## 状态

默认、悬停、聚焦、错误、禁用；遮蔽 / 显示两种模式，切换按钮的图标与名称（“显示密码”“隐藏密码”）随之变化。

## 事件

`Changed()`：内容改变。`Submitted()`：Enter。日志只写切换事件（`mujica: password_input reveal`，含是否显示），不写内容。

## 键盘操作

Tab 依次到输入框、切换按钮；切换按钮上 Space 或 Enter 切换；输入框中 Enter 提交，Cmd/Ctrl+V 粘贴。

## 限制

- MyGo 的 `Password()` 一旦设置无法撤销（`editor.go:791`），因此遮蔽与显示两种模式各用一个编辑器，切换时撤销历史和光标位置不保留，焦点留在切换按钮上。
- 显示模式下内容是普通文本，可复制，辅助技术可朗读——这是用户主动选择的结果。
