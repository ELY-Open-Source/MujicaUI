# 071 DataTable · 数据表格

## 用途

按列展示数据行，只构建可见行。提供排序、筛选、多选、列宽、列顺序、列显隐与自定义单元格。行与列都以稳定 ID 保存状态：排序、筛选或调整列之后，选择仍指向同一条数据。默认使用本地排序与筛选（仅在数据、排序或筛选条件变化时计算一次并缓存）；`Controlled` 模式下按调用方顺序显示，调用方在 `SortChanged()` 时自行查询。

## 最小示例

```go
type relic struct {
	ID   int
	Name string
	Age  int
}
rows := []relic{{1, "Censer", 300}, {2, "Altar", 120}, {3, "Bell", 80}}
var state mujica.DataTableState[int]
opts := mujica.DataTableOptions[relic, int]{
	Key: func(r relic) int { return r.ID },
	Columns: []mujica.DataColumn[relic]{
		{ID: "name", Title: "Name", Sortable: true,
			Text:    func(r relic) string { return r.Name },
			Compare: func(a, b relic) int { return cmp.Compare(a.Name, b.Name) }},
		{ID: "age", Title: "Age", Width: 80, Align: ui.End, Sortable: true,
			Text:    func(r relic) string { return strconv.Itoa(r.Age) },
			Compare: func(a, b relic) int { return cmp.Compare(a.Age, b.Age) }},
	},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.DataTable(c, &state, rows, opts).Size(360, 240)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func DataTable[R any, K comparable](c *ui.Context, s *DataTableState[K], rows []R, opts DataTableOptions[R, K]) DataTableView`。另有 `DataTableColumnPicker(c, s, cols)`（“列”按钮与勾选弹层，控制显隐）和纯函数 `DataTableRows(rows, cols, sort, filter) []int`（本地排序筛选辅助，返回行下标）。

`DataTableState` 字段：`Sort ui.SortOrder`、`Filter string`（不区分大小写包含匹配全部列的 `Text`）、`Hidden map[string]bool`（按列 ID 隐藏）、`Layout ui.TableLayout`（用户拖动得到的列顺序与列宽，可保存恢复）、`Selection ui.Selection[K]`（多选）。方法：`Selected()`、`Rows()`（筛选后行数）。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Columns` | `[]DataColumn[R]` | 必填 | 列定义：`ID`（必填、唯一）、`Title`、`Width`（0 表示分配剩余宽度）、`MinWidth`、`Align`、`Sortable`、`Compare`（本地排序比较函数）、`Text`（必填，纯文本，用于筛选与默认单元格）、`Cell`（自定义单元格）。 |
| `Key` | `func(R) K` | 必填 | 行的稳定 ID；重复时 panic。 |
| `Multiple` | `bool` | false | 多选：首列显示勾选框，点击勾选框增减，Cmd/Ctrl、Shift 点击与 Cmd/Ctrl+A。 |
| `Controlled` | `bool` | false | 受控模式：不在本地排序筛选。 |
| `Version` | `int` | 0 | 原地修改 `rows` 后递增，通知缓存失效；换新切片无需递增。 |
| `Status` | `DataStatus` | `DataReady` | 加载状态：`DataLoading` 首次加载（无行时显示加载标记）、`DataLoadingMore` 追加加载（行下方显示“正在加载更多”）、`DataFailed` 失败（显示错误与“重试”）。 |
| `Error` | `string` | 内置“加载失败” | `DataFailed` 时显示的具体错误。 |
| `Empty` | `string` | 内置“暂无数据” | 无数据时的说明文字。 |

## 状态

默认、表头悬停、排序箭头（升/降）、行悬停、选中行（Accent 填充 + OnAccent 文字；单选首列菱形，多选首列勾选框）、聚焦（表格外焦点环）、加载、追加加载、失败、空、无匹配。数据行高度不低于当前密度的数据行高（28/36/44 DIP）。

## 事件

`Changed()`：选择变化；`Submitted()`：Enter 或双击；`SortChanged()`：点击表头改变了 `State.Sort`；`Retried()`：点击“重试”。

## 键盘操作

- ↑/↓、Home/End：移动选择；Enter：提交；输入字母按首个可见列定位。
- 多选：Shift/Cmd(Ctrl) 配合点击或方向键；Cmd/Ctrl+A 全选。
- 表头：点击排序、再次点击反向；拖动表头边缘调整列宽，双击适配内容；拖动表头移动列。
- 列选择器：Tab 进入勾选项，Space 切换；Escape 关闭。

## 限制

- MyGo `ui.Table` 只按调用方给出的顺序显示（docs/ui/table.md “Sorting”），因此本组件在本地完成排序，或在受控模式下交给调用方。
- 表头内容只能是 `Title` 文字（ui/table.go `table` 构建表头），无法放置全选勾选框；全选使用 Cmd/Ctrl+A。
- 单选的跟随由组件按键重新定位，跨越任意距离都有效；多选的键盘锚点在大规模重排后由 MyGo 重置（ui/list.go:30）。
- 删除的数据行对应的选择会在数据变化时清除；被筛选隐藏的行保留选择。
- 10 万行 × 10 列：测试统计“最后一帧构建过单元格的行数”不超过 16，并断言重复帧不再调用比较函数（data_table_test.go）。
