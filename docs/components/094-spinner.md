# 094 Spinner 加载标记

## 用途

表示未知长度的加载：四个菱形绕中心旋转。可调尺寸，可显示说明文字。Button 的加载状态也使用它。标记基于 MyGo `ui.Spinner` 元素重绘，读屏视为不确定进度。

## 最小示例

```go
core.Spinner(c, core.SpinnerOptions{Size: 20, Label: "Syncing", ShowLabel: true})
```

签名：`func Spinner(c *ui.Context, opts SpinnerOptions) ui.Element`

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Size` | `float32` | 16 | 标记边长（DIP）；负数、NaN 或无穷 panic。 |
| `Label` | `string` | 内置“Loading / 加载中” | 无障碍名称；`ShowLabel` 时也显示在旁边。 |
| `ShowLabel` | `bool` | `false` | 在标记旁显示说明。 |
| `Color` | `ui.Color` | `AccentText` | 标记颜色。 |

## 状态

- 只在被绘制时请求重绘：不再构建（隐藏或停止加载）或滚出视口后动画立即停止，不留定时器。
- 减少动画模式下静止显示，保留名称与说明。

## 事件

无。

## 键盘操作

不可聚焦。

## 限制

- 每帧重绘，而非 MyGo 自带 spinner 的分步重绘；只在显示时产生开销。
