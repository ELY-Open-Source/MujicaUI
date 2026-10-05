# 043 InputGroup

输入组合。

## 用途

把单行输入框与前后缀文字、按钮放进同一个边框，例如网址前缀、单位、清空或搜索按钮。每个按钮保留独立焦点。

## 最小示例

```go
site := ""
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	g := input.InputGroup(c, &site, input.InputGroupOptions{
		Prefix: "https://", Label: "Site",
		Buttons: []input.InputGroupButton{{Label: "Open"}},
	})
	if g.Clicked == 0 {
		c.OpenURL("https://" + site)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `*string` | 必填 | 输入框的值。 |
| `Prefix / Suffix` | `string` | "" | 输入框前后的固定文字。 |
| `Buttons` | `[]InputGroupButton` | nil | 按钮：Label（必填，文字或无障碍名称）、Icon、IconOnly、Leading（放在输入框前）、Disabled。 |
| `Placeholder` | `string` | "" | 占位文字。 |
| `Label` | `string` | "" | 不在 FormField 内时为输入框命名。 |
| `Disabled` | `bool` | false | 禁用整个组合。 |
| `ReadOnly` | `bool` | false | 值显示为可选择、可复制、可聚焦的文字，不能修改；按钮仍可用。 |
| `Invalid` | `bool` | false | 边框显示危险色。 |

## 状态

默认、悬停（边框加深）、聚焦（强调色边框与焦点环）、禁用、只读、错误。

## 事件

`InputGroupParts.Changed()` 值被编辑；`Submitted()` 在输入框中按 Enter；`Clicked` 为被点击按钮在 Buttons 中的下标，未点击为 -1。`Input` 为可聚焦的输入元素。

## 键盘操作

输入框与每个按钮都是独立的 Tab 停靠点；按钮响应 Enter 和 Space；输入框保留选择、剪贴板、撤销与输入法。

## 限制

- MyGo 没有公开的只读输入框，ReadOnly 改用可选择文字实现，宽度与输入框一致。
- IconOnly 按钮必须提供 Icon，Label 作为无障碍名称；缺失时 panic。
