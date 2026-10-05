# 430 RecurrenceEditor · 重复规则编辑器

## 用途

选择每天、每周、每月、每年或自定义间隔，每周可选星期，并设置结束条件（永不、次数、日期）。

## 最小示例

```go
rule := datetime.Recurrence{Unit: datetime.RepeatWeekly, Interval: 1, Weekdays: []time.Weekday{time.Monday}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := datetime.RecurrenceEditor(c, &rule, datetime.RecurrenceEditorOptions{Start: core.Date{Year: 2026, Month: time.October, Day: 5}})
	if r.Changed() && r.Problem() == "" {
		fmt.Println(datetime.RecurrenceSummary(c, rule))
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func RecurrenceEditor(c *ui.Context, rule *Recurrence, opts RecurrenceEditorOptions) RecurrenceResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Start` | `Date` | 零值 | 首次发生日期；结束日期不能早于它，每周规则默认取其星期。 |
| `WeekStart` | `WeekStart` | `WeekStartLocale` | 周起始日：跟随语言或周日、周一、周六。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

默认、自定义（显示间隔与单位）、每周（星期按钮）、结束条件、问题提示（无星期、次数小于 1、结束日期过早）、禁用。

## 事件

`Changed()` 报告规则变化；`Problem()` 返回当前不可保存的原因；`RecurrenceSummary` 生成本地化描述。

## 键盘操作

- Tab 在重复选择、间隔、星期与结束条件之间移动，各控件保留自身按键。

## 限制

- 只描述规则，不展开具体日期。
