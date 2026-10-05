# 241 ConversationList · 会话列表

## 用途

按“置顶、今天、昨天、7 天内、更早”分组显示历史会话，只构建可见行，适合十万条以上的会话历史。

## 最小示例

```go
var state chat.ConversationListState
convs := []chat.ChatConversation{
	{ID: "1", Title: "Court etiquette", Updated: time.Now()},
	{ID: "2", Title: "Royal letters", Updated: time.Now().AddDate(0, 0, -30), Pinned: true},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.ConversationList(c, &state, convs, chat.ConversationListOptions{Label: "History"}, nil).Changed() {
		log.Println("open", state.Selected)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ConversationList(c *ui.Context, s *ConversationListState, items []ChatConversation, opts ConversationListOptions, row func(conv ChatConversation)) ConversationListView`。`row` 为 nil 时显示标题；传入时可放 `ConversationItem` 等自定义内容。

`ChatConversation`：`ID`（唯一）、`Title`、`Updated`、`Pinned`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Now` | `time.Time` | 零值（窗口时钟） | 分组用的当前时间。 |
| `Label` | `string` | "" | 列表的无障碍名称。 |
| `Disabled` | `bool` | false | 禁用：不能选择，退出 Tab 序列。 |

`ConversationListState.Selected` 为当前会话 ID，调用方可直接设置。

## 状态

分组标题（衬线小字、旧金色）、普通行、悬停、选中（Selection 底色 + 左侧菱形）、聚焦（列表焦点环）、禁用。

## 事件

`Changed()`：选中了另一会话，ID 在 `State.Selected`。`Submitted()`：在选中行上按 Enter 或双击。

## 键盘操作

- Tab：聚焦列表。
- ↑/↓、Home、End：在会话间移动并选择，跳过分组标题；输入首字母定位。
- Enter：打开当前会话。

## 限制

- 分组与排序只在 `items` 切片（首元素地址或长度）或当天日期变化时重算。原地修改切片元素后，请传入新切片以触发重新分组。
- 选中行由 MyGo 列表报告选中状态；分组标题由 `ListState.Header` 标记，不可选中。
- 自定义 `row` 中若含按钮（如 `ConversationItem`），点击由按钮处理，不会改变列表选择；请在按钮事件里设置 `State.Selected`。
