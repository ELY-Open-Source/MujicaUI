# 421 CalendarMonthView · 月视图

## 用途

按月显示事件格。每天最多显示若干行事件，其余折叠为“+N”，点击后在浮层中列出；事件可拖到其他日期改期。同一天时间重叠的事件并排放在同一行。

## 最小示例

```go
loc := time.Local
events := exampleEvents(loc)
month := datetime.YearMonth{Year: 2026, Month: time.October}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := datetime.CalendarMonthView(c, &month, events, datetime.CalendarMonthViewOptions{Location: loc})
	if ch, ok := r.Changed(); ok {
		fmt.Println("move", ch.ID, "to", ch.Start)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func CalendarMonthView(c *ui.Context, month *YearMonth, events []CalendarEvent, opts CalendarMonthViewOptions) CalendarViewResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Location` | `*time.Location` | 必填 | 切分日期与显示时间所用的时区；为 nil 时 panic，不做隐式 UTC 转换。 |
| `Version` | `int` | 0 | 事件在原切片内被修改时须改变；换新切片会被自动发现。 |
| `WeekStart` | `WeekStart` | `WeekStartLocale` | 周起始日：跟随语言或周日、周一、周六。 |
| `MaxPerDay` | `int` | 0（即 3） | 每天显示的事件行数，超出部分折叠为“+N”。 |
| `Today` | `Date` | 零值（系统时钟） | 标记为今天的日期。 |
| `Hour12` | `bool` | false | 12 小时制显示时间。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

默认、悬停、聚焦（焦点环）、今日（圆形底色 + 加粗）、非本月（弱化）、拖动经过的日期（选中底色 + 内嵌框）、禁用；改期失败时下方显示错误。

## 事件

`Opened()` 返回被激活的事件 ID；`DayChosen()` 返回点击的日期；`Changed()` 返回 `CalendarMove` 变更（保持墙上时间与时长）。上一月/下一月直接修改 `*month`。

## 键盘操作

- Tab 到达日期数字与事件块，方向键在其间移动。
- Enter 或 Space 打开事件或选择日期。
- 事件块上 Alt+←/→ 移动一天，Alt+↑/↓ 移动一周。
- Escape 关闭“+N”浮层。

## 限制

- 拖动改期使用 MyGo 的 `Drag`/`Drop`；拖动预览是 MyGo 绘制的元素副本。
- 时间按调用方给定的 `Location` 计算与显示；夏令时跳过或重复的墙上时间会报告，不自动移动。
- 事件由调用方持有：视图只报告移动、调整与新建（`Changed()`），由调用方写回数据。
