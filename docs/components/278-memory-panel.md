# 278 MemoryPanel 记忆面板

## 用途

列出 Agent 的记忆条目，支持搜索、行内编辑和确认后删除。

## 最小示例

```go
entries := []agent.MemoryEntry{{ID: "1", Text: "Lady Sakiko prefers the east wing.", Source: "Chat"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if r := agent.MemoryPanel(c, &entries, agent.MemoryPanelOptions{}); r.Deleted != "" {
		fmt.Println("forget", r.Deleted)
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func MemoryPanel(c *ui.Context, entries *[]MemoryEntry, opts MemoryPanelOptions) MemoryPanelResult`。`entries` 为 nil、条目 `ID` 为空或 `MaxHeight` 为负时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `MaxHeight` | `float32` | 360 | 列表最大高度。 |
| `Disabled` | `bool` | `false` | 只读，禁止搜索、编辑与删除。 |

## 状态

- `MemoryEntry` 字段：`ID`、`Text`、`Source`。
- 搜索不区分大小写，匹配内容与来源；只在搜索词或条目变化时重新筛选。
- 编辑时出现文本框和保存、取消；删除前用 Popconfirm 确认，确认按钮为危险样式且默认聚焦取消。
- 条目超过 200 条时虚拟化。

## 事件

- `Edited`：本帧保存的条目 ID，`*entries` 已更新。
- `Deleted`：本帧删除的条目 ID，已从 `*entries` 移除。

## 键盘操作

- Tab 到搜索框、各条的编辑与删除；删除确认中 Escape 取消。

## 限制

- 调用方原地修改条目文本而不改变切片长度时，筛选结果到下次搜索变化才更新。
