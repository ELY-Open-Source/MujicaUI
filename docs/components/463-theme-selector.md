# 463 ThemeSelector

主题选择。

## 用途

以三张预览卡选择跟随系统、浅色或深色；预览用对应主题的设计变量绘制。

## 最小示例

```go
mode := mujica.System
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{Mode: mode})
	mujica.ThemeSelector(c, &mode, mujica.ThemeSelectorOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `mode` | `*Mode` | 必填 | System、Light 或 Dark；其他值 panic。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

默认、悬停、聚焦、选中（选中底色、边框与勾选）、禁用。

## 事件

`ChoiceResult.Changed()`。

## 键盘操作

Tab 进入单选组，←/→ 移动，Space 选择。

## 限制

- 无。
