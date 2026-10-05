# 195 NewsFeed

## 用途

新闻流。

## 最小示例

```go
var state finance.NewsFeedState
rows := []finance.NewsItem{{ID: "a", Title: "Apple unveils new chips", Source: "Wire", Time: time.Date(2026, 10, 5, 14, 30, 0, 0, time.UTC), Symbols: []string{"AAPL"}}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := finance.NewsFeed(c, &state, rows, finance.NewsFeedOptions{})
	if r.Opened != "" || r.Symbol != "" || r.Retry {
		fmt.Println(r.Opened, r.Symbol)
	}
}
ui.Render(view, 960, 440, 1)
```

## 参数

`Version` 在调用方原地修改新闻或已读映射时递增；`ReadChanged` 报告已读状态变化。

`s.Read` 以 ID 保存已读；`items` 按调用方给定顺序显示。`Today` 默认窗口日期；`Status/Error/Disabled` 默认零值。

## 状态

默认、聚焦、禁用；异步结果支持加载、失败和空。

## 事件

if r.Opened != "" || r.Symbol != "" || r.Retry { 调用方回调(r.Opened, r.Symbol) }

## 键盘操作

Tab 聚焦列表或相关代码；上下选择；Enter 打开并标记已读。

## 限制

虚拟化构建可见行；读取状态以稳定 ID 保存；链接打开由调用方负责。
