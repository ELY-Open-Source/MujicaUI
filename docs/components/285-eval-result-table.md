# 285 EvalResultTable 评测结果表

## 用途

列出评测用例、得分与通过/失败，顶部显示通过数，每行可查看差异。

## 最小示例

```go
var state data.DataTableState[string]
cases := []feedback.EvalCase{{ID: "c1", Name: "Greets by title", Score: 0.98, Passed: true}, {ID: "c2", Name: "Refuses forgery", Score: 0.41}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := feedback.EvalResultTable(c, &state, cases, feedback.EvalResultTableOptions{})
	r.Table.Element.Height(260)
	if id, ok := r.Diff(); ok {
		_ = id // open expected and actual output side by side
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func EvalResultTable(c *ui.Context, s *DataTableState[string], cases []EvalCase, opts EvalResultTableOptions) EvalResultTableResult`

## 参数

`EvalCase`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| ID | string | 必填 | 唯一标识。 |
| Name | string | "" | 用例名。 |
| Score | float64 | 0 | 得分；非有限值 panic。 |
| Passed | bool | false | 是否通过。 |

`EvalResultTableOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Decimals | int | 0（即 2） | 得分小数位；为负 panic。 |
| Status | DataStatus | DataReady | 加载、失败状态。 |
| Error | string | "" | 失败说明。 |
| Version | int | 0 | 原地修改数据时递增；新切片或长度变化自动刷新。 |

返回 `EvalResultTableResult`，含命名字段 `Element` 与 `Table`（DataTableView）。给 `Element` 一个高度。通过数与数据校验按数据版本缓存。

## 状态

- 结果列为图标加“通过/失败”文字与颜色。
- 顶部“通过 n / m”与比例条。

## 事件

- `Diff()`：要查看差异的用例 ID。
- `Table` 的 `Changed()`、`SortChanged()`。

## 键盘操作

- ↑/↓ 选择行；差异按钮用 Tab 到达，Enter 或 Space 触发。

## 限制

- 差异内容由调用方显示，例如用 DiffViewer。
