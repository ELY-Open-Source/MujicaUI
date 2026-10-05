# StorageUsage

## 用途

存储分段条、分类明细及键盘读数。

## 最小示例

```go
segments := []files.StorageSegment{{Label: "Documents", Bytes: 24 << 30}, {Label: "Photos", Bytes: 40 << 30}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	files.StorageUsage(c, segments, files.StorageUsageOptions{Capacity: 128 << 30, Label: "Court disk"})
}
ui.Render(view, 500, 240, 1)
```

## 参数

Capacity 必须为正；Label 为业务名称；分类总量须处于容量内；原地更新时递增 Revision，分类明细采用虚拟列表。

## 状态

占用、剩余、焦点读数。

## 事件

方向键读数变化会朗读并记录状态。

## 键盘操作

Tab 聚焦分段条；左右循环读取分类与剩余空间。

## 限制

统计由调用方提供，组件展示分类大小与比例。
