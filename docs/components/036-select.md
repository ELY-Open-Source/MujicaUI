# 036 Select · 选择器

## 用途

从弹出列表中选择一个值。支持占位、清空、禁用项、空状态与输入首字母定位。

## 最小示例

```go
plan := ""
options := []mujica.SelectOption[string]{
	{Value: "free", Label: "Free"},
	{Value: "pro", Label: "Pro"},
	{Value: "team", Label: "Team", Disabled: true},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Select(c, &plan, options, mujica.SelectOptions{Clearable: true, Label: "Plan"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Select[T comparable](c *ui.Context, value *T, options []SelectOption[T], opts SelectOptions) ChoiceResult`。`*value` 不等于任何选项时显示占位；清空把值设为 `T` 的零值。

`SelectOption[T]`：`Value T`（稳定值）、`Label string`、`Disabled bool`（显示但不可选，方向键跳过）。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Placeholder` | `string` | 内置“请选择…” | 无值时的提示。 |
| `Clearable` | `bool` | false | 有值时显示清空按钮，并允许 Backspace/Delete 清空。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Empty` | `string` | 内置“暂无选项” | `options` 为空时弹层中的说明。 |
| `Label` | `string` | "" | 无障碍名称。 |

## 状态

默认、悬停（边框加深）、展开（酒红边框、箭头朝上）、聚焦、禁用；弹层中高亮项（浅底）、已选项（玫瑰底 + 勾号 + 酒红文字）、键盘高亮项另有左侧竖条、禁用项、空状态。

## 事件

`Changed()` 在值变化（选择或清空）的那一帧返回 true。

## 键盘操作

- Enter / Space / ↑ / ↓：打开弹层。
- 弹层中 ↑/↓、Home、End 移动高亮，Enter 或 Space 选择，Escape 关闭；关闭后焦点回到选择器。
- 弹层关闭时输入字母：选中以该前缀开头的选项（1 秒内连续输入组成前缀，重复同一字母则循环）；支持输入法提交的中文。
- Backspace / Delete：`Clearable` 时清空。

## 限制

- 基于 `ui.SelectBase`。其高亮索引没有公开的设置接口（ui/base.go:260），所以弹层打开时输入字母不移动高亮，只在弹层关闭时定位。
- 清空按钮不占 Tab 位（与 MyGo 搜索框一致），键盘用 Backspace/Delete 清空；它的命中区域为 28×28 DIP。
- 选项行的角色为 ListItem，选中状态由 `SelectBase.Item` 报告。`SelectParts` 没有关闭弹层的接口（ui/base.go:310），所以选择器在弹层打开时被禁用后不再构建弹层；重新启用时弹层会再次显示。
