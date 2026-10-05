# 352 PlaybackSpeedControl 倍速选择

## 用途

在 0.5× 到 2× 之间选择播放速度。

## 最小示例

```go
rate := 1.0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.PlaybackSpeedControl(c, &rate, mujica.PlaybackSpeedControlOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func PlaybackSpeedControl(c *ui.Context, rate *float64, opts PlaybackSpeedControlOptions) ChoiceResult`

`*rate` 必须是 `Rates` 中的一项，否则 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Rates` | `[]float64` | 0.5, 0.75, 1, 1.25, 1.5, 2 | 升序的正数。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

基于 Select，状态与键盘同 Select。

## 事件

`Changed()` 在选择新倍速后为 true。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Enter / Space | 打开列表 |
| ↑ / ↓ | 移动 |
| Enter | 确认 |
| Escape | 关闭 |

## 限制

- 无。
