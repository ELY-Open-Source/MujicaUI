# 238 SourceCard 来源卡片

## 用途

来源标题、站点、摘要、时间；打开回调。

## 最小示例

```go
var cites chat.CitationState
src := chat.Source{Title: "Hall of Mirrors", URL: "https://example.org/hall", Snippet: "Seventeen arches."}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.SourceCard(c, 1, src, &cites).Clicked() {
		fmt.Println("open", src.URL)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`n` 为引用编号；`src.Title`、`src.URL` 必填，`Site`、`Snippet`、`Published` 可选；`st` 与正文引用共享。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

Element.Clicked() 表示打开来源请求。

## 键盘操作

Tab 到卡片；Enter、Space 请求打开。

## 限制

打开 URL 由调用方处理；组件返回事件。
