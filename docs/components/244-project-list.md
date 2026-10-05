# 244 ProjectList · 项目列表

## 用途

按文件夹列出项目，展开项目可见其中的会话；把会话拖到项目上或用“移动”菜单即可改变归属。

## 最小示例

```go
var state chat.ProjectListState
projects := []chat.ChatProject{{ID: "p1", Name: "Coronation", Folder: "Court"}}
convs := func(id string) []chat.ChatConversation {
	return []chat.ChatConversation{{ID: "1", Title: "Menu"}}
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if v := chat.ProjectList(c, &state, projects, convs, chat.ProjectListOptions{}); v.Moved {
		log.Println("move", v.MoveID, "to", v.MoveTo)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ProjectList(c *ui.Context, s *ProjectListState, projects []ChatProject, conversations func(projectID string) []ChatConversation, opts ProjectListOptions) ProjectListView`。

`ChatProject`：`ID`、`Name`、`Folder`（文件夹名，空表示不在文件夹中）。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | "" | 树的无障碍名称。 |
| `Disabled` | `bool` | false | 禁用选择、展开与移动。 |

## 状态

文件夹（展开/收起）、项目行（悬停、选中：Selection 底 + 菱形）、拖入目标（角饰边框）、会话行（拖动时半透明）、禁用。

## 事件

`ProjectListView`：`Changed`（选中项目，ID 在 `State.Selected`）、`Opened`（打开的会话 ID）、`Moved` 与 `MoveID`/`MoveTo`（把会话移入项目，由调用方更新数据）。

## 键盘操作

- Tab：进入树；↑/↓ 在文件夹、项目、会话间移动。
- 项目行上 →/← 展开/收起会话；Enter 选择项目。
- 文件夹上 Enter/Space 展开或收起。
- 会话行的菜单按钮：选择“移入项目”下的目标完成移动（拖放的键盘替代）。

## 限制

- 项目行通过 `Description` 读出“已展开/已收起”；MyGo 不向自定义元素公开 expanded 状态（ui/element.go:282），文件夹使用 `CollapsibleBase`，其展开状态可被辅助技术读到。
- 列表未虚拟化，适合数百个项目以内。
