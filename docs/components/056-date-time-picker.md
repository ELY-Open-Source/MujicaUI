# 056 DateTimePicker · 日期时间选择器

## 用途

在调用方指定的时区中编辑日期和时间，确认后写入 `time.Time`。

## 最小示例

```go
tokyo, err := time.LoadLocation("Asia/Tokyo")
if err != nil {
	panic(err)
}
var meeting time.Time
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	opts := mujica.DateTimePickerOptions{Label: "Meeting", Location: tokyo, MinuteStep: 15}
	if mujica.DateTimePicker(c, &meeting, opts).Changed() {
		fmt.Println("meeting", meeting)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func DateTimePicker(c *ui.Context, value *time.Time, opts DateTimePickerOptions) DateTimeResult`，返回值嵌入 `*ui.Element`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Location` | `*time.Location` | 必填 | 用户编辑所用时区；为 nil 时 panic。 |
| `Layout` | `string` | `2006-01-02` | 日期部分的 Go 时间布局，用于显示与错误说明。 |
| `Min` | `time.Time` | 零值 | 最早时刻。 |
| `Max` | `time.Time` | 零值 | 最晚时刻。 |
| `Seconds` | `bool` | false | 显示秒。 |
| `Hour12` | `bool` | false | 12 小时制。 |
| `MinuteStep` | `int` | 0（即 1） | 分钟步长。 |
| `WeekStart` | `WeekStart` | `WeekStartLocale` | 周起始日。 |
| `Today` | `Date` | 零值（系统时钟） | 今日标记与空值时的默认日期。 |
| `ReadOnly` | `bool` | false | 只读：可聚焦、可复制（Cmd/Ctrl+C），不能修改，不打开弹层。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | 内置占位文字 | 辅助技术名称。 |

## 状态

触发器显示该时区下的时间与时区缩写。弹层：日历、时间段、时区名、错误说明、“取消/确定”。

## 事件

`Changed()`：确认且时刻改变时为 true，写入的值位于 `Location`。`Cancelled()`：取消、Escape 或点击外部关闭时为 true。

## 键盘操作

- Enter/Space/↓：打开。
- 日历与时间段按键同 DatePicker、TimePicker。
- Tab 到“取消”“确定”。
- Escape：取消。

## 限制

- 夏令时跳过的时间（如纽约 2026-03-08 02:30）和重复的时间（如 2026-11-01 01:30）都会提示，不会自动移到相邻时刻或选择其中之一。
- 空值（零 `time.Time`）打开时默认今天 00:00。
- 日期与时间值不带时区：`Date`、`YearMonth`、`TimeOfDay` 只保存日历字段，组件内部计算不做 UTC 转换。
- `Today` 为零值时读取 `c.Now()`（系统本地时区）；需要固定“今天”（测试、跨时区应用）时显式传入。
- 无障碍：MyGo v0.2.10 没有公开设置弹出按钮展开状态（`expanded`）、数值范围（`hasRange`/`accRange`）与值（`accValue`）的接口，这些字段只在 MyGo 内部组件中赋值（ui/element.go:286-295，例如 ui/base.go:297、ui/indicators.go:61）。因此触发器只能声明为 `RolePopUpButton` 而无法报告展开状态；时间段以 `Description` 读出当前值与范围，错误通过 `Error()` 关联。
