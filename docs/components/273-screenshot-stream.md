# 273 ScreenshotStream Agent 画面流

## 用途

查看 Agent 的画面帧序列：默认跟随最新帧，可暂停并逐帧前后查看。帧由调用方提供。

## 最小示例

```go
var st mujica.ScreenshotStreamState
var frames []mujica.ScreenshotFrame // appended as the agent captures them
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ScreenshotStream(c, &st, frames, mujica.ScreenshotStreamOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func ScreenshotStream(c *ui.Context, st *ScreenshotStreamState, frames []ScreenshotFrame, opts ScreenshotStreamOptions) ChoiceResult`。`st` 为 nil、帧没有图像或 `Height` 为负时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Height` | `float32` | 280 | 画面区高度。 |
| `Disabled` | `bool` | `false` | 禁止暂停与逐帧。 |

## 状态

- `ScreenshotFrame` 字段：`Image`（`ui.ImageSource`）、`Caption`、`At`。
- `ScreenshotStreamState` 字段：`Index`、`Paused`；未暂停时 `Index` 总是最后一帧。
- 逐帧操作会自动暂停；“回到实时”恢复跟随。
- 没有帧时显示“等待画面”。

## 事件

- `Changed()`：用户暂停、恢复或移动了帧。

## 键盘操作

- Tab 到画面区：← → 逐帧，Space 暂停或恢复；上一帧、下一帧按钮同样可用。

## 限制

- 不解码视频，帧图像由调用方提供。
