# 055 TimeRangePicker · 时间范围选择器

## 用途

编辑开始与结束时间，校验顺序；跨午夜必须显式开启。

## 最小示例

```go
shift := datetime.TimeRange{Start: datetime.TimeOfDay{Hour: 22}, End: datetime.TimeOfDay{Hour: 6}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	opts := datetime.TimeRangePickerOptions{Label: "Night shift", AllowOvernight: true}
	if datetime.TimeRangePicker(c, &shift, opts).Changed() {
		fmt.Println("shift", shift.Start, shift.End, shift.Overnight())
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func TimeRangePicker(c *ui.Context, value *TimeRange, opts TimeRangePickerOptions) DateTimeResult`，返回值的 `Element` 字段是组件元素。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Seconds` | `bool` | false | 显示秒。 |
| `Hour12` | `bool` | false | 12 小时制。 |
| `MinuteStep` | `int` | 0（即 1） | 分钟步长。 |
| `SecondStep` | `int` | 0（即 1） | 秒步长。 |
| `Min` | `TimeOfDay` | 00:00:00 | 两端最早时间。 |
| `Max` | `TimeOfDay` | 零值（即 23:59:59） | 两端最晚时间。 |
| `AllowOvernight` | `bool` | false | 允许结束早于开始，表示次日结束，并显示“次日”。 |
| `ReadOnly` | `bool` | false | 只读：可聚焦、可复制（Cmd/Ctrl+C），不能修改，不打开弹层。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | 空 | 整组名称。 |

## 状态

同 TimePicker；范围不合法时两端都显示错误边框；跨午夜时显示“+1 day / 次日”。

## 事件

`Changed()`：两端都合法且顺序正确、与原值不同时为 true。

## 键盘操作

- 两端各自按 TimePicker 操作；Tab 从开始移到结束。

## 限制

- 未开启 `AllowOvernight` 时传入跨午夜的值会 panic（编程错误）。
- 开始等于结束视为错误。
- `TimeRange.Overnight()` 报告是否跨午夜。
- 日期与时间值不带时区：`Date`、`YearMonth`、`TimeOfDay` 只保存日历字段，组件内部计算不做 UTC 转换。
- 无障碍：MyGo v0.2.10 没有公开设置弹出按钮展开状态（`expanded`）、数值范围（`hasRange`/`accRange`）与值（`accValue`）的接口，这些字段只在 MyGo 内部组件中赋值（ui/element.go:286-295，例如 ui/base.go:297、ui/indicators.go:61）。因此触发器只能声明为 `RolePopUpButton` 而无法报告展开状态；时间段以 `Description` 读出当前值与范围，错误通过 `Error()` 关联。
