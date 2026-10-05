# 467 LoginForm

登录表单。

## 用途

账号、密码与“记住我”；提交中禁用输入；显示调用方给出的错误（含频率限制、锁定）。密码在私有临时状态中编辑，经提交回调传出，完成后清除；账号校验由调用方完成。

## 最小示例

```go
var state mujica.LoginFormState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.LoginForm(c, &state, mujica.LoginFormOptions{OnSubmit: func(account, password string, remember bool) {
		// Submit credentials through the caller-provided auth service.
		state.EndSubmit("")
	}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `s` | `*LoginFormState` | 必填 | Account、Remember 可读写；密码不可读取。 |
| `OnSubmit` | `func(account, password string, remember bool)` | 必填 | 唯一接收凭据的出口。 |
| `AccountLabel` | `string` | 内置“Email or username” | 账号标签。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认、必填错误、提交中、失败（错误文字）、禁用。

## 事件

`OnSubmit` 每次提交调用一次；完成后调用 `s.EndSubmit(err)`，密码随即清除。

## 键盘操作

Tab 顺序：账号、密码、显示密码、记住我、登录；任一字段中 Enter 提交；缺项时聚焦第一个空字段。

## 限制

- 凭据只经调用方回调传出；密码、验证码与密钥不会写入日志、Label 或 Description。
- 提交失败后需重新输入密码。
