# 051 Calendar · 日历

## 用途

在页面中直接显示一个月份，用于选择单个日期。支持今日标记、不可选日期与周起始日配置。

## 最小示例

```go
day := mujica.Date{Year: 2026, Month: time.October, Day: 5}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.Calendar(c, &day, mujica.CalendarOptions{WeekStart: mujica.WeekStartMonday}).Changed() {
		fmt.Println("chosen", day)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Calendar(c *ui.Context, value *Date, opts CalendarOptions) DateTimeResult`，返回值的 `Element` 字段是组件元素。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Min` | `Date` | 零值 | 最早可选日期；零值不限制。 |
| `Max` | `Date` | 零值 | 最晚可选日期；零值不限制。 |
| `DisableDate` | `func(Date) bool` | nil | 返回 true 的日期不可选，显示删除线。 |
| `WeekStart` | `WeekStart` | `WeekStartLocale` | 首列星期：跟随语言（zh-CN 周一、en-US 周日），或 `WeekStartSunday`/`WeekStartMonday`/`WeekStartSaturday`。 |
| `Today` | `Date` | 零值（系统时钟） | 标记为今天的日期：加粗并在下方绘制菱形。 |
| `Disabled` | `bool` | false | 禁用整个日历：不能翻页、不能选择，退出 Tab 序列。 |
| `Label` | `string` | 内置“日历” | 辅助技术读出的名称。 |

## 状态

默认、悬停、按下、聚焦（键盘焦点环）、选中（酒红填充 + 内嵌细框 + 单选语义）、今日（加粗 + 菱形）、不可选（删除线并由 MyGo 置灰）、非本月（弱化文字）、禁用。

## 事件

`Calendar(...).Changed()` 在用户选择了与原值不同的日期的那一帧返回 true，`*value` 已更新。翻月不触发事件。

## 键盘操作

- Tab：进入日期格，落在已选日期或上次聚焦的日期。
- ←/→：前后一天；↑/↓：前后一周（自动跳过不可选日期，最多 31 天）。
- PageUp/PageDown：上/下一月；Shift+PageUp/PageDown：上/下一年。
- Home/End：本月第一天/最后一天。
- Enter/Space：选择当前日期。

## 限制

- 不可选日期不能获得焦点；若方向键 31 天内找不到可选日期，视图会翻页但焦点停在原处。
- MyGo 的 `ui.Calendar` 固定周一起始、只有英文且不支持禁用日期（ui/date.go:85-160），因此本组件基于 `RadioBase` 自绘网格，不复用 `ui.Calendar`。
- 日期与时间值不带时区：`Date`、`YearMonth`、`TimeOfDay` 只保存日历字段，组件内部计算不做 UTC 转换。
- `Today` 为零值时读取 `c.Now()`（系统本地时区）；需要固定“今天”（测试、跨时区应用）时显式传入。
- 无障碍：MyGo v0.2.10 没有公开设置弹出按钮展开状态（`expanded`）、数值范围（`hasRange`/`accRange`）与值（`accValue`）的接口，这些字段只在 MyGo 内部组件中赋值（ui/element.go:286-295，例如 ui/base.go:297、ui/indicators.go:61）。因此触发器只能声明为 `RolePopUpButton` 而无法报告展开状态；时间段以 `Description` 读出当前值与范围，错误通过 `Error()` 关联。
