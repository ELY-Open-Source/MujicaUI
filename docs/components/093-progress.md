# 093 Progress 进度

## 用途

显示任务进度（确定或不确定，线形或环形），以及计量值（如磁盘占用）及其阈值。归并 MyGo 的 Progress bar 与 Meter：`Progress` 为进度模式，`Meter` 为计量模式。两者分别基于 `ui.Progress` 与 `ui.Meter` 元素重绘，保留其无障碍角色与数值。

## 最小示例

```go
feedback.Progress(c, 0.62, feedback.ProgressOptions{Label: "Uploading", ShowValue: true})
feedback.Progress(c, 0, feedback.ProgressOptions{Shape: feedback.ProgressRing, Indeterminate: true, Label: "Connecting"})
feedback.Meter(c, 93, feedback.MeterOptions{
	Max: 100, Label: "Disk", ShowValue: true,
	Levels: &ui.MeterLevels{Warning: 70, Critical: 90},
})
```

签名：

- `func Progress(c *ui.Context, value float64, opts ProgressOptions) ui.Element`
- `func Meter(c *ui.Context, value float64, opts MeterOptions) ui.Element`

## 参数

`ProgressOptions`：

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `float64` | 必填 | 0 到 1；超出或 NaN 时 panic（不确定模式忽略）。 |
| `Shape` | `ProgressShape` | `ProgressLinear` | `ProgressLinear` 或 `ProgressRing`；其他值 panic。 |
| `Indeterminate` | `bool` | `false` | 未知长度的进度。 |
| `Label` | `string` | `""` | 显示在旁边并作为无障碍名称；为空时名称为“Loading”。 |
| `ShowValue` | `bool` | `false` | 显示百分比。 |
| `Size` | `float32` | 环形 48、线形填满 | 环形直径或线形宽度；负数、NaN 或无穷 panic。 |

`MeterOptions`：

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `value` | `float64` | 必填 | 须在范围内，否则 panic。 |
| `Min`, `Max` | `float64` | 0, 1（两者都为 0 时） | 范围；非有限值或 `Max <= Min` panic。 |
| `Levels` | `*ui.MeterLevels` | `nil` | 阈值：达到 `Warning` 为偏高、达到 `Critical` 为严重；`Critical < Warning` 时低值为坏（如电量）。须在范围内。 |
| `Label` | `string` | `""` | 名称。 |
| `ShowValue` | `bool` | `false` | 显示“值 / 最大值”。 |
| `Width` | `float32` | 填满 | 宽度；负数、NaN 或无穷 panic。 |

## 状态

- 确定进度变化时以 120 ms 过渡，减少动画时直接跳到新值。
- 不确定进度只在绘制时移动（只重绘、不重建视图）；减少动画时显示满格静态填充（对比度不低于 3:1）和“In progress”文字；切换到减少动画时进行中的过渡直接跳到目标值。
- 计量条按级别使用 Success / Warning / Danger 颜色，同时显示“Normal / High / Critical / Low”文字与阈值刻度，并把级别写入无障碍名称。

## 事件

无。

## 键盘操作

不可聚焦；读屏读取名称与数值。

## 限制

- 百分比按缩放后的文字尺寸测量，放不进环内时显示在环旁边。
- MyGo 不提供自定义无障碍数值的公开接口，因此必须基于 `ui.Progress` / `ui.Meter` 元素绘制。
