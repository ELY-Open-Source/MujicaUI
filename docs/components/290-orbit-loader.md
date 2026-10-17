# 290 OrbitLoader 环绕加载

## 用途

一枚菱形沿一圈小菱形环绕，与 Spinner 同一视觉语言。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.OrbitLoader(c, feedback.LoaderOptions{Size: 24, Label: "Connecting"})
}
ui.Render(view, 400, 300, 1)
```

签名：`func OrbitLoader(c *ui.Context, opts LoaderOptions) ui.Element`

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

- 移动的菱形坐标吸附像素。
- 减少动画：菱形停在顶端并显示文字。

## 事件

- 无。

## 键盘操作

- 无。

## 限制

- 无。
