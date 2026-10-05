# 253 StopGeneratingButton · 停止生成按钮

## 用途

生成回复时显示的停止按钮，并提供 Cmd+.（其他平台 Ctrl+.）快捷键；生成结束后不显示。

## 最小示例

```go
generating := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.StopGeneratingButton(c, chat.StopGeneratingButtonOptions{Generating: generating}).Changed() {
		generating = false
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func StopGeneratingButton(c *ui.Context, opts StopGeneratingButtonOptions) ChoiceResult`。不显示时返回的 `Element` 为 nil。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Generating` | `bool` | false | 正在生成时显示。 |
| `Disabled` | `bool` | false | 禁用按钮与快捷键。 |

## 状态

默认（方块图标 + 文字 + 快捷键提示）、悬停、聚焦、禁用、隐藏。

## 事件

`Changed()`：请求停止。

## 键盘操作

- Enter/Space：停止。
- Cmd+.（macOS）/ Ctrl+.（Windows、Linux）：窗口级快捷键，焦点所在控件未占用该键时生效。

## 限制

- 快捷键用 MyGo 窗口快捷键注册，焦点控件若注册了同一组合键会优先获得。
