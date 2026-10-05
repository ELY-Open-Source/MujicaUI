# 464 AccentColorPicker

主色选择。

## 用途

从候选色中选择主色；承载浅色文字对比度不足 4.5:1 的候选带警示标记，选中时给出提示。

## 最小示例

```go
accent := ui.Hex("#7D2034")
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.AccentColorPicker(c, &accent, mujica.AccentColorPickerOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `color` | `*ui.Color` | 必填 | 当前主色。 |
| `Candidates` | `[]ui.Color` | 内置 8 色 | 候选色。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认、选中（外框与勾选）、对比度不足（警示图标与说明）、禁用。

## 事件

`ChoiceResult.Changed()`。

## 键盘操作

Tab 进入单选组，←/→ 移动，Space 选择。

## 限制

- 对比度按 OnAccent 文字计算；只标记，不阻止选择。
