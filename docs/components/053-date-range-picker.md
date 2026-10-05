# 053 DateRangePicker · 日期范围选择器

## 用途

在弹层中先选开始日期、再选结束日期，悬停时预览范围，确认后写入。

## 最小示例

```go
var stay mujica.DateRange
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.DateRangePicker(c, &stay, mujica.DateRangePickerOptions{Label: "Stay"})
	if r.Changed() {
		fmt.Println("stay", stay.Start, stay.End, stay.Days())
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func DateRangePicker(c *ui.Context, value *DateRange, opts DateRangePickerOptions) DateTimeResult`，返回值嵌入 `*ui.Element`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Layout` | `string` | `2006-01-02` | 触发器与摘要中的日期格式。 |
| `Min` | `Date` | 零值 | 两端最早日期。 |
| `Max` | `Date` | 零值 | 两端最晚日期。 |
| `DisableDate` | `func(Date) bool` | nil | 不可选日期；范围内包含此类日期时拒绝并提示。 |
| `WeekStart` | `WeekStart` | `WeekStartLocale` | 周起始日。 |
| `Today` | `Date` | 零值（系统时钟） | 今日标记与空值时的起始月份。 |
| `ReadOnly` | `bool` | false | 只读：可聚焦、可复制（Cmd/Ctrl+C），不能修改，不打开弹层。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | 内置“日期范围” | 辅助技术名称。 |

## 状态

触发器：默认、悬停、聚焦、禁用。弹层：起止日期酒红填充，中间日期浅色带；选择结束前随指针或焦点预览范围；错误说明；“确定”在两端选齐前禁用。

## 事件

`Changed()`：点击“确定”或“清空”使值改变时为 true。`Cancelled()`：按 Escape、点击“取消”或点击弹层外关闭、草稿被丢弃时为 true。

## 键盘操作

- Enter/Space/↓：打开弹层。
- 日历内方向键移动，Enter 依次选择开始与结束。
- 早于开始日期的选择会重新开始范围。
- Tab 到“清空”“取消”“确定”。
- Escape：放弃草稿。

## 限制

- 弹层一次显示一个月；跨月范围通过翻页完成。
- `DateRange` 为闭区间，`Days()` 计入两端。
- 日期与时间值不带时区：`Date`、`YearMonth`、`TimeOfDay` 只保存日历字段，组件内部计算不做 UTC 转换。
- `Today` 为零值时读取 `c.Now()`（系统本地时区）；需要固定“今天”（测试、跨时区应用）时显式传入。
- 无障碍：MyGo v0.2.10 没有公开设置弹出按钮展开状态（`expanded`）、数值范围（`hasRange`/`accRange`）与值（`accValue`）的接口，这些字段只在 MyGo 内部组件中赋值（ui/element.go:286-295，例如 ui/base.go:297、ui/indicators.go:61）。因此触发器只能声明为 `RolePopUpButton` 而无法报告展开状态；时间段以 `Description` 读出当前值与范围，错误通过 `Error()` 关联。
