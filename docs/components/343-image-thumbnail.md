# 343 ImageThumbnail 缩略图

## 用途

固定尺寸、按比例裁切的可点击缩略图，带选中、加载中和加载失败状态。

## 最小示例

```go
photo := ui.NewBitmap(examplePhoto)
selected := false
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if media.ImageThumbnail(c, photo, media.ImageThumbnailOptions{Alt: "Grand hall", Aspect: 4.0 / 3, Selected: selected}).Changed() {
		selected = !selected
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ImageThumbnail(c *ui.Context, src ui.ImageSource, opts ImageThumbnailOptions) ChoiceResult`

`src` 为 nil 时显示占位图标。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Width` | `float32` | 0（即 96） | 宽度。 |
| `Aspect` | `float32` | 0（即 1） | 宽除以高，图片按 Cover 裁切填满。 |
| `Selected` | `bool` | false | 选中：2 DIP 强调边框加右上角菱形。 |
| `Loading` | `bool` | false | 显示加载标记，尺寸不变。 |
| `Err` | `error` | nil | 显示失败图标；错误文字作为提示与描述。 |
| `Alt` | `string` | 必填 | 替代说明，为空时 panic。 |
| `Disabled` | `bool` | false | 禁用点击。 |

## 状态

默认、悬停（边框加深）、聚焦、选中、加载、失败、禁用，所有状态尺寸相同。

## 事件

`Changed()` 表示本帧被点击，或聚焦时按了 Enter / Space。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 聚焦缩略图 |
| Enter / Space | 激活 |

## 限制

- 选中状态由调用方保存；组件只报告点击。
