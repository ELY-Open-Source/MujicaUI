# 434 AttendeeList · 参与者列表

## 用途

列出参与者、头像、邮箱与回复状态（文字 + 图标），顶部汇总各状态人数。

## 最小示例

```go
people := []datetime.Attendee{{Name: "Sakiko", Status: datetime.AttendeeAccepted, Organizer: true}, {Name: "Mutsumi"}}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if i, ok := datetime.AttendeeList(c, people, datetime.AttendeeListOptions{Removable: true}).Removed(); ok {
		people = append(people[:i], people[i+1:]...)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func AttendeeList(c *ui.Context, attendees []Attendee, opts AttendeeListOptions) AttendeeListResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Removable` | `bool` | false | 除组织者外显示移除按钮。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

接受、待定、拒绝、未回复（各有图标与文字）、组织者、可选、空状态、禁用。

## 事件

`Removed()` 返回被移除者的下标，由调用方从切片中删除。

## 键盘操作

- Tab 到达各移除按钮，Enter 移除。

## 限制

- 回复状态不由组件修改。
