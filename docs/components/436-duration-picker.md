# 436 DurationPicker · 时长选择器

## 用途

以时、分、秒三个数字框编辑时长，并检查范围。

## 最小示例

```go
length := 90 * time.Minute
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	opts := mujica.DurationPickerOptions{Min: 15 * time.Minute, Max: 8 * time.Hour}
	if mujica.DurationPicker(c, &length, opts).Changed() {
		fmt.Println("length", length)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func DurationPicker(c *ui.Context, value *time.Duration, opts DurationPickerOptions) ChoiceResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Seconds` | `bool` | false | 显示秒字段；关闭时值须为整分钟。 |
| `Min / Max` | `time.Duration` | 0 / 99 时 59 分 59 秒 | 允许的范围。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

默认、聚焦、超出范围（错误说明，值不写入）、禁用。

## 事件

`Changed()` 报告合法的新时长。

## 键盘操作

- Tab 在三个字段间移动；↑/↓ 步进；键入数字后 Enter。

## 限制

- 超出范围的输入保留并报告，不会被夹取。
