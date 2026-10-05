# 047 ColorPicker

颜色选择器。

## 用途

选择带透明度的颜色：色板按钮弹出面板，包含饱和度/亮度方块、色相与透明度滑块、HEX 与 RGB、透明度输入以及预设色。透明部分显示棋盘格。

## 最小示例

```go
tint := ui.Hex("#7D2034")
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	input.ColorPicker(c, &tint, input.ColorPickerOptions{Label: "Tint"})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `color` | `*ui.Color` | 必填 | 当前颜色。 |
| `Presets` | `[]ui.Color` | 内置 12 色宫廷色板 | 预设色。 |
| `Label` | `string` | 内置“Color/颜色” | 按钮的无障碍名称。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认、悬停、聚焦、展开、输入无效（对应输入框边框变红并显示格式说明，颜色不变）、预设选中（边框加勾选标记）、禁用。

## 事件

`ColorPickerParts.Changed()` 用户选择了新颜色；`Trigger` 为按钮元素。

## 键盘操作

Enter/Space 打开；打开后焦点在方块上，方向键调整饱和度与亮度；Tab 依次到达色相、透明度滑块（方向键、Home、End）、HEX、R、G、B、A 输入框与预设色；Escape 或点击外部关闭并把焦点还给按钮。

## 限制

- HEX 接受 RGB、RRGGBB、RRGGBBAA，可带或不带 #；RGB 为 0–255 整数，透明度为 0–100。无效输入保留在输入框中直到失焦，不修改颜色。
- 色谱渐变是 PRD 允许的例外。
- 弹层位置由 MyGo PopoverBase 决定，空间不足时显示在按钮上方。
