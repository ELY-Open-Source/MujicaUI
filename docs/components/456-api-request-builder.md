# 456 ApiRequestBuilder · API 请求构建器

## 用途

编辑方法、URL、参数、请求头和正文，报告发送请求。

## 最小示例

```go
request := mujica.ApiRequest{Method: "POST", URL: "https://example.com/api", Headers: []mujica.KeyValue{{Key: "Accept", Value: "application/json", Enabled: true}}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.ApiRequestBuilder(c, &request, mujica.ApiRequestBuilderOptions{})
	if r.Send {
		fmt.Println("send", request.Method)
	}
	if r.Changed {
		fmt.Println("request changed")
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

Sending 防止重复发送；Disabled 禁用；Method 默认为 GET；URL 须为 HTTP 或 HTTPS；启用且键非空的参数进入预览。

## 状态

空地址、有效地址、无效地址、发送中、禁用。

## 事件

Send 表示请求发送；Changed 为编辑变化。

## 键盘操作

Tab 遍历控件；URL 框 Enter 发送；左右切换标签页。

## 限制

网络请求由调用方执行；支持 GET、POST、PUT、PATCH、DELETE、HEAD、OPTIONS。
