# 237 CitationBadge 引用编号

## 用途

行内引用编号；悬停预览来源；点击定位来源。

## 最小示例

```go
var cites chat.CitationState
src := chat.Source{Title: "The Palace of Versailles", URL: "https://example.com/versailles"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if chat.CitationBadge(c, 1, src, &cites).Changed() {
		fmt.Println("show source", cites.Active)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`n` 为从 1 开始的引用编号；`src.Title`、`src.URL` 必填；`st` 为与来源卡片共享的 CitationState；`Active=0` 表示尚未选中。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

每次点击或键盘激活触发 ChoiceResult.Changed() 并定位来源；重复激活同一引用也会重新定位。

## 键盘操作

Tab 到引用；Enter、Space 定位；悬停或键盘焦点 500 毫秒后预览；Escape 关闭预览。

## 限制

编号按调用方引用顺序映射；预览内容来自调用方。
