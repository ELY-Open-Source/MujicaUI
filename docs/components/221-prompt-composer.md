# 221 PromptComposer 提示编辑器

## 用途

多行输入；Enter 发送、Shift+Enter 换行可配置；输入法组合期间不发送；附件与工具插槽。

## 最小示例

```go
var prompt string
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := chat.PromptComposer(c, &prompt, chat.PromptComposerOptions{Placeholder: "Ask the herald"})
	if r.Submitted() {
		fmt.Println("send", prompt)
		prompt = ""
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`value` 为调用方持有的文本指针。`Placeholder`、`Label` 配置占位与名称；`CmdEnterSends=false` 使用 Enter 发送；`Composing=false`，调用方在组合输入时设为 true；`MaxHeight=0` 使用 200 DIP。`Attachments`、`Tools`、`Actions` 为内容插槽；`Disabled=false`。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

InputResult.Changed() 表示文本编辑；Submitted() 表示发送请求；Input 暴露编辑器。

## 键盘操作

Enter 发送，Shift+Enter 换行；CmdEnterSends 下使用 Command/Ctrl+Enter；组合保护期间 Enter 留给输入法。

## 限制

MyGo v0.2.10 的 ui/input.go:688–708 将组合事件直接送到私有编辑器，公开 HandleInput 无法观察组合状态。自动检测组合输入仍需上游 API；调用方须设置 Composing。
