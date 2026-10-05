# 294 Typewriter 打字机

## 用途

逐字显示文本，可调速度、可跳过、完成时通知；文本增长时从当前位置继续。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	r := feedback.Typewriter(c, "The court assembles at dusk.", feedback.TypewriterOptions{Speed: 30, Skippable: true})
	if r.Changed() {
		// show the next message
	}
}
ui.Render(view, 400, 300, 1)
```

签名：`func Typewriter(c *ui.Context, text string, opts TypewriterOptions) ChoiceResult`

## 参数

`TypewriterOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Speed | float64 | 0（即每秒 40 字） | 速度；须为有限值，并产生 1 ns 到 time.Duration 上限之间的正间隔。非法值 panic。 |
| Skippable | bool | false | 显示“跳过”按钮。 |

返回 `ChoiceResult`，含 `Element`。

## 状态

- 按定时请求帧，不连续重绘。
- 整段文本作为名称朗读，不逐字朗读。
- 减少动画：直接显示全文。

## 事件

- `Changed()`：全文显示完的那一帧（含跳过）。

## 键盘操作

- “跳过”按钮用 Tab 到达，Enter 或 Space 跳过。

## 限制

- 文本被替换为不以旧文本开头的内容时从头开始。
