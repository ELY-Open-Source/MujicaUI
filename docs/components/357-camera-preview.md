# 357 CameraPreview 摄像头预览

## 用途

显示调用方提供的摄像头帧，无帧时显示“无信号”，可切换镜像。

## 最小示例

```go
frame := mujica.CameraFrame{Image: examplePhoto, Seq: 1} // a new Seq for each captured frame
mirror := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.CameraPreview(c, frame, &mirror, mujica.CameraPreviewOptions{})
}
ui.Render(view, 400, 300, 1)
```

## 参数

`func CameraPreview(c *ui.Context, frame CameraFrame, mirror *bool, opts CameraPreviewOptions) ChoiceResult`

`CameraFrame`：`Image`（标准库 `image.Image`，nil 为无信号）、`Seq`（帧序号，复用同一图像缓冲时用新序号触发重绘）。

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Label` | `string` | "Camera" | 无障碍名称。 |
| `Disabled` | `bool` | false | 禁用镜像切换。 |

## 状态

16:9 黑底；镜像按钮开启时有选中底色和下划线。

## 事件

`Changed()` 在镜像切换后为 true。

## 键盘操作

| 按键 | 作用 |
| --- | --- |
| Tab | 聚焦预览 |
| M | 切换镜像 |

## 限制

- MyGo 的 Painter 不能水平翻转位图（ui/paint.go:834），所以镜像时组件每个新帧复制一份翻转后的像素；高分辨率帧会占用 CPU。
- 组件不打开摄像头。
