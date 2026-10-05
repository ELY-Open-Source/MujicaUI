# 256 SharedConversationView · 分享会话视图

## 用途

只读显示他人分享的会话，带分享者与时间，提供复制链接回调。

## 最小示例

```go
msgs := []mujica.SessionMessage{{Role: mujica.SessionUser, Text: "Hello"}, {Role: mujica.SessionAssistant, Text: "Good day."}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.SharedConversationView(c, msgs, mujica.SharedConversationViewOptions{Title: "Greeting", SharedBy: "Ann"}).Changed() {
		c.WriteClipboard("https://example.com/share/1")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func SharedConversationView(c *ui.Context, messages []SessionMessage, opts SharedConversationViewOptions) ChoiceResult`。

`SessionMessage`：`Role`（`SessionUser`、`SessionAssistant`、`SessionSystem`）、`Author`、`Text`、`Time`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 必填 | 会话标题。 |
| `SharedBy` | `string` | "" | 分享者，空显示“未知”。 |
| `SharedAt` | `time.Time` | 零值 | 分享时间，零值不显示。 |
| `Disabled` | `bool` | false | 禁用复制链接。 |

## 状态

标题区、只读标记（锁形图标）、用户消息（描边气泡）、助手消息、禁用。

## 事件

`Changed()`：点击“复制链接”，由调用方生成链接并写入剪贴板。

## 键盘操作

- Tab：复制链接按钮，再到消息列表；消息文字可选择复制。

## 限制

- 消息以纯文本显示，不渲染 Markdown（Markdown 渲染属于 MarkdownView 组件）。
- 消息列表虚拟化，只构建可见消息。
