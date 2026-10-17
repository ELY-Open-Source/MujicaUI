# 298 LayoutTransition 布局过渡

## 用途

内容的尺寸与位置变化时从旧布局过渡到新布局（FLIP）。

## 最小示例

```go
open := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.LayoutTransition(c, feedback.LayoutTransitionOptions{Clip: true}, func() {
		ui.Text(c, "Program")
		if open {
			ui.Text(c, "Overture · Masquerade · Finale")
		}
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func LayoutTransition(c *ui.Context, opts LayoutTransitionOptions, content func()) ui.Element`

## 参数

`LayoutTransitionOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Duration | time.Duration | 0（即 120 ms） | 时长；为负 panic。 |
| Row | bool | false | 横向排列内容。 |
| Clip | bool | false | 尺寸变化时裁剪溢出内容。 |

基于 MyGo `Transition` 的 Position 与 Size。

## 状态

- 尺寸过渡中内容按每一帧的尺寸重新布局。
- 窗口缩放引起的变化不动画；减少动画时直接到位。

## 事件

- 无。

## 键盘操作

- 无。

## 限制

- 被同一变化推动的兄弟元素需要各自的 LayoutTransition 才会同步移动。
