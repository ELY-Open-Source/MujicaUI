# 260 ChatEmptyState · 会话空状态

## 用途

会话界面的三种空状态：新会话、无搜索结果、离线，并给出下一步操作。

## 最小示例

```go
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.ChatEmptyState(c, mujica.ChatEmptyStateOptions{Kind: mujica.ChatEmptyOffline}).Action {
		log.Println("reconnect")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ChatEmptyState(c *ui.Context, opts ChatEmptyStateOptions) EmptyStateResult`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Kind` | `ChatEmptyKind` | `ChatEmptyNew` | `ChatEmptyNew`、`ChatEmptyNoResults`、`ChatEmptyOffline`；越界 panic。 |
| `Query` | `string` | "" | 无结果时在说明中显示的查询。 |
| `Title`, `Description` | `string` | 内置 | 替换文案。 |
| `Action` | `string` | 内置 | 按钮文字：新建对话、清除搜索、重试。 |

## 状态

新会话（星形图标）、无结果（搜索图标）、离线（断网图标）。

## 事件

返回值 `Action`：点击按钮。

## 键盘操作

- Tab 聚焦按钮；Enter/Space 触发。

## 限制

- 基于 EmptyState 组件，布局与其一致。
