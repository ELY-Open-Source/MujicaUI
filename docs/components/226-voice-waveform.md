# 226 VoiceWaveform 输入波形

## 用途

输入电平波形；静音提示。

## 最小示例

```go
levels := []float32{0.2, 0.6, 0.9, 0.4, 0.1}
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.VoiceWaveform(c, levels, mujica.VoiceWaveformOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`levels` 为 0–1 的电平样本，最新样本放在末尾。`Silence=0` 使用 0.05；`SilentFor=0` 使用 20 个样本；`Width=0` 使用 200 DIP。

## 状态

展示由输入数据驱动；边界值与警告状态带文字说明。

## 事件

展示组件，样本变化由调用方更新。

## 键盘操作

展示组件，沿用外层 Tab 顺序。

## 限制

电平采集由调用方提供；只绘制适合宽度的最新样本。
