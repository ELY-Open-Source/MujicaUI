# 242 ConversationItem · 会话条目

## 用途

一行会话：标题、时间、置顶标记，以及菜单中的重命名（行内编辑）、置顶、导出、删除。删除与导出只回调，由调用方执行。

## 最小示例

```go
conv := mujica.ChatConversation{ID: "1", Title: "Court etiquette", Updated: time.Now()}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	v := mujica.ConversationItem(c, &conv, mujica.ConversationItemOptions{})
	if v.Deleted {
		log.Println("delete", conv.ID)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ConversationItem(c *ui.Context, conv *ChatConversation, opts ConversationItemOptions) ConversationItemView`。`conv` 不能为 nil，`ID` 不能为空。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Today` | `Date` | 零值（窗口时钟） | 决定时间显示为时刻、“昨天”或日期。 |
| `Selected` | `bool` | false | 显示为当前会话。 |
| `Disabled` | `bool` | false | 禁用全部操作。 |
| `Draggable` | `bool` | false | 允许拖到 `ProjectList` 的项目上，携带 `ConversationDrag`。 |

## 状态

默认、悬停、选中（Selection 底色 + 酒红标题）、置顶（星形图标）、行内编辑（输入框酒红边框）、禁用。

## 事件

`ConversationItemView`：`Opened`（点击或 Enter）、`Renamed`（新标题已写入 `conv.Title`）、`PinToggled`、`Exported`、`Deleted`。

## 键盘操作

- Enter/Space：打开会话。
- F2：行内重命名；输入框中 Enter 确认，Escape 取消；空白标题不会被接受。
- Delete：请求删除。
- 菜单按钮：Enter、Space 或 ↓ 打开，方向键移动，Escape 关闭。

## 限制

- 删除没有内置确认；需要确认时由调用方弹出 AlertDialog。
- 菜单按钮不在条目的 Tab 位之外另设快捷键。
