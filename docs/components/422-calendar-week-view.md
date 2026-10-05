# 422 CalendarWeekView · 周视图

## 用途

一周的时间网格，顶部为全天行。在空白处拖动新建事件，拖动事件块移动，拖动底边调整时长；重叠事件按列并排。

## 最小示例

```go
loc := time.Local
events := exampleEvents(loc)
week := mujica.Date{Year: 2026, Month: time.October, Day: 14}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.CalendarWeekView(c, &week, events, mujica.CalendarWeekViewOptions{Location: loc, SlotMinutes: 15})
	if ch, ok := r.Changed(); ok && ch.Kind == mujica.CalendarCreate {
		fmt.Println("new event", ch.Start, ch.End)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func CalendarWeekView(c *ui.Context, week *Date, events []CalendarEvent, opts CalendarWeekViewOptions) CalendarViewResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Location` | `*time.Location` | 必填 | 切分日期与显示时间所用的时区；为 nil 时 panic，不做隐式 UTC 转换。 |
| `Version` | `int` | 0 | 事件在原切片内被修改时须改变；换新切片会被自动发现。 |
| `WeekStart` | `WeekStart` | `WeekStartLocale` | 周起始日：跟随语言或周日、周一、周六。 |
| `SlotMinutes` | `int` | 0（即 30） | 拖动与键盘的时间格长度，须整除一小时。 |
| `Today` | `Date` | 零值（系统时钟） | 标记为今天的日期。 |
| `Now` | `time.Time` | 零值（系统时钟） | 当前时间线的位置。 |
| `Hour12` | `bool` | false | 12 小时制显示时间。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

默认、悬停、聚焦列（时间格光标：选中底色 + 内嵌框）、拖动中（新建区域虚线框、被拖事件内嵌框）、今日表头、当前时间线、禁用；无法表示的时间显示错误。

## 事件

`Changed()` 返回 `CalendarCreate`、`CalendarMove` 或 `CalendarResize`；`Opened()` 返回点击的事件；`DayChosen()` 返回点击的表头日期。

## 键盘操作

- Tab 聚焦某一天的列，↑/↓ 移动时间格光标，Enter 新建一格长的事件。
- 事件块上 Alt+↑/↓ 按格移动，Alt+←/→ 移动一天，Alt+Shift+↑/↓ 调整结束时间。

## 限制

- 调整时长的底边手柄只有 6 DIP 高，低于 28 DIP 命中区域要求；同样的操作可用事件块上的 Alt+Shift+↑/↓ 完成。
- 网格高 24 小时，默认滚动到 8:00；跨越 24 小时以上的事件放在全天行。
- 时间按墙上时间排布，夏令时当天的重复或缺失小时不改变格高。
- 时间按调用方给定的 `Location` 计算与显示；夏令时跳过或重复的墙上时间会报告，不自动移动。
- 事件由调用方持有：视图只报告移动、调整与新建（`Changed()`），由调用方写回数据。
