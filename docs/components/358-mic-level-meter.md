# 358 MicLevelMeter 麦克风电平

## 用途

分段电平条，带峰值保持与 dB 读数。

## 最小示例

```go
level := float32(0.4) // the input's peak, 0 to 1 of full scale
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.MicLevelMeter(c, level, media.MicLevelMeterOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func MicLevelMeter(c *ui.Context, level float32, opts MicLevelMeterOptions) ui.Element`

`level` 为满刻度的 0 到 1，超出 panic。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Hold` | `time.Duration` | 1.5 秒 | 峰值保持时长。 |
| `Label` | `string` | "Microphone level" | 无障碍名称；描述中读出当前与峰值 dB。 |

## 状态

24 段：前 70% 为成功色，接着警告色，最后 10% 为危险色；峰值段为空心标记；右侧显示 dB 数值，不只靠颜色。

## 事件

无。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| — | 不获取焦点 |

## 限制

- 保持结束时只请求一次重绘，空闲时不持续重绘。
