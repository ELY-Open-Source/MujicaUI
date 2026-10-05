# 077 DescriptionList · 描述列表

## 用途

以名称与值的形式展示详情，支持多列、跨列与复制值。容器宽度不足时自动减少列数，最窄时为单列。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.DescriptionList(c, []data.DescriptionItem{
		{Label: "Archivist", Value: "Sister Agnès", Copyable: true},
		{Label: "Shelf", Value: "North vault, row 7"},
		{Label: "Seal", Value: "Wine-red wax", Span: -1},
	}, data.DescriptionListOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func DescriptionList(c *ui.Context, items []DescriptionItem, opts DescriptionListOptions) DescriptionListView`。`DescriptionItem` 字段：`Label`、`Value`、`Span`（跨列数，负数跨整行）、`Copyable`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Columns` | `int` | 2 | 最多列数。 |
| `MinColumnWidth` | `float32` | 220 | 每列最小宽度（DIP）；容器放不下时减少列数。 |

## 状态

默认、复制按钮悬停与聚焦、已复制（勾选图标 1.5 秒，并向读屏软件播报“已复制”）。

## 事件

`DescriptionListView.Copied()` 返回本帧被复制项的下标。

## 键盘操作

- Tab 到达每个复制按钮（名称为“复制 + 标签”），Enter/Space 复制。

## 限制

- 列数根据上一帧的宽度计算，宽度变化后在下一帧重新排列（组件会立即请求该帧）。
- 复制使用 MyGo 剪贴板接口，只写入纯文本。
