# 329 FindWidget · 查找替换

## 用途

在文本中查找与替换：区分大小写、整词、正则表达式，显示匹配数并上下跳转。

## 最小示例

```go
text := "the seal and the seal"
state := code.FindState{Query: "seal", Replacement: "signet", ShowReplace: true}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if code.FindWidget(c, &text, &state, code.FindWidgetOptions{}).Replaced > 0 {
		fmt.Println(text)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func FindWidget(c *ui.Context, text *string, s *FindState, opts FindWidgetOptions) FindWidgetResult`。匹配在文本、查询或选项变化时重算；`s.Matches()` 返回各匹配的字节范围。替换直接修改 `*text`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `FindState.Query`, `Replacement` | `string` | "" | 查找与替换文字；正则模式下替换支持 `$1` 分组。 |
| `MatchCase`, `WholeWord`, `Regex` | `bool` | false | 区分大小写、整词、正则。 |
| `ShowReplace` | `bool` | false | 展开替换行。 |
| `Current` | `int` | 0 | 当前匹配。 |
| `ReadOnly` | `bool` | false | 隐藏替换。 |
| `Disabled` | `bool` | false | 不能输入与替换。 |

## 状态

选项按钮按下状态、匹配计数、无匹配、正则错误（Danger 文字说明具体错误）、替换按钮在无匹配时禁用、禁用。

## 事件

`Moved`：查询、选项或当前匹配变化；`Replaced`：本帧替换的数量。

## 键盘操作

- Enter / Shift+Enter：下一项 / 上一项；macOS Cmd+G，其他平台 F3。
- 替换输入框中 Enter：替换当前项。
- Tab 依次到达展开按钮、查找框、上下按钮、三个选项与替换控件。

## 限制

- 使用 Go 的 RE2 正则语法，不支持回溯引用。
- 不在编辑器中绘制高亮，调用方用 `Matches()` 自行标记。
