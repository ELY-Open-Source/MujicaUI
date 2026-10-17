# 426 EventChip · 事件块

## 用途

显示单个事件：左侧色条、标题与时间，超出时截断为省略号。

## 最小示例

```go
loc := time.Local
ev := exampleEvents(loc)[0]
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if datetime.EventChip(c, ev, datetime.EventChipOptions{Location: loc}).Clicked() {
		fmt.Println("open", ev.ID)
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func EventChip(c *ui.Context, ev CalendarEvent, opts EventChipOptions) ui.Element`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Location` | `*time.Location` | 必填 | 切分日期与显示时间所用的时区；为 nil 时 panic，不做隐式 UTC 转换。 |
| `Hour12` | `bool` | false | 12 小时制显示时间。 |
| `Inline` | `bool` | false | 开始时间与标题放在一行。 |
| `Selected` | `bool` | false | 选中：选中底色 + 内嵌框。 |
| `Disabled` | `bool` | false | 禁用：不响应点击、拖动与快捷键，弹层随之关闭。 |

## 状态

默认、悬停与按下（色调加深）、聚焦、选中、禁用。

## 事件

返回的元素可用 `.Clicked()` 判断激活。

## 键盘操作

- Tab 聚焦，Enter 或 Space 激活。

## 限制

- 事件颜色来自主题：AccentText、Info、Success、Warning、Danger、TextMuted，均在两套主题下对 Surface 达到 3:1；Ornament 在浅色主题下仅 2.8:1，不纳入。
