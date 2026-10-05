# 439 Countdown · 倒计时

## 用途

显示距截止时间的剩余时间，归零时报告一次并朗读提示。

## 最小示例

```go
deadline := time.Date(2026, time.December, 31, 23, 59, 59, 0, time.Local)
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if datetime.Countdown(c, deadline, datetime.CountdownOptions{}).Finished() {
		fmt.Println("happy new year")
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Countdown(c *ui.Context, deadline time.Time, opts CountdownOptions) CountdownResult`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Format` | `CountdownFormat` | `CountdownClock` | 时钟格式 01:02:03 或文字格式“1 小时 2 分 3 秒”。 |
| `Now` | `time.Time` | 零值（系统时钟） | 固定当前时间并停止刷新。 |
| `Label` | `string` | 内置文案 | 辅助技术读出的名称。 |

## 状态

计时中、归零（强调色）。

## 事件

`Finished()` 在每个截止时间归零的那一帧报告一次。

## 键盘操作

- 不可聚焦；归零时通过 `Announce` 通知读屏。

## 限制

- 仅在可见时每秒用 `c.After` 预约重绘，隐藏或归零后停止。
