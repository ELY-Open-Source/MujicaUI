# 293 CountUp 数字过渡

## 用途

数值变化时从旧值过渡到新值。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.CountUp(c, 12840, feedback.CountUpOptions{Duration: 800 * time.Millisecond})
}
ui.Render(view, 400, 300, 1)
```

签名：`func CountUp(c *ui.Context, value float64, opts CountUpOptions) ui.Element`

## 参数

`CountUpOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Duration | time.Duration | 0（即 120 ms） | 过渡时长；为负 panic。 |
| Decimals | int | 0 | 小数位；为负 panic。 |
| Format | func(float64) string | nil | 自定义显示，优先于 Decimals。 |

value 须为有限值。首次显示直接为目标值。

## 状态

- 缓出过渡，过渡中途再变化时从当前显示值继续。
- 等宽数字防止抖动；朗读最终值。
- 减少动画：直接显示新值。过渡结束或屏外隐藏后停止重绘；恢复动画设置时保持已显示的最终值。

## 事件

- 无。

## 键盘操作

- 无。

## 限制

- 无。
