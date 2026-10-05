# 246 WelcomeScreen · 欢迎页

## 用途

新会话的欢迎页：衬线问候、菱形分隔、副标题、能力卡片，以及放置起始提示的插槽。

## 最小示例

```go
cards := []chat.CapabilityCard{{Title: "Write", Examples: []string{"Draft a decree"}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if v := chat.WelcomeScreen(c, chat.WelcomeScreenOptions{Greeting: "Good evening", Cards: cards}); v.Example != "" {
		log.Println("ask", v.Example)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func WelcomeScreen(c *ui.Context, opts WelcomeScreenOptions) CapabilityCardsView`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Greeting` | `string` | 必填 | 问候语，空时 panic。 |
| `Subtitle` | `string` | "" | 副标题。 |
| `Cards` | `[]CapabilityCard` | nil | 能力卡片，最多三列。 |
| `Prompts` | `func()` | nil | 起始提示插槽，例如 `SuggestionChips`。 |
| `Disabled` | `bool` | false | 禁用卡片示例与插槽内容。 |

## 状态

默认、禁用（由 MyGo 置灰）。卡片与示例按钮的状态见 249 CapabilityCards。

## 事件

返回值的 `Example`：点击的示例提示。插槽内组件的事件由调用方自行读取。

## 键盘操作

- Tab 依次经过各卡片的示例与插槽内控件；卡片内 ↑/↓ 在示例间移动；Enter/Space 选择。

## 限制

- 问候语按纯文本显示，不解析 Markdown。
