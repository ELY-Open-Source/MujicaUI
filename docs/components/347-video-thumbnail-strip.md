# 347 VideoThumbnailStrip 缩略帧条

## 用途

横向缩略帧条：标出当前帧，点击或按键跳转。

## 最小示例

```go
frames := []media.StripFrame{{At: 0}, {At: 10 * time.Second}, {At: 20 * time.Second}}
current := 0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if media.VideoThumbnailStrip(c, &current, frames, media.VideoThumbnailStripOptions{}).Changed() {
		_ = frames[current].At // seek the player there
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func VideoThumbnailStrip(c *ui.Context, current *int, frames []StripFrame, opts VideoThumbnailStripOptions) ChoiceResult`

`StripFrame` 有 `Image`（可为 nil）和 `At`。`*current` 越界 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `FrameWidth` | `float32` | 0（即 96） | 每格宽度，高度按 16:9。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

- 只构建放得下的那几格，窗口跟随当前帧移动；10 万帧也只建十几个元素。
- 当前帧：强调边框加菱形；每格左下角显示时间。

## 事件

`Changed()` 在当前帧变化后为 true；调用方据此跳转播放器。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 聚焦帧条 |
| ← / → | 上一帧 / 下一帧 |
| Home / End | 第一帧 / 最后一帧 |

## 限制

- 不提供横向滚动条，用两侧按钮或键盘移动。
