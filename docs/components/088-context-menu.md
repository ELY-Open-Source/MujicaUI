# 088 ContextMenu 右键菜单

## 用途

为一块区域提供应用内右键菜单，由 MujicaUI 自绘。右键在指针处打开，键盘在焦点区域下方打开，并避开窗口边缘。菜单项与 DropdownMenu 相同（`PopupMenu`）。

## 最小示例

```go
view := func(c *ui.Context) {
	core.Use(c, core.Settings{})
	row := ui.Row(c).Padding(12).Children(func() {
		ui.Text(c, "Letter to the Countess")
	})
	overlay.ContextMenu(c, row, overlay.ContextMenuOptions{}, func(m *overlay.PopupMenu) {
		if m.Item("Reply", overlay.PopupMenuItemOptions{}) {
			fmt.Println("reply")
		}
		m.Item("Delete", overlay.PopupMenuItemOptions{})
	})
}
ui.Render(view, 400, 300, 1)
```

签名：`func ContextMenu(c *ui.Context, target ui.Element, opts ContextMenuOptions, build func(m *PopupMenu)) ui.Element`，返回 `target`。`target` 或 `build` 为 nil 时 panic。

## 参数

| 字段 | 类型 | 默认值 | 含义 |
| --- | --- | --- | --- |
| `Disabled` | `bool` | `false` | 不响应右键与菜单键，并关闭已打开的菜单。 |

`PopupMenu` 的方法与选项见 [087 DropdownMenu](087-dropdown-menu.md)。

## 状态

与 DropdownMenu 菜单面板相同。右下角附近打开时，菜单向左、向上翻转，保持在窗口内。

## 事件

`Item` / `Check` 的返回值；每次选择只报告一次。

## 键盘操作

| 按键 | 行为 |
| --- | --- |
| Menu 键 / Shift+F10 | 焦点在 `target` 内时，在 `target` 左下方打开，焦点在第一项。 |
| 菜单内按键 | 同 DropdownMenu。 |
| Escape | 关闭，焦点回到打开前的控件。 |

## 限制

- `target` 必须是容器：菜单会向它追加一个不可见的定位子元素，文字元素不能有非文字子元素（MyGo 会 panic）。
- 菜单通过 `target.HandleInput` 接收右键按下，因此会替换调用方在同一元素上设置的 `HandleInput`（`ui/handler.go:73`）。
- 不使用 MyGo 的 `Element.ContextMenu`：那是系统菜单，PRD 不允许把它当作哥特样式菜单。
- 打开期间外部点击只用于关闭，不传给下方控件：`PopoverBase` 用覆盖整窗的背板接收点击（`ui/base.go:373`）。
- Escape 只关闭最上层的 MujicaUI 浮层（窗口内按打开顺序记录）。MyGo 自带的 `ui.Popover`、`ui.Select` 弹层不在这份记录中；它们打开在非模态 MujicaUI 浮层之上时，Escape 可能先关闭下面的 MujicaUI 浮层。
