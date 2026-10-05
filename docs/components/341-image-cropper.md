# 341 ImageCropper 图片裁剪

## 用途

在图片上拖出裁剪框，可锁定宽高比、按 90 度旋转，输出裁剪矩形。

## 最小示例

```go
crop := mujica.ImageCrop{} // empty: starts fitted to the image
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	if mujica.ImageCropper(c, examplePhoto, &crop, mujica.ImageCropperOptions{Aspect: 1}).Changed() {
		_ = crop.Rect // in pixels of the photo turned by crop.Rotation
	}
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func ImageCropper(c *ui.Context, src image.Image, crop *ImageCrop, opts ImageCropperOptions) ChoiceResult`

`src` 必填，是标准库 `image.Image`（组件要自己转动像素）。`ImageCrop.Rect` 是旋转 `Rotation` 度之后那张图上的像素矩形；`Rotation` 只能是 0、90、180、270。`Rect` 为空时自动设为最大的居中框。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Aspect` | `float64` | 0 | 锁定宽除以高；0 为自由比例。`Rect` 不符合比例时 panic。 |
| `Height` | `float32` | 0（即 320） | 舞台高度。 |
| `Disabled` | `bool` | false | 禁用拖动、按键与旋转。 |
| `Label` | `string` | "Crop area" | 舞台的无障碍名称。 |

## 状态

- 默认：框外区域压暗，框线与三分线为象牙色，四角有菱形手柄。
- 聚焦：舞台显示 2 DIP 焦点环；下方读数同时作为无障碍描述。
- 禁用：MyGo 统一调暗，不响应输入。

## 事件

`Changed()` 在框移动、缩放或旋转后为 true。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| 方向键 | 移动裁剪框，每次图片宽度的 1% |
| Shift+方向键 | 右、下放大，左、上缩小；锁定比例时保持比例 |
| R / Shift+R | 向右 / 向左旋转 90 度 |

## 限制

- MyGo 的 Painter 不能旋转位图（只有单色图标能转，ui/svg.go:98），所以组件在图片或角度变化时复制一份转好的像素；超大图片会在那一帧花时间。
- 锁定比例时旋转会重新取最大的居中框，原选区不保留。
- 拖动角点时，按下点距角 14 DIP 以内才算抓住手柄。
