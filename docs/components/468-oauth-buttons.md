# 468 OAuthButtons

第三方登录按钮。

## 用途

按调用方提供的名称显示“使用 X 继续”按钮组；不使用商标图形。

## 最小示例

```go
providers := []account.OAuthProvider{{ID: "sso", Name: "Company SSO"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if id, ok := account.OAuthButtons(c, providers, account.OAuthButtonsOptions{}).Chosen(); ok {
		fmt.Println("start sign-in with", id)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `providers` | `[]OAuthProvider` | 必填 | ID、Name 必填；Icon 可选，应为非商标图形。 |
| `Pending` | `string` | "" | 正在进行的提供方：显示加载，其余禁用。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认、进行中、禁用。

## 事件

`OAuthButtonsResult.Chosen()` 返回提供方 ID。

## 键盘操作

Tab 在按钮间移动，Enter/Space 选择。

## 限制

- 登录流程由调用方完成。
