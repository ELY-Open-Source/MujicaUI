# 099 NotificationCenter 通知中心

## 用途

应用内消息列表：未读标记、已读/未读切换、单条清除、全部已读、全部清除和每条消息的操作。消息列表与持久化由调用方持有；组件就地修改调用方的切片，并报告发生了什么，供调用方保存。

## 最小示例

```go
r := feedback.NotificationCenter(c, &items, feedback.NotificationCenterOptions{})
if r.Action == "a" {
	// open the build log
}
```

签名：`func NotificationCenter(c *ui.Context, items *[]Notification, opts NotificationCenterOptions) NotificationCenterResult`

## 参数

`Notification`：

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `ID` | `string` | 必填 | 列表内唯一；为空或重复时 panic。 |
| `Title` | `string` | `""` | 标题，未读时加粗。 |
| `Body` | `string` | `""` | 正文。 |
| `Time` | `string` | `""` | 按原样显示的时间，如“5 min ago”。 |
| `Read` | `bool` | `false` | 是否已读。 |
| `Action` | `string` | `""` | 操作按钮文字。 |

`NotificationCenterOptions.Title`：面板标题，默认“Notifications / 通知”。`items` 为 nil 时 panic。

## 状态

- 未读：左侧酒红菱形、粗体标题、无障碍名称前缀“Unread”。
- 行悬停加浅底；按钮有悬停、按下与焦点状态；无未读或无消息时头部按钮禁用。
- 空列表显示 EmptyState。
- 面板标题下为旧金细线。列表超出容器高度时滚动。

## 事件

`NotificationCenterResult`：`Action`、`Toggled`（已读状态切换）、`Opened`（点击正文后标为已读）、`Removed` 为对应 ID；`MarkedAllRead`、`ClearedAll` 为布尔值。

## 键盘操作

列表是一个 Tab 停靠点；方向键在各条消息及其按钮之间移动，Home/End 到首尾。焦点在消息正文上时 Enter 或 Space 标为已读；Delete 或 Backspace 清除当前消息。

## 限制

- 不分页，也不虚拟化；面向几十到几百条的应用内消息。
- 时间文字由调用方格式化。
