# 203 MessageBubble 消息气泡

## 用途

按角色框住一条消息：用户靠右、Selection 底色；助手靠左、Surface 底与细边；系统居中、灰色小字；错误带 Danger 边框与警示图标。头像与名称为插槽，宽度可限制。

## 最小示例

```go
chat.MessageBubble(c, chat.MessageAssistant, chat.MessageBubbleOptions{Name: "Herald"}, func() {
	ui.Text(c, "The ball begins at nine.")
})
```

签名：`func MessageBubble(c *ui.Context, role MessageRole, opts MessageBubbleOptions, content func()) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `role` | `MessageRole` | 必填 | `MessageUser`、`MessageAssistant`、`MessageSystem`、`MessageError`。 |
| `Name` | `string` | "" | 显示在气泡上方的名称。 |
| `Avatar` | `func()` | nil | 头像插槽，用户在右侧、其他角色在左侧；系统消息不显示。 |
| `MaxWidth` | `float32` | 560 | 气泡最大宽度 DIP；负数 panic。 |
| `content` | `func()` | 必填 | 气泡内容，如 MarkdownView 或 StreamingText。 |

## 状态

四种角色外观；错误角色同时使用图标、边框与文字，不只依赖颜色。

## 事件

无；交互由内容与 MessageActions 提供。

## 键盘操作

气泡本身不获取焦点；内容中的链接、按钮按 Tab 顺序访问。辅助技术读到“角色，名称”分组。

## 限制

- 气泡不自带操作栏；把 MessageActions 放在 content 中或气泡下方。
