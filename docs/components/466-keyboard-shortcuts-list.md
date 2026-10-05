# 466 KeyboardShortcutsList

快捷键列表。

## 用途

分组显示快捷键表，可按动作名或按键搜索。

## 最小示例

```go
query := ""
groups := []account.ShortcutGroup{{Title: "General", Bindings: []account.ShortcutBinding{
	{Action: "Command palette", Shortcut: account.Shortcut{Mods: ui.Cmd | ui.Shift, Key: ui.KeyP}},
}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	account.KeyboardShortcutsList(c, &query, groups, account.KeyboardShortcutsListOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `query` | `*string` | 必填 | 搜索文字。 |
| `groups` | `[]ShortcutGroup` | 必填 | Title 与 Bindings。 |
| `Label` | `string` | 内置“Search shortcuts” | 搜索框名称。 |

## 状态

全部、筛选、无结果。

## 事件

无。

## 键盘操作

Tab 进入搜索框；每行读出动作与按键。

## 限制

大集合使用虚拟列表，构建当前可见行；搜索与分组在数据或查询变化时更新。

- 按键按平台显示，macOS 使用符号。
