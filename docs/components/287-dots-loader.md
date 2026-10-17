# 287 DotsLoader 三点加载

## 用途

三个圆点依次变亮的加载标记，只在加载期间显示。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.DotsLoader(c, feedback.LoaderOptions{Label: "Thinking", ShowLabel: true})
}
ui.Render(view, 400, 300, 1)
```

签名：`func DotsLoader(c *ui.Context, opts LoaderOptions) ui.Element`

## 参数

`LoaderOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Size | float32 | 0（即 16） | 标记高度；为负或无穷 panic。 |
| Label | string | 内置“加载中” | 朗读文字。 |
| ShowLabel | bool | false | 在标记旁显示文字。 |
| Color | ui.Color | AccentText | 标记颜色。 |

五个加载器共用 `LoaderOptions`。

## 状态

- 只在绘制时请求下一帧：不再构建或被隐藏后停止重绘。
- 减少动画：标记静止，并始终显示文字说明。
- 无弹跳、无发光。

## 事件

- 无。

## 键盘操作

- 无：不获取焦点，以进度角色朗读。

## 限制

- 无。
