# 251 DateSeparator · 日期分隔

## 用途

在消息流中标出新的一天：今天、昨天，其余显示本地化的完整日期。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	chat.DateSeparator(c, core.DateOf(time.Now()), chat.DateSeparatorOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func DateSeparator(c *ui.Context, day Date, opts DateSeparatorOptions) ui.Element`。`day` 必须是有效日期。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Today` | `Date` | 零值（窗口时钟） | 判断“今天/昨天”的当天。 |

## 状态

细线 + 旧金色菱形 + 说明文字；无交互状态。

## 事件

无。

## 键盘操作

不获取焦点；辅助技术读出日期文字。

## 限制

- 日期格式沿用 Calendar 的文案（`calendar.day_label`），en-US 为“Tuesday, September 1, 2026”，zh-CN 为“2026年9月1日 星期二”。
