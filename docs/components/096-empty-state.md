# 096 EmptyState 空状态

## 用途

说明视图为何为空并给出下一步。区分“无数据”（集合本身为空）与“无搜索结果”（筛选没有匹配），两者使用不同的默认图标与文案。

## 最小示例

```go
r := feedback.EmptyState(c, feedback.EmptyStateOptions{
	Kind:   feedback.EmptyNoResults,
	Query:  query,
	Action: "Clear search",
})
if r.Action {
	query = ""
}
```

签名：`func EmptyState(c *ui.Context, opts EmptyStateOptions) EmptyStateResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Kind` | `EmptyKind` | `EmptyNoData` | `EmptyNoData` 或 `EmptyNoResults`。 |
| `Title` | `string` | 按类型的内置文案 | 衬线标题。 |
| `Description` | `string` | 按类型的内置文案 | 说明。 |
| `Query` | `string` | `""` | 仅 `EmptyNoResults`：默认说明写出该关键词。用于 `EmptyNoData` 时 panic。 |
| `Icon` | `*ui.SVG` | inbox / search | 替换图标。 |
| `NoIcon` | `bool` | `false` | 不显示图标；与 `Icon` 同时设置时 panic。 |
| `Action` | `string` | `""` | 次按钮文字。 |

## 状态

静态展示；操作按钮有悬停、按下与焦点状态。

## 事件

`EmptyStateResult.Action`：操作按钮被触发。

## 键盘操作

Tab 聚焦操作按钮；Enter 或 Space 触发。

## 限制

- 只有一个操作；需要多个下一步时使用 Result。
