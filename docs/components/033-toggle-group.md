# 033 ToggleGroup · 切换按钮组

## 用途

一排相连的切换按钮，单选（如对齐方式）或多选（如粗体、斜体、下划线）。可设置是否允许全部松开。

## 最小示例

```go
align := []string{"left"}
items := []mujica.ToggleItem[string]{
	{Value: "left", Label: "Left"},
	{Value: "center", Label: "Center"},
	{Value: "right", Label: "Right"},
}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ToggleGroup(c, &align, items, mujica.ToggleGroupOptions{Label: "Align"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ToggleGroup[T comparable](c *ui.Context, value *[]T, items []ToggleItem[T], opts ToggleGroupOptions) ChoiceResult`。单选模式下 `*value` 最多一个元素，否则 panic。

`ToggleItem[T]`：

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Value` | `T` | — | 稳定的选项值。 |
| `Label` | `string` | "" | 文字。 |
| `Icon` | `*ui.SVG` | nil | 图标。 |
| `Name` | `string` | "" | 仅图标时的无障碍名称。 |
| `Disabled` | `bool` | false | 单项禁用，键盘导航跳过它。 |

`ToggleGroupOptions`：

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Multiple` | `bool` | false | 允许同时按下多项。 |
| `AllowEmpty` | `bool` | false | 允许松开最后一个已按下项。 |
| `Vertical` | `bool` | false | 纵向排列，↑/↓ 导航。 |
| `Disabled` | `bool` | false | 禁用整组。 |
| `Label` | `string` | "" | 组的无障碍名称。 |

## 状态

未按下、按下（玫瑰底 + 酒红文字 + 底部横条）、悬停、按住、聚焦、单项禁用、整组禁用。

## 事件

返回的 `ChoiceResult.Changed()` 在值变化的那一帧返回 true。不允许空选时点击唯一按下项不产生事件。

## 键盘操作

- Tab：进入组。首次落在第一个可用项（与哪项按下无关）；之后回到组内上次聚焦的项。再按 Tab 离开。
- ←/→（纵向为 ↑/↓）、Home、End：在项之间移动焦点，跳过禁用项，首尾循环。
- Space：按下或松开当前项。

## 限制

- 每项基于 `ui.ToggleBase`，组用 MyGo `FocusGroup` 实现方向键导航；MyGo 的 `Key` 不能用于自带输入处理的控件（ui/element.go:474），所以项以位置区分，`items` 在帧之间增删时焦点可能落到相邻项。
