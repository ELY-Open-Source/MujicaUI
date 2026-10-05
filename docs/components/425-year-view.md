# 425 YearView · 年视图

## 用途

十二个月的小日历，每天下方用 1～3 个点表示事件密度，读屏会读出事件数量。

## 最小示例

```go
loc := time.Local
events := exampleEvents(loc)
year := 2026
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if day, ok := datetime.YearView(c, &year, events, datetime.YearViewOptions{Location: loc}).DayChosen(); ok {
		fmt.Println("go to", day)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func YearView(c *ui.Context, year *int, events []CalendarEvent, opts YearViewOptions) CalendarViewResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Location` | `*time.Location` | 必填 | 切分日期与显示时间所用的时区；为 nil 时 panic，不做隐式 UTC 转换。 |
| `Version` | `int` | 0 | 事件在原切片内被修改时须改变；换新切片会被自动发现。 |
| `WeekStart` | `WeekStart` | `WeekStartLocale` | 周起始日：跟随语言或周日、周一、周六。 |
| `Today` | `Date` | 零值（系统时钟） | 标记为今天的日期。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

默认、悬停、聚焦、今日（底色 + 加粗）、密度点（1～2 项一点、3～5 项两点、6 项以上三点）、禁用。

## 事件

`DayChosen()` 返回被激活的日期；上一年/下一年及跨年移动修改 `*year`。

## 键盘操作

- Tab 进入日期，←/→ 移动一天，↑/↓ 移动一周，PageUp/PageDown 移动一月，可跨年。
- Enter 选择。

## 限制

- 密度不只靠颜色：用点数区分，并写入标签。
- 时间按调用方给定的 `Location` 计算与显示；夏令时跳过或重复的墙上时间会报告，不自动移动。
- 事件由调用方持有：视图只报告移动、调整与新建（`Changed()`），由调用方写回数据。
