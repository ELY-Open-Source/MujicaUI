# 214 MessageActions 消息操作

## 用途

消息下方的操作工具栏：复制、重新生成、编辑、赞、踩、朗读、分享，按角色显示：用户消息为复制、编辑；助手消息为复制、重新生成、赞、踩、朗读、分享；系统消息为复制；错误消息为复制、重新生成。

## 最小示例

```go
if a, ok := mujica.MessageActions(c, mujica.MessageAssistant, mujica.MessageActionsOptions{Feedback: &rating}).Action(); ok {
	fmt.Println("action", a)
}
```

签名：`func MessageActions(c *ui.Context, role MessageRole, opts MessageActionsOptions) MessageActionsResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `role` | `MessageRole` | 必填 | 决定显示哪些操作。 |
| `Feedback` | `*MessageFeedback` | nil | 赞踩状态；nil 时不显示赞与踩。再次点击已选的赞或踩会取消。 |

## 状态

赞或踩选中时为按下状态：Selection 底色、AccentText 边框与图标；辅助技术读到切换按钮的按下状态。

## 事件

`Action()` 返回本帧的操作：`ActionCopy`、`ActionRegenerate`、`ActionEdit`、`ActionLike`、`ActionDislike`、`ActionSpeak`、`ActionShare`。复制、朗读与分享由调用方执行。

## 键盘操作

整个工具栏是一个 Tab 停靠点，Left/Right 在操作间移动，Home、End 到首末；Enter 或 Space 触发。

## 限制

- 不内置复制内容与朗读：调用方收到操作后自行处理。
