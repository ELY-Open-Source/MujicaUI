# 349 AudioWaveform 音频波形

## 用途

绘制调用方提供的采样波形，覆盖播放进度，可拖动选择区间。

## 最小示例

```go
samples := make([]float32, 48000) // from your decoder, -1 to 1
pos, sel := time.Duration(0), media.MediaRange{}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.AudioWaveform(c, samples, media.AudioWaveformOptions{
		Duration: time.Second, Position: pos, Selection: &sel,
		OnSeek: func(d time.Duration) { pos = d },
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func AudioWaveform(c *ui.Context, samples []float32, opts AudioWaveformOptions) ChoiceResult`

`samples` 取值 -1 到 1。数据就绪时 `Duration` 须大于 0，`OnSeek` 必填。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Duration` | `time.Duration` | — | 音频长度。 |
| `Position` | `time.Duration` | 0 | 已播放部分。 |
| `OnSeek` | `func(time.Duration)` | — | 点击或方向键跳转。 |
| `Selection` | `*MediaRange` | nil | 非 nil 时可拖动或 Shift+方向键选择区间。 |
| `Height` | `float32` | 0（即 96） | 高度。 |
| `Status / Error` | `DataStatus / string` | DataReady | 加载中、失败（须给 Error）与无数据时显示同尺寸占位。 |
| `Label` | `string` | "Waveform" | 无障碍名称。 |

## 状态

- 已播放部分用强调色，其余用说明色，并有播放线和菱形标记，不只靠颜色区分。
- 选区用 Selection 底色加两条边线。

## 事件

`Changed()` 在选区变化后为 true；跳转通过 `OnSeek` 报告。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| ← / → | 跳转 ±5 秒 |
| Shift+← / → | 从播放位置开始或调整选区终点 |
| Escape | 清除选区 |

## 限制

- 峰值按像素列计算，只在数据或宽度变化时重算，不会每帧计算。
- 数据以首元素地址和长度识别；原地修改采样请换一个切片。
