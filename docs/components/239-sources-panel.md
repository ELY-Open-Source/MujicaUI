# 239 SourcesPanel 来源面板

## 用途

来源列表；按引用顺序；与正文引用联动高亮。

## 最小示例

```go
var cites mujica.CitationState
sources := []mujica.Source{{Title: "Hall of Mirrors", URL: "https://example.org/hall"}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if n, ok := mujica.SourcesPanel(c, sources, &cites).Opened(); ok {
		fmt.Println("open", sources[n-1].URL)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`sources` 按引用顺序提供；`st.Active` 为 0 或列表内的引用编号；外层容器提供高度。

## 状态

默认、悬停、键盘焦点与祖先禁用；业务状态由调用方提供。

## 事件

Opened() 返回从 1 开始的来源编号。

## 键盘操作

Tab 到可见来源；Enter、Space 请求打开；滚轮滚动；Active 变化定位来源。

## 限制

使用虚拟列表，只构建可见来源；调用方保证引用编号和来源顺序一致。
