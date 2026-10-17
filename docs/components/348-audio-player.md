# 348 AudioPlayer 音频播放器

## 用途

音频播放外壳：播放、暂停、进度与时长，播放本身由调用方驱动。

## 最小示例

```go
st := media.AudioState{Title: "Overture", Duration: 3 * time.Minute}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.AudioPlayer(c, st, media.AudioPlayerOptions{
		OnPlay:  func() { st.Playing = true },
		OnPause: func() { st.Playing = false },
		OnSeek:  func(d time.Duration) { st.Position = d },
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func AudioPlayer(c *ui.Context, st AudioState, opts AudioPlayerOptions) ui.Element`

`AudioState`：`Title`、`Position`、`Duration`、`Playing`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `OnPlay / OnPause / OnSeek` | `func` | 必填 | 缺少任一回调时 panic。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

播放中显示暂停按钮，暂停时显示播放按钮；`Duration` 为 0 时不显示进度条。

## 事件

通过回调报告，组件不改动 `st`。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 播放按钮、进度条 |
| Space | 组件内任意位置播放或暂停（按钮聚焦时由按钮处理） |
| ← / → | 进度条聚焦时跳转 5 秒 |

## 限制

- 组件不播放声音。
