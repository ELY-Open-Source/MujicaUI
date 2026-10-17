# 292 Shimmer 占位扫光

## 用途

加载中显示带扫光的占位块，加载完成后构建内容。

## 最小示例

```go
loading := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.Shimmer(c, loading, feedback.ShimmerOptions{Height: 48}, func() {
		ui.Text(c, "Report ready")
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func Shimmer(c *ui.Context, loading bool, opts ShimmerOptions, content func()) ui.Element`

## 参数

`ShimmerOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Width | float32 | 0（填满） | 宽度；为负 panic。 |
| Height | float32 | 0（即 ControlHeight） | 高度。 |
| Label | string | 内置“加载中” | 加载时朗读的文字。 |

`loading` 为 false 时只构建 content。

## 状态

- 扫光为一条半透明 Surface 色带，位置吸附像素，只在绘制时动画。
- 减少动画：静态填充。
- 加载结束后不再重绘。

## 事件

- 无。

## 键盘操作

- 无。

## 限制

- 只用于加载中，不作装饰。
