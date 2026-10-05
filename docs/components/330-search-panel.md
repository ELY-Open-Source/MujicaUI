# 330 SearchPanel · 多文件搜索

## 用途

展示多文件搜索结果：按文件分组、可折叠，匹配预览带高亮；提供包含与排除模式输入。

## 最小示例

```go
results := []code.SearchFileResult{{Path: "court/seal.go", Matches: []code.SearchMatch{{Line: 12, Text: "func Seal() {", Start: 5, End: 9}}}}
state := code.SearchPanelState{Query: "seal"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := code.SearchPanel(c, results, &state, code.SearchPanelOptions{})
	r.Element.Height(280)
	if r.Opened != nil {
		fmt.Println("open", r.Opened.Path, r.Opened.Line)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func SearchPanel(c *ui.Context, results []SearchFileResult, s *SearchPanelState, opts SearchPanelOptions) SearchPanelResult`。搜索由调用方执行：查询或模式变化时 `Changed` 为 true，完成后传入新的 `results`。行列表虚拟化。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `SearchFileResult.Path` | `string` | 必填 | 文件路径，也是折叠键。 |
| `SearchMatch.Line`, `Text` | `int`, `string` | — | 行号与行文字。 |
| `SearchMatch.Start`, `End` | `int` | — | 匹配在 `Text` 中的字节范围；越界 panic。 |
| `SearchPanelState.Query`, `Include`, `Exclude` | `string` | "" | 查询与包含、排除模式（格式由调用方解释）。 |
| `Status`, `Error` | `DataStatus`, `string` | `DataReady` | 搜索中或失败。 |
| `Disabled` | `bool` | false | 不能输入、折叠与打开。 |

## 状态

结果摘要、文件标题（折叠箭头与匹配数）、匹配行（Selection 背景与粗体下划线标出匹配）、选中行（Selection 与菱形）、搜索中、失败（可重试）、未输入、无结果、禁用。

## 事件

`Changed`：查询或模式变化，或点击重试；`Opened`：用户打开的匹配（路径与行号）。

## 键盘操作

- Tab 依次到达查询、包含、排除与结果列表。
- ↑/↓ 选择；Enter 或双击：打开匹配，或折叠/展开文件；文件行上 ← 折叠、→ 展开。

## 限制

- 长行只显示匹配前 40 个字节。
