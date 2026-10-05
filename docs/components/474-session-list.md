# 474 SessionList

登录会话。

## 用途

列出已登录设备、位置与最近活跃时间；可注销单个会话或经确认注销其他全部会话。

## 最小示例

```go
sessions := []account.Session{{ID: "s1", Device: "MacBook Pro", Current: true}, {ID: "s2", Device: "iPhone", Location: "Osaka"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if id, ok := account.SessionList(c, sessions, account.SessionListOptions{}).Revoked(); ok {
		fmt.Println("sign out", id)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `sessions` | `[]Session` | 必填 | ID、Device 必填；不包含会话令牌。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

本设备标记、最近活跃未知、确认对话框、禁用。

## 事件

`Revoked()` 返回会话 ID；`OthersRevoked()` 在确认后为 true。

## 键盘操作

Tab 到各注销按钮；确认框默认聚焦取消。

## 限制

大集合使用虚拟列表，仅构建当前可见行。

- 注销由调用方执行。
