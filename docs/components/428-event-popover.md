# 428 EventPopover · 事件浮层

## 用途

在锚点旁显示事件摘要：时间、地点、重复、参与者与提醒，并提供编辑与删除入口。

## 最小示例

```go
loc := time.Local
ev := exampleEvents(loc)[0]
open := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chip := datetime.EventChip(c, ev, datetime.EventChipOptions{Location: loc})
	if chip.Clicked() {
		open = true
	}
	r := datetime.EventPopover(c, chip, &open, ev, datetime.EventPopoverOptions{Location: loc})
	if r.Delete() {
		fmt.Println("delete", ev.ID)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func EventPopover(c *ui.Context, anchor ui.Element, open *bool, ev CalendarEvent, opts EventPopoverOptions) EventPopoverResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Location` | `*time.Location` | 必填 | 切分日期与显示时间所用的时区；为 nil 时 panic，不做隐式 UTC 转换。 |
| `Hour12` | `bool` | false | 12 小时制显示时间。 |
| `ReadOnly` | `bool` | false | 隐藏编辑与删除。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |

## 状态

打开（160 ms 淡入，减少动画时直接显示）、只读、禁用（自动关闭）。

## 事件

`Edit()`、`Delete()` 报告点击并关闭浮层；删除只报告，确认由调用方决定。

## 键盘操作

- Tab 到达关闭、删除与编辑；Escape 或点击外部关闭。

## 限制

- 浮层尺寸受窗口限制，内容超出时滚动。
