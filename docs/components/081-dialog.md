# 081 Dialog 对话框

## 用途

在窗口上方显示模态对话框：标题、说明、内容与操作区。打开时限制焦点并屏蔽背景，关闭后把焦点还给打开它的控件。视觉回归基准组件之一。

## 最小示例

```go
open := true
name := "Ophelia"
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if input.Button(c, "Rename…", input.ButtonOptions{}).Clicked() {
		open = true
	}
	overlay.Dialog(c, &open, overlay.DialogOptions{
		Title:       "Rename the chamber",
		Description: "The new name shows in every hall.",
		Decorated:   true,
		Actions: func() {
			if input.Button(c, "Cancel", input.ButtonOptions{Variant: input.Secondary}).Clicked() {
				open = false
			}
			if input.Button(c, "Rename", input.ButtonOptions{}).Clicked() {
				open = false
			}
		},
	}, func() {
		ui.TextInput(c, &name)
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func Dialog(c *ui.Context, open *bool, opts DialogOptions, content func()) ui.Element`，关闭时返回 nil。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 必填 | 衬线标题，也是对话框的无障碍名称；为空时 panic。 |
| `Description` | `string` | `""` | 标题下的说明文字，同时作为无障碍描述。 |
| `Width` | `float32` | `440` | 面板宽度（DIP），始终不超过窗口宽度减 32。 |
| `Decorated` | `bool` | `false` | 开启双线边框、四角菱形角饰和标题下的装饰线。 |
| `Actions` | `func()` | `nil` | 页脚操作区，靠右排列，可换行。 |
| `KeepOnBackdrop` | `bool` | `false` | 点击遮罩不关闭。 |
| `KeepOnEscape` | `bool` | `false` | Escape 不关闭。 |
| `HideClose` | `bool` | `false` | 不显示右上角关闭按钮。 |

`open` 为 false 时不构建任何元素。`content` 可为 nil；说明与内容一起滚动，标题与操作区始终留在窗口内（含 200% 文字缩放）。

## 状态

- 打开：遮罩渐入，面板 160 ms 淡入并上移 8 DIP；减少动画时直接显示。
- 关闭按钮：悬停、按下、键盘焦点环（2 DIP，留 2 DIP 间隔）。
- 装饰：只出现在边缘与标题区，不占焦点，不影响点击区域。

## 事件

- 关闭按钮、遮罩点击、Escape 把 `*open` 设为 false（各自可关闭）。
- 业务事件由 `Actions` 和 `content` 中的控件自己报告。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| Tab / Shift+Tab | 在对话框内循环，不会离开。 |
| Escape | 先交给输入法与正在编辑的控件，再关闭最上层浮层；嵌套 Popover 或菜单先关。 |
| Enter / Space | 触发聚焦的按钮。 |

## 限制

- 输入法组合期间按 Escape 仍会关闭浮层：MyGo 没有公开组合状态（`ui/editor.go:73` 的 `compose` 未导出，组合事件不交给带编辑器元素的 `HandleInput`，`ui/input.go:689`），编辑器也不接收 Escape（`ui/editor.go:163`），macOS 上按键先于输入法送达应用（`internal/darwin/surface.go:704`）。
- 关闭后焦点优先由 MyGo 还给打开者；打开者已消失、焦点落空时，焦点停在组件调用处的不可见占位元素上，Tab 从这里继续。不能回到打开者的父级：MyGo 不公开当前焦点或打开者（`ui/scope.go:123`、`ui/scope.go:137`）。
- 没有关闭过渡：`*open` 变为 false 的下一帧浮层元素即被移除。
