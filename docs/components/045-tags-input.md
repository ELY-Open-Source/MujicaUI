# 045 TagsInput

标签输入。

## 用途

编辑一组不重复的标签，对应 MyGo 的 Token field。支持 Enter 或分隔符添加、整段粘贴拆分、去重、数量上限与添加前校验。

## 最小示例

```go
tags := []string{"gothic"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.TagsInput(c, &tags, input.TagsInputOptions{Max: 5, Placeholder: "Add a tag", Label: "Tags"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `tags` | `*[]string` | 必填 | 标签列表。 |
| `Max` | `int` | 0（不限） | 最多标签数；负数 panic。 |
| `Separators` | `string` | ",;，；" 与换行、制表符 | 输入或粘贴时用于拆分的字符。 |
| `Validate` | `func(string) string` | nil | 添加前校验，返回错误文字或 ""。 |
| `Placeholder` | `string` | "" | 占位文字。 |
| `Label` | `string` | "" | 不在 FormField 内时为输入框命名。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Invalid` | `bool` | false | 边框显示危险色。 |

## 状态

默认、聚焦、错误（重复、超出上限、校验失败时显示具体原因并保留输入文字）、禁用。

## 事件

`TagsInputParts.Changed()` 添加或移除了标签；`Input` 为输入框元素。

## 键盘操作

Enter 或分隔符添加；输入框为空时 Backspace 移除最后一个标签；Tab 可到达每个标签的移除按钮，Enter/Space 移除。

## 限制

- 粘贴含分隔符的文本时，未通过检查的片段保留在输入框中，供用户修改。
- MyGo 单行输入框会处理粘贴文本中的换行，换行能否作为分隔符取决于平台粘贴行为；逗号、分号总是有效。
