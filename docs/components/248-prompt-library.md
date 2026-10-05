# 248 PromptLibrary · 提示模板库

## 用途

按分类与搜索浏览提示模板，预览模板，填写 `{{变量}}` 后插入。

## 最小示例

```go
var state chat.PromptLibraryState
templates := []chat.PromptTemplate{{ID: "1", Category: "Letters", Title: "Invitation", Body: "Invite {{guest}}."}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if v := chat.PromptLibrary(c, &state, templates, chat.PromptLibraryOptions{}); v.Inserted != "" {
		log.Println("insert", v.Inserted)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func PromptLibrary(c *ui.Context, s *PromptLibraryState, templates []PromptTemplate, opts PromptLibraryOptions) PromptLibraryView`。另有 `PromptLibraryVariables(body)` 与 `PromptLibraryFill(body, values)`。

`PromptTemplate`：`ID`、`Category`、`Title`（ID 与 Title 必填）、`Body`（`{{name}}` 为变量，名称可含字母、数字、`_.-`）。

`PromptLibraryState`：`Category`（空为全部）、`Query`、`Selected`、`Values`（变量值）。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | "" | 无障碍名称。 |
| `Disabled` | `bool` | false | 禁用全部操作。 |

## 状态

分类与模板行（选中：Selection 底 + 菱形）、无匹配、未选模板、预览（未填变量以 [名称] 标出，已填值酒红加粗）、插入按钮（变量未填全时禁用）、禁用。

## 事件

`Inserted`：变量已填好的模板正文，非空时由调用方插入输入框。

## 键盘操作

- Tab 依次进入分类列、搜索框、模板列、变量输入框、插入按钮；列内 ↑/↓ 移动；Enter/Space 选择。

## 限制

- 过滤只在模板切片、查询或分类变化时重算。
- 变量为简单占位替换，不支持条件或循环。
