# 211 AudioMessage 语音消息

## 用途

语音条：播放/暂停按钮、波形、播放进度与时长。播放由调用方驱动：组件只报告播放切换与跳转请求，调用方更新 Position 与 Playing。波形基于 MyGo `ui.SliderBase`，可拖动与键盘跳转。

## 最小示例

```go
r := mujica.AudioMessage(c, mujica.AudioMessageOptions{Duration: time.Minute, Position: pos})
if d, ok := r.Seek(); ok {
	pos = d
}
```

签名：`func AudioMessage(c *ui.Context, opts AudioMessageOptions) AudioMessageResult`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Duration` | `time.Duration` | 必填 | 时长，必须大于 0。 |
| `Position` | `time.Duration` | 0 | 当前进度，0 到 Duration 之间，否则 panic。 |
| `Playing` | `bool` | false | 按钮显示暂停（播放中）或播放。 |
| `Waveform` | `[]float32` | nil | 0 到 1 的响度采样，按宽度重采样为竖条；为空画平线；越界 panic。 |

## 状态

已播放的竖条为 AccentText，未播放为半透明 TextMuted，当前位置另有一条竖线，进度不只靠颜色表达；右侧显示 `m:ss / m:ss`。

## 事件

`Toggled()`：按下播放或暂停。`Seek()`：点击、拖动或按键请求的新位置。

## 键盘操作

Tab 到播放按钮（Enter 或 Space）与波形；波形上 Left/Right（Up/Down）按时长的 1% 移动，Home、End 到开头与结尾。

## 限制

- 不解码音频、不生成波形：采样由调用方提供。
