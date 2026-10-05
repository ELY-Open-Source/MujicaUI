# 431 AvailabilityPicker · 空闲时段选择

## 用途

按天与时间的网格，拖动涂抹自己的空闲时段，并显示其他人在每个时段的空闲人数。

## 最小示例

```go
days := []mujica.Date{{Year: 2026, Month: time.October, Day: 12}, {Year: 2026, Month: time.October, Day: 13}}
mine := [][]bool{make([]bool, 18), make([]bool, 18)}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.AvailabilityPicker(c, &mine, mujica.AvailabilityPickerOptions{Days: days}).Changed() {
		fmt.Println("availability changed")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func AvailabilityPicker(c *ui.Context, mine *[][]bool, opts AvailabilityPickerOptions) ChoiceResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Days` | `[]Date` | 必填 | 网格的列。 |
| `From / To` | `int` | 9 / 18 | 起止小时，0 ≤ From < To ≤ 24。 |
| `SlotMinutes` | `int` | 0（即 30） | 时段长度。 |
| `Others` | `[]AvailabilityPerson` | nil | 其他人的空闲表，维度须与值一致。 |
| `Hour12` | `bool` | false | 12 小时制显示时间。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

默认、自己空闲（选中底色 + 勾）、他人空闲（数量 + 条形）、全员可行（加粗强调）、键盘光标（内嵌框）、禁用。

## 事件

`Changed()` 报告自己的空闲表变化。

## 键盘操作

- Tab 聚焦网格，方向键移动光标，Space 或 Enter 切换。

## 限制

- 值的维度为 [天][时段]，与选项不符时 panic。
