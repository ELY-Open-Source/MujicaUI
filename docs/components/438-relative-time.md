# 438 RelativeTime · 相对时间

## 用途

把时间显示为“3 分钟前”“2 天后”，并在提示与标签中给出精确时间。

## 最小示例

```go
posted := time.Date(2026, time.October, 5, 9, 0, 0, 0, time.Local)
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	datetime.RelativeTime(c, posted, datetime.RelativeTimeOptions{Location: time.Local})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func RelativeTime(c *ui.Context, t time.Time, opts RelativeTimeOptions) ui.Element`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Location` | `*time.Location` | 必填 | 切分日期与显示时间所用的时区；为 nil 时 panic，不做隐式 UTC 转换。 |
| `Now` | `time.Time` | 零值（系统时钟） | 固定当前时间并停止刷新。 |
| `Hour12` | `bool` | false | 12 小时制显示时间。 |

## 状态

刚刚、分钟、小时、天、月、年。

## 事件

无事件。

## 键盘操作

- 不可聚焦；精确时间写在标签与提示中。

## 限制

- 刷新频率随距离变化：2 分钟内每秒、2 小时内每分钟、更远每小时；仅在可见时用 `c.After` 预约，隐藏即停止。
