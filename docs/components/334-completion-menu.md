# 334 CompletionMenu · 补全菜单

## 用途

在编辑器下方显示补全候选：类型徽记、详情与说明；方向键选择，Tab 或 Enter 接受，焦点始终留在编辑器。

## 最小示例

```go
text := "se"
open := true
items := []code.Completion{{Label: "Seal", Kind: code.SymbolMethod, Detail: "func()"}, {Label: "Search", Kind: code.SymbolFunction}}
var state code.CompletionMenuState
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	in := input.TextInput(c, &text, input.TextInputOptions{Label: "Code"})
	if r := code.CompletionMenu(c, in.Input, &open, items, &state, code.CompletionMenuOptions{}); r.Accepted != nil {
		text = r.Accepted.Label
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func CompletionMenu(c *ui.Context, anchor ui.Element, open *bool, items []Completion, s *CompletionMenuState, opts CompletionMenuOptions) CompletionMenuResult`。`anchor` 是编辑器元素（如 `InputResult.Input`）；打开期间组件替换它的 `HandleInput` 以先拿到菜单按键。候选为空、禁用或编辑器禁用时自动关闭。列表虚拟化，最多显示 8 行后滚动；窗口下方空间不足时在上方弹出。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Completion.Label` | `string` | 必填 | 候选文字。 |
| `Completion.Kind` | `SymbolKind` | `SymbolFunction` | 类型徽记；越界 panic。 |
| `Completion.Detail`, `Doc` | `string` | "" | 右侧详情与当前候选的说明。 |
| `CompletionMenuState.Selected` | `int` | 0 | 当前候选。 |
| `Disabled` | `bool` | false | 关闭菜单。 |

## 状态

当前候选 Selection 背景与菱形；弹层 160 ms 淡入（减少动画时直接显示）。

## 事件

`Accepted`：被接受的候选，同时 `*open` 置为 false。

## 键盘操作

- ↓/↑：循环移动；Tab 或 Enter：接受；Escape：关闭。点击候选也接受。

## 限制

- 打开期间编辑器自身设置的 `HandleInput` 会被覆盖。
