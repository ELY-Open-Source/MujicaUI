# 288 BarsLoader 条形加载

## 用途

四根竖条依次升降的加载标记。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.BarsLoader(c, feedback.LoaderOptions{Size: 20, Label: "Indexing"})
}
ui.Render(view, 400, 300, 1)
```

签名：`func BarsLoader(c *ui.Context, opts LoaderOptions) ui.Element`

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

- 条高吸附像素；只在绘制时动画。
- 减少动画：静止并显示文字。

## 事件

- 无。

## 键盘操作

- 无。

## 限制

- 无。
