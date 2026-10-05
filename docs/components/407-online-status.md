# 407 OnlineStatus 在线状态

## 用途

用圆形、月牙、横杠圆形和空心圆配合文字表示四种在线状态。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.OnlineStatus(c, messaging.PresenceAway, messaging.OnlineStatusOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`p PresenceStatus` 的零值为离线，支持 `PresenceOnline`、`PresenceAway`、`PresenceBusy`、`PresenceOffline`。`HideText` 默认 false；隐藏文字时保留名称和提示。`Size` 默认 10 DIP，必须 ≥ 0。

## 状态

状态由调用方数据控制；浅深主题、密度与文字缩放随窗口设置变化。适用插槽操作沿用控件焦点和禁用状态。

## 事件

纯展示，在线状态由调用方更新。

## 键盘操作

展示组件保持自然焦点序列。

## 限制

在线状态由调用方提供；隐藏文字仍保留状态形状和无障碍名称。
