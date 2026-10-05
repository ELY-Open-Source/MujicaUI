# 448 TraceWaterfall · 调用链瀑布图

## 用途

按跨度父子层级与耗时展示调用链，并读取跨度详情。

## 最小示例

```go
spans := []mujica.TraceSpan{{ID: "root", Name: "GET /api", Service: "gateway", Duration: 200 * time.Millisecond}, {ID: "db", Parent: "root", Name: "Query", Service: "database", Start: 20 * time.Millisecond, Duration: 80 * time.Millisecond}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.TraceWaterfall(c, spans, mujica.TraceWaterfallOptions{}).Changed() {
		fmt.Println("span selected")
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

Height 默认 320 DIP；TraceSpan.ID 唯一，Name 必填；Parent 关联父跨度；Start、Duration 为非负时间；Error 标记失败。

## 状态

空数据、跨度选择、失败跨度、父级禁用。

## 事件

Changed() 表示选择了跨度。

## 键盘操作

Tab 聚焦跨度；上下方向键移动；Enter 读取详情。

## 限制

更新调用链提交新切片；跨度父子关系须构成树。
