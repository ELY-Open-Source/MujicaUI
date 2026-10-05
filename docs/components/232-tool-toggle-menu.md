# 232 ToolToggleMenu 工具开关菜单

## 用途

联网、代码执行等工具开关；需授权项标记。

## 最小示例

```go
web := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	tools := []mujica.ToolItem{{ID: "web", Label: "Web search", On: &web}}
	if id, ok := mujica.ToolToggleMenu(c, tools).Changed(); ok {
		fmt.Println("toggled", id)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`tools` 每项 `ID`、`On` 必填；`Label`、`Description` 由调用方提供；`NeedsAuth=false`；需要授权时由 `Authorized` 控制开关可用性。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

Changed() 返回切换工具 ID；Authorize() 返回授权请求 ID。

## 键盘操作

Tab 到菜单按钮；Enter/Space 打开；Tab 到开关与授权；Escape 关闭。

## 限制

授权与工具执行由调用方完成。
