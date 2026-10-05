# 409 CallControls 通话控制

## 用途

提供静音、摄像头、屏幕共享和挂断操作，状态由调用方回传。

## 最小示例

```go
st := messaging.CallState{}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	messaging.CallControls(c, st, messaging.CallControlsOptions{
		OnMute:   func(on bool) { st.Muted = on },
		OnCamera: func(off bool) { st.CameraOff = off },
		OnShare:  func(on bool) { st.Sharing = on },
		OnHangUp: func() {}, // leave the call
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`CallState` 包含 `Muted`、`CameraOff`、`Sharing`。`OnMute`、`OnCamera`、`OnShare`、`OnHangUp` 均必填。摄像头回调参数表示是否关闭。`Disabled` 默认 false。

## 状态

默认、聚焦和悬停状态沿用主题、密度与文字缩放；适用操作在自身或祖先禁用时停止处理。

## 事件

每次操作调用对应回调一次；调用方回传更新后的 CallState。

## 键盘操作

Tab 进入工具栏；左右方向键移动；Enter / Space 激活。Cmd/Ctrl+D 静音，Cmd/Ctrl+E 摄像头，Cmd/Ctrl+Shift+S 共享，Cmd/Ctrl+Shift+H 挂断。

## 限制

通话、设备和屏幕采集由调用方处理；快捷键作用于当前窗口。
