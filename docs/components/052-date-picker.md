# 052 DatePicker · 日期选择器

## 用途

可键入的日期输入框，带日历弹层。用于表单中录入单个日期。

## 最小示例

```go
due := mujica.Date{Year: 2026, Month: time.October, Day: 5}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	opts := mujica.DatePickerOptions{Label: "Due date", Clearable: true,
		Min: mujica.Date{Year: 2026, Month: time.January, Day: 1}}
	if mujica.DatePicker(c, &due, opts).Changed() {
		fmt.Println("due", due)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func DatePicker(c *ui.Context, value *Date, opts DatePickerOptions) DateTimeResult`，返回值嵌入 `*ui.Element`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Layout` | `string` | `2006-01-02` | Go 时间布局，用于显示与解析输入。 |
| `Min` | `Date` | 零值 | 最早可接受日期。 |
| `Max` | `Date` | 零值 | 最晚可接受日期。 |
| `DisableDate` | `func(Date) bool` | nil | 不可选日期；键入此类日期会报告错误。 |
| `WeekStart` | `WeekStart` | `WeekStartLocale` | 弹层日历的周起始日。 |
| `Today` | `Date` | 零值（系统时钟） | 弹层日历中的今日标记；空值时弹层打开到此月份。 |
| `Clearable` | `bool` | false | 显示清空按钮；清空输入后提交表示“无日期”。 |
| `Placeholder` | `string` | 内置“选择日期” | 空输入框的占位文字。 |
| `ReadOnly` | `bool` | false | 只读：可聚焦、可复制（Cmd/Ctrl+C），不能修改，不打开弹层。 |
| `Disabled` | `bool` | false | 禁用输入与弹层。 |
| `Label` | `string` | 占位文字 | 辅助技术名称。 |

## 状态

默认、悬停（边框加深）、聚焦（焦点色边框）、只读（背景与边框减弱，可选择复制）、错误（危险色边框 + 下方说明，输入保留）、禁用。弹层内日历状态同 Calendar。

## 事件

`Changed()`：键入后按 Enter 或离开输入框并解析成功、在日历中选择、或点击清空导致值变化时为 true。

## 键盘操作

- 键入日期后 Enter 提交；Tab 离开时也会提交。
- ↓ 或 Alt+↓：打开日历；日历内按键同 Calendar，Enter 选择后关闭并写入。
- Escape：关闭日历，不改值。

## 限制

- 无效或越界的输入（如 2026-02-30）显示具体错误并保留原文，不会被修正为相邻日期。
- MyGo 的 `ui.DateInput` 固定格式且不可键入（ui/date.go:13），因此本组件组合 `TextInputBase` 与自绘日历弹层（`PopoverBase`）。
- 只读模式沿用 TextInput 的做法：编辑副本并在每帧恢复，文本仍可选择与复制。
- 日期与时间值不带时区：`Date`、`YearMonth`、`TimeOfDay` 只保存日历字段，组件内部计算不做 UTC 转换。
- `Today` 为零值时读取 `c.Now()`（系统本地时区）；需要固定“今天”（测试、跨时区应用）时显式传入。
- 无障碍：MyGo v0.2.10 没有公开设置弹出按钮展开状态（`expanded`）、数值范围（`hasRange`/`accRange`）与值（`accValue`）的接口，这些字段只在 MyGo 内部组件中赋值（ui/element.go:286-295，例如 ui/base.go:297、ui/indicators.go:61）。因此触发器只能声明为 `RolePopUpButton` 而无法报告展开状态；时间段以 `Description` 读出当前值与范围，错误通过 `Error()` 关联。
