# 331 ProblemsPanel · 问题面板

## 用途

列出错误、警告与信息，按文件分组并按行列排序，可按级别过滤，打开问题时回调跳转。

## 最小示例

```go
problems := []code.Problem{{Severity: code.ProblemError, File: "court/seal.go", Line: 12, Col: 5, Message: "undefined: Decree"}}
var state code.ProblemsPanelState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := code.ProblemsPanel(c, problems, &state, code.ProblemsPanelOptions{})
	r.Element.Height(240)
	if r.Opened != nil {
		fmt.Println("go to", r.Opened.File, r.Opened.Line)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ProblemsPanel(c *ui.Context, problems []Problem, s *ProblemsPanelState, opts ProblemsPanelOptions) ProblemsPanelResult`。分组排序只在数据、过滤或折叠变化时进行。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Problem.Severity` | `ProblemSeverity` | `ProblemError` | `ProblemError`、`ProblemWarning`、`ProblemInfo`；越界 panic。 |
| `Problem.File`, `Line`, `Col` | `string`, `int` | — | 位置。 |
| `Problem.Message`, `Source` | `string` | — | 说明与来源（如编译器名）。 |
| `ProblemsPanelState.Hidden` | `[3]bool` | 全显示 | 被隐藏的级别。 |
| `Disabled` | `bool` | false | 不能过滤、折叠与打开。 |

## 状态

级别按钮带图标与计数；问题行用不同形状的图标（圆叉、三角、圆 i）与颜色共同标识级别，并有无障碍名称；选中行 Selection 与菱形；无问题；禁用。

## 事件

`Opened`：用户打开的问题。

## 键盘操作

- Tab 到达级别按钮与列表；↑/↓ 选择；Enter 或双击打开问题，在文件行上折叠/展开；← 折叠、→ 展开。

## 限制

- 同一文件的问题按行、列排序，文件按路径排序。
