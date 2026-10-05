# 410 VideoCallGrid 视频通话宫格

## 用途

按宫格显示调用方提供的视频帧，标记正在发言和静音的参与者。

## 最小示例

```go
people := []messaging.CallParticipant{
	{ID: "1", Name: "Sakiko", Speaking: true}, // Frame: the latest decoded frame
	{ID: "2", Name: "Uika", Muted: true},
}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.VideoCallGrid(c, people, messaging.VideoCallGridOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`people []CallParticipant` 的 ID 唯一、Name 必填；`Frame ui.ImageSource` 提供视频帧，`Image` 提供头像。`Muted`、`CameraOff`、`Speaking` 控制状态。`Columns` 默认 0，自动选择列数，必须 ≥ 0。

## 状态

状态由调用方数据控制；浅深主题、密度与文字缩放随窗口设置变化。适用插槽操作沿用控件焦点和禁用状态。

## 事件

纯展示，帧和参与者状态由调用方更新。

## 键盘操作

Tab 进入宫格；方向键移动参与者焦点并读取状态。

## 限制

视频解码、采集与帧更新由调用方处理，更新后调用 c.Invalidate()。
