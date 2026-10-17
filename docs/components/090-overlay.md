# 090 Overlay 覆盖层

## 用途

在窗口内显示一层自定义内容：九个位置、可选遮罩、模态或非模态。模态时复用 MyGo `DialogBase` 的焦点限制与背景屏蔽，可与 MyGo 自带的 Select、Popover、菜单嵌套使用。

## 最小示例

```go
open := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	overlay.Overlay(c, &open, overlay.OverlayOptions{Place: overlay.OverlayBottomRight, Label: "Now playing"}, func() {
		ui.Text(c, "Symbol I: △")
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func Overlay(c *ui.Context, open *bool, opts OverlayOptions, content func()) ui.Element`。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Place` | `OverlayPlace` | `OverlayCenter` | `OverlayCenter`、`Top`、`Bottom`、`Left`、`Right`、`TopLeft`、`TopRight`、`BottomLeft`、`BottomRight`；越界 panic。 |
| `Inset` | `float32` | 16 | 面板与窗口边缘的距离。 |
| `Modal` | `bool` | `false` | 限制焦点、屏蔽背景，点击背景关闭。 |
| `Dim` | `bool` | `false` | 模态时加暗色遮罩；未开 `Modal` 时 panic。 |
| `Label` | `string` | `""` | 覆盖层的无障碍名称。 |

## 状态

弹层外观（Surface 底、1 DIP 边框、6 DIP 圆角）与 160 ms 淡入；减少动画时直接显示。面板不超过窗口减去 `Inset`，内容超出时滚动。

## 事件

Escape 与（模态时）背景点击把 `*open` 设为 false。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| Escape | 关闭；先交给编辑中的控件和更上层浮层。 |
| Tab | 模态时在覆盖层内循环；非模态时覆盖层元素排在窗口焦点顺序末尾。 |

## 限制

- 输入法组合期间按 Escape 仍会关闭浮层：MyGo 没有公开组合状态（`ui/editor.go:73` 的 `compose` 未导出，组合事件不交给带编辑器元素的 `HandleInput`，`ui/input.go:689`），编辑器也不接收 Escape（`ui/editor.go:163`），macOS 上按键先于输入法送达应用（`internal/darwin/surface.go:704`）。
- Escape 只关闭最上层的 MujicaUI 浮层（窗口内按打开顺序记录）。MyGo 自带的 `ui.Popover`、`ui.Select` 弹层不在这份记录中；它们打开在非模态 MujicaUI 浮层之上时，Escape 可能先关闭下面的 MujicaUI 浮层。
- 关闭后焦点优先由 MyGo 还给打开者；打开者已消失、焦点落空时，焦点停在组件调用处的不可见占位元素上，Tab 从这里继续。不能回到打开者的父级：MyGo 不公开当前焦点或打开者（`ui/scope.go:123`、`ui/scope.go:137`）。
- 非模态面板首帧尚未测量尺寸，该帧透明，下一帧定位后淡入。
