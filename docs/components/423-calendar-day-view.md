# 423 CalendarDayView · 日视图

## 用途

单日时间网格；重叠事件按列并排，交互与周视图相同。

## 最小示例

```go
loc := time.Local
events := exampleEvents(loc)
day := core.Date{Year: 2026, Month: time.October, Day: 14}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if id, ok := datetime.CalendarDayView(c, &day, events, datetime.CalendarDayViewOptions{Location: loc}).Opened(); ok {
		fmt.Println("open", id)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func CalendarDayView(c *ui.Context, day *Date, events []CalendarEvent, opts CalendarDayViewOptions) CalendarViewResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Location` | `*time.Location` | 必填 | 切分日期与显示时间所用的时区；为 nil 时 panic，不做隐式 UTC 转换。 |
| `Version` | `int` | 0 | 事件在原切片内被修改时须改变；换新切片会被自动发现。 |
| `SlotMinutes` | `int` | 0（即 30） | 拖动与键盘的时间格长度，须整除一小时。 |
| `Today` | `Date` | 零值（系统时钟） | 标记为今天的日期。 |
| `Now` | `time.Time` | 零值（系统时钟） | 当前时间线的位置。 |
| `Hour12` | `bool` | false | 12 小时制显示时间。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

同周视图。

## 事件

同周视图：`Changed()`、`Opened()`、`DayChosen()`；前一天/后一天修改 `*day`。

## 键盘操作

- 同周视图。

## 限制

- 与周视图共用网格实现。
- 时间按调用方给定的 `Location` 计算与显示；夏令时跳过或重复的墙上时间会报告，不自动移动。
- 事件由调用方持有：视图只报告移动、调整与新建（`Changed()`），由调用方写回数据。
