# 249 CapabilityCards · 能力卡片

## 用途

用卡片说明助手能做什么，每张卡片附示例提示，点击示例交给调用方。

## 最小示例

```go
cards := []chat.CapabilityCard{{Title: "Plan", Icon: icons.Must("calendar"), Examples: []string{"Plan a ball"}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if v := chat.CapabilityCards(c, cards, chat.CapabilityCardsOptions{}); v.Example != "" {
		log.Println("ask", v.Example)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func CapabilityCards(c *ui.Context, cards []CapabilityCard, opts CapabilityCardsOptions) CapabilityCardsView`。

`CapabilityCard`：`Title`（必填）、`Description`、`Icon`、`Examples`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Columns` | `int` | 0 → 3 | 每行卡片数；负数 panic。 |
| `Disabled` | `bool` | false | 禁用示例按钮。 |

## 状态

卡片（衬线标题、酒红图标）、示例按钮（悬停、按下、聚焦）、禁用。

## 事件

`Example`：点击的示例文字。

## 键盘操作

- Tab 在各卡片示例组间移动；组内 ↑/↓；Enter/Space 选择。

## 限制

- 卡片等宽分栏，不随窗口宽度自动减少列数；窄窗口时请调用方传入较小的 `Columns`。
