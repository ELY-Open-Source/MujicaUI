# 402 ChatMessage 聊天消息

## 用途

展示连续聊天消息、发送时间、编辑标记，以及调用方提供的回复和回应插槽。

## 最小示例

```go
prev := mujica.ChatEntry{ID: "1", Author: "Uika", Time: time.Now(), Text: "Rehearsal at nine."}
m := mujica.ChatEntry{ID: "2", Author: "Uika", Time: time.Now(), Text: "Bring the masks.", Edited: true}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ChatMessage(c, prev, mujica.ChatMessageOptions{})
	mujica.ChatMessage(c, m, mujica.ChatMessageOptions{
		Continued: mujica.ChatContinues(prev, m),
		Reactions: func() { mujica.Tag(c, "👍 2", mujica.TagOptions{}) },
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`ChatEntry` 包含 `ID`、必填的 `Author`、`Image`、`Time`、`Text`、`Edited`、`Own`。`ChatMessageOptions.Continued` 默认 false；`Reply` 与 `Reactions` 为可选插槽。`ChatContinues(prev, cur)` 要求同一作者、相同 Own、间隔非负且小于五分钟。

## 状态

状态由调用方数据控制；浅深主题、密度与文字缩放随窗口设置变化。适用插槽操作沿用控件焦点和禁用状态。

## 事件

插槽中的按钮、链接和回应事件由调用方处理。

## 键盘操作

可选文字支持选择复制；插槽控件按 Tab 访问。

## 限制

调用方计算 Continued 或使用 ChatContinues；回复与回应的业务逻辑在插槽中实现。
