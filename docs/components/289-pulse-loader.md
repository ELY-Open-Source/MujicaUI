# 289 PulseLoader 呼吸加载

## 用途

单个圆点在 60% 到 100% 之间缓慢呼吸的加载标记，强度受限，不发光。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.PulseLoader(c, feedback.LoaderOptions{Label: "Listening", ShowLabel: true})
}
ui.Render(view, 400, 300, 1)
```

签名：`func PulseLoader(c *ui.Context, opts LoaderOptions) ui.Element`

## 参数

`LoaderOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Size | float32 | 0（即 16） | 标记高度。 |
| Label | string | 内置“加载中” | 朗读文字。 |
| ShowLabel | bool | false | 显示文字。 |
| Color | ui.Color | AccentText | 颜色。 |

见 DotsLoader。

## 状态

- 透明度 0.55–1、尺寸 0.6–1，无光晕。
- 减少动画：静止并显示文字。

## 事件

- 无。

## 键盘操作

- 无。

## 限制

- 无。
