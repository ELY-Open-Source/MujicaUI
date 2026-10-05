# 360 ImageCompare 前后对比

## 用途

前后两张图片叠放，用滑块分隔对比。

## 最小示例

```go
before, after := ui.NewBitmap(examplePhoto), ui.NewBitmap(examplePhoto)
split := 0.5
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	media.ImageCompare(c, before, after, &split, media.ImageCompareOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ImageCompare(c *ui.Context, before, after ui.ImageSource, split *float64, opts ImageCompareOptions) ChoiceResult`

`*split` 取 0 到 1，表示分隔线位置；超出 panic。两张图都须为 `*ui.Bitmap` 或 `*ui.SVG`，盒子按 `before` 的宽高比。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `BeforeLabel / AfterLabel` | `string` | "Before" / "After" | 两侧标签，也组成无障碍名称。 |
| `Disabled` | `bool` | false | 禁用。 |

## 状态

分隔线与菱形手柄；两侧标签始终显示，不只靠画面区分。

## 事件

`Changed()` 在分隔位置变化后为 true。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| ← / → | 移动 1% |
| Home / End | 最左 / 最右 |
| 拖动 | 任意位置拖动分隔线 |

## 限制

- 两张图都按 Cover 铺满盒子，尺寸不同时会裁切。
