# 471 AccountSwitcher

账户切换。

## 用途

显示当前账户，菜单中切换已登录账户或添加账户。

## 最小示例

```go
current := "a"
accounts := []account.Account{{ID: "a", Name: "Sakiko Togawa", Email: "sakiko@court"}, {ID: "b", Name: "Oblivionis"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if account.AccountSwitcher(c, &current, accounts, account.AccountSwitcherOptions{}).Add() {
		fmt.Println("sign in another account")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `current` | `*string` | 必填 | 当前账户 ID，须在列表中。 |
| `accounts` | `[]Account` | 必填 | ID、Name 必填；Email、Avatar 可选。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

关闭、展开、当前账户勾选、禁用。

## 事件

`Changed()` 切换；`Add()` 请求添加账户。

## 键盘操作

Enter/Space/↓ 打开菜单，方向键移动，Enter 选择，Escape 关闭。

## 限制

大集合使用虚拟列表，仅构建当前可见行。

- 菜单项只显示文字，头像只在按钮上。
