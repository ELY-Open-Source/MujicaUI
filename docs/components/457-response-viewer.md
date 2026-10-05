# 457 ResponseViewer · 响应查看器

## 用途

显示状态码、耗时、大小、JSON 或文本正文与响应头。

## 最小示例

```go
response := &mujica.ApiResponse{Status: 200, Duration: 182 * time.Millisecond, Size: 32, Body: `{"ready":true}`, Headers: []mujica.KeyValue{{Key: "Content-Type", Value: "application/json"}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.ResponseViewer(c, response, mujica.ResponseViewerOptions{}).Changed() {
		fmt.Println("retry")
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

resp 为 nil 时显示等待；Status、Error 控制请求状态；Height 默认 320 DIP；HTTP 状态码须在 100 到 599；Size 为非负字节数。

## 状态

等待请求、加载、失败、JSON 正文、文本正文、空正文、空响应头、父级禁用。

## 事件

Changed() 表示重试。

## 键盘操作

Tab 到标签页；左右切换正文与响应头；JSON 树沿用查看器键盘操作。

## 限制

响应来自调用方；正文为有效 JSON 时使用树视图；其余正文使用文本。
