# 243 ConversationSearch · 会话搜索

## 用途

按标题与内容搜索会话。查询交给调用方异步执行，结果回传后高亮匹配文字。

## 最小示例

```go
query := "seat"
hits := []chat.ConversationSearchHit{{ID: "1", Title: "Ballroom seating plan"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	v := chat.ConversationSearch(c, &query, hits, chat.ConversationSearchOptions{})
	if v.QueryChanged {
		log.Println("search", query)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ConversationSearch(c *ui.Context, query *string, hits []ConversationSearchHit, opts ConversationSearchOptions) ConversationSearchView`。

`ConversationSearchHit`：`ID`、`Title`、`Snippet`（匹配的消息片段）。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Loading` | `bool` | false | 调用方正在搜索。 |
| `Error` | `string` | "" | 搜索失败的说明，附重试按钮。 |
| `Placeholder` | `string` | 内置“搜索” | 输入框占位。 |
| `Label` | `string` | "" | 无障碍名称。 |
| `Disabled` | `bool` | false | 禁用输入与选择。 |

## 状态

空查询、搜索中（旋转标记）、有结果（匹配处加粗、下划线、Selection 底）、无结果、失败（错误文字 + 重试）、禁用。

## 事件

`QueryChanged`：查询变化，调用方据此发起搜索，结果到达后传入新的 `hits` 并调用 `c.Invalidate()`。`Chosen`：点击或 Enter 打开的会话 ID。`Retry`：点击重试。

## 键盘操作

- 输入框中 ↓：进入结果列表。
- 列表中 ↑/↓ 移动，Enter 打开。
- Escape：清空查询（SearchInput 行为）。

## 限制

- 匹配高亮为不区分大小写的逐字符匹配，不做分词、拼音或模糊匹配；排序完全由调用方决定。
- 结果行固定高度（标题一行、片段一行），片段超长时截断。
