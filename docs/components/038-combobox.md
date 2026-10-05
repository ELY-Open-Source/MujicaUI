# 038 Combobox · 组合框

## 用途

输入文字筛选选项并选择，覆盖 MyGo 的 Autocomplete（自由输入）与 Combobox（选项选择）。可允许自定义值；选项可由调用方异步提供。

## 最小示例

```go
font := ""
fonts := []string{"Garamond", "Georgia", "Didot"}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := input.Combobox(c, &font, fonts, input.ComboboxOptions{AllowCustom: true, Label: "Font"})
	if r.QueryChanged() {
		log.Println("search:", r.Query())
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Combobox(c *ui.Context, value *string, options []string, opts ComboboxOptions) ComboboxResult`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Placeholder` | `string` | "" | 空输入时的提示。 |
| `AllowCustom` | `bool` | false | 允许不在选项中的文字成为值：列表首项为“使用“xxx””，Enter 或离开字段时保留输入。 |
| `ExternalFilter` | `bool` | false | 选项已由调用方按 `Query` 过滤，组件不再筛选。 |
| `Loading` | `bool` | false | 调用方正在获取选项：字段内显示旋转标记，无选项时弹层显示“加载中…”。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | "" | 无障碍名称。 |

`ComboboxResult` 嵌入 `ChoiceResult`，另有 `Query() string`（当前输入）与 `QueryChanged() bool`（用户本帧输入了文字，此时发起请求）。

## 状态

默认、悬停、聚焦（酒红边框 + 光标）、展开、加载、无匹配（“无匹配项”）、禁用；列表高亮项、当前值（勾号）。

## 事件

- `Changed()`：值变化（选择、自定义提交）。
- `QueryChanged()`：输入变化。异步流程：在此发起请求，结果到达后（经 `Window.Update`）通过 `options` 传入并清除 `Loading`。

## 键盘操作

- 输入：筛选（前缀匹配优先）并打开列表，高亮第一项。
- ↓ / ↑：打开或移动高亮；Enter 选择；Escape 关闭。
- Tab 离开：不允许自定义时输入恢复为当前值；允许时提交输入。

## 限制

- 基于 `ui.ComboboxBase`。它在没有选项时隐藏弹层（ui/combobox.go:191），所以“加载中”“无匹配”用一个不接收输入的独立浮层显示。
- 异步结果在打字后才到达时 MyGo 不会自动高亮首项，组件让 Enter 选择首个结果。
- MyGo 替换输入框文字时保留原光标位置（ui/editor.go:117），且不公开光标设置；组件在选择后重建输入框，使光标落在新文字末尾并保持焦点。
