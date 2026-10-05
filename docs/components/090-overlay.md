# 090 Overlay 覆盖层

## 用途

在窗口内显示一层自定义内容：九个位置、可选遮罩、模态或非模态。模态时复用 MyGo `DialogBase` 的焦点限制与背景屏蔽，可与 MyGo 自带的 Select、Popover、菜单嵌套使用。

## 最小示例

```go
open := true
view := func(c *ui.Context) {
	mujica.Use(c, mujica.Settings{})
	mujica.Overlay(c, &open, mujica.OverlayOptions{Place: mujica.OverlayBottomRight, Label: "Now playing"}, func() {
		ui.Text(c, "Symbol I: △")
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func Overlay(c *ui.Context, open *bool, opts OverlayOptions, content func()) *ui.Element`。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Place` | `OverlayPlace` | `OverlayCenter` | `OverlayCenter`、`Top`、`Bottom`、`Left`、`Right`、`TopLeft`、`TopRight`、`BottomLeft`、`BottomRight`；越界 panic。 |
| `Inset` | `float32` | 16 | 面板与窗口边缘的距离。 |
| `Modal` | `bool` | `false` | 限制焦点、屏蔽背景，点击背景关闭。 |
| `Dim` | `bool` | `false` | 模态时加暗色遮罩；未开 `Modal` 时 panic。 |
| `Label` | `string` | `""` | 覆盖层的无障碍名称。 |

## 状态

弹层外观（Surface 底、1 DIP 边框、6 DIP 圆角）与 160 ms 淡入；减少动画时直接显示。

## 事件

Escape 与（模态时）背景点击把 `*open` 设为 false。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| Escape | 关闭；先交给编辑中的控件和更上层浮层。 |
| Tab | 模态时在覆盖层内循环；非模态时覆盖层元素排在窗口焦点顺序末尾。 |

## 限制

- 非模态覆盖层的 Escape 是窗口级快捷键（MyGo 浮层专用快捷键 `overlayShortcut` 不公开，`ui/scope.go:109`），同时打开多个非模态浮层时由先构建的那个响应。
- 非模态覆盖层关闭时不自动恢复焦点（`ui/scope.go:123` openOverlay 不公开）。
- 非模态面板首帧尚未测量尺寸，该帧透明，下一帧定位后淡入。
