# 086 HoverCard 悬停卡片

## 用途

悬停或聚焦锚点时预览更多信息，例如人物资料。指针可以移入卡片并操作其中的控件，键盘可以 Tab 进入卡片，离开后延迟关闭。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	b := input.Button(c, "Lady Sakiko", input.ButtonOptions{Variant: input.Ghost})
	overlay.HoverCard(c, b, overlay.HoverCardOptions{}, func() {
		ui.Text(c, "Keyboardist of the court.")
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func HoverCard(c *ui.Context, anchor ui.Element, opts HoverCardOptions, content func()) ui.Element`，返回卡片，未显示时为 nil。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `OpenDelay` | `time.Duration` | 500 ms | 悬停或聚焦多久后打开；负数 panic。 |
| `CloseDelay` | `time.Duration` | 300 ms | 指针与焦点都离开后多久关闭；负数 panic。 |
| `Width` | `float32` | 280 | 卡片宽度。 |
| `Label` | `string` | `""` | 卡片的无障碍名称。 |

## 状态

- 打开、关闭都按延迟进行；减少动画时取消淡入。
- 指针在锚点或卡片上、或键盘焦点在锚点或卡片内时保持打开。
- 下方空间不足时翻到上方，左右限制在窗口内。

## 事件

返回值由 nil 变为非 nil 即打开；卡片内控件自行报告事件。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| Tab 聚焦锚点 | `OpenDelay` 后打开。 |
| Tab | 从锚点进入卡片内的控件，再继续到后面的元素。 |
| Escape | 立即关闭，焦点回到锚点，直到焦点离开前不再打开。 |

## 限制

- 卡片由 `PopoverBase` 承载，使焦点顺序紧跟锚点，在 Dialog 内也属于对话框的焦点范围；代价是打开期间外部点击只用于关闭，不传给下方控件：`PopoverBase` 用覆盖整窗的背板接收点击（`ui/base.go:373`）。
- 指针在卡片外但焦点仍在卡片内时卡片保持打开，点击外部或按 Escape 关闭。
- Escape 只关闭最上层的 MujicaUI 浮层（窗口内按打开顺序记录）。MyGo 自带的 `ui.Popover`、`ui.Select` 弹层不在这份记录中；它们打开在非模态 MujicaUI 浮层之上时，Escape 可能先关闭下面的 MujicaUI 浮层。
