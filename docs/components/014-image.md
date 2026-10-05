# 014 Image 图片

## 用途

显示位图（`*ui.Bitmap`）或 SVG（`*ui.SVG`），可选包含或裁切填充；提供占位、加载中、加载失败三种状态和替代说明。

## 最小示例

```go
var photo *ui.Bitmap // decode once with ui.DecodeBitmap, outside the view
loading, loadErr := false, errors.New("timeout")
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	img := display.Image(c, photo, display.ImageOptions{
		Width: 160, Height: 120, Crop: true, Alt: "Hall of mirrors",
		Loading: loading, Err: loadErr,
	})
	if img.Retry {
		loading, loadErr = true, nil // start loading again here
	}
}
```

## 参数

`func Image(c *ui.Context, src ui.ImageSource, opts ImageOptions) ImageView`

`src` 为 `nil`（未赋值的接口）时显示占位。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Width`, `Height` | `float32` | `0` | 盒子尺寸（DIP）。为 0 时使用 `src` 的尺寸（只给一边时按比例补齐）；`src` 为 nil 时为 160×120。占位、加载、错误状态与已加载图片使用相同的宽高比和伸缩规则，父容器变窄时一起按比例缩小，所以各状态的盒子完全相同。负数 panic。 |
| `Crop` | `bool` | `false` | `true` 覆盖并裁切（`ui.Cover`）；`false` 完整放入（`ui.Contain`）。 |
| `Alt` | `string` | `""` | 替代说明，供辅助技术读取；在其他状态下作为描述。 |
| `Loading` | `bool` | `false` | 调用方正在加载。 |
| `Err` | `error` | `nil` | 调用方加载失败；优先于 `Loading`。 |

返回值 `ImageView`：`Element` 为外层元素；`Retry` 表示本帧点击了重试。

## 状态

优先级：错误 > 加载中 > 占位 > 图片。

- 占位：浅色框内显示图片图标，名称“暂无图片”。
- 加载中：框内显示 Spinner，名称“正在加载图片”。
- 各状态共用同一个外层盒子：图片元素始终参与布局（非加载状态下隐藏），所以在默认拉伸的 Column、Row 和固定尺寸的父容器里，盒子大小完全相同。
- 错误：盒子中央是 28×28 DIP 的仅图标“重试”按钮（`refresh-cw`），`Err.Error()` 作为按钮的悬停提示和描述，也是盒子的无障碍错误说明。盒子足够宽时，按钮旁边还会显示“图片加载失败”文字。

## 事件

- `ImageView.Retry`：用户点击了重试，或在按钮上按 Enter/Space。

## 键盘操作

只有错误状态的“重试”按钮可以聚焦（Tab），用 Enter 或 Space 触发。

## 限制

- “重试”按钮画在不裁切的图层上，以盒子为中心，始终保持 28×28 DIP 的点击区域和完整的焦点环。盒子矮于 28 DIP 时，按钮会超出盒子的上下边缘，可能盖住相邻内容；需要按钮完全落在盒子内时，请让盒子至少 28 DIP 高。
- “图片加载失败”文字只在盒子宽度足够时显示。宽度按上一帧的布局判断，盒子尺寸变化时组件会多布局一帧来修正。

- 组件不读取文件，也不访问网络。解码（`ui.DecodeBitmap`、`ui.ParseSVG`）和异步加载由调用方完成，结果通过 `Window.Update` 回到界面线程，再设置 `Loading` 与 `Err`。
- SVG 的绘制范围以 MyGo 支持的部分为准：不支持文本、嵌入图片、pattern、marker 和 filter。
- 接口里装着 nil 指针（如 `(*ui.Bitmap)(nil)`）时不算 `nil`，不会显示占位。
