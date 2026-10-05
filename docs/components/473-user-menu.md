# 473 UserMenu

用户菜单。

## 用途

头像按钮打开菜单：在线状态、设置与退出登录。头像角标以形状与颜色共同表示状态。

## 最小示例

```go
status := account.UserStatusOnline
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if account.UserMenu(c, &status, account.UserMenuOptions{Name: "Sakiko Togawa"}).SignOut() {
		fmt.Println("sign out")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `status` | `*UserStatus` | 必填 | UserStatusOnline、UserStatusAway、UserStatusBusy、UserStatusInvisible。 |
| `Name` | `string` | 必填 | 用户名。 |
| `Email` | `string` | "" | 邮箱。 |
| `Avatar` | `*ui.Bitmap` | nil | 头像。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

关闭、展开、四种状态、禁用。

## 事件

`StatusChanged()`、`Settings()`、`SignOut()`。

## 键盘操作

Enter/Space/↓ 打开，方向键移动，Enter 选择，Escape 关闭。

## 限制

- 按钮名称包含用户名与状态。
