# 300 BlinkHighlight 变化高亮

## 用途

数值变化时单次高亮：上升 Success、下降 Danger，并留下方向箭头。

## 最小示例

```go
price := 101.5
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.BlinkHighlight(c, price, feedback.BlinkHighlightOptions{})
}
ui.Render(view, 400, 300, 1)
```

签名：`func BlinkHighlight(c *ui.Context, value float64, opts BlinkHighlightOptions) ui.Element`

## 参数

`BlinkHighlightOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Format | func(float64) string | 两位小数 | 显示格式。 |
| Duration | time.Duration | 0（即 120 ms） | 高亮淡出时长；为负 panic。 |

value 须为有限值；首次显示不高亮。

## 状态

- 背景从 22% 透明度淡出一次，结束后停止重绘。
- 方向箭头与朗读的“上升/下降”在变化后一直保留，不只靠颜色。
- 减少动画：不闪烁，只显示静态箭头。

## 事件

- 无。

## 键盘操作

- 无。

## 限制

- 开启减少动画会结束当前高亮；恢复动画时保持静态方向标记。
