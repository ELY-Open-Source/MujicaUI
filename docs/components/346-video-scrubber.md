# 346 VideoScrubber 视频时间轴

## 用途

可拖动的时间轴：显示已缓冲区间，悬停显示时间和缩略图插槽。

## 最小示例

```go
pos := 30 * time.Second
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.VideoScrubber(c, &pos, mujica.VideoScrubberOptions{
		Duration: 2 * time.Minute,
		Buffered: []mujica.MediaRange{{Start: 0, End: time.Minute}},
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func VideoScrubber(c *ui.Context, pos *time.Duration, opts VideoScrubberOptions) ChoiceResult`

`*pos` 必须在 0 到 `Duration` 之间，否则 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Duration` | `time.Duration` | 必填 | 时间轴长度，须大于 0。 |
| `Buffered` | `[]MediaRange` | nil | 已缓冲区间，以浅色填充。 |
| `Preview` | `func(at time.Duration)` | nil | 悬停时在插槽里构建该时间的缩略图。 |
| `PreviewWidth / PreviewHeight` | `float32` | 160×90 | 插槽尺寸。 |
| `Step` | `time.Duration` | 5 秒 | 方向键每次移动的时长。 |
| `Disabled` | `bool` | false | 禁用。 |
| `Label` | `string` | "Timeline" | 无障碍名称；以滑块角色报告范围与数值。 |

## 状态

默认、悬停（显示时间气泡）、拖动、聚焦（菱形手柄外圈）、禁用。

## 事件

`Changed()` 在拖动、点击或按键改变 `*pos` 后为 true。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| ← / ↓ | 后退一个 Step |
| → / ↑ | 前进一个 Step |
| Home / End | 开头 / 结尾 |

## 限制

- 位置按毫秒取整。
- 悬停插槽只在指针位于时间轴上时显示，键盘操作时不显示。
