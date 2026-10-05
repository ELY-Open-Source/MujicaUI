# 235 ContextWindowMeter 上下文容量

## 用途

上下文使用率条；分段（系统、历史、附件）；阈值警告。

## 最小示例

```go
segments := []chat.ContextSegment{{Label: "System", Tokens: 1200}, {Label: "History", Tokens: 58000}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chat.ContextWindowMeter(c, segments, 128000, chat.ContextWindowMeterOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`segments` 每项 `Label` 必填、`Tokens` 非负；`limit` 为正数；`Warn=0` 默认阈值 0.8，允许配置 0–1。

## 状态

展示由输入数据驱动；边界值与警告状态带文字说明。

## 事件

展示组件，调用方更新各段用量。

## 键盘操作

展示组件，沿用外层 Tab 顺序。

## 限制

上下文数据由调用方统计，分段同时提供文字与图形标记。
