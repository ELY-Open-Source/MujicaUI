# 420 SnoozePicker 稍后提醒

## 用途

选择快捷提醒时间，或通过日期时间输入设定自定义提醒。

## 最小示例

```go
var until time.Time
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if messaging.SnoozePicker(c, &until, messaging.SnoozePickerOptions{}).Changed() {
		_ = until // hide the mail until then
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`until *time.Time` 保存提醒时间。`Location` 默认 `time.Local`；`MorningHour` 默认 9，允许 0–23，其中 0 使用默认值；`Now` 默认使用当前帧时间，可传固定时间；`Disabled` 默认 false。`SnoozeChoices` 的 morning 参数直接接受 0–23。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

返回 `ChoiceResult`；`Changed()` 表示提醒时间改变。

## 键盘操作

Tab 进入；上下方向键选择快捷项；Enter / Space 激活；自定义时间沿用 DateTimePicker 的输入、确认和取消操作。

## 限制

组件返回时间，提醒调度与邮件隐藏由调用方处理；本地日期按日历构造以保留夏令时后的早晨时刻。
