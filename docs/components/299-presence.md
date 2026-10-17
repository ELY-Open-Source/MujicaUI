# 299 Presence 进出场

## 用途

内容随 show 进出场：淡入淡出、滑动或折叠，退出完成后移除。

## 最小示例

```go
show := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.Presence(c, show, feedback.PresenceOptions{Kind: feedback.PresenceSlide}, func() {
		ui.Text(c, "Saved")
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func Presence(c *ui.Context, show bool, opts PresenceOptions, content func()) ui.Element`

## 参数

`PresenceOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Kind | PresenceKind | PresenceFade | PresenceFade、PresenceSlide、PresenceCollapse；未知值 panic。 |
| Duration | time.Duration | 0（即 PopupDuration） | 时长；为负 panic。 |

基于 MyGo `Transition` 的 Enter 与 Exit：退出时保留上一帧的副本直到动画结束，不占布局、不接收指针与焦点。

## 状态

- 减少动画：直接出现与消失。动画结束后停止重绘。

## 事件

- 无。

## 键盘操作

- 退出中的内容不接收焦点。

## 限制

- 不提供缩放：MyGo `ui.Motion` 只有位移、透明度与折叠（ui/transition.go:40）。
- 父元素在上一帧不存在时，内容随父元素一起出现而不进场。
