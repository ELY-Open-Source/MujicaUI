# 413 PinnedMessages 置顶消息

## 用途

展示虚拟化置顶消息列表，发出跳转和取消置顶事件。

## 最小示例

```go
pins := []messaging.ChatEntry{{ID: "a", Author: "Sakiko", Text: "Rehearsal moves to nine."}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := messaging.PinnedMessages(c, pins, messaging.PinnedMessagesOptions{})
	if r.Unpinned != "" {
		pins = pins[:0] // unpin it
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`pins []ChatEntry` 的 ID 与 Author 必填。`Height` 默认 320 DIP，必须 ≥ 0；`Disabled` 默认 false。返回 `PinnedMessagesResult`，包含命名字段 `Element`、`Jump`、`Unpinned`。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

`Jump` 和 `Unpinned` 返回本帧操作的消息 ID；调用方更新消息定位和置顶数据。

## 键盘操作

Tab 访问跳转与取消置顶；方向键移动；Enter / Space 激活。

## 限制

调用方维护置顶消息；列表行只构建可见部分。
