# 427 EventEditor · 事件编辑器

## 用途

编辑事件的标题、时间、地点、参与者、提醒与重复规则。编辑的是副本，点击保存并校验通过后才写回。

## 最小示例

```go
loc := time.Local
ev := exampleEvents(loc)[0]
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.EventEditor(c, &ev, mujica.EventEditorOptions{Location: loc})
	if r.Saved() {
		fmt.Println("saved", ev.Title, ev.Start)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func EventEditor(c *ui.Context, ev *CalendarEvent, opts EventEditorOptions) EventEditorResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Location` | `*time.Location` | 必填 | 切分日期与显示时间所用的时区；为 nil 时 panic，不做隐式 UTC 转换。 |
| `Hour12` | `bool` | false | 12 小时制显示时间。 |
| `WeekStart` | `WeekStart` | `WeekStartLocale` | 周起始日：跟随语言或周日、周一、周六。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

默认、字段错误（保存后显示：缺少标题、结束早于开始、重复规则问题）、全天模式（只选日期）、禁用。

## 事件

`Saved()`：校验通过并写回 `*ev`；`Canceled()`：丢弃草稿。

## 键盘操作

- Tab 在字段间移动，各字段保留自身按键；保存与取消为按钮。

## 限制

- 全天事件保存时规整为当天零点到次日零点。
- 参与者按姓名添加，新参与者状态为“未回复”。
