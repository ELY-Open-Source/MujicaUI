# 037 MultiSelect · 多选选择器

## 用途

从列表中选择多个值，弹层内可搜索、全选当前匹配项，可限制数量；已选值以标签显示，超出部分折叠为“+N”。

## 最小示例

```go
tags := []int{1}
options := []mujica.SelectOption[int]{
	{Value: 1, Label: "Rose"},
	{Value: 2, Label: "Lily"},
	{Value: 3, Label: "Iris"},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.MultiSelect(c, &tags, options, mujica.MultiSelectOptions{Max: 2, Label: "Flowers"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func MultiSelect[T comparable](c *ui.Context, values *[]T, options []SelectOption[T], opts MultiSelectOptions) ChoiceResult`。`*values` 按选择顺序保存；组件每次修改都会生成新切片。`Max` 或 `MaxTags` 为负数时 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Placeholder` | `string` | 内置“请选择…” | 无选择时的提示。 |
| `Max` | `int` | 0 | 最多选择数，0 不限；达到上限后其余项禁用，并显示“已选 n / Max”。 |
| `MaxTags` | `int` | 0 → 3 | 字段中显示的标签数，其余折叠为“+N”。 |
| `Disabled` | `bool` | false | 禁用（标签不可移除）。 |
| `Label` | `string` | "" | 无障碍名称。 |

## 状态

字段：占位、标签、折叠、展开、悬停、聚焦、禁用。弹层：搜索框、全选（全选/半选/未选）、选中项（酒红方框 + 白勾）、禁用项、达到上限、无结果。

## 事件

`Changed()` 在勾选、取消、全选或移除标签的那一帧返回 true。

## 键盘操作

- Enter / Space / ↓：打开弹层，焦点在搜索框。
- 搜索框中 ↓：进入列表（“全选”）；列表内 ↑/↓、Home、End 移动；Space 勾选。
- Escape：关闭，焦点回到字段。
- 字段上 Backspace：移除最后一个值。

## 限制

- 选项行基于 `ui.CheckboxBase` 自绘方框；全选的“半选”只在画面上表示，MyGo 没有为自定义元素设置 mixed 状态的公开接口（ui/element.go:282 的 checked 字段仅由内部设置），屏幕阅读器读作未选中。
- 标签上的移除按钮不占 Tab 位，键盘用 Backspace 移除。
- 列表未虚拟化，适合数百项以内。
