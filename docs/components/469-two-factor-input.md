# 469 TwoFactorInput

两步验证码。

## 用途

六位验证码输入（基于 PinInput），支持整段粘贴；重新发送带倒计时；显示调用方错误。

## 最小示例

```go
code := ""
resendAt := time.Now().Add(30 * time.Second)
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.TwoFactorInput(c, &code, mujica.TwoFactorInputOptions{ResendAt: resendAt})
	if r.Completed() {
		fmt.Println("verify the code")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `code` | `*string` | 必填 | 验证码数字。 |
| `ResendAt` | `time.Time` | 零值 | 此时间前重新发送按钮倒计时并禁用。 |
| `Error` | `string` | "" | 调用方错误，如过期。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

输入中、完成、错误、倒计时、可重发、禁用。

## 事件

`TwoFactorResult.Completed()` 在填满第六位时；`Resend()` 请求重发。

## 键盘操作

数字键填格，Backspace 删除，Cmd/Ctrl+V 粘贴；Tab 到重新发送按钮。

## 限制

- 验证码与密钥通过调用方状态接入；组件事件与日志仅包含完成或失败标记。
- 倒计时每秒刷新一次，结束后不再重绘。
