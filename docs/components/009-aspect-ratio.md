# 009 AspectRatio 宽高比

## 用途

在父容器的宽高约束内放置一个保持指定宽高比的框，尽可能大并居中；子内容在框内居中，不被拉伸，超出部分裁剪。适合预览、视频占位、封面图。

## 最小示例

```go
ui.Box(c).Size(320, 240).Children(func() {
	layout.AspectRatio(c, layout.AspectRatioOptions{Ratio: 16.0 / 9}, func() {
		ui.Text(c, "16 : 9")
	})
})
```

签名：`func AspectRatio(c *ui.Context, opts AspectRatioOptions, child func()) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Ratio` | `float32` | 1 | 宽 / 高；负数 panic。 |
| `child` | `func()` | 必填 | 框内内容；nil 时 panic。需要铺满框时让子元素 `Fill()`。 |

## 状态

布局组件，无交互状态。返回的元素占满父级内容区，比例框位于其中央。

## 事件

无。

## 键盘操作

不获取焦点。

## 限制

- 父级高度不确定（随内容增长）时只按宽度计算。
- 框尺寸依据上一帧的父级尺寸计算，窗口缩放后组件请求下一帧完成调整（MyGo 布局的 `AspectRatio` 由已定的宽推出高，或由固定的高推出宽，见 `ui/layout.go:345` 与 `:364`；但不会在父级宽高同时约束时取两者中较小的一边，因此组件自行计算）。
