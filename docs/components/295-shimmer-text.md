# 295 ShimmerText 文字扫光

## 用途

处理中的状态文字（如“思考中”）上有光扫过。

## 最小示例

```go
thinking := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	feedback.ShimmerText(c, "Thinking…", feedback.ShimmerTextOptions{Active: thinking})
}
ui.Render(view, 400, 300, 1)
```

签名：`func ShimmerText(c *ui.Context, text string, opts ShimmerTextOptions) ui.Element`

## 参数

`ShimmerTextOptions`

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| Active | bool | false | 扫光；否则为静态 TextMuted 文字。 |

text 为空时 panic。

## 状态

- 逐字由 TextMuted 渐亮到 Text，再回落；只在绘制时动画。
- 以状态角色朗读文字。
- 减少动画或未激活：静态文字。

## 事件

- 无。

## 键盘操作

- 无。

## 限制

- 扫光时文字保持单行。
