# 301 DiffViewer 差异查看器

## 用途

按行比较两段文本，统一或并排显示：行号、增删高亮、折叠较长的未变更区；大文件只构建可见行。

## 最小示例

```go
var state git.DiffViewerState
from := "func main() {\n\tprintln(\"hi\")\n}\n"
to := "func main() {\n\tprintln(\"hello\")\n}\n"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	git.DiffViewer(c, &state, from, to, git.DiffViewerOptions{Language: "go", SideBySide: true}).Height(240)
}
ui.Render(view, 400, 300, 1)
```

签名：`func DiffViewer(c *ui.Context, s *DiffViewerState, from, to string, opts DiffViewerOptions) ui.Element`

## 参数

`DiffViewerOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| SideBySide | bool | false | 并排两栏；否则统一一栏。 |
| Context | int | 0（即 3） | 变更前后保留的未变更行数；为负 panic。 |
| Language | string | "" | 语法高亮语言，如 go；空为纯文本。 |

`DiffViewerState` 由调用方持有，保存差异、已展开的折叠与滚动位置。差异只在 from 或 to 改变时重新计算（Myers 行差异）。给返回元素一个高度。

## 状态

- 删除行：Danger 浅色底与 − 号；新增行：Success 浅色底与 + 号；不只靠颜色。
- 统一视图显示新旧两列行号；并排视图每侧一列，缺行的一侧为空底。
- 超过 2×Context 的未变更区折叠为“N 行未更改”按钮，点击展开，悬停 120 ms 过渡。

## 事件

- 无；展开折叠保存在状态中。

## 键盘操作

- 列表可获得焦点并滚动；折叠按钮用 Tab 到达，Enter 或 Space 展开。

## 限制

- 只做行级差异；词级差异见 InlineDiff。
- 长行不换行，超出部分裁剪。
- 文本换成新字符串后折叠全部重置。
