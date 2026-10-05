# 227 SlashCommandMenu 斜杠命令菜单

## 用途

输入“/”弹出命令；筛选；方向键选择；Tab 补全。

## 最小示例

```go
var prompt string
commands := []chat.SlashCommand{{Name: "summarize", Description: "Summarize the thread"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	composer := chat.PromptComposer(c, &prompt, chat.PromptComposerOptions{})
	if cmd, ok := chat.SlashCommandMenu(c, composer, &prompt, commands).Chosen(); ok {
		fmt.Println("run", cmd.Name)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`composer` 为本帧构建的 PromptComposer 结果；`value` 为同一文本指针；`commands` 每项 `Name` 为必填单词，`Description` 为调用方说明。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

Chosen() 返回选中的命令；Tab 补全文本，Enter 选择命令。

## 键盘操作

编辑器内 Up/Down 选择，Enter 选中，Tab 补全，Escape 关闭。

## 限制

菜单按可见行虚拟化；筛选在查询或命令变化时更新。
