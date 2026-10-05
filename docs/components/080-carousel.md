# 080 Carousel · 走马灯

## 用途

逐张展示内容，提供上一项、下一项与位置指示。默认不自动播放；开启自动播放后提供暂停按钮，并在指针悬停或焦点位于组件内时暂停。

## 最小示例

```go
index := 0
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	data.Carousel(c, &index, 3, data.CarouselOptions{Label: "Windows", Autoplay: true}, func(i int) {
		ui.Text(c, fmt.Sprintf("Rose window %d", i+1)).Padding(24)
	}).Element.Size(360, 220)
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func Carousel(c *ui.Context, index *int, n int, opts CarouselOptions, slide func(i int)) ChoiceResult`。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | 空 | 无障碍名称。 |
| `Autoplay` | `bool` | false | 每隔 `Interval` 前进一项（循环）。 |
| `Interval` | `time.Duration` | 5s | 自动播放间隔；负数 panic。 |
| `Loop` | `bool` | false | 首尾相接；自动播放总是循环。 |

## 状态

默认、按钮悬停/按下/聚焦、首尾禁用（未循环时）、当前位置（实心菱形）、自动播放中、已暂停（显示继续按钮）。

## 事件

`ChoiceResult.Changed()`：当前项变化，`*index` 已更新；同时向读屏软件播报“第 N 项，共 M 项”。

## 键盘操作

- 焦点位于组件内时：←/→ 上一项、下一项；Home/End 第一项、最后一项。
- 位置指示为一个焦点组：方向键在指示点间移动，当前项随焦点切换；此时方向键不再由整个组件处理。
- Tab 依次到达上一项、下一项、位置指示与暂停按钮。

## 限制

- 切换不做滑动动画，直接显示新的一项。
- 不自动播放时没有计时器；自动播放只在可见帧中通过 `Context.After` 预约下一帧。
- `*index` 超出范围时会被夹到有效范围（记录 debug 日志），例如条目数减少后。
