# 236 SendButton 发送按钮

## 用途

发送与停止两态；禁用条件；快捷键提示。

## 最小示例

```go
streaming := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := chat.SendButton(c, chat.SendButtonOptions{Streaming: streaming, Shortcut: "Enter"})
	if r.Sent() || r.Stopped() {
		streaming = r.Sent()
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`Streaming=false` 显示发送；true 显示停止；`Disabled=false`，禁用发送时停止仍可用；`Shortcut` 仅展示在提示中。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

Sent() 表示发送请求；Stopped() 表示停止请求。

## 键盘操作

Tab 到按钮；Enter、Space 激活；快捷键提示由调用方绑定。

## 限制

生成和停止逻辑由调用方实现。
