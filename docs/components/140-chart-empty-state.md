# 140 ChartEmptyState 图表占位

## 用途

在图表无数据、加载或失败时保持图表尺寸，并在失败时提供重试。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chart.ChartEmptyState(c, data.DataFailed, chart.ChartEmptyStateOptions{Height: 180, Error: "The archive did not answer"})
}
ui.Render(view, 400, 240, 1)
```

## 参数

`status DataStatus`：DataReady 显示空数据，DataLoading 显示加载，DataFailed 显示错误和重试；其他有效加载状态同样显示加载。`Height` 默认 240 DIP，负值 panic；`Error`、`Empty` 默认内置本地化文案。

## 状态

空数据、加载、错误；祖先禁用保持展示并阻止重试；减少动画显示静态加载说明。

## 事件

返回 `ChoiceResult`，`Changed()` 表示点击重试，数据请求由调用方执行。

## 键盘操作

Tab 聚焦重试；Enter 或 Space 触发一次重试。

## 限制

只管理占位展示和事件；数据和异步状态由调用方维护。
