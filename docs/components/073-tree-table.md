# 073 TreeTable · 树形表格

## 用途

把树形层级与表格列结合：首列显示展开箭头与缩进，其余列显示属性。支持展开、选择、列宽与列顺序、排序状态（同级节点在本地排序，结果按父节点缓存，只在子项或排序变化时重新计算），以及外部加载子项。只构建可见行。

## 最小示例

```go
kids := map[string][]string{"Vestry": {"Stole", "Alb"}}
size := map[string]int{"Vestry": 9, "Stole": 2, "Alb": 5, "Organ": 7}
var state mujica.TreeTableState[string]
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.TreeTable(c, &state, mujica.TreeTableOptions[string]{
		Roots:    []string{"Vestry", "Organ"},
		Children: func(n string) []string { return kids[n] },
		Label:    func(n string) string { return n },
		Columns: []mujica.TreeColumn[string]{
			{ID: "name", Title: "Name", Cell: func(n string) { ui.Text(c, n) },
				Compare: cmp.Compare[string]},
			{ID: "size", Title: "Size", Width: 80, Align: ui.End,
				Cell: func(n string) { ui.Text(c, strconv.Itoa(size[n])) }},
		},
	}).Size(360, 240)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func TreeTable[K comparable](c *ui.Context, s *TreeTableState[K], opts TreeTableOptions[K]) DataTableView`。`TreeTableState` 字段：`Sort`、`Layout`、`Selection`；方法：`Selected()`、`SetOpen`、`IsOpen`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Columns` | `[]TreeColumn[K]` | 必填 | 列：`ID`（必填、唯一）、`Title`、`Width`、`MinWidth`、`Align`、`Compare`（设置后可点击排序同级节点）、`Cell`（必填）。首列在箭头后显示 `Cell`。 |
| `Roots` | `[]K` | 空 | 顶层节点。 |
| `Children` | `func(K) []K` | 必填 | 子节点：nil 为叶子，空切片为未加载。子项变化时请返回新切片（排序缓存按切片识别）。 |
| `Label` | `func(K) string` | 必填 | 无障碍名称与输入定位。 |
| `ItemStatus`, `Load` | | nil | 与 Tree 相同的外部加载协议。 |
| `Multiple` | `bool` | false | 多选。 |
| `Status` | `DataStatus` | `DataReady` | 加载状态：`DataLoading` 首次加载（无行时显示加载标记）、`DataLoadingMore` 追加加载（行下方显示“正在加载更多”）、`DataFailed` 失败（显示错误与“重试”）。 |
| `Error` | `string` | 内置“加载失败” | `DataFailed` 时显示的具体错误。 |
| `Empty` | `string` | 内置“暂无数据” | 无数据时的说明文字。 |

## 状态

默认、悬停、选中（Accent 填充 + 勾选标记）、聚焦、展开/折叠、排序箭头、子项加载中/失败、空。

## 事件

返回 `DataTableView`：`Changed()`、`Submitted()`、`SortChanged()`、`Retried()`。

## 键盘操作

- 与 Tree 相同的方向键与展开键；点击表头排序同级节点，再次点击反向；拖动表头边缘或表头调整列宽与顺序。

## 限制

- 基于 MyGo `ui.OutlineTable`；表头只能显示文字。
- 排序只作用于同级节点，层级关系保持不变。
