# 050 Transfer

穿梭框。

## 用途

在左右两个列表间移动条目，选出一个子集。两侧各有搜索、复选与全选，中间按钮批量移动。以稳定 ID 保存状态。

## 最小示例

```go
chosen := []string{"uika"}
items := []input.TransferItem{{ID: "uika", Label: "Uika"}, {ID: "mutsumi", Label: "Mutsumi"}, {ID: "nyamu", Label: "Nyamu", Disabled: true}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.Transfer(c, &chosen, items, input.TransferOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `target` | `*[]string` | 必填 | 右侧已选条目的 ID，按加入顺序；含未知 ID 时 panic。 |
| `items` | `[]TransferItem` | 必填 | 全部条目：ID（唯一、非空）、Label、Disabled。 |
| `SourceTitle / TargetTitle` | `string` | 内置“Available/可选”“Chosen/已选” | 两侧标题。 |
| `Height` | `float32` | 240 | 列表高度（DIP）。 |
| `Disabled` | `bool` | false | 禁用整个组件。 |

## 状态

默认、悬停、勾选（底色加勾选标记）、半选（全选框显示减号）、禁用条目（不可勾选、不随全选移动）、空列表、无搜索结果、禁用。

## 事件

`ChoiceResult.Changed()` 有条目移动。

## 键盘操作

Tab 进入列表，↑/↓ 在条目间移动，Space 勾选；Tab 到达全选框、搜索框与移动按钮。

## 限制

- 条目逐项创建，适合数百条以内；大数据量请使用虚拟化列表组合。
- 全选只作用于当前搜索结果中的可用条目。
