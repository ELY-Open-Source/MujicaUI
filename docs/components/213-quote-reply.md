# 213 QuoteReply 引用回复

## 用途

显示被引用消息的作者与单行摘要（多余空白折叠，过长省略）；点击跳转到原消息。

## 最小示例

```go
if chat.QuoteReply(c, "Marie", "Let them eat cake.").Clicked() {
	fmt.Println("jump to the quoted message")
}
```

签名：`func QuoteReply(c *ui.Context, author, excerpt string) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `author` | `string` | 必填 | 原消息作者；为空 panic。 |
| `excerpt` | `string` | "" | 原消息摘要；换行与连续空白合并为一个空格。 |

## 状态

左侧 3 DIP AccentText 竖线标明引用；悬停时底色变为 Selection。

## 事件

`Clicked()`：点击或 Enter、Space；跳转由调用方完成（如 MessageList 的 ScrollTo）。

## 键盘操作

Tab 聚焦，Enter 或 Space 触发。辅助技术读到“回复 作者：摘要”。

## 限制

- 只显示一行摘要，不渲染原消息的格式。
