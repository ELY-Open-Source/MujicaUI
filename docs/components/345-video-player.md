# 345 VideoPlayer 视频播放器

## 用途

播放器外壳：画面由调用方提供，组件绘制状态并通过回调报告播放、暂停、跳转、倍速、字幕、全屏和画中画。

## 最小示例

```go
st := media.VideoState{Duration: 4 * time.Minute, Rate: 1}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.VideoPlayer(c, st, media.VideoPlayerOptions{
		OnPlay:  func() { st.Playing = true },  // start your decoder
		OnPause: func() { st.Playing = false }, // pause it
		OnSeek:  func(d time.Duration) { st.Position = d },
		OnRate:  func(r float64) { st.Rate = r },
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func VideoPlayer(c *ui.Context, st VideoState, opts VideoPlayerOptions) ui.Element`

`VideoState`：`Frame`（当前画面，nil 显示“暂无视频”）、`Position`、`Duration`、`Buffered`、`Playing`、`Rate`、`Subtitle`、`SubtitlesOn`、`Fullscreen`、`PiP`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `OnPlay / OnPause / OnSeek` | `func` | 必填 | 播放、暂停、跳转。 |
| `OnRate` | `func(float64)` | nil | 倍速；nil 时不显示倍速选择。 |
| `OnSubtitles / OnFullscreen / OnPiP` | `func(bool)` | nil | 开关回调；nil 时隐藏对应按钮。 |
| `Disabled` | `bool` | false | 禁用全部控件。 |
| `Label` | `string` | "Video" | 无障碍名称。 |

## 状态

- 画面区 16:9，黑底；字幕开启且有字幕时在底部显示。
- 字幕、画中画、全屏按钮开启时有选中底色和下划线。

## 事件

全部通过回调报告；组件不改动 `st`。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 先聚焦画面，再到进度条和各按钮 |
| Space / K | 画面聚焦时播放或暂停 |
| ← / → | 画面聚焦时后退 / 前进 5 秒 |
| F | 切换全屏 |
| C | 切换字幕 |

## 限制

- 组件不解码视频，也不实现全屏或画中画窗口，只报告用户的请求。
- 帧率取决于调用方更新 `Frame` 并调用 `c.Invalidate()` 的频率。
