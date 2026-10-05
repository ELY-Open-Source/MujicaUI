# 264 PermissionPrompt 权限请求

## 用途

请求权限：允许一次、始终允许或拒绝，并说明“始终允许”覆盖的范围。

## 最小示例

```go
open := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.PermissionPrompt(c, &open, mujica.PermissionPromptOptions{
		Title: "Write to the ledger?", Scope: "Write access to ~/court/ledger",
	})
	if r.Choice == mujica.PermissionAllowAlways {
		fmt.Println("remember the grant")
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func PermissionPrompt(c *ui.Context, open *bool, opts PermissionPromptOptions) PermissionResult`。`Title` 或 `Scope` 为空时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 必填 | 问题。 |
| `Description` | `string` | `""` | 补充说明。 |
| `Scope` | `string` | 必填 | “始终允许”覆盖的范围。 |

## 状态

- 范围说明放在带锁形图标的框内。
- 打开时焦点在“拒绝”，点击遮罩无效。

## 事件

- `Choice`：本帧的选择，`PermissionAllowOnce`、`PermissionAllowAlways`、`PermissionDeny`；无选择为 `PermissionNone`。Escape 记为拒绝。

## 键盘操作

- 打开时焦点在拒绝；Tab 依次到“始终允许”“允许一次”；Enter 触发聚焦按钮；Escape 拒绝。

## 限制

- 组件不保存授权，记住“始终允许”由调用方负责。
