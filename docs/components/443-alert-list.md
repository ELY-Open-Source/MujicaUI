# 443 AlertList · 告警列表

## 用途

按级别、来源与发生时间展示告警；报告确认和静默请求。

## 最小示例

```go
alerts := []mujica.AlertItem{{ID: "cpu", Title: "CPU usage high", Level: mujica.AlertWarning, Source: "worker", At: time.Date(2026, 10, 5, 9, 0, 0, 0, time.UTC)}}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	r := mujica.AlertList(c, alerts, mujica.AlertListOptions{})
	if r.Ack != "" {
		fmt.Println("ack", r.Ack)
	}
	if r.Silence != "" {
		fmt.Println("silence", r.Silence)
	}
}
ui.Render(view, 720, 480, 1)
```

## 参数

MaxHeight 为列表最大高度，默认 360 DIP；Disabled 禁用操作。AlertItem.ID 唯一，Title 必填，Level 为有效枚举；Acked、Silence 展示处理状态。

## 状态

未处理、已确认、静默中、空列表、禁用。

## 事件

Ack、Silence 返回告警 ID；调用方更新 Acked 或 Silence。

## 键盘操作

Tab 聚焦按钮；Enter 或 Space 激活。

## 限制

静默请求固定为一小时；告警数据由调用方管理。
