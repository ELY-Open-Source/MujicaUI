# 453 JsonViewer · JSON 查看器

## 用途

显示 JSON 树、折叠节点、复制路径并搜索键和值。

## 最小示例

```go
data := `{"server":{"port":8080},"ready":true}`
var state mujica.JsonViewerState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.JsonViewer(c, data, &state, mujica.JsonViewerOptions{})
	if r.Copied != "" {
		fmt.Println("path copied")
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

state 必填；Height 默认 360 DIP；Expanded 首次展开所有节点；Query 为搜索文本，Current 为匹配索引。

## 状态

展开、折叠、搜索匹配、无匹配、解析失败、父级禁用。

## 事件

Copied 为写入剪贴板的路径。

## 键盘操作

Tab 到搜索、复制和树；上下移动；左右展开折叠；Enter 切换节点。

## 限制

JSONPath 使用点号、数组索引和带引号的键；解析结果按数据变化缓存。
