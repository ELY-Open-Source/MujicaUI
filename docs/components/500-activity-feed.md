# 500 ActivityFeed 活动流

## 用途

按日分组显示“参与者 动作 对象 时间”，今天与昨天用文字，其余显示日期；只构建可见行。

## 最小示例

```go
events := []feedback.Activity{
	{ID: "1", Actor: "Lin", Action: "commented on", Object: "Quarterly report", Time: time.Now()},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.ActivityFeed(c, events, feedback.ActivityFeedOptions{}).Size(380, 280)
}
```

签名：`func ActivityFeed(c *ui.Context, items []Activity, opts ActivityFeedOptions) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `items` | `[]Activity` | 必填 | ID、Actor、Action、Object、Time；缺 ID、Actor、Action 或 Time panic。 |
| `Location` | `*time.Location` | `time.Local` | 按哪个时区划分日期。 |
| `Today` | `Date` | 时钟 | 固定“今天”。 |

## 状态

日期标题使用衬线字与旧金细线；无数据时显示空状态。

## 事件

无。

## 键盘操作

滚动容器按 MyGo 规则滚动。

## 限制

分组只在条数、首尾 ID 或时区变化时重算；中间条目原地修改需换新切片。
