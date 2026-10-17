# 095 Skeleton 骨架屏

## 用途

内容加载前占位，形状为文字行、圆形或矩形。尺寸固定，加载完成后替换为真实内容时不跳动。

## 最小示例

```go
ui.Row(c).Padding(16).Gap(12).Children(func() {
	feedback.Skeleton(c, feedback.SkeletonOptions{Shape: feedback.SkeletonCircle})
	ui.Column(c).Grow(1).Children(func() {
		feedback.Skeleton(c, feedback.SkeletonOptions{Lines: 2})
	})
})
```

签名：`func Skeleton(c *ui.Context, opts SkeletonOptions) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Shape` | `SkeletonShape` | `SkeletonText` | `SkeletonText`、`SkeletonCircle`、`SkeletonRect`。 |
| `Width` | `float32` | 填满父级；圆形 40 | 宽度；圆形为直径。 |
| `Height` | `float32` | 当前密度控件高度 | 矩形高度；文字形状由行数决定（每行 21 DIP）。 |
| `Lines` | `int` | 1 | 文字行数，多行时最后一行为 60% 宽。 |

负值、NaN、无穷或未知形状 panic。文字行高随系统文本缩放。

## 状态

明暗缓慢脉动（1.6 秒周期），只重绘不重建视图；减少动画模式使用静态填充。

## 事件

无。

## 键盘操作

不可聚焦，不进入无障碍树。

## 限制

- 骨架本身不向读屏报告加载；请在容器上配合 Spinner 或 LoadingOverlay 的说明文字。
