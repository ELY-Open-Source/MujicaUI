# 074 Accordion · 手风琴

## 用途

多个可展开的区段放在一个带边框的盒子中。支持单项或多项展开、是否允许全部关闭，展开状态由调用方持有（受控）。

## 最小示例

```go
open := []string{"choir"}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Accordion(c, &open, []mujica.AccordionItem{
		{ID: "choir", Title: "Choir", Content: func() { ui.Text(c, "Twelve voices.") }},
		{ID: "nave", Title: "Nave", Content: func() { ui.Text(c, "Black oak pews.") }},
	}, mujica.AccordionOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Accordion(c *ui.Context, open *[]string, items []AccordionItem, opts AccordionOptions) AccordionView`。`*open` 是展开区段的 ID 列表；`AccordionItem` 字段：`ID`（必填、唯一）、`Title`、`Content`（必填）、`Disabled`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Multiple` | `bool` | false | 允许同时展开多项；否则展开一项会关闭其他项。 |
| `KeepOne` | `bool` | false | 最后一个展开项不能被关闭。 |

## 状态

默认、悬停、聚焦、展开（箭头翻转，内容 120 ms 展开；减少动画时直接切换）、禁用（不展开，退出 Tab 序列）。

## 事件

`AccordionView.Changed()`：某区段因用户操作而展开或关闭，`*open` 已更新。外部直接修改 `*open` 不触发事件。

## 键盘操作

- Tab：到达标题；标题打开时 Tab 进入其内容。
- ↑/↓：上/下一个标题（跳过禁用项）；Home/End：第一个/最后一个。
- Enter/Space：展开或关闭。

## 限制

- 基于 MyGo `ui.CollapsibleBase`；MyGo 的 `ui.Accordion` 外观固定且用内部状态移动焦点（ui/collapsible.go `accordionKeys`），因此键盘移动在组件内重新实现：向上移动时焦点在下一帧到达。
- 关闭动画进行的 120 ms 内内容仍在界面树中。
