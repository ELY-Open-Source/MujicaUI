# 076 Timeline · 时间线

## 用途

纵向列出事件：时间、标题、内容与状态。标记之间以细线连接，长内容自动换行。状态同时用形状、图标和文字表示，不只依赖颜色。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.Timeline(c, []data.TimelineItem{
		{Time: "06:00", Title: "Matins sung", Status: data.TimelineDone},
		{Time: "09:00", Title: "Relics moved", Content: "From the crypt to the altar.", Status: data.TimelineActive},
		{Time: "18:00", Title: "Vespers"},
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Timeline(c *ui.Context, items []TimelineItem) ui.Element`。`TimelineItem` 字段：`Time`、`Title`、`Content`、`Status`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Status` | `TimelineStatus` | `TimelineDefault` | `TimelineDefault`（空心菱形）、`TimelineDone`（勾选图标 + “已完成”）、`TimelineActive`（实心菱形 + “进行中”）、`TimelineWarning`（警告图标 + “警告”）、`TimelineError`（错误图标 + “错误”）。 |

## 状态

仅展示：默认、完成、进行中、警告、错误。

## 事件

无。

## 键盘操作

不获取焦点。

## 限制

- 仅纵向排列；内容为纯文本。
- 不虚拟化，适合几十条以内的事件。
