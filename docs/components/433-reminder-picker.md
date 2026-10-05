# 433 ReminderPicker · 提醒选择

## 用途

多选事件前的提醒时间，按选项顺序保存。

## 最小示例

```go
reminders := []time.Duration{15 * time.Minute}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.ReminderPicker(c, &reminders, mujica.ReminderPickerOptions{Max: 3}).Changed() {
		fmt.Println("reminders", reminders)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ReminderPicker(c *ui.Context, value *[]time.Duration, opts ReminderPickerOptions) ChoiceResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Choices` | `[]time.Duration` | 0、5 分、15 分、30 分、1 时、1 天、1 周 | 可选的提前量，须为整分钟且不重复。 |
| `Max` | `int` | 0（不限） | 最多选择的数量，达到后其余选项禁用。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

默认、选中（选中底色 + 下划线标记）、达到上限的禁用项、禁用。

## 事件

`Changed()` 报告选择变化。

## 键盘操作

- Tab 进入列表，↑/↓ 移动，Space 切换。

## 限制

- 只支持事件开始前的提醒；提前量须为整分钟，否则 panic。
