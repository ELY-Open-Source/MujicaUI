# 202 MessageList 消息列表

## 用途

虚拟化的消息列表，基于 MyGo `ui.List`：只构建可见消息，以稳定消息 ID 保持位置；从底部开始并跟随新消息，用户上翻后停止跟随，并显示新消息提示；可按日期插入分隔。

## 最小示例

```go
r := mujica.MessageList(c, &list, len(ids), mujica.MessageListOptions{
	ID: func(i int) string { return ids[i] },
}, func(i int) { ui.Text(c, "Message "+ids[i]) })
fmt.Println("unread", r.Unread())
```

签名：`func MessageList(c *ui.Context, s *MessageListState, n int, opts MessageListOptions, row func(i int)) MessageListResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `ID` | `func(i int) string` | 必填 | 消息的稳定 ID；在上方插入旧消息时视野中的消息保持不动。ID 须唯一。 |
| `Date` | `func(i int) time.Time` | nil | 设置后在每个新日期的第一条消息前显示“今天 / 昨天 / 日期”分隔。 |
| `Label` | `func(i int) string` | nil | 消息的无障碍名称，也用于键入定位。 |
| `n` | `int` | 必填 | 消息数；负数 panic。 |

## 状态

跟随底部；离开底部时底部居中出现 ScrollToBottomButton（含未读数）；点击后回到底部并清零。

## 事件

`Unread()`：离开底部后新增的消息数。`MessageListState.ScrollToEnd()` 与 `AtEnd()` 供调用方控制。

## 键盘操作

列表获得焦点后，滚轮、方向键与 PageUp/PageDown 滚动；回到最新按钮可用 Tab 聚焦，Enter 或 Space 触发。

## 限制

- 不检查 ID 是否重复（逐帧检查 10 万条代价过高）；重复 ID 会让两条消息共享状态。
- 消息高度由 MyGo 在显示时测量，未显示消息按已测高度估算，长距离跳转时滚动条可能微调。
