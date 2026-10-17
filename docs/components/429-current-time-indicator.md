# 429 CurrentTimeIndicator · 当前时间线

## 用途

填满父元素，按当前时间在当天中的比例画一条横线，用于时间网格。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	ui.Box(c).Size(200, 240).Children(func() {
		datetime.CurrentTimeIndicator(c, datetime.CurrentTimeIndicatorOptions{Location: time.Local})
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func CurrentTimeIndicator(c *ui.Context, opts CurrentTimeIndicatorOptions) ui.Element`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Location` | `*time.Location` | 必填 | 切分日期与显示时间所用的时区；为 nil 时 panic，不做隐式 UTC 转换。 |
| `Day` | `Date` | 零值（今天） | 父元素自上而下代表的日期；不是今天时不显示。 |
| `Now` | `time.Time` | 零值（系统时钟） | 当前时间；固定后不再定时刷新，便于测试。 |
| `Hour12` | `bool` | false | 12 小时制显示时间。 |

## 状态

显示（Danger 色线与菱形）、不在当天时隐藏。

## 事件

无事件；读屏通过标签读出当前时间。

## 键盘操作

- 不可聚焦。

## 限制

- 仅在被构建（可见）时用 `c.After` 预约下一分钟的重绘；隐藏后不再预约，自动停止。
