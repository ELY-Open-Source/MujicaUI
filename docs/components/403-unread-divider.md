# 403 UnreadDivider 未读分隔线

## 用途

标记新消息开始位置，并朗读完整未读计数。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.UnreadDivider(c, 3)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`count int` 必须 ≥ 0；0 显示“新消息”，正值显示完整计数。

## 状态

状态由调用方数据控制；浅深主题、密度与文字缩放随窗口设置变化。适用插槽操作沿用控件焦点和禁用状态。

## 事件

纯展示，计数由调用方更新。

## 键盘操作

展示组件保持自然焦点序列。

## 限制

组件展示调用方提供的计数。
