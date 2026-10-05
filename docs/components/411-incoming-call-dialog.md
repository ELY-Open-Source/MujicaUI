# 411 IncomingCallDialog 来电对话框

## 用途

显示语音或视频来电，提供接听、拒绝与可配置的初始焦点。

## 最小示例

```go
ringing := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.IncomingCallDialog(c, &ringing, messaging.IncomingCall{Caller: "Sakiko", Video: true}, messaging.IncomingCallOptions{
		OnAccept:  func(video bool) {}, // join, with the camera on when video
		OnDecline: func() {},
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`open *bool` 控制显示；`IncomingCall.Caller` 必填，`Image` 为头像，`Video` 默认 false。`OnAccept func(video bool)` 与 `OnDecline func()` 必填。`FocusAccept` 默认 false，初始聚焦拒绝；`Disabled` 默认 false，禁用时关闭。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

接听调用 `OnAccept`，视频接听参数为 true；拒绝和 Escape 调用 `OnDecline`。完成后设置 `*open=false`。

## 键盘操作

默认焦点为拒绝；Tab 循环访问操作；Enter / Space 激活；Escape 拒绝；关闭后恢复焦点。

## 限制

来电数据和通话连接由调用方处理；点击遮罩保持来电显示；禁用关闭时回调保持静默。
