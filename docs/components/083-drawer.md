# 083 Drawer 抽屉

## 用途

从窗口四边之一滑出的面板，用于筛选、详情和设置。可选模态或非模态。

## 最小示例

```go
open := true
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	if input.Button(c, "Filters", input.ButtonOptions{}).Clicked() {
		open = true
	}
	overlay.Drawer(c, &open, overlay.DrawerOptions{Title: "Filters", Side: overlay.DrawerRight, Size: 320}, func() {
		ui.Text(c, "Guests are placed by rank.")
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func Drawer(c *ui.Context, open *bool, opts DrawerOptions, content func()) ui.Element`。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Title` | `string` | 必填 | 标题与无障碍名称；为空时 panic。 |
| `Side` | `DrawerSide` | `DrawerRight` | `DrawerRight`、`DrawerLeft`、`DrawerTop`、`DrawerBottom`；越界时 panic。 |
| `Size` | `float32` | 左右 360，上下 280 | 宽度或高度；限制在 160 DIP 与窗口的 90% 之间；负数 panic。 |
| `NonModal` | `bool` | `false` | 非模态：不加遮罩、不限制焦点，窗口其余部分仍可操作。 |
| `Actions` | `func()` | `nil` | 页脚操作区。 |

## 状态

- 打开时从所在边滑入并淡入 160 ms；减少动画时直接出现。
- 面板贴边，只在朝内一侧画 1 DIP 分隔线。
- 内容超出时在内容区滚动。

## 事件

关闭按钮与 Escape 把 `*open` 设为 false；模态抽屉点击遮罩同样关闭。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| Tab / Shift+Tab | 模态：在抽屉内循环。非模态：按窗口焦点顺序，抽屉内容排在窗口内容之后。 |
| Escape | 关闭（先交给正在编辑的控件和更上层浮层）。 |

## 限制

- 输入法组合期间按 Escape 仍会关闭浮层：MyGo 没有公开组合状态（`ui/editor.go:73` 的 `compose` 未导出，组合事件不交给带编辑器元素的 `HandleInput`，`ui/input.go:689`），编辑器也不接收 Escape（`ui/editor.go:163`），macOS 上按键先于输入法送达应用（`internal/darwin/surface.go:704`）。
- Escape 只关闭最上层的 MujicaUI 浮层（窗口内按打开顺序记录）。MyGo 自带的 `ui.Popover`、`ui.Select` 弹层不在这份记录中；它们打开在非模态 MujicaUI 浮层之上时，Escape 可能先关闭下面的 MujicaUI 浮层。
- 关闭后焦点优先由 MyGo 还给打开者；打开者已消失、焦点落空时，焦点停在组件调用处的不可见占位元素上，Tab 从这里继续。不能回到打开者的父级：MyGo 不公开当前焦点或打开者（`ui/scope.go:123`、`ui/scope.go:137`）。
- 非模态抽屉打开时不移动焦点，其内容排在窗口焦点顺序末尾。
