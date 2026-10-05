# 316 ReviewComment 评审评论

## 用途

行级评审讨论：显示文件与行号、评论列表，支持回复、解决与重新打开。

## 最小示例

```go
thread := mujica.ReviewThread{File: "main.go", Line: 12, Messages: []mujica.ReviewMessage{
	{Author: "Uika", When: time.Now(), Body: "Can this return early?"},
}}
draft := ""
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.ReviewComment(c, thread, &draft, mujica.ReviewCommentOptions{})
	if text, ok := r.Replied(); ok {
		thread.Messages = append(thread.Messages, mujica.ReviewMessage{Author: "Sakiko", When: time.Now(), Body: text})
	}
	if r.Resolved() || r.Reopened() {
		thread.Resolved = !thread.Resolved
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func ReviewComment(c *ui.Context, thread ReviewThread, draft *string, opts ReviewCommentOptions) ReviewCommentResult`

## 参数

`ReviewThread`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| File | string | 必填 | 文件路径。 |
| Line | int | 必填 | 行号，须为正。 |
| Resolved | bool | false | 已解决：隐藏回复框，显示“重新打开”。 |
| Messages | []ReviewMessage | 必填 | 至少一条评论（作者、时间、正文）。 |

`ReviewCommentOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| TimeFormat | string | "2006-01-02 15:04" | 时间格式。 |
| Disabled | bool | false | 禁用回复、解决与重新打开。 |

`draft` 保存回复草稿，发送后清空。

## 状态

- 已解决：标题栏显示 Success 色图标与“已解决”。
- 回复按钮在草稿为空时禁用。

## 事件

- `Replied()`：发送的回复文本（去除首尾空白）。
- `Resolved()`、`Reopened()`：对应按钮。

## 键盘操作

- Tab 经过回复框与按钮；回复框内 Cmd+Enter（非 macOS 为 Ctrl+Enter）发送。

## 限制

- 评论正文为纯文本，不渲染 Markdown。
- 线程状态由调用方更新。
