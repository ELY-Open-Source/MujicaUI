# 291 WaveLoader 波形加载

## 用途

一条向右流动的正弦线加载标记。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.WaveLoader(c, feedback.LoaderOptions{Label: "Speaking", ShowLabel: true})
}
ui.Render(view, 400, 300, 1)
```

签名：`func WaveLoader(c *ui.Context, opts LoaderOptions) ui.Element`

## 参数

`LoaderOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Size | float32 | 0（即 16） | 标记高度，宽为 2.5 倍。 |
| Label | string | 内置“加载中” | 朗读文字。 |
| ShowLabel | bool | false | 显示文字。 |
| Color | ui.Color | AccentText | 颜色。 |

见 DotsLoader。

## 状态

- 路径点吸附像素。
- 减少动画：静止的波形并显示文字。

## 事件

- 无。

## 键盘操作

- 无。

## 限制

- 无。
