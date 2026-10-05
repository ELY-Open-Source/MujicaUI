# 057 MonthPicker · 月份选择器

## 用途

按年份分页选择月份。

## 最小示例

```go
var month mujica.YearMonth
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	opts := mujica.MonthPickerOptions{Label: "Billing month", Clearable: true}
	if mujica.MonthPicker(c, &month, opts).Changed() {
		fmt.Println("month", month)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func MonthPicker(c *ui.Context, value *YearMonth, opts MonthPickerOptions) DateTimeResult`，返回值嵌入 `*ui.Element`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Min` | `YearMonth` | 零值 | 最早月份。 |
| `Max` | `YearMonth` | 零值 | 最晚月份。 |
| `DisableMonth` | `func(YearMonth) bool` | nil | 不可选月份。 |
| `Clearable` | `bool` | false | 有值时显示清空按钮。 |
| `Today` | `Date` | 零值（系统时钟） | 标记当前月份。 |
| `ReadOnly` | `bool` | false | 只读：可聚焦、可复制（Cmd/Ctrl+C），不能修改，不打开弹层。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | 内置“选择月份” | 辅助技术名称。 |

## 状态

默认、悬停、聚焦、选中（酒红填充 + 内嵌细框）、当前月（加粗 + 菱形）、不可选（删除线）、禁用。

## 事件

`Changed()`：选择或清空导致值变化时为 true。

## 键盘操作

- Enter/Space/↓：打开。
- ←/→：前后一月；↑/↓：前后三月（一行）。
- PageUp/PageDown：上/下一年。
- Home/End：一月/十二月。
- Enter：选择并关闭；Escape：关闭。

## 限制

- 键盘移到不可选月份时焦点保持在原处，可继续移动。
- 日期与时间值不带时区：`Date`、`YearMonth`、`TimeOfDay` 只保存日历字段，组件内部计算不做 UTC 转换。
- `Today` 为零值时读取 `c.Now()`（系统本地时区）；需要固定“今天”（测试、跨时区应用）时显式传入。
- 无障碍：MyGo v0.2.10 没有公开设置弹出按钮展开状态（`expanded`）、数值范围（`hasRange`/`accRange`）与值（`accValue`）的接口，这些字段只在 MyGo 内部组件中赋值（ui/element.go:286-295，例如 ui/base.go:297、ui/indicators.go:61）。因此触发器只能声明为 `RolePopUpButton` 而无法报告展开状态；时间段以 `Description` 读出当前值与范围，错误通过 `Error()` 关联。
