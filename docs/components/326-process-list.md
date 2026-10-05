# 326 ProcessList · 进程列表

## 用途

以表格列出进程，可按名称、PID、CPU、内存排序；结束进程前必须确认，确认后通过结果交给调用方执行。

## 最小示例

```go
procs := []mujica.Process{{PID: 7, Name: "archivist", CPU: 42, Memory: 50 << 20}, {PID: 9, Name: "herald", CPU: 3, Memory: 8 << 20}}
var state mujica.ProcessListState
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.ProcessList(c, procs, &state, mujica.ProcessListOptions{})
	r.Element.Height(240)
	if r.End != nil {
		fmt.Println("kill", r.End.PID)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ProcessList(c *ui.Context, procs []Process, s *ProcessListState, opts ProcessListOptions) ProcessListResult`。排序只在切片或排序方式变化时进行；刷新数据请传入新切片。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Process.PID`, `Name` | `int`, `string` | — | 进程号与名称；PID 是行的稳定键。 |
| `Process.CPU` | `float64` | 0 | 占用百分比。 |
| `Process.Memory` | `uint64` | 0 | 常驻内存字节数，显示为 KB/MB/GB。 |
| `ProcessListState.Sort` | `ui.SortOrder` | CPU 降序 | 列为 `name`、`pid`、`cpu`、`memory`，其他值 panic。 |
| `Disabled` | `bool` | false | 不能排序与结束进程。 |

## 状态

选中行（Selection 背景与菱形）、排序列箭头、确认对话框（默认聚焦“取消”）、无进程、禁用。

## 事件

`End`：用户确认要结束的进程，没有则为 nil。

## 键盘操作

- Tab 到达表格；↑/↓、Home/End 选择；Delete 结束所选进程（先确认）。
- 点击列头排序，再点反向；每行结束按钮可用 Tab 到达。

## 限制

- 组件不发送信号，只回报确认结果。
