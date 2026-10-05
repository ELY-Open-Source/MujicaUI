# 201 ChatContainer 对话容器

## 用途

整页对话布局：左侧可选侧栏插槽，右侧为消息区与底部输入区。消息区就是 MessageList：新消息到来时自动滚到底部，用户向上翻看时停止跟随，并在底部显示回到最新按钮。

## 最小示例

```go
mujica.ChatContainer(c, &chat, len(msgs), mujica.ChatContainerOptions{
	List:  mujica.MessageListOptions{ID: func(i int) string { return msgs[i] }},
	Input: func() { ui.Text(c, "Composer goes here") },
}, func(i int) { ui.Text(c, msgs[i]) })
```

签名：`func ChatContainer(c *ui.Context, s *MessageListState, n int, opts ChatContainerOptions, row func(i int)) MessageListResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `s` | `*MessageListState` | 必填 | 消息区的位置与已读状态，由调用方持有。 |
| `n / row` | `int / func(i int)` | 必填 | 消息数量与第 i 条消息的构建函数；只构建可见消息。 |
| `List` | `MessageListOptions` | `ID` 必填 | 传给消息列表的选项，见 202。 |
| `Input` | `func()` | nil | 输入区插槽，位于消息区下方，带上边线。 |
| `Sidebar` | `func()` | nil | 侧栏插槽；nil 时不显示侧栏。 |
| `SidebarWidth` | `float32` | 260 | 侧栏宽度 DIP；负数 panic。 |

## 状态

消息区跟随底部 / 用户上翻后停止跟随；离开底部时出现带未读数的回到最新按钮。

## 事件

返回的 `MessageListResult.Unread()` 给出离开底部后新到的消息数。

## 键盘操作

Tab 在侧栏、消息与输入区之间移动；消息区获得焦点后滚轮与 PageUp/PageDown（Linux、Windows）滚动。

## 限制

- 侧栏宽度固定，不提供拖动调整；需要时用 SplitPane（008）包裹。
