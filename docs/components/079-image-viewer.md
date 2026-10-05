# 079 ImageViewer · 图片查看器

## 用途

在窗口上方的模态层中查看图片：放大、缩小、平移、适应窗口、上一张/下一张，以及键盘操作与关闭。

## 最小示例

```go
picture := ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="240" height="160"><rect width="240" height="160" fill="#AD9164"/><path d="M120 20 L200 80 L120 140 L40 80 Z" fill="#7D2034"/></svg>`))
open, index := true, 0
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.ImageViewer(c, &open, &index, []mujica.ViewerImage{{Source: picture, Alt: "Rose window"}})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ImageViewer(c *ui.Context, open *bool, index *int, images []ViewerImage) ImageViewerView`。`ViewerImage` 字段：`Source`（`*ui.Bitmap` 或 `*ui.SVG`，必填）、`Alt`（说明文字与无障碍名称）。关闭时返回的 `Element` 为 nil。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `open` | `*bool` | 必填 | 为 true 时显示；Escape、点击遮罩或关闭按钮设为 false。 |
| `index` | `*int` | 必填 | 当前图片下标；越界 panic。 |
| `images` | `[]ViewerImage` | 必填 | 至少一张。 |

## 状态

缩放 25%–800%（相对适应窗口），按钮在边界处禁用；适应窗口时“适应”按钮禁用；首张/末张时上一张/下一张禁用。

## 事件

`ImageViewerView.Changed()`：切换了图片，`*index` 已更新。关闭通过 `*open` 变为 false 观察。

## 键盘操作

- + / -：放大、缩小；0：适应窗口。
- ←/→：上一张、下一张；Shift+方向键：平移。
- 滚轮：缩放；拖动图片：平移。
- Escape：关闭；Tab 在工具按钮与图片之间移动。

## 限制

- 基于 MyGo `ui.DialogBase`，焦点限制与关闭后的焦点恢复由 MyGo 处理。
- 图片解码由调用方完成（`ui.DecodeBitmap`），组件不读取文件或网络。
- 不支持旋转与触控板捏合手势。
