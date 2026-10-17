# 350 AudioSpectrum 音频频谱

## 用途

把调用方提供的频带数值画成柱状频谱；组件自己不做动画，隐藏时不重绘。

## 最小示例

```go
bands := []float32{0.2, 0.6, 0.9, 0.4, 0.1} // set by your analyser, then c.Invalidate()
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.AudioSpectrum(c, bands, media.AudioSpectrumOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func AudioSpectrum(c *ui.Context, bands []float32, opts AudioSpectrumOptions) ui.Element`

每个频带取 0 到 1，超出范围 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Height` | `float32` | 0（即 120） | 高度。 |
| `BandLabel` | `func(i int) string` | nil | 读数中频带的名称，如“1 kHz”；nil 时编号。 |
| `Status / Error` | `DataStatus / string` | DataReady | 加载、失败与无数据占位，尺寸相同。 |
| `Label` | `string` | "Spectrum" | 无障碍名称。 |

## 状态

聚焦后用方向键选中的频带以正文色显示，带虚线框和读数，不只靠颜色。

## 事件

无事件；方向键读数会通过 `Announce` 朗读。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 聚焦频谱 |
| ← / → | 上一个 / 下一个频带 |
| Home / End | 第一个 / 最后一个频带 |

## 限制

- 刷新节奏由调用方决定：更新数据后调用 `c.Invalidate()`。
