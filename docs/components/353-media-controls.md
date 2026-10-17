# 353 MediaControls 媒体控制

## 用途

上一首、播放/暂停、下一首、循环、随机，一组可用方向键切换焦点的按钮。

## 最小示例

```go
st := media.MediaControlsState{}
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.MediaControls(c, st, media.MediaControlsOptions{
		OnPrevious: func() {}, OnNext: func() {},
		OnPlay: func() { st.Playing = true }, OnPause: func() { st.Playing = false },
		OnLoop: func(m media.LoopMode) { st.Loop = m }, OnShuffle: func(on bool) { st.Shuffle = on },
	})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func MediaControls(c *ui.Context, st MediaControlsState, opts MediaControlsOptions) ui.Element`

`MediaControlsState`：`Playing`、`Loop`（`LoopOff`、`LoopAll`、`LoopOne`）、`Shuffle`、`NoPrevious`、`NoNext`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `OnPrevious / OnPlay / OnPause / OnNext` | `func()` | 必填 | 按钮回调。 |
| `OnLoop` | `func(LoopMode)` | 必填 | 收到下一个循环模式：不循环 → 列表 → 单曲。 |
| `OnShuffle` | `func(bool)` | 必填 | 随机开关。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

循环与随机开启时有选中底色和下划线，单曲循环换用带“1”的图标；循环按钮的名称说明当前模式。

## 事件

全部通过回调报告。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 进入按钮组 |
| ← / → | 在按钮间移动焦点 |
| Enter / Space | 按下 |

## 限制

- 无。
